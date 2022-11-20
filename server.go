package goloadbalancer

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"
)

//reverse proxy
// func Serve() {
// 	director := func(request *http.Request) {
// 		request.URL.Scheme = "http"
// 		request.URL.Host = ":8081"
// 	}
// 	rp := &httputil.ReverseProxy{
// 		Director: director,
// 	}

// 	s := http.Server{
// 		Addr:     ":8080",
// 		Handler: rp,
// 	}

// 	if err := s.ListenAndServe(); err != nil {
// 		log.Fatal(err.Error())
// 	}
// }

// config 
type Config struct {
	Proxy Proxy 
	Backends []Backends
}
type Proxy struct {
	Port string
}
type Backends struct {
	URL string
	IsDead bool
	mu *sync.RWMutex
}
var mu sync.Mutex
var idx int = 0

func lbHandler(w http.ResponseWriter, r *http.Request){
	maxLen := len(cfg.Backends)
	mu.Lock()
	currentBackend := cfg.Backends[idx%maxLen]
	targetURL, err := url.Parse(currentBackend.URL)
	if err != nil {
		log.Fatal(err.Error())
	}
	idx++
	mu.Unlock()
	reverseProxy := httputil.NewSingleHostReverseProxy(targetURL)
	reverseProxy.ServeHTTP(w, r)
}

var cfg Config

func Serve() {
	data, err := ioutil.ReadFile("./config.json")
	if err != nil {
        log.Fatal(err.Error())
    }
    json.Unmarshal(data, &cfg)
	
	s := http.Server{
		Addr: ":" + cfg.Proxy.Port,
		Handler: http.HandlerFunc(lbHandler),
	}
	if err = s.ListenAndServe(); err != nil {
        log.Fatal(err.Error())
    }
}