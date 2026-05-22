package server

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"pswitch/internal/adminui"
	"pswitch/internal/config"
	"pswitch/internal/pool"
	"pswitch/internal/protocol/anthropic"
	"pswitch/internal/protocol/openai"
	pruntime "pswitch/internal/runtime"
)

const AdminPrefix = "/dashboard"

func NewRouter(manager *pruntime.Manager, adminToken string) http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.Recoverer)

	adminHandler := adminui.New(manager, adminToken)
	router.Handle(AdminPrefix, http.RedirectHandler(AdminPrefix+"/", http.StatusPermanentRedirect))
	router.Handle(AdminPrefix+"/", http.StripPrefix(AdminPrefix, adminHandler))
	router.Handle(AdminPrefix+"/*", http.StripPrefix(AdminPrefix, adminHandler))

	// Health check endpoint for container orchestrators (K8s, Docker, etc.)
	router.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	router.Handle("/*", &dispatcher{manager: manager})

	return router
}

type dispatcher struct {
	manager *pruntime.Manager
}

func (d *dispatcher) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cfg, providerPool := d.manager.Snapshot()
	route, ok := matchRoute(cfg.Routes, r.URL.Path)
	if !ok {
		http.NotFound(w, r)
		return
	}

	openAIHandler := openai.NewHandler(providerPool, openai.Options{
		Mode:    pool.Mode(cfg.Mode),
		Metrics: d.manager.Metrics(),
	})

	var handler http.Handler
	switch route.Kind {
	case "openai":
		handler = openAIHandler
	case "anthropic":
		handler = anthropic.NewHandler(anthropic.Options{
			Model:         route.Model,
			UpstreamModel: route.UpstreamModel,
			Upstream:      openAIHandler,
		})
	default:
		http.NotFound(w, r)
		return
	}

	if route.Prefix != "/" {
		handler = http.StripPrefix(route.Prefix, handler)
	}
	handler.ServeHTTP(w, r)
}

func matchRoute(routes []config.Route, requestPath string) (config.Route, bool) {
	var (
		match    config.Route
		matchLen = -1
	)
	for _, route := range routes {
		if !pathMatches(route.Prefix, requestPath) {
			continue
		}
		if len(route.Prefix) > matchLen {
			match = route
			matchLen = len(route.Prefix)
		}
	}
	return match, matchLen >= 0
}

func pathMatches(prefix, requestPath string) bool {
	if prefix == "/" {
		return true
	}
	return requestPath == prefix || strings.HasPrefix(requestPath, prefix+"/")
}
