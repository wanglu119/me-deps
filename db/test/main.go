package main

import (
	"fmt"

	depDbConfig "github.com/wanglu119/me-deps/db/common/config"
	depOrm "github.com/wanglu119/me-deps/db/orm"
)

func init() {
	depDbConfig.Database.Type = "sqlite3"
	depDbConfig.Database.Path = "./bin/test.db"

	/* test mysql
	depDbConfig.Database.Type = "mysql"
	depDbConfig.Database.Host = "10.130.17.177:3306"
	depDbConfig.Database.User = "root"
	depDbConfig.Database.Password = "root"
	depDbConfig.Database.Name = "test"
	*/

	Init()
}

func Init() {
	log.Info(fmt.Sprintf("database type: %s", depDbConfig.Database.Type))

	if depDbConfig.Database.Type == "mongo" {
		panic("not support")
	} else {
		depOrm.AddTable(new(Test))
		err := depOrm.NewEngine()
		if err != nil {
			panic(err)
		}

		Tests = &tests{
			CommonOp: &depOrm.Common{},
		}
	}
}

func main() {
	bean := &Test{
		Username: "xxx",
	}
	Tests.InsertOne(bean)

	res := []*Test{}
	Tests.Find(&res, &Test{})
	for i, r := range res {
		log.Info(i, r.Id, r.Username)
	}

	Tests.DeleteMany(bean)
	Tests.DeleteMany(&Test{
		Username: "xxx",
	})
}
