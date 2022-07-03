package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Account struct {
	Id       primitive.ObjectID `bson:"_id"`
	Mobile   string
	Password string
	Status   uint8
	RealName string
	IdNum    string
}
