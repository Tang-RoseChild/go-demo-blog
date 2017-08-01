package store

import (
	"fmt"
	"time"

	"github.com/Tang-RoseChild/go-demo-blog/utils/db"
	"github.com/Tang-RoseChild/go-demo-blog/utils/id"

	"github.com/jinzhu/gorm"
)

type Comment struct {
	ID        string
	Content   string
	BlogID    string
	Timestamp time.Time
	Status    int
	UserID    string
}

type ListCommentsReq struct {
	BlogID string
	UserID string
	Limit  int
	From   int
}
type ListCommentsResp struct {
	Comments []*Comment
	HasMore  bool
	Err      error
}

type CommentService interface {
	Create(userID, blogID string, content string) *Comment
	ListComments(req *ListCommentsReq) *ListCommentsResp
}

var DefaultService = &service{}

type service struct{}

func (s *service) Create(userID, blogID string, content string) *Comment {
	c := &Comment{
		UserID:    userID,
		BlogID:    blogID,
		Content:   content,
		Timestamp: time.Now().UTC(),
		ID:        idutils.DefaultGenerator.GetID(),
	}
	dbutils.DB.Create(c)

	return c
}

func (s *service) ListComments(req *ListCommentsReq) *ListCommentsResp {
	var (
		comments []*Comment
		err      error
		count    int
	)
	switch {
	case req.BlogID != "":
		err = dbutils.DB.Order("timestamp desc").Find(&comments, "blog_id = ?", req.BlogID).Count(&count).Error
	case req.UserID != "":
		err = dbutils.DB.Order("timestamp desc").Find(&comments, "user_id = ?", req.UserID).Count(&count).Error
	default:
		err = dbutils.DB.Order("timestamp desc").Find(&comments).Count(&count).Error
	}
	fmt.Println("list comments >>> , ")
	if err != nil && err != gorm.ErrRecordNotFound {
		fmt.Println("err >>> in panic", err)
		panic(err)
	}

	resp := &ListCommentsResp{
		Comments: comments,
		Err:      err,
		HasMore:  count > req.Limit+req.From,
	}
	fmt.Println("resp >> ", resp)
	return resp
}
