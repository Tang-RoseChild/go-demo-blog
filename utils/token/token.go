package tokenutils

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	expireSeconds = int64(600)
)

func GenToken(id string) string {
	claim := NewStandClaim()
	claim.ID = id
	claim.IssueAt = uint64(time.Now().Unix())
	// claim.ExpireAt = uint64(time.Now().Unix() + 600)
	claim.ExpireAt = uint64(time.Now().Unix() + expireSeconds)

	return "bear " + claim.SignedToString(hmacKey)
}

func IssueToken(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHead := r.Header.Get("Authorization")
		fields := strings.Fields(authHead)
		fmt.Println("fields", fields, " authHead >>> ", authHead)
		if len(fields) != 2 || strings.ToLower(fields[0]) != "bear" {

			http.Error(w, "invalid token: request head", http.StatusBadRequest)
			return
		}
		claim := NewStandClaim()

		if err := claim.Decode(fields[1], hmacKey); err != nil {
			fmt.Printf("claim >>> %#v \n", claim)
			http.Error(w, "invalid token: "+err.Error(), http.StatusBadRequest)
			return
		}
		// claim.ExpireAt = uint64(time.Now().Unix() + int64(600))
		claim.ExpireAt = uint64(time.Now().Unix() + int64(expireSeconds))
		w.Header().Set("Authorization", "bear "+claim.SignedToString(hmacKey))
		r = r.WithContext(context.WithValue(r.Context(), "uid", claim.ID))

		if handler != nil {
			handler(w, r)
		}

	}
}

var secret []byte

func init() {
	key, err := base64.StdEncoding.DecodeString("QTJERDNzZjNEM2FzODNLOXNr") // it's not used, just a demo
	if err != nil {
		panic(err)
	}

	secret = key
}
func hmacKey() []byte {
	return secret
}

func NewStandClaim() *StandClaim {
	return &StandClaim{}
}

type Tokener interface {
	Validate() bool
}

type StandClaim struct {
	ID       string
	IssueAt  uint64
	ExpireAt uint64
}

func (c *StandClaim) SignedToString(keyFn func() []byte) string {
	payload, err := json.Marshal(c)
	if err != nil {
		panic(err)
	}

	payloadSum := sha256.Sum256(payload)

	mac := hmac.New(sha256.New, keyFn())
	mac.Write(payloadSum[:])

	return base64.StdEncoding.EncodeToString(payload) + "." + base64.StdEncoding.EncodeToString(mac.Sum(nil))
}
func (c *StandClaim) Decode(token string, keyFn func() []byte) error {
	splits := strings.Split(token, ".")
	if len(splits) != 2 {
		return errors.New("invalid token: stype wrong")
	}

	payloadData, err := base64.StdEncoding.DecodeString(splits[0])
	if err != nil {
		return err
	}

	err = json.Unmarshal(payloadData, c)
	if err != nil {
		return err
	}

	mac := hmac.New(sha256.New, keyFn())
	payloadSum := sha256.Sum256(payloadData)
	mac.Write(payloadSum[:])
	expectedMAC := mac.Sum(nil)

	sig := base64.StdEncoding.EncodeToString(expectedMAC)
	// fmt.Printf("%s  \n expectedMAC %s \n", splits[1], sig)
	// if base64.StdEncoding.EncodeToString(expectedMAC) != splits[1] {
	if sig != splits[1] {
		return errors.New("invalid token: signature invalid")
	}
	err = json.Unmarshal(payloadData, c)
	if err != nil {
		return err
	}
	if !c.Validate() {
		return errors.New("invalid token:validate failed")
	}

	return nil
}

func (c *StandClaim) Validate() bool {
	now := uint64(time.Now().Unix())
	if c.ExpireAt < c.IssueAt || c.IssueAt > now || c.ExpireAt < now {
		return false
	}
	return true
}

var initSecret []byte

func SetSecret(secret string) {
	initSecret = []byte(secret)
}

func GetSecret() []byte {
	return initSecret
}
