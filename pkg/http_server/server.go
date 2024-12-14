package http_server

import (
	"company-crud/pkg/logger"
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
	logger *logger.Logger
	server *http.Server
}

type GroupRouter interface {
	AddRoute(r *mux.Router)
}

func New(log *logger.Logger, withSwagger bool, groupRoutes ...GroupRouter) *Server {
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

	for _, route := range groupRoutes {
		route.AddRoute(router)
	}

	return &Server{
		router: router,
		server: srv,
		logger: log,
	}
}

func (s *Server) Start() {
	s.logger.Info(fmt.Sprintf("Listening on: %s", s.server.Addr))
	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.logger.Error(fmt.Sprintf("Server error: %s", err.Error()))
	}
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
