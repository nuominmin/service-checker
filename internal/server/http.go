package server

import (
	"bytes"
	"embed"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/gorilla/handlers"
	nethttp "net/http"
	pb "service-checker/api"
	"service-checker/internal/conf"
	"service-checker/internal/service"
	"time"
)

// frontend/services-status/dist/*
//
//go:embed *
var content embed.FS

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, svc *service.Service) *http.Server {
	var opts = []http.ServerOption{
		http.Timeout(time.Second * 30),
		http.Filter(handlers.CORS(
			handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization", "User-Agent", "Content-Length", "Access-Control-Allow-Credentials"}),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}),
			handlers.AllowedOrigins([]string{"*"}), handlers.AllowedOrigins([]string{"*"}),
			handlers.AllowCredentials(),
		)),
		http.Middleware(
			recovery.Recovery(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)

	srv.Handle("/static/", nethttp.StripPrefix("/static/", nethttp.FileServer(nethttp.FS(content))))
	srv.HandleFunc("/", func(w nethttp.ResponseWriter, r *nethttp.Request) {
		body, err := content.ReadFile("dist/index.html")
		if err != nil {
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		nethttp.ServeContent(w, r, "index.html", time.Now(), bytes.NewReader(body))
	})

	pb.RegisterV1HTTPServer(srv, svc)
	return srv
}
