package orm

import (
	"fmt"
	"os"
	"path"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"xorm.io/xorm"

	"github.com/wanglu119/me-deps/db/common/config"
	"github.com/wanglu119/me-deps/log"
)

var (
	x            *xorm.Engine
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
	if len(config.Database.Charset) == 0 {
		config.Database.Charset = "utf8"
	}
	switch config.Database.Type {
	case "sqlite3":
		if err := os.MkdirAll(path.Dir(config.Database.Path), os.ModePerm); err != nil {
			return nil, fmt.Errorf("create directories: %v", err)
		}
		connStr = "file:" + config.Database.Path + "?cache=shared&mode=rwc"
	case "mysql":
		//  "root:root@tcp(127.0.0.1:3306)/xorm?charset=utf8"
		connStr = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&loc=Local",
			config.Database.User, config.Database.Password,
			config.Database.Host, config.Database.Name, config.Database.Charset)
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

	if log.DEBUG {
		x.ShowSQL(true)
	}
	x.SetTZDatabase(time.Local)

	err = Init()

	if err != nil {
		return err
	}

	return nil
}

func GetEngine() *xorm.Engine {
	return x
}
