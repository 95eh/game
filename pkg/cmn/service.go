package cmn

import (
	"github.com/95eh/eg"
	"github.com/95eh/eg/svc/proto"
)

const (
	SvcAccount eg.TService = 50 + iota
	SvcRegion
)

const (
	// SvcGate 网关服务
	SvcGate = proto.SvcRegionBegin + iota
	// SvcCharacter 角色服务
	SvcCharacter
	// SvcScene 通用场景服务别名
	SvcScene
	// SvcSceneWorld 大世界场景服务
	SvcSceneWorld
)

func init() {
	eg.SetService(SvcAccount, "account")
	eg.SetService(SvcRegion, "region")
	eg.SetService(SvcGate, "gate")
	eg.SetService(SvcCharacter, "character")
	eg.SetService(SvcScene, "scene")
	eg.SetService(SvcSceneWorld, "scene_world")
}
