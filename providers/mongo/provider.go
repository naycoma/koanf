package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/knadh/koanf/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mongo struct {
	Uri         string
	Database    string
	Collection  string
	Filter      bson.M // どのドキュメントを読み込むか（任意）
	ReadTimeout time.Duration
}

func Provider(uri, database, collection string, filter ...bson.M) *Mongo {
	var f bson.M
	if len(filter) > 0 && filter[0] != nil {
		f = filter[0]
	} else {
		f = bson.M{} // デフォルトは全てのドキュメントを対象
	}
	return &Mongo{
		Uri:         uri,
		Database:    database,
		Collection:  collection,
		Filter:      f,
		ReadTimeout: 10 * time.Second,
	}
}

func (m *Mongo) Read() (map[string]any, error) {
	ctx, cancel := context.WithTimeout(context.Background(), m.ReadTimeout)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(m.Uri))
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(ctx)

	coll := client.Database(m.Database).Collection(m.Collection)

	var result map[string]any
	if err := coll.FindOne(ctx, m.Filter).Decode(&result); err != nil {
		return nil, fmt.Errorf("mongodb read error: %w", err)
	}

	delete(result, "_id")

	return result, nil
}

func (m *Mongo) ReadBytes() ([]byte, error) {
	return nil, fmt.Errorf("mongo provider does not support ReadBytes")
}

var _ koanf.Provider = (*Mongo)(nil)
