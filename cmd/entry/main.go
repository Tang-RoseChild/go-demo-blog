package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Tang-RoseChild/go-demo-blog/account"
	"github.com/Tang-RoseChild/go-demo-blog/account/service"
	"github.com/Tang-RoseChild/go-demo-blog/blog"
	"github.com/Tang-RoseChild/go-demo-blog/comments"
	"github.com/Tang-RoseChild/go-demo-blog/utils/db"
	_ "github.com/Tang-RoseChild/go-demo-blog/utils/db/sqlite3"
)

var defaultService = service.New(service.Options{})
var addr string

func init() {
	flag.StringVar(&addr, "addr", ":9991", "addr")
}

func main() {
	flag.Parse()

	dbutils.Connector.Connect()
	// dbutils.DB.LogMode(true) // for testing

	engine := NewGinEngine()
	engine.StaticFS("/", http.Dir("./static"))
	LoadHandler(engine)

	fmt.Println("listen ..... " + addr)
	if err := engine.Run(addr); err != nil {
		panic(err)
	}

}

func NewGinEngine() *gin.Engine {
	return gin.Default()
}

func LoadHandler(engine *gin.Engine) {
	r := engine.Group("/api")
	account.GinLoad(r)
	blog.GinLoad(r)
	comments.GinLoad(r)
}
