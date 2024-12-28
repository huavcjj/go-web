package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func getBoy(w http.ResponseWriter, r *http.Request) {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
		}
	}()
	time.Sleep(2 * time.Second)
	w.Write([]byte("Hi Boy"))
}

func getGirl(w http.ResponseWriter, r *http.Request) {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
		}
	}()
	time.Sleep(2 * time.Second)
	w.Write([]byte("Hi Girl"))
}

func timeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		begin := time.Now()
		next.ServeHTTP(w, r)
		timeElapsed := time.Since(begin)
		log.Printf("%s used time %v", r.URL.Path, timeElapsed.Milliseconds())
	})
}

func limitMiddleware(next http.Handler) http.Handler {
	limitCh := make(chan struct{}, 10)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		limitCh <- struct{}{}
		log.Printf("current limitCh length: %d", len(limitCh))
		next.ServeHTTP(w, r)
		<-limitCh
	})
}

// func main() {
// 	http.Handle("/boy", timeMiddleware(limitMiddleware(http.HandlerFunc(getBoy))))
// 	http.Handle("/girl", timeMiddleware(limitMiddleware(http.HandlerFunc(getGirl))))

// 	if err := http.ListenAndServe(":8080", nil); err != nil {
// 		log.Fatal(err)
// 	}

// }
