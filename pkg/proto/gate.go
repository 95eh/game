package proto

import (
	"game/pkg/cmn"
	"game/pkg/proto/pb"
	"github.com/95eh/eg"
)

const (
	CdGatePushToUser eg.TCode = iota
	CdGatePushToUsers
	CdGatePushToAllUsers
	CdGateAuth
	CdGateAuthRepeat
	CdGateSetCharacterId
	CdGateDisconnect
	CdGateUserOfflineNtc
)

func init() {
	eg.SetCodeNames(cmn.SvcGate, map[eg.TCode]string{
		CdGatePushToUser:     "push to user",
		CdGatePushToUsers:    "push to users",
		CdGatePushToAllUsers: "push to all users",
		CdGateAuth:           "auth",
		CdGateAuthRepeat:     "auth_repeat",
		CdGateSetCharacterId: "set_character_id",
		CdGateDisconnect:     "disconnect",
		CdGateUserOfflineNtc: "user_offline_ntc",
	})
}
func InitGateCodec() {
	c := eg.Codec()
	c.BindFac(cmn.SvcGate, CdGateAuth,
		func() interface{} {
			return &pb.ReqGateAuth{}
		},
		func() interface{} {
			return &pb.ResGateAuth{}
		})
	c.BindFac(cmn.SvcGate, CdGateAuthRepeat,
		nil,
		func() interface{} {
			return &pb.ResGateAuthRepeat{}
		})
	c.BindFac(cmn.SvcGate, CdGateSetCharacterId,
		func() interface{} {
			return &pb.ReqGateSetCharacterId{}
		},
		func() interface{} {
			return &pb.ResGateSetCharacterId{}
		})
	c.BindFac(cmn.SvcGate, CdGatePushToUser,
		func() interface{} {
			return &pb.PushToUser{}
		},
		nil)
	c.BindFac(cmn.SvcGate, CdGatePushToUsers,
		func() interface{} {
			return &pb.PushToUsers{}
		},
		nil)
	c.BindFac(cmn.SvcGate, CdGatePushToAllUsers,
		func() interface{} {
			return &pb.PushToAllUsers{}
		},
		nil)
	c.BindFac(cmn.SvcGate, CdGateDisconnect,
		func() interface{} {
			return &pb.ReqGateDisconnect{}
		},
		func() interface{} {
			return &pb.ResGateDisconnect{}
		})
	c.BindFac(cmn.SvcGate, CdGateUserOfflineNtc,
		nil,
		func() interface{} {
			return &pb.NtcGateUserOffline{}
		})
}
