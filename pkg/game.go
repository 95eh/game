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
		SvcOpts: []svc.Option{
			svc.ServiceAlias(),
		},
	}
	if conf.Gate.Enable {
		opt.Gate = true
		opt.GateOpts = []gate.Option{
			gate.BlackFilter(userBlackFilter),
			gate.DefaultMask(eg.GenMask(cmn.RoleGuest)),
			gate.DefaultRole(cmn.RoleSvc),
		}
		opt.ConnOpts = []gate.ConnOption{
			gate.ConnTcpPort(conf.Gate.Tcp.Port),
			gate.ConnHttpPort(conf.Gate.Http.Port),
			gate.ConnResOk(func(a eg.IAgent, svc eg.TService, code eg.TCode, o interface{}) {
				bytes, err := eg.Codec().Marshal(o)
				if err != nil {
					eg.Log().Error(err)
					return
				}
				buffer := eg.NewByteBufferWithLen(uint32(4 + len(bytes)))
				buffer.WUint16(svc)
				buffer.WUint16(code)
				buffer.Write(bytes)
				a.Send(buffer.All())
			}),
			gate.ConnResErr(func(a eg.IAgent, service eg.TService, code eg.TCode, e eg.TErrCode) {
				eg.Log().Error(eg.NewErr(e, eg.M{
					"server": service,
					"code":   code,
				}))
			}),
			gate.ConnPacker(func(service eg.TService, code eg.TCode, o interface{}) []byte {
				bytes, err := eg.Codec().Marshal(o)
				if err != nil {
					eg.Log().Error(err)
					return nil
				}
				buffer := eg.NewByteBufferWithLen(uint32(4 + len(bytes)))
				buffer.WUint16(service)
				buffer.WUint16(code)
				buffer.Write(bytes)
				return buffer.All()
			}),
			gate.ConnClosed(func(agentId string, err eg.IErr) {
				if err == nil {
					return
				}
				eg.Gate().DelUserData(agentId, func(id, alias string, ok bool) {
					if !ok {
						return
					}
					eg.Log().Info("user disconnect", eg.M{
						"id":      id,
						"alias":   alias,
						"agentId": agentId,
					})
					if alias == "" {
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
	if conf.ScnWorld.Enable {
		service.InitSceneWorld(conf.ScnWorld)
	}
	if conf.Character.Enable {
		service.InitCharacter(conf.Character)
	}

	eg.Start()
	eg.BeforeExit("game", func() {
		eg.Log().Info("exit game", nil)
		eg.Dispose()
	})
	eg.WaitExit()
}

func userBlackFilter(ip string) bool {
	//todo
	return false
}
