package comments

import (
	"net/http"

	"github.com/gin-gonic/gin"

	commentStore "github.com/Tang-RoseChild/go-demo-blog/comments/store"
	"github.com/Tang-RoseChild/go-demo-blog/utils/http"
	"github.com/Tang-RoseChild/go-demo-blog/utils/token"
)

func GinLoad(rootGroup *gin.RouterGroup) {
	g := rootGroup.Group("/comment")
	g.POST("/create", httputils.ToGinHandler(CreateCommentHandler))
	g.POST("/list_by_blog", httputils.ToGinHandler(ListCommentsByBlogHandler))
	g.POST("/list_by_user", httputils.ToGinHandler(ListCommentsByUserHandler))

}
func Load() {
	http.HandleFunc("/api/comment/create", tokenutils.IssueToken(CreateCommentHandler))
	http.HandleFunc("/api/comment/list_by_blog", ListCommentsByBlogHandler)
	http.HandleFunc("/api/comment/list_by_user", tokenutils.IssueToken(ListCommentsByUserHandler))
}

func CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Content string
		BlogID  string
	}

	httputils.MustUnmarshalReq(r, &req)
	comment := commentStore.DefaultService.Create(r.Context().Value("uid").(string), req.BlogID, req.Content)
	httputils.MustMarshalResp(w, comment)
}

func ListCommentsByBlogHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		BlogID string
	}
	httputils.MustUnmarshalReq(r, &req)
	resp := commentStore.DefaultService.ListComments(&commentStore.ListCommentsReq{BlogID: req.BlogID})
	httputils.MustMarshalResp(w, resp)
}

func ListCommentsByUserHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID string
	}
	httputils.MustUnmarshalReq(r, &req)
	resp := commentStore.DefaultService.ListComments(&commentStore.ListCommentsReq{UserID: req.UserID})
	httputils.MustMarshalResp(w, resp)
}
