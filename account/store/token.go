package accountstore

import (
	"errors"

	"github.com/dgrijalva/jwt-go"
)

var (
	invalidTokenError = errors.New("invalid token")
)

func NewClaim() *AccountClaim {
	return &AccountClaim{}
}

type AccountClaim struct {
	jwt.StandardClaims
}

func (ac *AccountClaim) GenerateToken(secret []byte) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, ac)
	ss, _ := token.SignedString(secret)
	return "Bear " + ss
}

func (ac *AccountClaim) Parse(tokenStr string, secret []byte) error {
	token, err := jwt.ParseWithClaims(tokenStr, ac, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err == nil {
		if claim, ok := token.Claims.(*AccountClaim); ok && token.Valid {
			*ac = *claim
			return nil
		}
	}

	return invalidTokenError
}
