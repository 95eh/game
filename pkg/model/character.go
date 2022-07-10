package model

import (
	"github.com/95eh/eg"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Character struct {
	Id        primitive.ObjectID `bson:"_id"`
	Cid       int64
	Uid       int64
	NickName  string `bson:"nick_name"`
	Gender    uint8
	Figure    uint8
	IsNovice  bool      `bson:"is_novice"`
	SceneType eg.TScene `bson:"scene_type"`
	SceneId   int64     `bson:"scene_id"`
	PosX      float32
	PosY      float32
}
