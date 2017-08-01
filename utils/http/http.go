package httputils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

func MustUnmarshalReq(r *http.Request, out interface{}) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(data, out)
	if err != nil {
		panic(err)
	}
}

func MustMarshalResp(w http.ResponseWriter, resp interface{}) {
	data, err := json.Marshal(resp)
	if err != nil {
		panic(err)
	}
	fmt.Println("string data")
	fmt.Fprint(w, string(data))
}

func ToGinHandler(h http.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		h(c.Writer, c.Request)
	}
}
