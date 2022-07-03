package model

import (
	"github.com/95eh/eg"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Character struct {
	Id        primitive.ObjectID `bson:"_id"`
	Uid       string
	NickName  string
	Gender    uint8
	Figure    uint8
	IsNovice  bool
	SceneType eg.TScene
	SceneId   int64
}
