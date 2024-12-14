package http_server

import (
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	swaggo "github.com/swaggo/http-swagger/v2"
	"net/http"
	"strconv"
)

type Server struct {
	router *mux.Router
	server *http.Server
}

type GroupRouter interface {
	AddRoute(r *mux.Router)
}

func New(withSwagger, withCors bool) *Server {
	srv := &http.Server{}
	srv.Addr = `:` + strconv.Itoa(8000)

	router := mux.NewRouter()

	if withSwagger {
		c := cors.New(cors.Options{
			AllowedOrigins: []string{"*"},
			AllowedMethods: []string{"GET", "HEAD", "POST", "PUT", "OPTIONS", "DELETE"},
		})

		router.PathPrefix("/swagger/").Handler(swaggo.Handler(
			swaggo.DeepLinking(true),
			swaggo.DocExpansion("none"),
			swaggo.DomID("swagger-ui"),
		)).Methods(http.MethodGet)

		srv.Handler = c.Handler(router)
	} else {
		srv.Handler = router
	}

	return &Server{
		router: router,
		server: srv,
	}
}

func (s *Server) CreateRoutes(groupRoutes ...GroupRouter) {
	for _, route := range groupRoutes {
		route.AddRoute(s.router)
	}
}

func (s *Server) Start() error {
	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("server error: %s", err.Error())
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
