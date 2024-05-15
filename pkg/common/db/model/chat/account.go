

package chat

import (
	"context"
	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/errs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"

	"github.com/openimsdk/chat/pkg/common/db/table/chat"
)

func NewAccount(db *mongo.Database) (chat.AccountInterface, error) {
	coll := db.Collection("account")
	_, err := coll.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{
			{Key: "user_id", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return &Account{coll: coll}, nil
}

type Account struct {
	coll *mongo.Collection
}

func (o *Account) Create(ctx context.Context, accounts ...*chat.Account) error {
	return mongoutil.InsertMany(ctx, o.coll, accounts)
}

func (o *Account) Take(ctx context.Context, userId string) (*chat.Account, error) {
	return mongoutil.FindOne[*chat.Account](ctx, o.coll, bson.M{"user_id": userId})
}

func (o *Account) Update(ctx context.Context, userID string, data map[string]any) error {
	if len(data) == 0 {
		return nil
	}
	return mongoutil.UpdateOne(ctx, o.coll, bson.M{"user_id": userID}, bson.M{"$set": data}, false)
}

func (o *Account) UpdatePassword(ctx context.Context, userId string, password string) error {
	return mongoutil.UpdateOne(ctx, o.coll, bson.M{"user_id": userId}, bson.M{"$set": bson.M{"password": password, "change_time": time.Now()}}, false)
}
