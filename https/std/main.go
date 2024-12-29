package main

import (
	"net/http"

	"github.com/unrolled/secure"
)

//openssl genrsa -out key.pem 2048
//openssl req -new -key key.pem -out cert.csr
//openssl x509 -req -days 365 -in cert.csr -signkey key.pem -out cert.pem

var myHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, TLS!\n"))
})

func main() {
	secureMiddleware := secure.New(secure.Options{
		FrameDeny:   true,
		SSLRedirect: true,
		SSLHost:     "localhost:4000",
	})

	app := secureMiddleware.Handler(myHandler)
	if err := http.ListenAndServeTLS(":4000", "../cert.pem", "../key.pem", app); err != nil {
		panic(err)
	}

}
