package blog

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Tang-RoseChild/go-demo-blog/blog/store"
	"github.com/Tang-RoseChild/go-demo-blog/utils/http"
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

	httputils.MustMarshalResp(w, map[string]interface{}{"blog": blog, "err": err})
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

	httputils.MustMarshalResp(w, map[string]interface{}{"blog": blog, "err": err})
}

func List(w http.ResponseWriter, r *http.Request) {
	// var req struct {
	// 	Limit int
	// 	From  int
	// }
	// httputils.MustUnmarshalReq(r, &req)

	list, hasMore, err := store.DefaultService.GetBlogList(&store.ListReq_V2{0, 0, "admin"})
	httputils.MustMarshalResp(w, map[string]interface{}{
		"blogs":   list,
		"hasMore": hasMore,
		"errro":   err,
	})
}

func Get(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID string
	}
	httputils.MustUnmarshalReq(r, &req)
	blog, err := store.DefaultService.GetBlog(req.ID)
	// commentsResp := commentStore.DefaultService.ListComments(&commentStore.ListCommentsReq{BlogID: req.ID})
	// if err == nil {
	// 	err = commentsResp.Err
	// }
	httputils.MustMarshalResp(w, map[string]interface{}{
		"blog":  blog,
		"error": err,
		// "comments": commentsResp.Comments,
	})
}

func ListByTag(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Tag string
		// 	Limit int
		// 	From  int
	}
	httputils.MustUnmarshalReq(r, &req)
	// list, hasMore, err := store.DefaultService.GetBlogListByTag(req.Tag)
	list := store.DefaultService.GetBlogsByTag(req.Tag)
	httputils.MustMarshalResp(w, map[string]interface{}{
		"blogs": list,
		// "hasMore": hasMore,
		// "errro":   err,
	})
}

func ListBySource(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Source string
		// 	Limit int
		// 	From  int
	}
	httputils.MustUnmarshalReq(r, &req)
	// list, hasMore, err := store.DefaultService.GetBlogListByTag(req.Tag)
	list := store.DefaultService.GetBlogsBySource(req.Source)
	httputils.MustMarshalResp(w, map[string]interface{}{
		"blogs": list,
		// "hasMore": hasMore,
		// "errro":   err,
	})
}
func GinLoad(rootGroup *gin.RouterGroup) {
	g := rootGroup.Group("/blog")
	g.POST("/", httputils.ToGinHandler(Get))
	g.POST("/upload", middleware.NeedLogin(), httputils.ToGinHandler(Upload))
	g.POST("/update", middleware.NeedLogin(), httputils.ToGinHandler(Update))
	g.POST("/list", httputils.ToGinHandler(List))
	g.POST("/list/tag", httputils.ToGinHandler(ListByTag))
	g.POST("/list/source", httputils.ToGinHandler(ListBySource))

}
