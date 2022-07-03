package service

import (
	"fmt"
	"game/pkg/cmn"
	"game/pkg/proto"
	"game/pkg/proto/pb"
	"github.com/95eh/eg"
	"github.com/gomodule/redigo/redis"
)

func InitGate(conf *cmn.GameConfig) {
	proto.InitGateCodec()
	initUser()
	svc := eg.Svc().BindService(cmn.SvcGate)
	{
		svc.Bind(proto.CdGateSetCharacterId, gateSetCharacterId)
		svc.Bind(proto.CdGateDisconnect, gateDisconnect)
		svc.Bind(proto.CdGatePushToUser, gatePushToUser)
		svc.Bind(proto.CdGatePushToUsers, gatePushToUsers)
		svc.Bind(proto.CdGatePushToAllUsers, gatePushToAllUser)
	}
	user := eg.Gate()
	{
		user.SetRole(cmn.RoleGuest, cmn.SvcGate,
			proto.CdGateAuth,
		)
		ug := user.BindService(cmn.SvcGate)
		{
			ug.BindRequest(proto.CdGateAuth, gateAuth)
		}
	}
	gateRegisterAddr(conf)
}

func gateDisconnect(ctx eg.ICtx) {
	id := ctx.GetId()
	eg.Gate().DelUserData(id, func(uid, alias string, ok bool) {
		ctx.Ok(&pb.ResGateDisconnect{})
		eg.Conn().Close(id)
		eg.Log().Info("user disconnect", eg.M{
			"uid":   uid,
			"alias": alias,
		})
		if alias == "" {
			return
		}
		eg.Ntc().NtcDelUserServiceAlias(0, alias)
		eg.Ntc().PushNotice(cmn.SvcGate, proto.CdGateUserOfflineNtc,
			0, alias, &pb.NtcGateUserOffline{})
	})
}

func gateSetCharacterId(ctx eg.ICtx) {
	req := ctx.Body().(*pb.ReqGateSetCharacterId)
	eg.Gate().SetAlias(ctx.GetId(), req.Cid, func(err eg.IErr) {
		if err != nil {
			ctx.Err(err)
			return
		}
		ctx.Ok(&pb.ResGateSetCharacterId{})
	})
}

func gateRegisterAddr(conf *cmn.GameConfig) {
	conn := eg.Redis().SpawnGlobalConn()
	defer conn.Close()

	serviceKey := getRedisGateKey(eg.Region())
	conn.Send(eg.MULTI)
	conn.Send(eg.HSET, serviceKey,
		"http_ip", conf.Gate.Http.Ip,
		"http_port", conf.Gate.Http.Port,
		"socket_ip", conf.Gate.Tcp.Ip,
		"socket_port", conf.Gate.Tcp.Port,
		"count", 0,
		"cap", conf.Gate.Cap,
		"status", 1)
	_, e := conn.Do(eg.EXEC)
	if e != nil {
		panic(e)
	}
	num, e := redis.Int(conn.Do(eg.HEXISTS, "region:name", fmt.Sprintf("%d", conf.App.Region)))
	if e != nil {
		panic(e)
	}
	if num == 0 {
		resultNum, e := redis.Int(conn.Do(eg.HSET, "region:name", fmt.Sprintf("%d", conf.App.Region), conf.App.Name))
		if e != nil {
			panic(e)
		}
		if resultNum == 0 {
			panic(e)
		}
	}
}

func gateAuth(agentId string, body interface{}, resOk eg.FnInt64Any, resErr eg.FnInt64Err) {
	req := body.(*pb.ReqGateAuth)
	_, claims, err := cmn.ParseJwt(req.Token)
	if err != nil {
		resErr(0, err)
		return
	}

	resOk(0, &pb.ResGateAuth{
		ErrCode: 0,
	})

	eg.Gate().SetUserData(agentId, claims.Uid, func(oldAgentId string, ok bool) {
		if !ok {
			return
		}
		eg.Log().Info("重复登陆", eg.M{
			"old agent id": oldAgentId,
		})
		eg.Conn().Send(oldAgentId, cmn.SvcGate, proto.CdGateAuthRepeat, &pb.ResGateAuthRepeat{})
		eg.Conn().Close(oldAgentId)
	})
}

func gatePushToUser(ctx eg.ICtx) {
	req := ctx.Body().(*pb.PushToUser)
	eg.Conn().SendByUserId(ctx.GetId(), req.Bytes)
}

func gatePushToUsers(ctx eg.ICtx) {
	req := ctx.Body().(*pb.PushToUsers)
	eg.Conn().SendByUserIds(req.Ids, req.Bytes)
}

func gatePushToAllUser(ctx eg.ICtx) {
	req := ctx.Body().(*pb.PushToAllUsers)
	eg.Conn().SendAll(eg.CopyBytes(req.Bytes))
}

func getRedisGateKey(region eg.TRegion) string {
	return fmt.Sprintf("region:%d:gate", region)
}

func initUser() {
	eg.Gate().SetRole(cmn.RoleGuest, cmn.SvcGate,
		proto.CdGateAuth)
	eg.Gate().SetRole(cmn.RolePlayer, cmn.SvcCharacter,
		proto.CdCharacterInfo,
		proto.CdCharacterCreate,
		proto.CdCharacterSelect,
		proto.CdCharacterEnterScene,
		proto.CdCharacterChangeScene,
		proto.CdCharacterExitScene,
	)
	eg.Gate().SetRole(cmn.RolePlayer, cmn.SvcScene)
}
