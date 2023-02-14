package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type DeviceCredentials struct {
	Id       primitive.ObjectID `bson:"_id,omitempty"`
	DeviceId string             `bson:"device_id"`
	Password string             `bson:"password"`
}
