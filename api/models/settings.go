package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Settings struct {
	Id                 primitive.ObjectID `bson:"_id"`
	DeviceId           string             `bson:"device_id"`
	DataRate           int                `bson:"data_rate" copier:"data_rate"`
	NotificationEmails []string           `bson:"email"`
}
