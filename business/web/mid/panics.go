package mid

import (
	"context"
	"fmt"
	"github.com/yakushou730/ardanlabs-ultimate-serice-v3/foundation/web"
	"net/http"
	"runtime/debug"
)

// Panics recovers from panics and converts the panic to an error, so it is
// reported in Metrics and handled in Errors
func Panics() web.Middleware {

	// This is the actual middleware function to be executed
	m := func(handler web.Handler) web.Handler {

		// Create the handler that will be attracted in the middleware chain
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) (err error) {

			// Defer a function to recover from a panic and set the error return
			// variable after the fact
			defer func() {
				if rec := recover(); rec != nil {

					// Stack trace will be provided
					trace := debug.Stack()
					err = fmt.Errorf("PANIC [%v] TRACE [%s]", rec, string(trace))
				}
			}()

			// Return the error, so it can be handled further up the chain
			return handler(ctx, w, r)
		}

		return h

	}

	return m
}
