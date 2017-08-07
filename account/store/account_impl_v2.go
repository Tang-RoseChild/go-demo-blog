package accountstore

import (
	"time"

	"github.com/Tang-RoseChild/go-demo-blog/utils/db"
	"github.com/Tang-RoseChild/go-demo-blog/utils/id"

	"github.com/jinzhu/gorm"
)

type Account_V2 struct {
	Mobile        string    `json:"mobile"`
	NO            string    `json:"no"`
	RegTime       time.Time `json:"reg_time"`
	LastLoginTime time.Time `json:"last_login_time"`
	ID            string    `json:"id"`
	Status        int       `json:"-"`
	IP            string    `json:"-"`
	Password      string    `json:"-"`
	Avatar        string    `json:"avatar"`
}

func (a *Account_V2) TableName() string { // for gorm: db name is accounts
	return "accounts"
}

func (a *Account_V2) Validate(plainPwd string) bool {
	return a.Password == plainPwd // simple one
}

type Service_V2 struct{}

func (s *Service_V2) GetAccountByMobile(no, mobile string) (*Account_V2, *Error) {
	var account Account_V2
	err := dbutils.DB.First(&account, "no = ? and mobile = ?", no, mobile).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &Error{notFound: true}
		}
		panic(err)
	}
	return &account, nil
}

func (s *Service_V2) CreateAccount(req *CreateAccountReq) *Account_V2 {
	now := time.Now().UTC()
	account := &Account_V2{
		ID:            idutils.DefaultGenerator.GetID(),
		NO:            req.No,
		Mobile:        req.Mobile,
		RegTime:       now,
		LastLoginTime: now,
		IP:            req.IP,
		Password:      req.Password,
	}
	if err := dbutils.DB.Save(account).Error; err != nil {
		panic(err)
	}
	return account
}

func (s *Service_V2) GetAccountByID(id string) (*Account_V2, *Error) {
	var account Account_V2
	err := dbutils.DB.First(&account, "id = ?", id).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			panic(err)
		} else {
			return nil, &Error{notFound: true}
		}
	}
	return &account, nil
}
