// Package handlers contains the full set of handler functions and routes
// supported by the web api
package handlers

import (
	"expvar"
	"github.com/yakushou730/ardanlabs-ultimate-serice-v3/app/services/sales-api/handlers/debug/checkgrp"
	"github.com/yakushou730/ardanlabs-ultimate-serice-v3/app/services/sales-api/handlers/v1/testgrp"
	"github.com/yakushou730/ardanlabs-ultimate-serice-v3/foundation/web"
	"go.uber.org/zap"
	"net/http"
	"net/http/pprof"
	"os"
)

func DebugStandardLibraryMux() *http.ServeMux {
	mux := http.NewServeMux()

	// Register all the standard library debug endpoints
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.Handle("/debug/vars", expvar.Handler())

	return mux
}

func DebugMux(build string, log *zap.SugaredLogger) http.Handler {
	mux := DebugStandardLibraryMux()

	// Register debug check endpoints
	cgh := checkgrp.Handlers{
		Build: build,
		Log:   log,
	}

	mux.HandleFunc("/debug/readiness", cgh.Readiness)
	mux.HandleFunc("/debug/liveness", cgh.Liveness)

	return mux
}

type APIMuxConfig struct {
	Shutdown chan os.Signal
	Log      *zap.SugaredLogger
}

func APIMux(cfg APIMuxConfig) *web.App {
	app := web.NewApp(cfg.Shutdown)

	tgh := testgrp.Handlers{
		Log: cfg.Log,
	}
	app.Handle(http.MethodGet, "/v1/test", tgh.Test)

	return app
}
