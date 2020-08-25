package orm

import (
	"os"
	"path"
	"fmt"
	"time"
	
	_ "github.com/mattn/go-sqlite3"
	"xorm.io/xorm"
	
	"github.com/wanglu119/me-deps/db/common/config"
)

var (
	x *xorm.Engine
	legacyTables []interface{}
)

func AddTable(table ...interface{}) {
	for _, t := range table {
		legacyTables = append(legacyTables, t)
	}
}

func Init() error {
	for _, t := range legacyTables {
		err := x.Sync2(t)
		if err != nil {
			return err
		}
	}
	
	return nil
}

func getEngine() (*xorm.Engine, error) {
	
	connStr := ""
	switch config.Database.Type {
		case "sqlite3":
			if err := os.MkdirAll(path.Dir(config.Database.Path), os.ModePerm); err != nil {
				return nil, fmt.Errorf("create directories: %v", err)
			}
			connStr = "file:" + config.Database.Path + "?cache=shared&mode=rwc"
			
		default:
			return nil, fmt.Errorf("unknow database type: %s", config.Database.Type)
	}
	return xorm.NewEngine(config.Database.Type, connStr)
}

func NewEngine() (err error) {
	x, err = getEngine()
	if err != nil {
		return fmt.Errorf("connect to database: %v", err)
	}
	
	x.SetConnMaxLifetime(time.Second)
	
	x.ShowSQL(true)
	
	err = Init()
	
	if err != nil {
		return err
	}
	
	return nil
}

func GetEngine() *xorm.Engine {
	return x
}




