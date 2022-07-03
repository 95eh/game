package proto

import (
	"game/pkg/cmn"
	"game/pkg/proto/pb"
	"github.com/95eh/eg"
)

const (
	CdAccountSignUp = iota
	CdAccountSignIn
	CdAccountSignOut
)

func init() {
	eg.SetCodeNames(cmn.SvcAccount, map[eg.TCode]string{
		CdAccountSignUp:  "sign_up",
		CdAccountSignIn:  "sign_in",
		CdAccountSignOut: "sign_out",
	})
}

func InitAccountCodec() {
	c := eg.Codec()
	c.BindFac(cmn.SvcAccount, CdAccountSignIn,
		func() interface{} {
			return &pb.ReqAccountSignIn{}
		},
		func() interface{} {
			return &pb.ResAccountSignIn{}
		})
	c.BindFac(cmn.SvcAccount, CdAccountSignUp,
		func() interface{} {
			return &pb.ReqAccountSignUp{}
		},
		func() interface{} {
			return &pb.ResAccountSignUp{}
		})
	c.BindFac(cmn.SvcAccount, CdAccountSignOut,
		func() interface{} {
			return &pb.ReqAccountSignOut{}
		},
		func() interface{} {
			return &pb.ResAccountSignOut{}
		})
}
