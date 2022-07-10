package service

import (
	"context"
	"game/pkg/cmn"
	"game/pkg/model"
	"game/pkg/proto"
	"game/pkg/proto/pb"
	"github.com/95eh/eg"
	"github.com/95eh/rpg"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitCharacter(conf cmn.CharacterConf) {
	proto.InitCharacterCodec()

	svc := eg.Svc().RegisterService(cmn.SvcCharacter)
	{
		svc.BindReq(proto.CdCharacterInfo, characterInfo)
		svc.BindReq(proto.CdCharacterCreate, characterCreate)
		svc.BindReq(proto.CdCharacterSelect, characterSelect)
		svc.BindReq(proto.CdCharacterEnterScene, characterEnterScene)
		svc.BindReq(proto.CdCharacterChangeScene, characterChangeScene)
		svc.BindReq(proto.CdCharacterExitScene, characterExitScene)
	}
	ntc := eg.Ntc().BindServiceNotice(cmn.SvcGate)
	{
		ntc.Bind(proto.CdGateUserOfflineNtc, cmn.SvcCharacter, characterOffLine)
	}
}

func characterColl() *mongo.Collection {
	return cmn.Db().Collection(model.MdCharacter)
}

func characterExitScene(ctx eg.ICtx) {
	var character model.Character
	e := characterColl().FindOne(context.TODO(), bson.D{
		{"cid", ctx.GetId()},
	}, options.FindOne().SetProjection(
		bson.D{
			{"scene_type", 1},
			{"scene_id", 1},
			{"posX", 1},
			{"posY", 1},
			{"nick_name", 1},
		}),
	).Decode(&character)
	if e != nil {
		ctx.Err(eg.WrapErr(eg.EcUnmarshallErr, e))
		return
	}

	eg.Req().Request(ctx.Tid(), ctx.GetId(), cmn.SvcScene, proto.CdSceneExit,
		&pb.ReqSceneExit{},
		func(tid int64, obj any) {
			ctx.Ok(&pb.ResCharacterExitScene{})
		}, func(tid int64, err eg.IErr) {
			ctx.Err(err)
		})
}

func characterChangeScene(ctx eg.ICtx) {
	ctx.Ok(&pb.ResCharacterChangeScene{})
	//req := ctx.Body().(*pb.ReqCharacterChangeScene)
	// todo 判断能否进入场景
	// todo 退出当前的场景
}

func characterEnterScene(ctx eg.ICtx) {
	var character model.Character
	e := characterColl().FindOne(context.TODO(), bson.D{
		{"cid", ctx.GetId()},
	}, options.FindOne().SetProjection(
		bson.D{
			{"scene_type", 1},
			{"scene_id", 1},
			{"posX", 1},
			{"posY", 1},
			{"nick_name", 1},
		}),
	).Decode(&character)
	if e != nil {
		ctx.Err(eg.WrapErr(eg.EcUnmarshallErr, e))
		return
	}

	bytes, err := eg.M{
		rpg.TnfPos: eg.NewVec3(character.PosX, 0, character.PosY),
		rpg.LfName: character.NickName,
	}.Gob()
	if err != nil {
		ctx.Err(err)
		return
	}
	eg.Req().Request(ctx.Tid(), ctx.GetId(), cmn.SvcScene, proto.CdSceneEnter,
		&pb.ReqSceneEnter{
			SceneId: character.SceneId,
			Type:    int32(cmn.ActorPlayer),
			Bytes:   bytes,
		}, func(tid int64, obj any) {
			ctx.Ok(&pb.ResCharacterEnterScene{})
		}, func(tid int64, err eg.IErr) {
			ctx.Err(err)
		})
}

func characterSelect(ctx eg.ICtx) {
	req := ctx.Body().(*pb.ReqCharacterSelect)
	var character model.Character
	err := characterColl().FindOne(context.TODO(), bson.D{
		{"cid", req.Cid},
	}, options.FindOne().SetProjection(bson.D{
		{"scene_type", 1},
	})).Decode(&character)
	if err != nil {
		ctx.Err(eg.WrapErr(eg.EcUnmarshallErr, err))
		return
	}
	if character.Id == primitive.NilObjectID {
		ctx.Ok(&pb.ResCharacterSelect{
			ErrCode: 1,
		})
		return
	}
	sceneService, e := cmn.SceneTypeToSceneService(character.SceneType)
	if e != nil {
		eg.Log().Error(e)
		return
	}
	ctx.Ok(&pb.ResCharacterSelect{})
	eg.Req().Request(0, ctx.GetId(), cmn.SvcGate, proto.CdGateSetCharacterId,
		&pb.ReqGateSetCharacterId{
			Cid: req.Cid,
		}, func(tid int64, obj any) {

		}, func(tid int64, err eg.IErr) {

		})
	eg.Ntc().NtcSetUserServiceAlias(ctx.Tid(), req.Cid, sceneService, cmn.SvcScene)
}

func characterInfo(ctx eg.ICtx) {
	//查询账号角色列表
	cur, err := characterColl().Find(context.TODO(), bson.D{
		{"uid", ctx.GetId()},
	})
	if err != nil {
		ctx.Ok(&pb.ResCharacter{
			ErrCode: eg.EcNotExistId,
		})
		return
	}
	var characters []model.Character
	err = cur.All(context.TODO(), &characters)
	if err != nil {
		ctx.Err(eg.NewErr(eg.EcServiceErr, eg.M{
			"uid": ctx.GetId(),
			"err": err.Error(),
		}))
		return
	}

	characterSlice := make([]*pb.Character, len(characters))
	for i, c := range characters {
		characterSlice[i] = &pb.Character{
			Cid:       c.Cid,
			NickName:  c.NickName,
			Gender:    0,
			Figure:    0,
			IsNovice:  false,
			SceneType: 0,
			SceneId:   0,
		}
	}

	ctx.Ok(&pb.ResCharacter{
		Characters: characterSlice})
}

func characterCreate(ctx eg.ICtx) {
	req := ctx.Body().(*pb.ReqCharacterCreate)

	// 检测敏感词
	ok, _ := cmn.SensitiveFilter().Validate(req.NickName)
	if !ok {
		ctx.Ok(&pb.ResCharacterCreate{ErrCode: 1})
		return
	}

	// 昵称长度
	if len([]rune(req.NickName)) < 2 || len([]rune(req.NickName)) > 6 {
		ctx.Ok(&pb.ResCharacterCreate{
			ErrCode: 2,
		})
		return
	}

	// 昵称性别
	if req.Gender != 0 && req.Gender != 1 {
		ctx.Ok(&pb.ResCharacterCreate{
			ErrCode: 3,
		})
		return
	}
	var character model.Character
	err := characterColl().FindOne(context.TODO(), bson.D{
		{"nick_name", req.NickName},
	}).Decode(&character)
	if err == nil {
		ctx.Ok(&pb.ResCharacterCreate{
			ErrCode: 4,
		})
		return
	}
	if character.Id != primitive.NilObjectID {
		ctx.Ok(&pb.ResCharacterCreate{
			ErrCode: 5,
		})
		return
	}

	cid := eg.SId().GetGlobalId()
	characterColl().InsertOne(context.TODO(), &model.Character{
		Id:        primitive.NewObjectID(),
		Cid:       cid,
		Uid:       ctx.GetId(),
		NickName:  req.NickName,
		Gender:    uint8(req.Gender),
		Figure:    uint8(req.Figure),
		IsNovice:  true,
		SceneType: cmn.SceneWorld,
		SceneId:   1,
	})
	ctx.Ok(&pb.ResCharacterCreate{})
}

func characterOffLine(ctx eg.ICtx) {
	eg.Log().Info("character offline", eg.M{
		"character_id": ctx.GetId(),
	})
}
