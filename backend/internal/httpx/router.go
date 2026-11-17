package httpx

import (
	"context"
	"net/http"
	"time"

	"github.com/complyhub/backend/internal/config"
	"github.com/complyhub/backend/internal/logger"
)

type Server struct {
	httpServer *http.Server
	router     *http.ServeMux
	log        *logger.Logger
	cfg        *config.Config
}

func NewServer(cfg *config.Config, log *logger.Logger) *Server {
	router := http.NewServeMux()
	s := &Server{
		router: router,
		log:    log,
		cfg:    cfg,
		httpServer: &http.Server{
			Addr:              ":" + cfg.Port,
			Handler:           router,
			ReadHeaderTimeout: 5 * time.Second,
			ReadTimeout:       15 * time.Second,
			WriteTimeout:      30 * time.Second,
			IdleTimeout:       60 * time.Second,
		},
	}
	s.setupMiddleware()
	return s
}

func (s *Server) setupMiddleware() {
	// Order: request ID -> recover -> logging -> security headers -> CORS -> body limit -> auth -> permissions
	s.router.HandleFunc("GET /healthz", s.Healthz)
	s.router.HandleFunc("GET /readyz", s.Readyz)
}

func (s *Server) Healthz(w http.ResponseWriter, r *http.Request) {
	WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"}, getRequestID(r.Context()))
}

func (s *Server) Readyz(w http.ResponseWriter, r *http.Request) {
	WriteJSON(w, http.StatusOK, map[string]string{"status": "ready"}, getRequestID(r.Context()))
}

func (s *Server) Router() *http.ServeMux {
	return s.router
}

func (s *Server) Start() error {
	s.log.Info("starting server", "port", s.cfg.Port)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.log.Info("shutting down server")
	return s.httpServer.Shutdown(ctx)
}

func (s *Server) WithAuth(authMiddleware func(http.Handler) http.Handler) *Server {
	// Wrap the router with auth middleware for protected routes
	// This is a placeholder - actual implementation will wrap specific routes
	return s
}

func (s *Server) RequirePermission(permission string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			permissions := GetPermissions(r.Context())
			hasPermission := false
			for _, p := range permissions {
				if p == permission || p == "admin" {
					hasPermission = true
					break
				}
			}
			if !hasPermission {
				WriteError(w, http.StatusForbidden, "FORBIDDEN", "insufficient permissions", nil, getRequestID(r.Context()))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}