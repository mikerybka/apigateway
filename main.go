package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"

	"golang.org/x/crypto/acme/autocert"
)

func main() {
	proxyPort := "8000"
	h := &httputil.ReverseProxy{
		Director: func(r *http.Request) {
			r.URL.Scheme = "http"
			r.URL.Host = "localhost:" + proxyPort
		},
	}
	http.Handle("/", h)
	go listen(":80")
	listenHTTPS()
}

func listen(addr string) {
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		panic(err)
	}
}

func listenHTTPS() {
	hosts := []string{"api.ryb.dev"}
	certsDir := "/certs"
	amceEmail := "merybka@gmail.com"
	m := autocert.Manager{
		Prompt: autocert.AcceptTOS,
		Cache:  autocert.DirCache(certsDir),
		HostPolicy: func(_ context.Context, host string) error {
			for _, h := range hosts {
				if h == host {
					return nil
				}
			}
			return fmt.Errorf("host %q not allowed", host)
		},
		Email: amceEmail,
	}
	err := http.Serve(m.Listener(), nil)
	if err != nil {
		panic(err)
	}
}
