package proto

import (
	"game/pkg/cmn"
	"game/pkg/proto/pb"
	"github.com/95eh/eg"
)

const (
	CdCharacterInfo = iota
	CdCharacterCreate
	CdCharacterSelect
	CdCharacterEnterScene
	CdCharacterChangeScene
	CdCharacterExitScene
)

func init() {
	eg.SetCodeNames(cmn.SvcCharacter, map[eg.TCode]string{
		CdCharacterInfo:        "info",
		CdCharacterCreate:      "create",
		CdCharacterSelect:      "select",
		CdCharacterEnterScene:  "enter_scene",
		CdCharacterChangeScene: "change_scene",
		CdCharacterExitScene:   "exit_scene",
	})
}

func InitCharacterCodec() {
	c := eg.Codec()
	c.BindFac(cmn.SvcCharacter, CdCharacterInfo,
		func() interface{} {
			return &pb.ReqCharacter{}
		},
		func() interface{} {
			return &pb.ResCharacter{}
		})
	c.BindFac(cmn.SvcCharacter, CdCharacterCreate,
		func() interface{} {
			return &pb.ReqCharacterCreate{}
		},
		func() interface{} {
			return &pb.ResCharacterCreate{}
		})
	c.BindFac(cmn.SvcCharacter, CdCharacterSelect,
		func() interface{} {
			return &pb.ReqCharacterSelect{}
		},
		func() interface{} {
			return &pb.ReqCharacterSelect{}
		})
	c.BindFac(cmn.SvcCharacter, CdCharacterEnterScene,
		func() interface{} {
			return &pb.ReqCharacterEnterScene{}
		},
		func() interface{} {
			return &pb.ResCharacterEnterScene{}
		})
	c.BindFac(cmn.SvcCharacter, CdCharacterChangeScene,
		func() interface{} {
			return &pb.ReqCharacterChangeScene{}
		},
		func() interface{} {
			return &pb.ResCharacterChangeScene{}
		})
	c.BindFac(cmn.SvcCharacter, CdCharacterExitScene,
		func() interface{} {
			return &pb.ReqCharacterExitScene{}
		},
		func() interface{} {
			return &pb.ResCharacterExitScene{}
		})
}
