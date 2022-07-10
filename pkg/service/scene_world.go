package service

import (
	"fmt"
	"game/pkg/cmn"
	"game/pkg/proto/pb"
	"github.com/95eh/eg"
	"github.com/95eh/eg/svc"
	"github.com/95eh/rpg"
)

func InitSceneWorld(conf cmn.ScnWorldConf) {
	service, err := cmn.SceneTypeToSceneService(cmn.SceneWorld)
	if err != nil {
		eg.Fatal(err)
	}
	InitScene(service)
	eg.Timer().StartQuickTicker(rpg.TnfMoveTicker, rpg.TnfMoveTickDur)
	eg.Scene().BindSceneFac(cmn.SceneWorld, worldSceneFac)
	rpg.TileScene().AddSceneTplConf(cmn.SceneWorld, conf.Tile)
}

func CreateSceneWorld() {
	//todo 读取数据
	_, err := eg.Scene().CreateScene(cmn.SceneWorld, 1, 1024, nil)
	if err != nil {
		eg.Fatal(err)
	}
}

func worldSceneFac(scene eg.IScene, o interface{}) eg.IErr {
	//eg.Timer().After(1000, func() {
	//	testActors()
	//	testPlayer()
	//})
	return nil
}

func testPlayer() {
	id := int64(1)
	bytes, _ := eg.M{
		rpg.TnfPos: eg.Vec3{
			X: 100,
			Y: 0,
			Z: 100,
		},
		rpg.TnfMoveSpeed: float32(6),
		rpg.LfName:       "95eh",
	}.Gob()
	reqEnter := &pb.ReqSceneEnter{
		SceneId: 1,
		Type:    int32(cmn.ActorPlayer),
		Bytes:   bytes,
	}
	tid := eg.Log().Sign(0, eg.M{
		"code":       "scene enter",
		"scene id":   reqEnter.SceneId,
		"actor type": reqEnter.Type,
	})
	ctx := svc.SpawnCtx(tid, id, reqEnter, func(tid int64, obj interface{}) {
		eg.Log().Info("scene enter ok", eg.M{
			"response": obj,
		})
	}, func(tid int64, err eg.IErr) {
		eg.Log().Info("scene enter err", eg.M{
			"error": err.String(),
		})
	})
	sceneEnter(ctx)

	reqMoveStart := &pb.ReqSceneMoveStart{
		ForX: 1,
		ForY: 0,
		ForZ: 1,
	}
	tid = eg.Log().Sign(0, eg.M{
		"move start": reqMoveStart,
	})
	moveStartCtx := svc.SpawnCtx(tid, id, reqMoveStart, func(tid int64, obj interface{}) {
		eg.Log().Info("move start ok", eg.M{
			"response": obj,
		})
	}, func(tid int64, err eg.IErr) {
		eg.Log().Info("move start err", eg.M{
			"error": err.String(),
		})
	})
	sceneMoveStart(moveStartCtx)
}

func testActors() {
	return
	tsc, _ := rpg.TileScene().GetSceneTplData(cmn.SceneWorld)
	conf := tsc.Conf()
	x := conf.Width / conf.TileSize
	y := conf.Length / conf.TileSize
	tsh := conf.TileSize / 2
	n := int32(8)
	v := x * y * n
	eg.Log().Debug("test total", eg.M{
		"v": v,
	})
	eg.Timer().After(1000, func() {
		for i := int32(0); i < x; i++ {
			for j := int32(0); j < y; j++ {
				for m := int32(0); m < n; m++ {
					//id := primitive.NewObjectID()
					id := eg.SId().GetRegionId()
					name := fmt.Sprintf("test_%d_%d_%d", i, j, m)
					px := float32((i * conf.TileSize) + tsh)
					py := float32((j * conf.TileSize) + tsh)
					bytes, _ := eg.M{
						rpg.TnfPos: eg.Vec3{
							X: px,
							Y: 0,
							Z: py,
						},
						rpg.TnfMoveSpeed: float32(6),
						rpg.LfName:       name,
					}.Gob()
					req := &pb.ReqSceneEnter{
						SceneId: 1,
						Type:    int32(cmn.ActorTest),
						Bytes:   bytes,
					}
					tid := eg.Log().Sign(0, eg.M{
						"code":       "scene enter",
						"scene id":   req.SceneId,
						"actor type": req.Type,
					})
					ctx := svc.SpawnCtx(tid, id, req, func(tid int64, obj interface{}) {

					}, func(tid int64, err eg.IErr) {

					})
					sceneEnter(ctx)
				}
			}
		}
	})
}
