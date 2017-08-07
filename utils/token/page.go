package tokenutils

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Pagination struct {
	Scope   string
	From    int
	Limit   int
	HasMore bool
	*jwt.StandardClaims
}

func NewPagination(scope string, from, limit int) *Pagination {
	return &Pagination{
		Scope: scope,
		From:  from,
		Limit: limit,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: time.Now().Add(2 * time.Minute).Unix(),
		},
	}
}
func (p *Pagination) GenerateToken(secret []byte) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, p)
	ss, _ := token.SignedString(secret)
	return "Bear " + ss
}

func (p *Pagination) Parse(tokenStr string, secret []byte) error {
	token, err := jwt.ParseWithClaims(tokenStr, p, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err == nil {
		if claim, ok := token.Claims.(*Pagination); ok && token.Valid {
			*p = *claim
			return nil
		}
	}

	return errors.New("invalid token")
}
