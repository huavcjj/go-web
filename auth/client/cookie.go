package main

import (
	"io"
	"net/http"
	"os"
)

func main() {
	if resp, err := http.Get("http://localhost:5678/login"); err != nil {
		panic(err)
	} else {
		io.Copy(os.Stdout, resp.Body)
		os.Stdout.WriteString("\n")
		loginCookies := resp.Cookies()
		resp.Body.Close()
		if req, err := http.NewRequest("GET", "http://localhost:5678/profile", nil); err != nil {
			panic(err)
		} else {
			for _, cookie := range loginCookies {
				req.AddCookie(cookie)
			}
			if resp, err := http.DefaultClient.Do(req); err != nil {
				panic(err)
			} else {
				io.Copy(os.Stdout, resp.Body)
				os.Stdout.WriteString("\n")
				resp.Body.Close()
			}
		}
	}
}
