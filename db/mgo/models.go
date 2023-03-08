package mgo

import (
	"sync"
	"context"
	
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	
	"github.com/wanglu119/me-deps/db/common/config"
)

type MTable interface {
	CreateIndex()
}

var (
	once sync.Once
	client *mongo.Client
	ctx context.Context
	cancelFunc context.CancelFunc
	legacyTables []interface{}
)

func AddTable(table ...interface{}) {
	for _, t := range table {
		legacyTables = append(legacyTables, t)
	}
}

func getEngine() (*mongo.Client, error) {
	var err error
	once.Do(func() {
		ctx, cancelFunc = context.WithCancel(context.Background())
		clientOpt := options.Client()
		clientOpt.Hosts = []string{config.Database.Host}
		client ,err = mongo.Connect(ctx, clientOpt)
	})
	
	return client, err
}

func NewEngine() (err error) {
	if _, err = getEngine(); err != nil {
		return err
	}
	
	for _, m := range legacyTables {
		m.(MTable).CreateIndex()
	}
	
	return nil
}

func GetEngine() *mongo.Client {
	return client
}





