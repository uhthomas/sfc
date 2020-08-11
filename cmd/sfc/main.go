package main

import (
	"context"
	"flag"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/uhthomas/sfc"
)

func Main(ctx context.Context) error {
	addr := flag.String("addr", ":80", "Listen address")
	www := flag.String("www", "www", "File server path")
	flag.Parse()

	return (&http.Server{
		Addr: *addr,
		Handler: &sfc.Handler{
			Client: &sfc.Client{
				C: &http.Client{Timeout: 10 * time.Second},
				BaseURL: &url.URL{
					Scheme: "https",
					Host:   "track.sendfromchina.com",
				},
			},
			FileServer: http.FileServer(http.Dir(*www)),
		},
		BaseContext: func(net.Listener) context.Context { return ctx },
	}).ListenAndServe()
}

func main() {
	if err := Main(context.Background()); err != nil {
		log.Fatal(err)
	}
}
