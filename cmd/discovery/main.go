package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

var host string
var defaultServer *server

func init() {
	flag.StringVar(&host, "host", ":9994", "host addr")
	flag.Parse()

	defaultServer = &server{}
}

func main() {
	http.HandleFunc("/", proxyHandler) // default handler
	http.HandleFunc("/_x/server/add", addServerHandler)
	http.HandleFunc("/_x/server/del", delServerHandler)
	http.HandleFunc("/_x/config", configHandler)

	fmt.Printf("listening   %s ...... \n", host)
	if err := http.ListenAndServe(host, nil); err != nil {
		panic(err)
	}
}

// weight :
// path的话太细了，如果一个server很多个接口，每个接口都添加个path会非常麻烦
// 用server name
type server struct {
	Route map[string]*route
	mux   sync.Mutex
}
type route struct {
	Addrs []string
	index int
}
type serverConfig struct {
	Name   string
	Addr   string
	Weight int
}

func (s *server) addServer(c *serverConfig) {
	s.mux.Lock()
	defer s.mux.Unlock()
	if s.Route == nil {
		s.Route = make(map[string]*route)
	}

	for i := 0; i < c.Weight; i++ {
		if s.Route[c.Name] == nil {
			s.Route[c.Name] = new(route)
		}
		s.Route[c.Name].Addrs = append(s.Route[c.Name].Addrs, c.Addr)
	}
}

func (s *server) delServer(c *serverConfig) {
	s.mux.Lock()
	defer s.mux.Unlock()
	if s.Route == nil {
		return
	}

	rout, ok := s.Route[c.Name]
	if !ok {
		return
	}
	var newAddrs []string

	for _, addr := range rout.Addrs {
		if addr != c.Addr {
			newAddrs = append(newAddrs, addr)
		}
	}

	s.Route[c.Name].Addrs = newAddrs
}

func (s *server) getServer(name string) string {
	s.mux.Lock()
	defer s.mux.Unlock()
	_route, ok := s.Route[name]
	if !ok {
		return ""
	}

	server := _route.Addrs[_route.index]
	_route.index++
	if _route.index == len(_route.Addrs) {
		_route.index = 0
	}
	return server
}

func addServerHandler(w http.ResponseWriter, r *http.Request) {
	var config serverConfig
	mustUnmarshalReq(r, &config)

	defaultServer.addServer(&config)
	mustMarshalResp(w, defaultServer)
}

func delServerHandler(w http.ResponseWriter, r *http.Request) {
	var config serverConfig
	mustUnmarshalReq(r, &config)

	defaultServer.delServer(&config)
	mustMarshalResp(w, defaultServer)
}

func configHandler(w http.ResponseWriter, r *http.Request) {
	mustMarshalResp(w, defaultServer)
}

// http get('host:port/serviceName/uri')
func proxyHandler(w http.ResponseWriter, r *http.Request) {
	url := strings.TrimPrefix(r.URL.Path, "/")
	index := strings.Index(url, "/")
	serviceName := url[:index]
	server := defaultServer.getServer(serviceName)
	if server == "" {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, `http://`+server+url[index:], http.StatusMovedPermanently)
}

func mustUnmarshalReq(r *http.Request, out interface{}) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(data, out)
	if err != nil {
		panic(err)
	}
}

func mustMarshalResp(w http.ResponseWriter, resp interface{}) {
	data, err := json.Marshal(resp)
	if err != nil {
		panic(err)
	}

	fmt.Fprint(w, string(data))
}
