package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Account struct {
	Id       primitive.ObjectID `bson:"_id"`
	Uid      int64
	Mobile   string
	Password string
	Status   uint8
	RealName string `bson:"real_name"`
	IdNum    string `bson:"id_num"`
	Mask     int64
}
