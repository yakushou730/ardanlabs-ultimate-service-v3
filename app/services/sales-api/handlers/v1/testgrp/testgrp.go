// Package testgrp contains all the test handlers
package testgrp

import (
	"context"
	"errors"
	"github.com/yakushou730/ardanlabs-ultimate-serice-v3/business/sys/validate"
	"github.com/yakushou730/ardanlabs-ultimate-serice-v3/foundation/web"
	"go.uber.org/zap"
	"math/rand"
	"net/http"
)

// Handlers manages the set of check endpoints
type Handlers struct {
	Log *zap.SugaredLogger
}

// Test handler is for development
func (h Handlers) Test(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	if n := rand.Intn(100); n%2 == 0 {
		//return errors.New("untrusted error")
		//return web.NewShutdownError("restart service")
		//panic("testing panic")
		return validate.NewRequestError(errors.New("trusted error"), http.StatusBadRequest)
	}

	status := struct {
		Status string
	}{
		Status: "OK",
	}

	return web.Respond(ctx, w, status, http.StatusOK)
}
