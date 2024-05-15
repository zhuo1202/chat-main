

package chat

import (
	"context"
	"github.com/openimsdk/tools/db/mongoutil"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"

	"github.com/openimsdk/chat/pkg/common/db/table/chat"
	"github.com/openimsdk/tools/errs"
)

func NewRegister(db *mongo.Database) (chat.RegisterInterface, error) {
	coll := db.Collection("register")
	_, err := coll.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "user_id", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
	})
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return &Register{coll: coll}, nil
}

type Register struct {
	coll *mongo.Collection
}

func (o *Register) Create(ctx context.Context, registers ...*chat.Register) error {
	return mongoutil.InsertMany(ctx, o.coll, registers)
}

func (o *Register) CountTotal(ctx context.Context, before *time.Time) (int64, error) {
	filter := bson.M{}
	if before != nil {
		filter["create_time"] = bson.M{"$lt": before}
	}
	return mongoutil.Count(ctx, o.coll, filter)
}
