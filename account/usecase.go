package account

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/Tang-RoseChild/go-demo-blog/account/service"
	"github.com/Tang-RoseChild/go-demo-blog/utils/http"
	"github.com/Tang-RoseChild/go-demo-blog/utils/token"
)

var defaultService = service.New(service.Options{})

func GinLoad(rootGroup *gin.RouterGroup) {
	accountGroup := rootGroup.Group("/account")
	accountGroup.POST("/login", httpToGin(LoginHandler))
	// accountGroup.POST("/register", httpToGin(RegisterHandler))
	// accountGroup.GET("/issue_token", httpToGin(IssueTokenHandler))
}

func Load() {

	http.HandleFunc("/api/login", httputils.IPLimit(10, httputils.LimitReq(5, LoginHandler)))
	http.HandleFunc("/api/register", RegisterHandler)
	http.HandleFunc("/api/issue_token", IssueTokenHandler)
}

func GinHandler(c *gin.Context) {
	LoginHandler(c.Writer, c.Request)
}

func httpToGin(h http.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		h(c.Writer, c.Request)
	}
}

// http
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		No       string
		Mobile   string
		Password string
	}
	httputils.MustUnmarshalReq(r, &req)
	req.No = "86"
	account, err := defaultService.GetAccountByMobile(req.No, req.Mobile)
	if err != nil {
		http.Error(w, "not found", http.StatusBadRequest)
		return
	}
	pwd, ok := account.(service.Passworder)
	if !ok {
		panic("not implement password")
	}
	if pwd.GetPassword() != req.Password {
		http.Error(w, "password wrong", http.StatusForbidden)
		return
	}

	w.Header().Set("Authorization", tokenutils.GenToken(account.GetID()))
	httputils.MustMarshalResp(w, map[string]interface{}{
		"account": account,
		"success": true,
	})
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		No       string
		Mobile   string
		Password string
	}

	httputils.MustUnmarshalReq(r, &req)
	fmt.Println("req in reg ", req)
	if _, err := defaultService.GetAccountByMobile(req.No, req.Mobile); err == nil {
		http.Error(w, "mobile exists", http.StatusBadRequest)
		return
	}

	ip := strings.Split(r.RemoteAddr, ":")[0]
	account, err := defaultService.CreateAccount(req.No, req.Mobile, req.Password, ip)
	if err != nil {
		http.Error(w, "internel error", http.StatusInternalServerError)
		return
	}

	token := tokenutils.GenToken(account.GetID())
	w.Header().Set("Authorization", token)
	httputils.MustMarshalResp(w, account)
}

func IssueTokenHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("issue totken")
	hanlder := tokenutils.IssueToken(nil)
	hanlder(w, r)
}
