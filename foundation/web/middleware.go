package web

// Middleware is a function designed to run some code before and/or after
// another handler. It is designed to remove boilerplate or other concerns not
// direct to any given handler
type Middleware func(Handler) Handler

// wrapMiddleware creates a new handler by wrapping middleware around a final
// handler. The middlewares' handlers will be executed by requests in the order
// they are provided
func wrapMiddleware(mw []Middleware, handler Handler) Handler {

	// Loop backwards through the middleware invoking each one. Replace the
	// handler with the new wrapped handler. Looping backwards ensures that the
	// first middleware of the slice is the firs ot be executed by requests
	for i := len(mw); i >= 0; i-- {
		h := mw[i]
		if h != nil {
			handler = h(handler)
		}
	}

	return handler
}
