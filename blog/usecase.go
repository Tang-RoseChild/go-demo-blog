package blog

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Tang-RoseChild/go-demo-blog/blog/store"
	"github.com/Tang-RoseChild/go-demo-blog/utils/http"
	"github.com/Tang-RoseChild/go-demo-blog/utils/token"
	"github.com/Tang-Rosechild/go-demo-blog/middleware"
)

func Upload(w http.ResponseWriter, r *http.Request) {
	var req store.CreateReq_V2
	httputils.MustUnmarshalReq(r, &req)
	var blog *store.Blog_V2
	var err error
	if req.ID != "" {
		blog, err = store.DefaultService.Update(&store.UpdateReq_V2{
			ID:          req.ID,
			Title:       req.Title,
			Content:     req.Content,
			Tag:         req.Tag,
			Source:      req.Source,
			Description: req.Description,
			Status:      store.StausPublished,
			Points:      req.Points,
		})
	} else {
		req.Status = store.StausPublished
		req.UserID = "admin"
		blog = store.DefaultService.Create(&req)
	}

	httputils.MustMarshalResp(r, w, map[string]interface{}{"blog": blog, "err": err})
}

func Update(w http.ResponseWriter, r *http.Request) {
	var req store.UpdateReq_V2
	httputils.MustUnmarshalReq(r, &req)
	var blog *store.Blog_V2
	var err error
	// fmt.Println("update req points .>> ", req.Points)
	if req.ID == "" {
		blog = store.DefaultService.Create(&store.CreateReq_V2{
			Title:       req.Title,
			Content:     req.Content,
			Description: req.Description,
			Tag:         req.Tag,
			UserID:      "admin",
			Status:      store.StatusSaved,
			Points:      req.Points,
		})
	} else {
		blog, err = store.DefaultService.Update(&req)
	}

	httputils.MustMarshalResp(r, w, map[string]interface{}{"blog": blog, "err": err})
}

func List(c *gin.Context) {
	pageClaim := c.MustGet(middleware.PaginationTokenKey).(*tokenutils.Pagination)
	list, total := store.DefaultService.GetBlogList(&store.ListReq_V2{Limit: pageClaim.Limit, From: pageClaim.From, UserID: "admin"})
	hasMore := pageClaim.From+pageClaim.Limit < total
	if hasMore {
		pageClaim.From += pageClaim.Limit
	}
	c.Writer.Header().Set("Authorization", pageClaim.GenerateToken(tokenutils.GetSecret()))
	httputils.MustMarshalResp(c.Request, c.Writer, map[string]interface{}{
		"blogs":   list,
		"hasMore": hasMore,
	})
}

func Get(c *gin.Context) {
	var req struct {
		ID string
	}
	httputils.MustUnmarshalReq(c.Request, &req)
	blog, err := store.DefaultService.GetBlog(req.ID)
	// commentsResp := commentStore.DefaultService.ListComments(&commentStore.ListCommentsReq{BlogID: req.ID})
	// if err == nil {
	// 	err = commentsResp.Err
	// }
	httputils.MustMarshalResp(c.Request, c.Writer, map[string]interface{}{
		"blog":  blog,
		"error": err,
		// "comments": commentsResp.Comments,
	})
}
func ListByTag(c *gin.Context) {
	var req store.ListReq_V2
	httputils.MustUnmarshalReq(c.Request, &req)

	pageClaim := c.MustGet(middleware.PaginationTokenKey).(*tokenutils.Pagination)
	req.Limit = pageClaim.Limit
	req.From = pageClaim.From
	list, total := store.DefaultService.GetBlogsByTag(&req)
	hasMore := pageClaim.From+pageClaim.Limit < total
	if hasMore {
		pageClaim.From += pageClaim.Limit
	}
	c.Writer.Header().Set("Authorization", pageClaim.GenerateToken(tokenutils.GetSecret()))
	httputils.MustMarshalResp(c.Request, c.Writer, map[string]interface{}{
		"blogs":   list,
		"hasMore": hasMore,
	})
}

func ListBySource(c *gin.Context) {
	var req store.ListReq_V2
	httputils.MustUnmarshalReq(c.Request, &req)
	pageClaim := c.MustGet(middleware.PaginationTokenKey).(*tokenutils.Pagination)
	req.Limit = pageClaim.Limit
	req.From = pageClaim.From
	list, total := store.DefaultService.GetBlogsBySource(&req)
	hasMore := pageClaim.From+pageClaim.Limit < total
	if hasMore {
		pageClaim.From += pageClaim.Limit
	}
	c.Writer.Header().Set("Authorization", pageClaim.GenerateToken(tokenutils.GetSecret()))
	httputils.MustMarshalResp(c.Request, c.Writer, map[string]interface{}{
		"blogs":   list,
		"hasMore": hasMore,
		// "errro":   err,
	})
}
func GinLoad(rootGroup *gin.RouterGroup) {
	g := rootGroup.Group("/blog")
	g.POST("/", Get)
	g.POST("/upload", middleware.NeedLogin(), httputils.ToGinHandler(Upload))
	g.POST("/update", middleware.NeedLogin(), httputils.ToGinHandler(Update))
	g.POST("/list", middleware.PaginationToken("list", 10), List)
	g.POST("/list/tag", middleware.PaginationToken("list", 10), ListByTag)
	g.POST("/list/source", middleware.PaginationToken("list", 10), ListBySource)

}
