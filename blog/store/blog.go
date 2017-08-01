package store

import (
	"errors"
	"time"

	"github.com/Tang-RoseChild/go-demo-blog/utils/db"
	"github.com/Tang-RoseChild/go-demo-blog/utils/id"

	"github.com/jinzhu/gorm"
)

const (
	StatusSaved    = 1
	StausPublished = 2
)

type pgBlog struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Content     string    `json:"content" sql:"type:text"`
	UserID      string    `json:"user_id"`
	Timestamp   time.Time `json:"timestamp"`
	Status      int       `json:"status"`
	Description string    `json:"description"`
	Tag         string    `json:"tag"`
	Source      string    `json:"source"`
}

func (b pgBlog) TableName() string {
	return "blogs"
}

type pgBlogList []*pgBlog

func (list pgBlogList) ToBlogList() BlogList {
	var b BlogList
	for _, item := range list {
		b = append(b, (*Blog)(item))
	}
	return b
}

var (
	NotFound = errors.New("not found")
)

type Blog pgBlog
type BlogList []*Blog

func (list BlogList) ToPGBlogList() []*pgBlog {
	var pglist []*pgBlog
	for _, item := range list {
		pglist = append(pglist, (*pgBlog)(item))
	}
	return pglist
}
func GetBlog(id string) (*Blog, error) {
	var blog pgBlog
	if err := dbutils.DB.First(&blog, "id=?", id).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			panic(err)
		}
		return nil, NotFound
	}

	return (*Blog)(&blog), nil
}

type ListReq struct {
	Limit  int
	From   int
	UserID string
}

func GetBlogList(req *ListReq) ([]*Blog, bool, error) {
	var blogs []*pgBlog
	// if err := dbutils.DB.Offset(req.From).Limit(req.Limit).Order("timestamp desc").Find(&blogs, "user_id=?", req.UserID).Error; err != nil {
	if err := dbutils.DB.Order("timestamp desc").Find(&blogs, "user_id=?", req.UserID).Error; err != nil {
		panic(err)
	}

	// var count int
	// dbutils.DB.Table(pgBlog{}.TableName()).Count(&count)
	var hasMore = false
	// if req.Limit+req.From < count {
	// 	hasMore = true
	// }
	return pgBlogList(blogs).ToBlogList(), hasMore, nil
}

type CreateReq struct {
	Title       string
	Content     string
	UserID      string
	Tag         string
	Source      string
	Description string
	Status      int
	ID          string
}

func Create(req *CreateReq) *Blog {
	blog := &pgBlog{
		ID:          idutils.DefaultGenerator.GetID(),
		Title:       req.Title,
		Content:     req.Content,
		UserID:      req.UserID,
		Timestamp:   time.Now().UTC(),
		Description: req.Description,
		Tag:         req.Tag,
		Status:      req.Status,
		Source:      req.Source,
	}

	if err := dbutils.DB.Create(blog).Error; err != nil {
		panic(err)
	}
	return (*Blog)(blog)
}

type UpdateReq struct {
	Title       string
	Content     string
	UserID      string
	Tag         string
	Source      string
	ID          string
	Description string
	Status      int
}

func Update(req *UpdateReq) (*Blog, error) {
	blog := &pgBlog{}
	switch {
	case req.Title != "":
		blog.Title = req.Title
		fallthrough
	case req.Content != "":
		blog.Content = req.Content
		fallthrough
	case req.Tag != "":
		blog.Tag = req.Tag
		fallthrough
	case req.Source != "":
		blog.Source = req.Source
		fallthrough
	case req.Description != "":
		blog.Description = req.Description
		fallthrough
	case req.Status != 0:
		blog.Status = req.Status
	}

	if err := dbutils.DB.Table(blog.TableName()).Where("id=?", req.ID).UpdateColumns(blog).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, NotFound
		}
		panic(err)
	}

	return GetBlog(req.ID)
}
