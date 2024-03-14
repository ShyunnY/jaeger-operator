package cmd

import (
	"fmt"
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/ShyunnY/jaeger-operator/internal/config"
)

func NewAdmin(cfg *config.Server) error {
	return StartAdmin(cfg)
}

func StartAdmin(cfg *config.Server) error {

	handlers := http.NewServeMux()
	address := fmt.Sprintf("0.0.0.0:15000")

	// register pprof endpoint
	handlers.HandleFunc("/debug/pprof/", pprof.Index)
	handlers.HandleFunc("/debug/pprof/profile", pprof.Profile)
	handlers.HandleFunc("/debug/pprof/trace", pprof.Trace)
	handlers.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	handlers.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)

	adminServer := &http.Server{
		Addr:              address,
		Handler:           handlers,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       15 * time.Second,
	}

	cfg.Logger.Info("start admin server", "address", address)
	go func() {
		if err := adminServer.ListenAndServe(); err != nil {
			cfg.Logger.Error(err, "failed to start admin server")
		}
	}()

	return nil
}
