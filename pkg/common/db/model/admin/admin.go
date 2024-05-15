

package admin

import (
	"context"
	"github.com/openimsdk/chat/pkg/common/constant"
	"github.com/openimsdk/chat/pkg/common/db/table/admin"
	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/db/pagination"
	"github.com/openimsdk/tools/errs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewAdmin(db *mongo.Database) (admin.AdminInterface, error) {
	coll := db.Collection("admin")
	_, err := coll.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{
			{Key: "account", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return &Admin{
		coll: coll,
	}, nil
}

type Admin struct {
	coll *mongo.Collection
}

func (o *Admin) Take(ctx context.Context, account string) (*admin.Admin, error) {
	return mongoutil.FindOne[*admin.Admin](ctx, o.coll, bson.M{"account": account})
}

func (o *Admin) TakeUserID(ctx context.Context, userID string) (*admin.Admin, error) {
	return mongoutil.FindOne[*admin.Admin](ctx, o.coll, bson.M{"user_id": userID})
}

func (o *Admin) Update(ctx context.Context, account string, update map[string]any) error {
	if len(update) == 0 {
		return nil
	}
	return mongoutil.UpdateOne(ctx, o.coll, bson.M{"user_id": account}, bson.M{"$set": update}, false)
}

func (o *Admin) Create(ctx context.Context, admins []*admin.Admin) error {
	return mongoutil.InsertMany(ctx, o.coll, admins)
}

func (o *Admin) ChangePassword(ctx context.Context, userID string, newPassword string) error {
	return mongoutil.UpdateOne(ctx, o.coll, bson.M{"user_id": userID}, bson.M{"$set": bson.M{"password": newPassword}}, false)

}

func (o *Admin) Delete(ctx context.Context, userIDs []string) error {
	if len(userIDs) == 0 {
		return nil
	}
	return mongoutil.DeleteMany(ctx, o.coll, bson.M{"user_id": bson.M{"$in": userIDs}})
}

func (o *Admin) Search(ctx context.Context, pagination pagination.Pagination) (int64, []*admin.Admin, error) {
	opt := options.Find().SetSort(bson.D{{"create_time", -1}})
	filter := bson.M{"level": constant.NormalAdmin}
	return mongoutil.FindPage[*admin.Admin](ctx, o.coll, filter, pagination, opt)
}
