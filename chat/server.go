package main

import (
	"flag"
	"fmt"
	"net/http"
)

func serveHome(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "home.html")
}

func main() {
	port := flag.String("port", "8080", "http service port")
	flag.Parse()
	hub := NewHub()
	go hub.run()

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		ServeWs(hub, w, r)
	})
	fmt.Printf("http server started on :%s\n", *port)
	if err := http.ListenAndServe(":"+*port, nil); err != nil {
		fmt.Printf("ListenAndServe: %v\n", err)
	}

}
