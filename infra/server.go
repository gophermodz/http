package infra

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	_ "net/http/pprof" //nolint: gosec

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Server struct {
	router         *mux.Router
	config         *config
	logger         *zap.Logger
	tracer         trace.Tracer
	tracerProvider *sdktrace.TracerProvider
}

func New(logger *zap.Logger, opts ...Option) (*Server, error) {
	cfg := NewDefaultConfig()
	for _, opt := range opts {
		opt.apply(cfg)
	}
	srv := &Server{
		router: mux.NewRouter(),
		config: cfg,
		logger: logger,
	}

	return srv, nil
}

// Run starts the HTTP server.
func (s *Server) Run(ctx context.Context) error {
	l := s.logger.With(zap.String("host", s.config.Host), zap.Int("port", s.config.Port))
	l.Info("[INFRA-HTTP] server starting")

	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", s.config.Host, s.config.Port),
		Handler:      s.router,
		ReadTimeout:  s.config.ReadTimeout,
		WriteTimeout: s.config.WriteTimeout,
		IdleTimeout:  s.config.IdleTimeout,
	}

	if s.config.TracerEnabled {
		s.initTracer(ctx)
		s.router.Use(NewOpenTelemetryMiddleware(s.config.OtelServiceName))
	}

	s.router.Handle("/metrics", promhttp.Handler())
	s.router.PathPrefix("/debug/pprof/").Handler(http.DefaultServeMux)
	s.router.HandleFunc("/api/info", s.infoHandler).Methods("GET")
	s.router.HandleFunc("/api/healthz", s.healthzHandler).Methods("GET")

	go func() {
		<-ctx.Done()
		l.Info("[INFRA-HTTP] server shutting down")
		if err := srv.Shutdown(ctx); err != nil {
			l.Error("[INFRA-HTTP] server shutdown error", zap.Error(err))
		}
	}()

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}

	return nil
}

// JSONResponse is a helper function to write a JSON response to the client.
func (s *Server) JSONResponse(w http.ResponseWriter, _ *http.Request, result interface{}) {
	body, err := json.Marshal(result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		s.logger.Error("JSON marshal failed", zap.Error(err))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(prettyJSON(body))
}

func prettyJSON(b []byte) []byte {
	var out bytes.Buffer
	_ = json.Indent(&out, b, "", "  ")
	return out.Bytes()
}
