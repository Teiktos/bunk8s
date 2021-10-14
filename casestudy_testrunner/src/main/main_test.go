package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"showcase/model"
	"strings"
	"testing"
	"time"

	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestDeviceServiceReceiveDeviceEvent(t *testing.T) {

	respGet, err := http.DefaultClient.Get("http://device.bunk8s-fe.svc.cluster.local:8080/device")

	if err != nil {
		assert.Nil(t, err)
		t.Fatalf("Error getting devices: %v", err)
	}
	defer respGet.Body.Close()

	body, err := io.ReadAll(respGet.Body)
	if err != nil {
		t.Fatalf("Error reading body: %v", err)
	}

	var availableDevices []model.Device

	err = json.Unmarshal(body, &availableDevices)
	if err != nil {
		t.Fatalf("Error unmarshalling body: %v", err)
	}

	deviceEvent := model.DeviceEvent{DeviceId: availableDevices[0].DeviceId, BadgeNumber: 1337}

	deviceEventJson, err := json.Marshal(deviceEvent)
	if err != nil {
		t.Fatalf("Error marshalling deviceEVent to JSON: %v", err)
	}

	req, err := http.NewRequest("POST", "http://device.bunk8s-fe.svc.cluster.local:8080/device/event", strings.NewReader(string(deviceEventJson)))

	if err != nil {
		t.Fatalf("Error creating new POST request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Acceppt", "application/json")

	_, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Error doing POST request: %v", err)
	}

}

func TestRoomServiceRoomEventResponse(t *testing.T) {

	connRabbit, err := amqp.Dial("amqp://rabbitmq.bunk8s-fe.svc.cluster.local:5672")

	if err != nil {
		t.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer connRabbit.Close()

	channelRabbit, err := connRabbit.Channel()
	if err != nil {
		t.Fatalf("Failed to create RabbitMQ channel: %v", err)
	}
	defer channelRabbit.Close()

	t.Run("Queue Dashboard", func(t *testing.T) {

		messagesDashboard, err := channelRabbit.Consume(
			"dashboard.queue.room.event", // queue name
			"",                           // consumer
			true,                         // auto-ack
			false,                        // exclusive
			false,                        // no local
			false,                        // no wait
			nil,                          // arguments
		)
		if err != nil {
			t.Fatalf("Failed to create consumer for RabbitMQ dashboard queue: %v", err)
		}

		ok := func() bool {
			for message := range messagesDashboard {

				var roomEvent model.RoomEvent

				json.Unmarshal(message.Body, &roomEvent)

				return roomEvent.CardNo == 1337

				// log.Printf(" > Received message from dashboard queue: %s\n", message.Body)
				// return
			}
			return false
		}()

		assert.True(t, ok)

	})

	t.Run("Queue Shout", func(t *testing.T) {

		messagesShout, err := channelRabbit.Consume(
			"shout.queue.room.event", // queue name
			"",                       // consumer
			true,                     // auto-ack
			false,                    // exclusive
			false,                    // no local
			false,                    // no wait
			nil,                      // arguments
		)
		if err != nil {
			t.Fatalf("Failed to create consumer for RabbitMQ dashboard queue: %v", err)
		}

		ok := func() bool {
			for message := range messagesShout {

				var roomEvent model.RoomEvent

				json.Unmarshal(message.Body, &roomEvent)

				return roomEvent.CardNo == 1337
				// log.Printf(" > Received message from shout queue: %s\n", message.Body)
				// return
			}
			return false
		}()
		assert.True(t, ok)
	})

}

func TestMongoDBRoomAllocationEntry(t *testing.T) {

	ctx := context.Background()

	mongoClient, err := ConnectToMongo(ctx, "mongodb://mongo.bunk8s-fe.svc.cluster.local:27017")

	if err != nil {
		t.Fatalf("Failed to create MongoDB client: %v", err)
	}

	cursor, err := mongoClient.Database("room").Collection("roomAllocation").Find(ctx, bson.M{"_id": 1337})
	if err != nil {
		t.Fatalf("Failed to search for entry on MongoDB: %v", err)
	}
	defer cursor.Close(ctx)

	var roomAllocations []*model.MongoRoomAllocation

	cursor.All(ctx, &roomAllocations)

	assert.Equal(t, 1, len(roomAllocations))

}

func TestMain(m *testing.M) {

	m.Run()
	os.Exit(0)

}

func ConnectToMongo(ctx context.Context, mongodbHost string) (*mongo.Client, error) {
	timedContext, cancelTimedContext := context.WithTimeout(ctx, (10 * time.Second))
	defer cancelTimedContext()

	client, err := mongo.NewClient(options.Client().ApplyURI(mongodbHost))
	if err != nil {
		return nil, fmt.Errorf("Failed to create client: %w", err)
	}

	err = client.Connect(timedContext)
	if err != nil {
		return nil, fmt.Errorf("Failed to initialize client: %w", err)
	}

	err = client.Ping(timedContext, nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to ping MongoDB: %w", err)
	}

	return client, nil
}
