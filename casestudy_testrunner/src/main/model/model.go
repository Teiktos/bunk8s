package model

type Device struct {
	DeviceId   int `json:"deviceId"`
	RoomNumber int `json:"roomNumber"`
}

type DeviceEvent struct {
	DeviceId    int `json:"deviceId"`
	BadgeNumber int `json:"badgeNumber"`
}

type RoomEvent struct {
	RoomNo   int    `json:"roomNo"`
	CardNo   int    `json:"cardNo"`
	Type     string `json:"type"`
	RoomName string `json:"roomName"`
}

type MongoRoomAllocation struct {
	Id      int       `bson:"_id"`
	Room    MongoRoom `bson:"room"`
	Version int       `bson:"version"`
	Class   string    `bson:"_class"`
}

type MongoRoom struct {
	Id       int    `bson:"_id"`
	RoomName string `bson:"roomName"`
}
