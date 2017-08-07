package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Tang-RoseChild/go-demo-blog/account"
	"github.com/Tang-RoseChild/go-demo-blog/blog"
	"github.com/Tang-RoseChild/go-demo-blog/comments"
	"github.com/Tang-RoseChild/go-demo-blog/utils/db"
	_ "github.com/Tang-RoseChild/go-demo-blog/utils/db/sqlite3"
	"github.com/Tang-RoseChild/go-demo-blog/utils/token"
)

var (
	addr   = flag.String("addr", ":9991", "addr")
	secret = flag.String("secret", "就不告诉你", "secret")
)

func main() {
	flag.Parse()
	tokenutils.SetSecret(*secret)
	dbutils.Connector.Connect()
	// dbutils.DB.LogMode(true) // for testing

	engine := NewGinEngine()
	engine.StaticFS("/", http.Dir("./static"))
	LoadHandler(engine)

	fmt.Println("listen ..... " + *addr)
	if err := engine.Run(*addr); err != nil {
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
