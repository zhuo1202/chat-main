

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

func NewIPForbidden(db *mongo.Database) (admin.IPForbiddenInterface, error) {
	coll := db.Collection("ip_forbidden")
	_, err := coll.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{
			{Key: "ip", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return &IPForbidden{
		coll: coll,
	}, nil
}

type IPForbidden struct {
	coll *mongo.Collection
}

func (o *IPForbidden) Take(ctx context.Context, ip string) (*admin.IPForbidden, error) {
	return mongoutil.FindOne[*admin.IPForbidden](ctx, o.coll, bson.M{"ip": ip})

}

func (o *IPForbidden) Find(ctx context.Context, ips []string) ([]*admin.IPForbidden, error) {
	return mongoutil.Find[*admin.IPForbidden](ctx, o.coll, bson.M{"ip": bson.M{"$in": ips}})

}

func (o *IPForbidden) Search(ctx context.Context, keyword string, state int32, pagination pagination.Pagination) (int64, []*admin.IPForbidden, error) {
	filter := bson.M{}

	switch state {
	case constant.LimitNil:
	case constant.LimitEmpty:
		filter = bson.M{"limit_register": 0, "limit_login": 0}
	case constant.LimitOnlyRegisterIP:
		filter = bson.M{"limit_register": 1, "limit_login": 0}
	case constant.LimitOnlyLoginIP:
		filter = bson.M{"limit_register": 0, "limit_login": 1}
	case constant.LimitRegisterIP:
		filter = bson.M{"limit_register": 1}
	case constant.LimitLoginIP:
		filter = bson.M{"limit_login": 1}
	case constant.LimitLoginRegisterIP:
		filter = bson.M{"limit_register": 1, "limit_login": 1}
	}

	if keyword != "" {
		filter["$or"] = []bson.M{
			{"ip": bson.M{"$regex": keyword, "$options": "i"}},
		}
	}
	return mongoutil.FindPage[*admin.IPForbidden](ctx, o.coll, filter, pagination)
}

func (o *IPForbidden) Create(ctx context.Context, ms []*admin.IPForbidden) error {
	return mongoutil.InsertMany(ctx, o.coll, ms)
}

func (o *IPForbidden) Delete(ctx context.Context, ips []string) error {
	if len(ips) == 0 {
		return nil
	}
	return mongoutil.DeleteMany(ctx, o.coll, bson.M{"ip": bson.M{"$in": ips}})
}
