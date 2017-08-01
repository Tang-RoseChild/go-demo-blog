package store

import (
	"time"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Account interface {
	GetMobile() (no, mobile string)
	GetRegTime() time.Time
	GetLastLoginTime() time.Time
	GetID() string
	GetStatus() int

	SetStatus(status int)
	Marshal() (typ, content string) // json,xml ...
}
