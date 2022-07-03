package cmn

import "github.com/95eh/eg"

type AppConf struct {
	Mode   string
	Name   string
	Ver    string
	Region eg.TRegion
}

type RedisConf struct {
	Addr     string
	Password string
	Db       int
}

type MongoDBConf struct {
	UserName string
	Password string
	Addr     string
	DbName   string
}

type AddrConf struct {
	Ip   string
	Port int
}
