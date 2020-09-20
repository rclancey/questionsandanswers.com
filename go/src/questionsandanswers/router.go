package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type Handler func(w http.ResponseWriter, r *http.Request, pathinfo string) (interface{}, error)

type RouteKey struct {
	Method string
	Path string
}

type Router struct {
	DocumentRoot string
	DefaultProxy *url.URL
	Routes map[RouteKey]Handler
}

func NewRouter(docroot, proxy string) (*Router, error) {
	docroot, err := filepath.Abs(docroot)
	if err != nil {
		return nil, err
	}
	docroot = filepath.Clean(docroot)
	st, err := os.Stat(docroot)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("document root %s does not exist", docroot)
		}
		return nil, err
	}
	if !st.IsDir() {
		return nil, fmt.Errorf("document root %s is not a directory", docroot)
	}
	var proxyUrl *url.URL
	if proxy != "" {
		proxyUrl, err = url.Parse(proxy)
		if err != nil {
			return nil, err
		}
	}
	return &Router{
		DocumentRoot: docroot,
		DefaultProxy: proxyUrl,
		Routes: map[RouteKey]Handler{},
	}, nil
}

func (router *Router) AddRoute(method, path string, handler Handler) {
	key := RouteKey{method, path}
	router.Routes[key] = handler
}

func (router *Router) Lookup(r *http.Request) (Handler, string) {
	pth := r.URL.Path
	pathinfo := ""
	var key RouteKey
	for len(pth) > 1 {
		pth = strings.TrimRight(pth, "/")
		key = RouteKey{r.Method, pth}
		handler, ok := router.Routes[key]
		if ok {
			return handler, pathinfo
		}
		var extra string
		pth, extra = path.Split(pth)
		if pathinfo == "" {
			pathinfo = extra
		} else {
			pathinfo = path.Join(extra, pathinfo)
		}
	}
	return nil, ""
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	iw := &InstrumentedResponseWriter{w: w, Status: 0, Bytes: 0}
	defer r.Body.Close()
	fn := filepath.Clean(filepath.FromSlash(strings.TrimLeft(r.URL.Path, "/")))
	fn = filepath.Join(router.DocumentRoot, fn)
	st, err := os.Stat(fn)
	if err == nil && st.IsDir() {
		fn = filepath.Join(fn, "index.html")
		_, err = os.Stat(fn)
	}
	if err == nil {
		log.Printf("serve static %s for %s", fn, r.URL)
		http.ServeFile(iw, r, fn)
		accessLogger.Log(iw, r, startTime)
		return
	}
	handler, pathinfo := router.Lookup(r)
	if handler == nil {
		if router.DefaultProxy != nil {
			router.Proxy(iw, r)
			accessLogger.Log(iw, r, startTime)
			return
		}
		SendError(iw, NotFound(nil))
		accessLogger.Log(iw, r, startTime)
		return
	}
	obj, err := handler(iw, r, pathinfo)
	if err != nil {
		SendError(iw, err)
		accessLogger.Log(iw, r, startTime)
		return
	}
	var data []byte
	if obj == nil {
		data = []byte("null")
	} else {
		data, err = json.Marshal(obj)
	}
	if err != nil {
		SendError(iw, err)
		accessLogger.Log(iw, r, startTime)
		return
	}
	iw.Header().Set("Content-Type", "application/json")
	iw.WriteHeader(http.StatusOK)
	iw.Write(data)
	accessLogger.Log(iw, r, startTime)
}
