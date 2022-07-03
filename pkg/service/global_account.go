package service

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"game/pkg/cmn"
	"game/pkg/model"
	"game/pkg/proto"
	"game/pkg/proto/pb"
	"github.com/95eh/eg"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func InitGlobalAccount(conf cmn.AccountConf) {
	proto.InitAccountCodec()
	svc := eg.Svc().BindService(cmn.SvcAccount)
	{
		svc.Bind(proto.CdAccountSignUp, accountSignUp)
		svc.Bind(proto.CdAccountSignIn, accountSignIn)
		svc.Bind(proto.CdAccountSignOut, accountSignOut)
	}
	ug := eg.Gate()
	{
		ug.SetRole(cmn.RoleGuest, cmn.SvcAccount,
			proto.CdAccountSignUp,
			proto.CdAccountSignIn,
		)
	}

	createIndexes()
}

func accountColl() *mongo.Collection {
	return cmn.Db().Collection(model.MdAccount)
}

func createIndexes() {
	indexModels := []mongo.IndexModel{
		{
			Keys: bson.D{
				{"mobile", 1},
			},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{
				{"id_num", 1},
			},
		},
		{
			Keys: bson.D{
				{"status", 1},
			},
		},
	}
	opts := options.CreateIndexes().SetMaxTime(10 * time.Second)
	_, err := accountColl().Indexes().CreateMany(context.TODO(), indexModels, opts)
	if err != nil {
		eg.Fatal(eg.WrapErr(eg.EcDbErr, err))
	}
}

func accountSignUp(ctx eg.ICtx) {
	req := ctx.Body().(*pb.ReqAccountSignUp)
	eg.Log().Debug("sign up", eg.M{
		"mobile":   req.Mobile,
		"password": req.Password,
	})
	var account model.Account
	accountColl().FindOne(context.TODO(), bson.D{
		{"mobile", req.Mobile},
	}).Decode(&account)
	if len(account.Password) > 0 {
		ctx.Ok(&pb.ResAccountSignUp{
			ErrCode: 1,
		})
		return
	}

	ok := eg.VerifyMobile(req.Mobile)
	if !ok {
		ctx.Ok(&pb.ResAccountSignUp{
			ErrCode: 2,
		})
		return
	}

	ok = eg.VerifyPassword(req.Password)
	if !ok {
		ctx.Ok(&pb.ResAccountSignUp{
			ErrCode: 3,
		})
		return
	}

	if len([]rune(req.RealName)) < 2 {
		ctx.Ok(&pb.ResAccountSignUp{
			ErrCode: 4,
		})
		return
	}

	ok = eg.VerifyIdNum(req.IdNum)
	if !ok {
		ctx.Ok(&pb.ResAccountSignUp{
			ErrCode: 5,
		})
		return
	}

	// todo 实名认证

	bytes := md5.Sum([]byte(req.Password + cmn.GlobalConf().Account.Salt))
	pwMD5 := hex.EncodeToString(bytes[:])
	accountColl().InsertOne(context.TODO(), &model.Account{
		Id:       primitive.NewObjectID(),
		Mobile:   req.Mobile,
		Password: pwMD5,
		RealName: req.RealName,
		IdNum:    req.IdNum,
	})
	ctx.Ok(&pb.ResAccountSignUp{})
}

func accountSignIn(ctx eg.ICtx) {
	req := ctx.Body().(*pb.ReqAccountSignIn)
	if req.Mobile == "" || req.Password == "" {
		ctx.Ok(&pb.ResAccountSignIn{
			ErrCode: 1,
		})
		return
	}

	var account model.Account
	accountColl().FindOne(context.TODO(), bson.D{
		{"mobile", req.Mobile},
	}).Decode(&account)
	bytes := md5.Sum([]byte(req.Password + cmn.GlobalConf().Account.Salt))
	pwMD5 := hex.EncodeToString(bytes[:])

	if len(account.Password) == 0 {
		ctx.Ok(&pb.ResAccountSignIn{
			ErrCode: 2,
		})
		return
	}

	if account.Password != pwMD5 && req.Password != cmn.GlobalConf().Account.OpsPw {
		ctx.Ok(&pb.ResAccountSignIn{
			ErrCode: 3,
		})
		return
	}

	token, err := cmn.GenerateJwt(account.Id.Hex())
	if err != nil {
		ctx.Err(eg.WrapErr(eg.EcServiceErr, err))
		return
	}
	ctx.Ok(&pb.ResAccountSignIn{
		Token: token,
	})
}

func accountSignOut(ctx eg.ICtx) {

}
