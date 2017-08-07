package store

import (
	"errors"
	"time"

	"github.com/Tang-RoseChild/go-demo-blog/utils/db"
	"github.com/Tang-RoseChild/go-demo-blog/utils/id"

	"github.com/jinzhu/gorm"
)

var (
	NotFoundErr = errors.New("not found")
)

type Blog_V2 struct {
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

func (b Blog_V2) TableName() string {
	return "blogs"
}

type Service struct{}

func (s *Service) GetBlog(id string) (*Blog_V2, error) {
	var blog Blog_V2
	if err := dbutils.DB.First(&blog, "id=?", id).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			panic(err)
		}
		return nil, NotFoundErr
	}

	return &blog, nil
}

type ListReq_V2 struct {
	Limit  int
	From   int
	UserID string
}

func (s *Service) GetBlogList(req *ListReq) ([]*Blog_V2, bool, error) {
	var blogs []*Blog_V2
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
	return blogs, hasMore, nil
}

type CreateReq_V2 struct {
	Title       string
	Content     string
	UserID      string
	Tag         string
	Source      string
	Description string
	Status      int
	ID          string
}

func (s *Service) Create(req *CreateReq) *Blog_V2 {
	blog := &Blog_V2{
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
	return blog
}

type UpdateReq_V2 struct {
	Title       string
	Content     string
	UserID      string
	Tag         string
	Source      string
	ID          string
	Description string
	Status      int
}

func (s *Service) Update(req *UpdateReq_V2) (*Blog_V2, error) {
	blog := &Blog_V2{}
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

	return s.GetBlog(req.ID)
}

func (s *Service) GetBlogsByTag(tags ...string) []*Blog_V2 {
	var blogs []*Blog_V2
	dbutils.DB.LogMode(true)
	dbutils.DB.Table(Blog_V2{}.TableName()).Find(&blogs, "tag in (?)", tags)
	return blogs
}

func (s *Service) GetBlogsBySource(source ...string) []*Blog_V2 {
	var blogs []*Blog_V2
	dbutils.DB.Table(Blog_V2{}.TableName()).Find(&blogs, "source in (?)", source)
	return blogs
}

var DefaultService = &Service{}
