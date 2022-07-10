package service

import (
	"game/pkg/cmn"
	"game/pkg/proto"
	"game/pkg/proto/pb"
	"github.com/95eh/eg"
	"github.com/gomodule/redigo/redis"
	"strconv"
	"sync"
)

var (
	_GlobalRegionMtx  sync.RWMutex
	_GlobalRegionList []*pb.Region
)

func InitGlobalRegion(conf cmn.RegionConf) {
	proto.InitRegionCodec()

	svc := eg.Svc().RegisterService(cmn.SvcRegion)
	{
		svc.BindReq(proto.CdRegionList, regionList)
	}
}

func regionList(ctx eg.ICtx) {
	_GlobalRegionMtx.RLock()
	ctx.Ok(&pb.ResRegionList{
		List: _GlobalRegionList,
	})
	_GlobalRegionMtx.RUnlock()
}

func GlobalTickUpdateRegions(secs int64) {
	globalUpdateRegions(0)
	eg.Timer().Tick(secs, 0, globalUpdateRegions)
}

func globalUpdateRegions(c int32) {
	conn := eg.Redis().SpawnGlobalConn()
	defer conn.Close()
	nm, e := redis.StringMap(conn.Do(eg.HGETALL, "region:name"))
	if e != nil {
		eg.Log().Error(eg.WrapErr(eg.EcRedisErr, e))
		return
	}
	regions := make([]*pb.Region, 0, len(nm))
	for id, name := range nm {
		rid, e := strconv.ParseInt(id, 10, 64)
		if e != nil {
			eg.Log().Error(eg.NewErr(eg.EcRedisErr, eg.M{
				"error": e.Error(),
			}))
			continue
		}
		gm, e := redis.StringMap(conn.Do(eg.HGETALL,
			getRedisGateKey(eg.TRegion(rid))))
		if e != nil {
			eg.Log().Error(eg.WrapErr(eg.EcRedisErr, e))
			continue
		}
		count, e := strconv.Atoi(gm["count"])
		if e != nil {
			eg.Log().Error(eg.WrapErr(eg.EcRedisErr, e))
			continue
		}
		cap, e := strconv.Atoi(gm["cap"])
		if e != nil {
			eg.Log().Error(eg.WrapErr(eg.EcRedisErr, e))
			continue
		}
		status, e := strconv.Atoi(gm["status"])
		if e != nil {
			eg.Log().Error(eg.WrapErr(eg.EcRedisErr, e))
			continue
		}
		regions = append(regions, &pb.Region{
			Id:            rid,
			Name:          name,
			HttpIp:        gm["http_ip"],
			HttpPort:      gm["http_port"],
			SocketIp:      gm["socket_ip"],
			SocketPort:    gm["socket_port"],
			WebsocketIp:   gm["websocket_ip"],
			WebsocketPort: gm["websocket_port"],
			UdpIp:         gm["udp_ip"],
			UdpPort:       gm["udp_port"],
			Count:         uint32(count),
			Cap:           uint32(cap),
			Status:        uint32(status),
		})
	}
	_GlobalRegionMtx.Lock()
	_GlobalRegionList = regions
	_GlobalRegionMtx.Unlock()
}
