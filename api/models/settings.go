package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Settings struct {
	Id                 primitive.ObjectID `bson:"_id,omitempty"`
	DeviceId           string             `bson:"device_id"`
	DataRate           int                `bson:"data_rate" copier:"DataRate"`
	NotificationEmails []string           `bson:"email" copier:"Email"`
}
