package pkg

import (
	"flag"
	"fmt"
	"game/pkg/cmn"
	"game/pkg/service"
	"github.com/95eh/eg"
	"github.com/95eh/eg/codec"
	"github.com/95eh/eg/gate"
	"github.com/95eh/eg/start"
	"github.com/95eh/eg/svc"
)

func StartGlobal() {
	var (
		folder string
		config string
	)
	defFolder := eg.ExeDir() + "/../configs"
	flag.StringVar(&folder, "f", defFolder, "启动目录，默认程序启动目录同级configs目录")
	flag.StringVar(&config, "c", eg.ExeName(), "配置文件名，默认程序名")
	flag.Parse()
	eg.SetConfRoot(folder)
	cmn.LoadGlobalConf(
		fmt.Sprintf("%s.yml", config),
		fmt.Sprintf("%s_dev.yml", config))

	conf := cmn.GlobalConf()
	cmn.InitMongo(conf.MongoDB)
	globalRedisFac, globalRedisPool := cmn.GetRedisFac(conf.Redis)
	eg.SetMode(conf.App.Mode)
	eg.SetAppName(conf.App.Name)
	eg.SetAppVer(conf.App.Ver)
	eg.SetRegion(conf.App.Region)
	deadline := int64(10)
	if eg.Mode() == eg.ModeDebug {
		deadline = 60 * 100
	}
	opt := start.Option{
		Codec: codec.NewPbCodec(),
		RedisOpts: []svc.RedisOption{
			svc.RedisGlobalConnFac(globalRedisFac),
			svc.RedisGlobalConnPool(globalRedisPool),
		},
		RegisterOpts: []svc.RegisterOption{
			svc.Register(svc.NewRedisRegister()),
		},
		NetOpts: []svc.NetOption{
			svc.NetIp(conf.Net.Ip),
			svc.NetPort(conf.Net.Port),
		},
		ReqOpts: []svc.ReqOption{
			svc.ServiceAlias(),
		},
		Gate: true,
		GateOpts: []gate.Option{
			gate.IpBlocker(userIpBlocker),
			gate.DefaultMask(eg.GenMask(cmn.RoleGuest)),
			gate.DefaultRole(cmn.RolePlayer),
			gate.Deadline(deadline),
		},
		ConnOpts: []gate.ConnOption{
			gate.ConnHttpPort(conf.Http.Port),
			gate.ConnHttpIdParser(func(bytes []byte) (int64, int64, eg.IErr) {
				if bytes == nil {
					return 0, 0, nil
				}
				_, c, err := cmn.ParseJwt(eg.BytesToStr(bytes))
				if err != nil {
					return 0, 0, err
				}
				return c.Id, c.Uid, nil
			}),
		},
	}
	start.Default(opt)

	service.InitGlobalAccount(conf.Account)
	service.InitGlobalRegion(conf.Region)

	eg.Start()

	service.GlobalTickUpdateRegions(conf.Region.CheckSecs)

	if conf.Account.Enable {
		eg.Register().RegisterService(cmn.SvcAccount)
	}
	if conf.Region.Enable {
		eg.Register().RegisterService(cmn.SvcRegion)
	}

	eg.BeforeExit("game", func() {
		eg.Log().Info("exit game", nil)
		eg.Dispose()
	})
	eg.WaitExit()
}
