package httputils

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

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

func MustMarshalResp(r *http.Request, w http.ResponseWriter, resp interface{}) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(resp)
	if err != nil {
		panic(err)
	}

	canCompress := strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")
	if !canCompress {
		fmt.Fprint(w, string(data))
		return
	}

	w.Header().Set("Content-Encoding", "gzip")
	gw, err := gzip.NewWriterLevel(w, gzip.BestCompression)
	if err != nil {
		panic(err)
	}
	n, err := gw.Write(data)
	if n != len(data) || err != nil {
		panic("write all data error")
	}
	gw.Flush()
}

func ToGinHandler(h http.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		h(c.Writer, c.Request)
	}
}
