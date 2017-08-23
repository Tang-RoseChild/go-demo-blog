package accountstore

import (
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Error struct {
	error
	notFound    bool
	mobileExist bool
}

func (e *Error) NotFound() bool {
	return e.notFound
}

func (e *Error) MobileExist() bool {
	return e.mobileExist
}
func (e *Error) Error() string {
	var str string
	switch {
	case e.notFound:
		str = "not found"
	case e.mobileExist:
		str = "mobile exists"
	default:
		str = e.error.Error()
	}
	return str
}

// 这样不适合interface的模式，相当于把interface的实现必须和这个包绑定
type CreateAccountReq struct {
	No       string
	Mobile   string
	Password string
	IP       string
}
type AccountService interface {
	GetAccountByID(id string) (*Account_V2, *Error)
	GetAccountByMobile(no, mobile string) (*Account_V2, *Error)
	CreateAccount(req *CreateAccountReq) *Account_V2
}
type Passworder interface {
	Validate(plainPwd string) bool
}

type Service struct {
	AccountService
	Passworder
}

func NewService() *Service {
	return &Service{
		AccountService: &Service_V2{},
	}
}
