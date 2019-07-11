package traceid

import (
	"net/http"

	"glint-backend/lib/go/jsonrpc2"
	"glint-backend/lib/go/jsonrpc2/metadata"
)

func CompatMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var newRequest *http.Request

		ctx := r.Context()

		if traceID, ok := FromContext(ctx); ok {
			// trace id is set, add jsonrpc metadata to context
			md := metadata.Metadata{}
			md = md.Add(jsonrpc2.TraceIdKey, traceID)
			newCtx := metadata.NewContext(ctx, md)
			newRequest = r.WithContext(newCtx)
		} else {
			// trace id isn't set, do nothing
			newRequest = r
		}

		next.ServeHTTP(w, newRequest)
	})
}
