package cmn

import (
	"github.com/95eh/eg"
)

var (
	_GlobalConf *GlobalConfig
)

func GlobalConf() *GlobalConfig {
	return _GlobalConf
}

func LoadGlobalConf(files ...string) {
	eg.LoadConf(&_GlobalConf, eg.ConvertConfLocalPath(files...)...)
}

type GlobalConfig struct {
	App     AppConf
	MongoDB MongoDBConf `yaml:"mongoDB"`
	Redis   RedisConf
	Net     AddrConf
	Http    AddrConf
	Region  RegionConf
	Account AccountConf
}

type RegionConf struct {
	Enable    bool
	CheckSecs int64
}

type AccountConf struct {
	Enable bool
	JwtKey string
	Salt   string
	OpsPw  string
}
