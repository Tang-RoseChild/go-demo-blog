package accountstore

import (
	"testing"

	"github.com/dgrijalva/jwt-go"
)

func TestToken(t *testing.T) {
	ac := &AccountClaim{
		StandardClaims: jwt.StandardClaims{
			Id: "test_id",
		},
	}
	secret := []byte("不告诉你")
	tokenStr, err := ac.GenerateToken(secret)
	if err != nil {
		t.Error(err)
	}
	err = ac.Parse(tokenStr, secret)
	if err != nil {
		t.Error(err)
	}
	if ac.Id != "test_id" {
		t.Error("id should eq")
	}

	ac1 := &AccountClaim{}
	err = ac1.Parse(tokenStr, secret)
	if err != nil {
		t.Error(err)
	}
	if ac1.Id != "test_id" {
		t.Error("id should eq")
	}
}
