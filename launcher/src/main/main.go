package main

import (
	"bunk8s/launcher/parser"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"bunk8s/launcher/pb"

	"crypto/tls"
	"crypto/x509"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {

	logLevel, _ := zerolog.ParseLevel("info")

	zerolog.SetGlobalLevel(logLevel)

	bunk8sConfig := parser.ParseConfig("/config.yaml")
	bunk8sLauncherConfig := bunk8sConfig.LauncherConfig
	bunk8sCoordinatorConfig := bunk8sConfig.CoordinatorConfig

	coordinatorConfigJson, err := json.Marshal(bunk8sCoordinatorConfig)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshall from coordinatorConfig struct to json")
	}
	log.Debug().Msg("Unmarshalled coordinatorConfig struct to Json")

	coordinatorConfigProto := pb.CoordinatorConfig{}

	err = json.Unmarshal(coordinatorConfigJson, &coordinatorConfigProto)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshall from coordinatorConfig json to proto")
	}
	log.Debug().Msgf("Unmarshalled Json to Proto: %+v\n ", &coordinatorConfigProto)

	var creds credentials.TransportCredentials

	if bunk8sLauncherConfig.CertFile == "" {

		creds = credentials.NewTLS(&tls.Config{InsecureSkipVerify: false})

	} else {

		pemServerCA, err := ioutil.ReadFile("/cert/" + bunk8sLauncherConfig.CertFile)
		if err != nil {
			log.Error().Err(err).Msg("Failed to read cert file")
		}

		certPool := x509.NewCertPool()

		ok := certPool.AppendCertsFromPEM(pemServerCA)
		if !ok {
			log.Error().Err(err).Msg("Failed to parse server certificate or adding them to certPool")
		}

		config := &tls.Config{
			RootCAs:            certPool,
			InsecureSkipVerify: false,
		}
		creds = credentials.NewTLS(config)
	}

	conn, err := grpc.Dial(bunk8sLauncherConfig.CoordinatorIp+":"+bunk8sLauncherConfig.CoordinatorPort, grpc.WithTransportCredentials(creds), grpc.WithBlock())
	if err != nil {
		log.Error().Err(err).Msg("Failed to Dial gRPC Server")
	}
	log.Debug().Msg("successfully dialed coordinator!")

	c := pb.NewTestRunClient(conn)

	serverReply, err := c.DeployTestRunner(context.Background(), &coordinatorConfigProto)
	if err != nil {
		log.Error().Err(err).Msg("Failed when calling DeployTestRunner")
	}
	log.Debug().Msgf("Server replied with: %+v\n ", serverReply)

	serverReplyJson, err := json.MarshalIndent(serverReply, "", "    ")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed marshalling server reply to json")
		os.Exit(1)
	}

	log.Debug().Msgf("Server replied: ", serverReplyJson)
	fmt.Println(string(serverReplyJson))

	for i := range serverReply.TestRunnerPods {
		switch serverReply.TestRunnerPods[i].ReturnCode {
		case 0:
			log.Debug().Msg("Test run finished successfully")
		case 1:
			log.Error().Msg("Namespace name invalid")
			os.Exit(1)
		case 2:
			log.Error().Msg("Test runner pod name invalid or already exists in given namespace")
			os.Exit(1)
		case 3:
			log.Error().Msg("Failed to create test runner pod")
			os.Exit(1)
		case 4:
			log.Error().Msg("Failed to watch test runner pod")
			os.Exit(1)
		case 5:
			log.Error().Msg("Test run timed out")
			os.Exit(1)
		default:
			log.Error().Msg("Unknown or empty return from server")
			os.Exit(1)
		}
	}

	conn.Close()
}
