package service

import "github.com/Tang-RoseChild/go-demo-blog/account/store"

type AccountService interface {
	GetAccountByID(id string) (store.Account, error)
	GetAccountByMobile(no, mobile string) (store.Account, error)
	CreateAccount(no, mobile string, password string, ip string) (store.Account, error)
}
type Passworder interface {
	GetPassword() string
}
type Service struct {
	AccountService
}

type Options struct {
}

func New(opts Options) *Service {
	return &Service{
		AccountService: &store.Service{},
	}
}
