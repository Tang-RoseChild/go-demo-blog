package store

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/Tang-RoseChild/go-demo-blog/utils/db"
	"github.com/Tang-RoseChild/go-demo-blog/utils/id"
	"github.com/jinzhu/gorm"
)

var (
	NotFound     = errors.New("not found")
	MobileExists = errors.New("mobile exists")
)

type dbAccount struct {
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

func (a *dbAccount) TableName() string {
	return "accounts"
}

func (a *dbAccount) GetMobile() (no, mobile string) {
	return a.NO, a.Mobile
}
func (a *dbAccount) GetRegTime() time.Time {
	return a.RegTime
}
func (a *dbAccount) GetLastLoginTime() time.Time {
	return a.LastLoginTime
}
func (a *dbAccount) GetID() string {
	return a.ID
}
func (a *dbAccount) GetStatus() int {
	return a.Status
}

func (a *dbAccount) SetStatus(status int) {
	a.Status = status
}
func (a *dbAccount) Marshal() (string, string) {
	data, err := json.Marshal(a)
	if err != nil {
		return "", err.Error()
	}
	return "json", string(data)
}

func (a *dbAccount) GetPassword() string {
	return a.Password
}

type Service struct{}

func (s *Service) GetAccountByMobile(no, mobile string) (Account, error) {
	var account dbAccount
	err := dbutils.DB.First(&account, "no = ? and mobile = ?", no, mobile).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			panic(err)
		} else {
			err = NotFound
		}
	}
	return &account, err
}
func (s *Service) GetAccountByID(id string) (Account, error) {
	var account dbAccount
	err := dbutils.DB.First(&account, "id = ?", id).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			panic(err)
		} else {
			err = NotFound
		}
	}
	return &account, err
}

func (s *Service) CreateAccount(no, mobile string, password string, ip string) (Account, error) {
	now := time.Now().UTC()
	account := &dbAccount{
		ID:            idutils.DefaultGenerator.GetID(),
		NO:            no,
		Mobile:        mobile,
		RegTime:       now,
		LastLoginTime: now,
		IP:            ip,
		Password:      password,
	}
	err := dbutils.DB.Save(account).Error
	return account, err
}
