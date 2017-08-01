package dbutils

import "github.com/jinzhu/gorm"

var DB *gorm.DB // 用 sql.DB 来替代: 只要支持sql.DB的，都可以替换，便于移植。缺点sql.DB更像底层的，需要自己写sql等
var Connector DBConnector

type DBConnector interface {
	Connect()
}
