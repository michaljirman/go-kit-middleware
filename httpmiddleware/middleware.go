package httpmiddleware

import "net/http"

type Middleware func(next http.Handler) http.Handler

func Chain(middleware ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		chain := NopMiddleware(next)

		// apply middlewares in reversed order
		for idx := len(middleware) - 1; idx >= 0; idx-- {
			chain = middleware[idx](chain)
		}

		return chain
	}
}

func Apply(handler http.Handler, middleware ...Middleware) http.Handler {
	chain := Chain(middleware...)
	return chain(handler)
}

func NopMiddleware(next http.Handler) http.Handler {
	return next
}
