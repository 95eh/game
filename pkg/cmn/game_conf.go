package cmn

import (
	"github.com/95eh/eg"
	"github.com/95eh/rpg"
)

var (
	_GameConf *GameConfig
)

func GameConf() *GameConfig {
	return _GameConf
}

func LoadGameConf(files ...string) {
	eg.LoadConf(&_GameConf, eg.ConvertConfLocalPath(files...)...)
}

type GameConfig struct {
	App         AppConf
	MongoDB     MongoDBConf `yaml:"mongoDB"`
	GlobalRedis RedisConf
	RegionRedis RedisConf
	Net         AddrConf
	Gate        GateConf
	Character   CharacterConf
	ScnWorld    ScnWorldConf
	Other       GameOther
}

type GateConf struct {
	Enable    bool
	Cap       int
	Tcp       AddrConf
	Http      AddrConf
	Websocket AddrConf
	Udp       AddrConf
}

type CharacterConf struct {
	Enable bool
}

type ScnWorldConf struct {
	Enable bool
	Tile   rpg.TileSceneConf
}

type GameOther struct {
	SensitiveFilterUrl string
}
