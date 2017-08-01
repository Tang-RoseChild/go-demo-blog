package postgres

import (
	"flag"
	"fmt"

	"github.com/Tang-RoseChild/go-demo-blog/utils/db"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var dsn string
var dbHost string
var dbUser string
var dbName string
var dbPassword string

func init() {
	flag.StringVar(&dbHost, "db-host", "localhost", "db host")
	flag.StringVar(&dbUser, "db-user", "demo", "")
	flag.StringVar(&dbName, "db-name", "demo", "")
	flag.StringVar(&dbPassword, "db-password", "demo123", "")
	dbutils.Connector = &connector{}
}

type connector struct{}

func (c *connector) Connect() {
	var err error
	dbutils.DB, err = gorm.Open("postgres", fmt.Sprintf("host=%s user=%s dbname=%s password=%s sslmode=disable", dbHost, dbUser, dbName, dbPassword))
	if err != nil {
		panic(err)
	}
}
