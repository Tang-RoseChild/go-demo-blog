package account

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/Tang-RoseChild/go-demo-blog/account/store"
	"github.com/Tang-RoseChild/go-demo-blog/utils/http"
	"github.com/Tang-RoseChild/go-demo-blog/utils/token"
)

var defaultService = accountstore.NewService()

func GinLoad(rootGroup *gin.RouterGroup) {
	accountGroup := rootGroup.Group("/account")
	accountGroup.POST("/login", httpToGin(LoginHandler))
	accountGroup.POST("/register", httpToGin(RegisterHandler))
	// accountGroup.GET("/issue_token", httpToGin(IssueTokenHandler))
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
	if req.No == "" {
		req.No = "86"
	}

	account, err := defaultService.GetAccountByMobile(req.No, req.Mobile)
	if err != nil {
		if err.NotFound() {
			http.Error(w, "not found", http.StatusBadRequest)
			return
		}
		panic(err)
	}

	if !account.Validate(req.Password) {
		http.Error(w, "password wrong", http.StatusForbidden)
		return
	}

	claim := accountstore.NewClaim()
	claim.Id = account.ID
	w.Header().Set("Authorization", claim.GenerateToken(tokenutils.GetSecret()))

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

	if _, err := defaultService.GetAccountByMobile(req.No, req.Mobile); err == nil {
		http.Error(w, "mobile exists", http.StatusBadRequest)
		return
	}

	ip := strings.Split(r.RemoteAddr, ":")[0]
	account := defaultService.CreateAccount(&accountstore.CreateAccountReq{req.No, req.Mobile, req.Password, ip})
	token := tokenutils.GenToken(account.ID)
	w.Header().Set("Authorization", token)
	httputils.MustMarshalResp(w, account)
}

func IssueTokenHandler(w http.ResponseWriter, r *http.Request) {
	hanlder := tokenutils.IssueToken(nil)
	hanlder(w, r)
}
