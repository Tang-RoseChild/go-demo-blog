package sqlite3

import (
	"flag"
	"fmt"

	"github.com/Tang-RoseChild/go-demo-blog/utils/db"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var dbName string

func init() {
	flag.StringVar(&dbName, "db-path", "/tmp/sqlite3.db", "sqlite db file path ")
	dbutils.Connector = &connector{}
}

type connector struct{}

func (c *connector) Connect() {
	fmt.Println("connet in sqlite3")
	var err error
	dbutils.DB, err = gorm.Open("sqlite3", dbName)
	if err != nil {
		panic(err)
	}
}
