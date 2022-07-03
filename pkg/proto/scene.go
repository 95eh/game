package proto

import (
	"game/pkg/cmn"
	"game/pkg/proto/pb"
	"github.com/95eh/eg"
)

const (
	CdSceneEnter = iota
	CdSceneEnterNtc
	CdSceneExit
	CdSceneChangePos
	CdSceneChangePosNtc
	CdSceneVisibleNtc
	CdSceneInvisibleNtc
	CdSceneMoveStart
	CdSceneMoveStartNtc
	CdSceneMoveStop
	CdSceneMoveStopNtc
	CdSceneCreate
	CdSceneDispose
)

func init() {
	eg.SetCodeNames(cmn.SvcScene, map[eg.TCode]string{
		CdSceneEnter:        "enter",
		CdSceneEnterNtc:     "enter_ntc",
		CdSceneExit:         "exit_curr",
		CdSceneChangePos:    "change_pos",
		CdSceneChangePosNtc: "change_pos_ntc",
		CdSceneVisibleNtc:   "visible_ntc",
		CdSceneInvisibleNtc: "invisible_ntc",
		CdSceneMoveStart:    "move_start",
		CdSceneMoveStartNtc: "move_start_ntc",
		CdSceneMoveStop:     "move_stop",
		CdSceneMoveStopNtc:  "move_stop_ntc",
		CdSceneCreate:       "create",
		CdSceneDispose:      "dispose",
	})
}

func InitSceneCodec() {
	c := eg.Codec()
	c.BindFac(cmn.SvcScene, CdSceneEnter,
		func() interface{} {
			return &pb.ReqSceneEnter{}
		},
		func() interface{} {
			return &pb.ResSceneEnter{}
		})
	c.BindFac(cmn.SvcScene, CdSceneEnterNtc,
		nil,
		func() interface{} {
			return &pb.NtcSceneEnter{}
		})
	c.BindFac(cmn.SvcScene, CdSceneExit,
		func() interface{} {
			return &pb.ReqSceneExit{}
		},
		func() interface{} {
			return &pb.ResSceneExit{}
		})
	c.BindFac(cmn.SvcScene, CdSceneChangePos,
		func() interface{} {
			return &pb.ReqSceneChangePos{}
		},
		func() interface{} {
			return &pb.ResSceneChangePos{}
		})
	c.BindFac(cmn.SvcScene, CdSceneChangePosNtc,
		nil,
		func() interface{} {
			return &pb.NtcSceneChangePos{}
		})
	c.BindFac(cmn.SvcScene, CdSceneVisibleNtc,
		nil,
		func() interface{} {
			return &pb.NtcSceneVisible{}
		})
	c.BindFac(cmn.SvcScene, CdSceneInvisibleNtc,
		nil,
		func() interface{} {
			return &pb.NtcSceneInvisible{}
		})
	c.BindFac(cmn.SvcScene, CdSceneMoveStart,
		func() interface{} {
			return &pb.ReqSceneMoveStart{}
		}, func() interface{} {
			return &pb.ResSceneMoveStart{}
		})
	c.BindFac(cmn.SvcScene, CdSceneMoveStartNtc,
		nil,
		func() interface{} {
			return &pb.NtcSceneMoveStart{}
		})
	c.BindFac(cmn.SvcScene, CdSceneMoveStop,
		func() interface{} {
			return &pb.ReqSceneMoveStop{}
		}, func() interface{} {
			return &pb.ResSceneMoveStop{}
		})
	c.BindFac(cmn.SvcScene, CdSceneMoveStopNtc,
		nil,
		func() interface{} {
			return &pb.NtcSceneMoveStop{}
		})
}
