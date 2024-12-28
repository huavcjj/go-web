package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

var limitCh = make(chan struct{}, 10)

func GetBoy(w http.ResponseWriter, r *http.Request) {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
		}
	}()
	time.Sleep(2 * time.Second)
	w.Write([]byte("Hi Boy"))
}

func GetGirl(w http.ResponseWriter, r *http.Request) {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
		}
	}()
	time.Sleep(2 * time.Second)
	w.Write([]byte("Hi Girl"))
}

func TimeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		begin := time.Now()
		next.ServeHTTP(w, r)
		timeElapsed := time.Since(begin)
		log.Printf("%s used time %v", r.URL.Path, timeElapsed.Milliseconds())
	})
}

func LimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		limitCh <- struct{}{}
		log.Printf("current limitCh length: %d", len(limitCh))
		next.ServeHTTP(w, r)
		<-limitCh
	})
}

type Middleware func(http.Handler) http.Handler

type Router struct {
	middlewareChain []Middleware
	routes          map[string]http.Handler
}

func NewRouter() *Router {
	return &Router{
		middlewareChain: []Middleware{},
		routes:          make(map[string]http.Handler),
	}
}

func (router *Router) Use(m Middleware) {
	router.middlewareChain = append(router.middlewareChain, m)
}

func (router *Router) Handle(method, pattern string, handler http.Handler) {
	wrappedHandler := handler
	for i := len(router.middlewareChain) - 1; i >= 0; i-- {
		wrappedHandler = router.middlewareChain[i](wrappedHandler)
	}
	router.routes[method+pattern] = wrappedHandler
}

func (router *Router) GET(pattern string, handler http.Handler) {
	router.Handle("GET", pattern, handler)
}

func (router *Router) POST(pattern string, handler http.Handler) {
	router.Handle("POST", pattern, handler)
}

func (router *Router) PUT(pattern string, handler http.Handler) {
	router.Handle("PUT", pattern, handler)
}

func (router *Router) DELETE(pattern string, handler http.Handler) {
	router.Handle("DELETE", pattern, handler)
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestPath := r.URL.Path
	if handler, ok := router.routes[r.Method+requestPath]; ok {
		handler.ServeHTTP(w, r)
		return
	}
	http.NotFound(w, r)
}

func main() {
	router := NewRouter()
	router.Use(TimeMiddleware)
	router.Use(LimitMiddleware)

	router.GET("/boy", http.HandlerFunc(GetBoy))
	router.GET("/girl", http.HandlerFunc(GetGirl))

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
