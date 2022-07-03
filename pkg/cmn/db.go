package cmn

import (
	"context"
	"fmt"
	"github.com/95eh/eg"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var (
	_Mongo *mongo.Database
)

func Db() *mongo.Database {
	return _Mongo
}

func InitMongo(conf MongoDBConf) {
	uri := fmt.Sprintf("mongodb://%s:%s@%s", conf.UserName, conf.Password, conf.Addr)
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		eg.Fatal(eg.WrapErr(eg.EcDbErr, err))
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		eg.Fatal(eg.WrapErr(eg.EcDbErr, err))
	}
	_Mongo = client.Database(conf.DbName)
}
