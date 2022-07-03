package proto

import (
	"game/pkg/cmn"
	"game/pkg/proto/pb"
	"github.com/95eh/eg"
)

const (
	CdRegionList = iota
)

func init() {
	eg.SetCodeNames(cmn.SvcRegion, map[eg.TCode]string{
		CdRegionList: "list",
	})
}

func InitRegionCodec() {
	c := eg.Codec()
	c.BindFac(cmn.SvcRegion, CdRegionList,
		func() interface{} {
			return &pb.ReqRegionList{}
		},
		func() interface{} {
			return &pb.ResRegionList{}
		})
}
