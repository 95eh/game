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

	svc := eg.Svc().RegisterService(cmn.SvcGate)
	{
		svc.BindReq(proto.CdGateSetCharacterId, gateSetCharacterId)
		svc.BindReq(proto.CdGatePushToUser, gatePushToUser)
		svc.BindReq(proto.CdGatePushToUsers, gatePushToUsers)
		svc.BindReq(proto.CdGatePushToAllUsers, gatePushToAllUser)
	}

	gateRegisterAddr(conf)
}

func gateSetCharacterId(ctx eg.ICtx) {
	req := ctx.Body().(*pb.ReqGateSetCharacterId)
	if req.Cid == ctx.GetId() {
		ctx.Ok(&pb.ResGateSetCharacterId{})
		return
	}
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
		"websocket_ip", conf.Gate.Websocket.Ip,
		"websocket_port", conf.Gate.Websocket.Port,
		"udp_ip", conf.Gate.Udp.Ip,
		"udp_port", conf.Gate.Udp.Port,
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

func gateAuth(agentId int64, body interface{}, resOk eg.FnInt64Any, resErr eg.FnInt64Err) {
	req := body.(*pb.ReqGateAuth)
	_, claims, err := cmn.ParseJwt(req.Token)
	if err != nil {
		resErr(0, err)
		return
	}

	resOk(0, &pb.ResGateAuth{
		ErrCode: 0,
	})

	eg.Gate().SetAgentData(agentId, claims.Mask, claims.Uid, claims.Expire, func(oldAgentId int64, ok bool) {
		if !ok {
			resOk(0, &pb.ResGateAuth{
				ErrCode: 1,
			})
			eg.Conn().Close(agentId)
			return
		}
		if oldAgentId == 0 {
			return
		}
		eg.Log().Info("????????????", eg.M{
			"old agent id": oldAgentId,
		})
		eg.Conn().Send(oldAgentId, cmn.SvcGate, proto.CdGateAuth, &pb.ResGateAuth{
			ErrCode: 2,
		})
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
