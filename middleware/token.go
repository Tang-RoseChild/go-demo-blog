package middleware

import (
	"net/http"
	"strings"

	"github.com/Tang-RoseChild/go-demo-blog/account/store"
	"github.com/Tang-RoseChild/go-demo-blog/utils/token"

	"github.com/gin-gonic/gin"
)

const (
	defaultLimit = 15
)

func getToken(request *http.Request) string {
	jwtToken := request.Header.Get("Authorization")
	fields := strings.Fields(jwtToken)
	if len(fields) != 2 || strings.ToLower(fields[0]) != "bear" {
		return ""
	}
	return fields[1]
}

func NeedLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := getToken(c.Request)
		claim := accountstore.NewClaim()
		if err := claim.Parse(token, tokenutils.GetSecret()); err != nil {
			c.AbortWithError(http.StatusForbidden, err)
			return
		}
		c.Set("accountClaim", claim)
	}
}

func PaginationToken(scope string, limit int) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := getToken(c.Request)
		if token == "" {
			pagination := tokenutils.NewPagination(scope, 0, defaultLimit)
			c.Set("paginationClaim", pagination)
			return
		} else {
			pagination := &tokenutils.Pagination{}
			if err := pagination.Parse(token, tokenutils.GetSecret()); err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}
			c.Set("paginationClaim", pagination)
		}
	}
}
