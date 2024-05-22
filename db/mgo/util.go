package mgo

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/wanglu119/me-deps/db/common/config"
)

func toTableName(clsName string) string {
	reg := regexp.MustCompile("[A-Z]*[a-z0-9]+")
	src := []byte(clsName)
	sarr := reg.FindAllString(string(src), -1)
	tableName := strings.Join(sarr, "_")
	tableName = strings.ToLower(tableName)

	return tableName
}

func GetCollection(t reflect.Type) *mongo.Collection {
	tableName := toTableName(t.Name())
	coll := client.Database(config.Database.Name).Collection(tableName)

	return coll
}

func CreateUniqueIndex(t reflect.Type, keys ...string) {
	tableName := toTableName(t.Name())
	col := client.Database(config.Database.Name).Collection(tableName)
	indexOpts := options.CreateIndexes().SetMaxTime(10 * time.Second)

	indexView := col.Indexes()
	keysDoc := bson.D{}

	for _, key := range keys {
		if strings.HasPrefix(key, "-") {
			// mongo-driver v1.15.0
			keysDoc = append(keysDoc, bson.E{Key: strings.TrimLeft(key, "-"), Value: -1})
			// mongo-driver v1.11.3
			// keysDoc = keysDoc.Append(strings.TrimLeft(key, "-"), bsonx.Int32(-1))
		} else {
			// mongo-driver v1.15.0
			keysDoc = append(keysDoc, bson.E{Key: key, Value: 1})
			// mongo-driver v1.11.3
			// keysDoc = keysDoc.Append(key, bsonx.Int32(1))
		}
	}

	result, err := indexView.CreateOne(ctx,
		mongo.IndexModel{
			Keys:    keysDoc,
			Options: options.Index().SetUnique(true),
		},
		indexOpts,
	)
	if result == "" || err != nil {
		log.Error(fmt.Sprintf("EnsureIndex error: %v", err))
	}
}

func CreateIndex(t reflect.Type, keys ...string) {
	tableName := toTableName(t.Name())
	col := client.Database(config.Database.Name).Collection(tableName)
	indexOpts := options.CreateIndexes().SetMaxTime(10 * time.Second)

	indexView := col.Indexes()
	// mongo-driver v1.11.3
	// keysDoc := bsonx.Doc{}
	// mongo-driver v1.15.0
	keysDoc := bson.D{}

	for _, key := range keys {
		if strings.HasPrefix(key, "-") {
			// mongo-driver v1.15.0
			keysDoc = append(keysDoc, bson.E{Key: strings.TrimLeft(key, "-"), Value: -1})
			// mongo-driver v1.11.3
			// keysDoc = keysDoc.Append(strings.TrimLeft(key, "-"), bsonx.Int32(-1))
		} else {
			// mongo-driver v1.15.0
			keysDoc = append(keysDoc, bson.E{Key: key, Value: 1})
			// mongo-driver v1.11.3
			// keysDoc = keysDoc.Append(key, bsonx.Int32(1))
		}
	}

	result, err := indexView.CreateOne(ctx,
		mongo.IndexModel{
			Keys:    keysDoc,
			Options: options.Index(),
		},
		indexOpts,
	)
	if result == "" || err != nil {
		log.Error(fmt.Sprintf("EnsureIndex error: %v", err))
	}
}
