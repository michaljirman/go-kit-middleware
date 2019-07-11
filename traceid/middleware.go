package traceid

import (
	"net/http"
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var newRequest *http.Request

		ctx := r.Context()

		if _, ok := FromContext(ctx); ok {
			// trace id is already set, do nothing
			newRequest = r
		} else {
			// trace id isn't set, create new one
			traceID := Generate()
			newCtx := NewContext(ctx, traceID)
			newRequest = r.WithContext(newCtx)
		}

		next.ServeHTTP(w, newRequest)
	})
}
