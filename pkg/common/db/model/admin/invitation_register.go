

package admin

import (
	"context"
	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/db/pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/openimsdk/chat/pkg/common/constant"
	"github.com/openimsdk/chat/pkg/common/db/table/admin"
	"github.com/openimsdk/tools/errs"
)

func NewInvitationRegister(db *mongo.Database) (admin.InvitationRegisterInterface, error) {
	coll := db.Collection("invitation_register")
	_, err := coll.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{
			{Key: "invitation_code", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return &InvitationRegister{
		coll: coll,
	}, nil
}

type InvitationRegister struct {
	coll *mongo.Collection
}

func (o *InvitationRegister) Find(ctx context.Context, codes []string) ([]*admin.InvitationRegister, error) {
	return mongoutil.Find[*admin.InvitationRegister](ctx, o.coll, bson.M{"invitation_code": bson.M{"$in": codes}})
}

func (o *InvitationRegister) Del(ctx context.Context, codes []string) error {
	if len(codes) == 0 {
		return nil
	}
	return mongoutil.DeleteMany(ctx, o.coll, bson.M{"invitation_code": bson.M{"$in": codes}})
}

func (o *InvitationRegister) Create(ctx context.Context, v []*admin.InvitationRegister) error {
	return mongoutil.InsertMany(ctx, o.coll, v)
}

func (o *InvitationRegister) Take(ctx context.Context, code string) (*admin.InvitationRegister, error) {
	return mongoutil.FindOne[*admin.InvitationRegister](ctx, o.coll, bson.M{"code": code})
}

func (o *InvitationRegister) Update(ctx context.Context, code string, data map[string]any) error {
	if len(data) == 0 {
		return nil
	}
	return mongoutil.UpdateOne(ctx, o.coll, bson.M{"invitation_code": code}, bson.M{"$set": data}, false)
}

func (o *InvitationRegister) Search(ctx context.Context, keyword string, state int32, userIDs []string, codes []string, pagination pagination.Pagination) (int64, []*admin.InvitationRegister, error) {
	filter := bson.M{}
	switch state {
	case constant.InvitationCodeUsed:
		filter = bson.M{"user_id": bson.M{"$ne": ""}}
	case constant.InvitationCodeUnused:
		filter = bson.M{"user_id": ""}
	}

	if len(userIDs) > 0 {
		filter["user_id"] = bson.M{"$in": userIDs}
	}
	if len(codes) > 0 {
		filter["invitation_code"] = bson.M{"$in": codes}
	}
	if keyword != "" {
		filter["$or"] = []bson.M{
			{"invitation_code": bson.M{"$regex": keyword, "$options": "i"}},
			{"user_id": bson.M{"$regex": keyword, "$options": "i"}},
		}
	}
	return mongoutil.FindPage[*admin.InvitationRegister](ctx, o.coll, filter, pagination)
}
