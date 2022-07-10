package pkg

import (
	"flag"
	"fmt"
	"game/pkg/cmn"
	"game/pkg/proto"
	"game/pkg/proto/pb"
	"game/pkg/service"
	"github.com/95eh/eg"
	"github.com/95eh/eg/codec"
	"github.com/95eh/eg/gate"
	"github.com/95eh/eg/start"
	"github.com/95eh/eg/svc"
)

func StartGame() {
	var (
		folder string
		config string
	)
	defFolder := eg.ExeDir() + "/../configs"
	flag.StringVar(&folder, "f", defFolder, "启动目录，默认程序启动目录同级configs目录")
	flag.StringVar(&config, "c", eg.ExeName(), "配置文件名，默认程序名")
	flag.Parse()
	eg.SetConfRoot(folder)
	cmn.LoadGameConf(
		fmt.Sprintf("%s.yml", config),
		fmt.Sprintf("%s_dev.yml", config))

	conf := cmn.GameConf()
	cmn.InitMongo(conf.MongoDB)

	globalRedisFac, globalRedisPool := cmn.GetRedisFac(conf.GlobalRedis)
	regionRedisFac, regionRedisPool := cmn.GetRedisFac(conf.RegionRedis)
	eg.SetMode(conf.App.Mode)
	eg.SetAppName(conf.App.Name)
	eg.SetAppVer(conf.App.Ver)
	eg.SetRegion(conf.App.Region)
	opt := start.Option{
		Codec: codec.NewPbCodec(),
		RedisOpts: []svc.RedisOption{
			svc.RedisGlobalConnFac(globalRedisFac),
			svc.RedisGlobalConnPool(globalRedisPool),
			svc.RedisRegionConnFac(regionRedisFac),
			svc.RedisRegionConnPool(regionRedisPool),
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
	}
	deadline := int64(30)
	if eg.Mode() == eg.ModeDebug {
		deadline = 60 * 100
	}
	if conf.Gate.Enable {
		opt.Gate = true
		opt.GateOpts = []gate.Option{
			gate.IpBlocker(userIpBlocker),
			gate.DefaultMask(eg.GenMask(cmn.RoleGuest)),
			gate.DefaultRole(cmn.RolePlayer),
			gate.Deadline(deadline),
		}
		opt.ConnOpts = []gate.ConnOption{
			gate.ConnTcpPort(conf.Gate.Tcp.Port),
			gate.ConnHttpPort(conf.Gate.Http.Port),
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
			gate.ConnClosed(func(agentId int64, err eg.IErr) {
				eg.Gate().DelAgentData(agentId, func(id, alias int64, ok bool) {
					if !ok {
						return
					}
					eg.Log().Info("user disconnect", eg.M{
						"id":      id,
						"alias":   alias,
						"agentId": agentId,
					})
					if alias == 0 {
						return
					}
					eg.Ntc().NtcDelUserServiceAlias(0, alias)
					eg.Ntc().PushNotice(cmn.SvcGate, proto.CdGateUserOfflineNtc,
						0, alias, &pb.NtcGateUserOffline{})
				})
			}),
		}
	}
	start.Default(opt)

	cmn.InitSensitiveWords(conf.Other.SensitiveFilterUrl)

	if conf.Gate.Enable {
		service.InitGate(conf)
	}
	service.InitSceneWorld(conf.ScnWorld)
	service.InitCharacter(conf.Character)

	eg.Start()

	service.CreateSceneWorld()

	if conf.ScnWorld.Enable {
		eg.Register().RegisterService(cmn.SvcSceneWorld)
	}
	if conf.Character.Enable {
		eg.Register().RegisterService(cmn.SvcCharacter)
	}
	if conf.Gate.Enable {
		eg.Register().RegisterService(cmn.SvcGate)
	}

	eg.BeforeExit("game", func() {
		eg.Log().Info("exit game", nil)
		eg.Dispose()
	})
	eg.WaitExit()
}

func userIpBlocker(ip string, fn eg.FnBool) {
	//todo
	fn(true)
}
