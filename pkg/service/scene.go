package service

import (
	"game/pkg/cmn"
	"game/pkg/proto"
	"game/pkg/proto/pb"
	"github.com/95eh/eg"
	"github.com/95eh/eg/scene"
	"github.com/95eh/rpg"
)

func InitScene(sceneService eg.TService) {
	proto.InitSceneCodec()
	svc := eg.Svc().BindService(sceneService)
	{
		svc.Bind(proto.CdSceneEnter, sceneEnter)
		svc.Bind(proto.CdSceneExit, sceneExit)
		svc.Bind(proto.CdSceneChangePos, sceneChangePos)
		svc.Bind(proto.CdSceneMoveStart, sceneMoveStart)
		svc.Bind(proto.CdSceneMoveStop, sceneMoveStop)
		svc.Bind(proto.CdSceneCreate, sceneCreate)
		svc.Bind(proto.CdSceneDispose, sceneDispose)
	}
	if eg.Scene() == nil {
		eg.AddModule(
			scene.NewMScene(),
			rpg.NewMTileScene(),
		)
	}
	eg.Scene().BindActorFac(cmn.ActorPlayer,
		nil,
		rpg.Ac_Transform, rpg.Ac_Life)
	eg.Scene().BindActorFac(cmn.ActorTest,
		func(actor eg.IActor, o *eg.Object) {
			actor.BindEventProcessor(rpg.Evt_Visible, DebugVisible)
			actor.BindEventProcessor(rpg.Evt_Invisible, DebugInvisible)
		},
		rpg.Ac_Transform, rpg.Ac_Life, rpg.Ac_Rand_Move) //
	eg.Scene().BindActorComponentFac(rpg.Ac_Transform, rpg.NewAcTransform)
	eg.Scene().BindActorComponentFac(rpg.Ac_Life, rpg.NewAcLife)
	eg.Scene().BindActorComponentFac(rpg.Ac_Rand_Move, rpg.NewAcAutoMove)
}

func sceneDispose(ctx eg.ICtx) {
	req := ctx.Body().(*pb.ReqSceneDispose)
	eg.Scene().SpawnScheduler(ctx.Tid(), nil).
		GetScene(req.SceneId).
		DisposeScene().
		Do(nil, func(object *eg.Object, err eg.IErr) {
			if err != nil {
				ctx.Err(err)
			} else {
				ctx.Ok(&pb.ResSceneDispose{})
			}
		})
}

func sceneCreate(ctx eg.ICtx) {
	req := ctx.Body().(*pb.ReqSceneCreate)
	eg.Scene().SpawnScheduler(ctx.Tid(), nil).
		CreateScene(eg.TScene(req.SceneType), req.SceneId, req.ActorCap, req.Data).
		Do(nil, func(object *eg.Object, err eg.IErr) {
			if err != nil {
				ctx.Err(err)
			} else {
				ctx.Ok(&pb.ResSceneCreate{})
			}
		})
}

func sceneEnter(ctx eg.ICtx) {
	req := ctx.Body().(*pb.ReqSceneEnter)
	m, err := eg.GobToM(req.Bytes)
	if err != nil {
		ctx.Err(err)
		return
	}
	o := eg.NewObject(m)
	actorId := ctx.GetId()
	eg.Scene().SpawnScheduler(ctx.Tid(), nil).
		GetScene(req.SceneId).
		SpawnActor(eg.TActor(req.Type), actorId, o).
		SpawnActorEvent(rpg.NewEvtVisible).
		GetTagsTag(func(actor eg.IActor, s eg.ISceneWorkerScheduler) ([]string, string, eg.IErr) {
			c, _ := actor.GetComponent(rpg.Ac_Transform)
			tags, tag := c.(*rpg.AcTransform).GetVisionTileTags()
			return tags, tag, nil
		}).
		GetTagsActors().
		ResetActorEvents(1).
		FnActors(func(actor eg.IActor, s eg.ISceneWorkerScheduler) eg.IErr {
			if actor.Id() == actorId {
				return nil
			}
			actor.ProcessEvent(s.ActorEvent())
			s.PushEvent(rpg.NewEvtVisible(actor))
			return nil
		}).
		ActorProcessEvents().
		Do(nil, func(object *eg.Object, err eg.IErr) {
			if err != nil {
				ctx.Err(err)
			} else {
				ctx.Ok(&pb.ResSceneEnter{})
			}
		})
}

func sceneExit(ctx eg.ICtx) {
	actorId := ctx.GetId()
	eg.Scene().SpawnScheduler(ctx.Tid(), nil).
		GetActorAndScene(actorId).
		SpawnActorEvent(rpg.NewEvtInvisible).
		GetTagsTag(func(actor eg.IActor, s eg.ISceneWorkerScheduler) ([]string, string, eg.IErr) {
			c, _ := actor.GetComponent(rpg.Ac_Transform)
			tags, tag := c.(*rpg.AcTransform).GetVisionTileTags()
			return tags, tag, nil
		}).
		GetTagsActors().
		ResetActorEvents(1).
		FnActors(func(actor eg.IActor, s eg.ISceneWorkerScheduler) eg.IErr {
			if actor.Id() == actorId {
				return nil
			}
			actor.ProcessEvent(s.ActorEvent())
			s.PushEvent(rpg.NewEvtInvisible(actor))
			return nil
		}).
		ActorProcessEvents().
		Do(nil, func(object *eg.Object, err eg.IErr) {
			if err != nil {
				ctx.Err(err)
			} else {
				ctx.Ok(&pb.ResSceneExit{})
			}
		})
}

func sceneChangePos(ctx eg.ICtx) {

}

func sceneMoveStart(ctx eg.ICtx) {
	req := ctx.Body().(*pb.ReqSceneMoveStart)
	eg.Scene().SpawnScheduler(ctx.Tid(), nil).
		GetActorAndScene(ctx.GetId()).
		SpawnActorEvent(rpg.NewEvtMoveStart).
		GetTags(func(actor eg.IActor, s eg.ISceneWorkerScheduler) (tags []string, e eg.IErr) {
			c, _ := actor.GetComponent(rpg.Ac_Transform)
			tags = c.(*rpg.AcTransform).StartMove(eg.NewVec3(req.ForX, req.ForY, req.ForZ))
			return
		}).
		GetTagsActors().
		ActorsProcessEvent().
		Do(nil, func(object *eg.Object, err eg.IErr) {
			if err != nil {
				ctx.Err(err)
			} else {
				ctx.Ok(&pb.ResSceneMoveStart{})
			}
		})
}

func sceneMoveStop(ctx eg.ICtx) {
	eg.Scene().SpawnScheduler(ctx.Tid(), nil).
		GetActorAndScene(ctx.GetId()).
		SpawnActorEvent(rpg.NewEvtMoveStop).
		GetTags(func(actor eg.IActor, s eg.ISceneWorkerScheduler) (tags []string, e eg.IErr) {
			c, _ := actor.GetComponent(rpg.Ac_Transform)
			tags = c.(*rpg.AcTransform).StopMove()
			return
		}).
		GetTagsActors().
		ActorsProcessEvent().
		Do(nil, func(object *eg.Object, err eg.IErr) {
			if err != nil {
				ctx.Err(err)
			} else {
				ctx.Ok(&pb.ResSceneMoveStop{})
			}
		})
}

func DebugVisible(actor eg.IActor, e eg.IActorEvent) {
	return
	evt := e.(*rpg.EvtVisible)
	if evt.ActorId == actor.Id() || evt.ActorType == cmn.ActorTest {
		return
	}
	//c, _ := actor.GetComponent(rpg.Ac_Transform)
	//tnf := c.(rpg.IActorComTransform)
	//pos := tnf.Position()
	//dis := eg.Vec3Distance(pos, evt.Pos)
	eg.Log().Debug("visible", eg.M{
		"actor": actor.Id(),
		"event": evt,
		//"dis":   dis,
	})
}

func DebugInvisible(actor eg.IActor, e eg.IActorEvent) {
	return
	evt := e.(*rpg.EvtInvisible)
	if evt.ActorId == actor.Id() || evt.ActorType == cmn.ActorTest {
		return
	}
	//c, _ := actor.GetComponent(rpg.Ac_Transform)
	//tnf := c.(rpg.IActorComTransform)
	//dis := eg.Vec3Distance(tnf.Position(), evt.Pos)
	eg.Log().Debug("invisible", eg.M{
		"actor": actor.Id(),
		"event": evt,
		//"dis":   dis,
	})
}
