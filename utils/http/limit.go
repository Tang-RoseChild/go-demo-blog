package httputils

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
)

func LimitReq(max int, handler http.HandlerFunc) http.HandlerFunc {
	ch := make(chan struct{}, max)

	return func(w http.ResponseWriter, r *http.Request) {
		ch <- struct{}{}
		defer func() {
			<-ch
		}()
		handler(w, r)
	}
}

type ipmap struct {
	m   map[string]int
	mux sync.RWMutex
	max int
}

func (m *ipmap) Pass(ip string) bool {

	m.mux.Lock()
	defer m.mux.Unlock()
	m.m[ip]++
	fmt.Println("ip in pass >>> ", ip, "map >>> ", m.m, " max >>> ", m.max)
	return m.m[ip] <= m.max
}

func (m *ipmap) SetMaxReq(max int) {
	m.max = max
}

var _ipmap = &ipmap{
	m: make(map[string]int),
}

func IPLimit(max int, handler http.HandlerFunc) http.HandlerFunc {
	_ipmap.SetMaxReq(max)
	return func(w http.ResponseWriter, r *http.Request) {
		splits := strings.Split(r.RemoteAddr, ":")
		fmt.Println("remote address , host and ip ", splits, " splits[0] >>> ", splits[0])
		if !_ipmap.Pass(splits[0]) {
			fmt.Println("not passed >>>>> too frequently ")
			http.Error(w, "request too frequently", http.StatusBadRequest)
			return
		}
		if handler != nil {
			handler(w, r)
		}
	}
}
