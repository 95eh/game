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
)

func InitCharacter(conf cmn.CharacterConf) {
	proto.InitCharacterCodec()

	svc := eg.Svc().BindService(cmn.SvcCharacter)
	{
		svc.Bind(proto.CdCharacterInfo, characterInfo)
		svc.Bind(proto.CdCharacterCreate, characterCreate)
		svc.Bind(proto.CdCharacterSelect, characterSelect)
		svc.Bind(proto.CdCharacterEnterScene, characterEnterScene)
		svc.Bind(proto.CdCharacterChangeScene, characterChangeScene)
		svc.Bind(proto.CdCharacterExitScene, characterExitScene)
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
	// todo 判断能否退出当前的场景
	ctx.Ok(&pb.ResCharacterExitScene{})
}

func characterChangeScene(ctx eg.ICtx) {
	ctx.Ok(&pb.ResCharacterChangeScene{})
	//req := ctx.Body().(*pb.ReqCharacterChangeScene)
	// todo 判断能否进入场景
	// todo 退出当前的场景
}

func characterEnterScene(ctx eg.ICtx) {
	sceneId := int64(1)
	//todo 获取角色数据
	bytes, err := eg.M{
		rpg.TnfPos: eg.NewVec3(100, 0, 100),
		rpg.LfName: "95eh",
	}.Gob()
	if err != nil {
		ctx.Err(err)
	}
	eg.Svc().Request(ctx.Tid(), ctx.GetId(), cmn.SvcScene, proto.CdSceneEnter,
		&pb.ReqSceneEnter{
			SceneId: sceneId,
			Type:    int32(cmn.ActorPlayer),
			Bytes:   bytes,
		}, func(tid int64, obj interface{}) {
			ctx.Ok(&pb.ResCharacterEnterScene{})
		}, func(tid int64, err eg.IErr) {
			ctx.Err(err)
		})
}

func characterSelect(ctx eg.ICtx) {
	req := ctx.Body().(*pb.ReqCharacterSelect)
	id, err := primitive.ObjectIDFromHex(req.Cid)
	if err != nil {
		ctx.Ok(&pb.ResCharacterSelect{
			ErrCode: 1,
		})
		return
	}
	var t model.Character
	err = characterColl().FindOne(context.TODO(), bson.D{
		{"_id", id},
	}).Decode(t)
	if err != nil {
		ctx.Ok(&pb.ResCharacterSelect{
			ErrCode: 2,
		})
		return
	}
	if t.Id == primitive.NilObjectID {
		ctx.Ok(&pb.ResCharacterSelect{
			ErrCode: 3,
		})
		return
	}
	ctx.Ok(&pb.ResCharacterSelect{})
	sceneService, iErr := cmn.SceneTypeToSceneService(t.SceneType)
	if iErr != nil {
		eg.Log().Error(iErr)
		return
	}
	eg.Svc().Request(0, ctx.GetId(), cmn.SvcGate, proto.CdGateSetCharacterId,
		&pb.ReqGateSetCharacterId{
			Cid: req.Cid,
		}, func(tid int64, obj interface{}) {

		}, func(tid int64, err eg.IErr) {

		})
	eg.Ntc().NtcSetUserServiceAlias(ctx.Tid(), req.Cid, sceneService, cmn.SvcScene)
}

func characterInfo(ctx eg.ICtx) {
	//查询账号角色列表
	characters := make([]model.Character, 0)
	cur, err := characterColl().Find(context.TODO(), bson.D{
		{"uid", ctx.GetId()},
	})
	if err != nil {
		ctx.Ok(&pb.ResCharacter{
			ErrCode: eg.EcNotExistId,
		})
		return
	}
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
			Id:        c.Id.Hex(),
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
	character := &model.Character{}
	err := characterColl().FindOne(context.TODO(), bson.D{
		{"nickname", req.NickName},
	}).Decode(character)
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

	characterColl().InsertOne(context.TODO(), &model.Character{
		Id:       primitive.NewObjectID(),
		Uid:      ctx.GetId(),
		NickName: req.NickName,
		Gender:   uint8(req.Gender),
		Figure:   uint8(req.Figure),
		IsNovice: true,
	})
	ctx.Ok(&pb.ResCharacterCreate{})
}

func characterOffLine(ctx eg.ICtx) {
	eg.Log().Info("character offline", eg.M{
		"character_id": ctx.GetId(),
	})
}
