package mgo

import (
    "time"
    "strings"
    "fmt"
    "reflect"
    "regexp"
    
    "go.mongodb.org/mongo-driver/x/bsonx"
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
	indexOpts := options.CreateIndexes().SetMaxTime(10*time.Second)
	
	indexView := col.Indexes()
	keysDoc := bsonx.Doc{}
	
	for _, key := range keys {
		if strings.HasPrefix(key, "-") {
			keysDoc = keysDoc.Append(strings.TrimLeft(key, "-"), bsonx.Int32(-1))	
		} else {
			keysDoc = keysDoc.Append(key, bsonx.Int32(1))
		}
	}
	
	result, err := indexView.CreateOne(ctx, 
		mongo.IndexModel {
			Keys: keysDoc,
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
	indexOpts := options.CreateIndexes().SetMaxTime(10*time.Second)
	
	indexView := col.Indexes()
	keysDoc := bsonx.Doc{}
	
	for _, key := range keys {
		if strings.HasPrefix(key, "-") {
			keysDoc = keysDoc.Append(strings.TrimLeft(key, "-"), bsonx.Int32(-1))	
		} else {
			keysDoc = keysDoc.Append(key, bsonx.Int32(1))
		}
	}
	
	result, err := indexView.CreateOne(ctx, 
		mongo.IndexModel {
			Keys: keysDoc,
			Options: options.Index(),
		},
		indexOpts,
	)
	if result == "" || err != nil {
		log.Error(fmt.Sprintf("EnsureIndex error: %v", err))
	}
}


