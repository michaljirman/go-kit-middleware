package httprequestlogger

import (
	"context"
	"net/http"

	"glint-backend/lib/go/httpmiddleware"
)

type MemoMiddlewareHandler func(ctx context.Context, request *http.Request, requestContent, responseContent []byte,
	responseCode int, responseHeader http.Header)

func NewMemoMiddleware(handler MemoMiddlewareHandler) httpmiddleware.Middleware {
	return (&memoMiddleware{handler: handler}).do
}

type memoMiddleware struct {
	handler MemoMiddlewareHandler
}

func (m *memoMiddleware) do(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reader := NewMemoReadCloser(r.Body)
		r.Body = reader

		writer := NewMemoResponseWriter(w)

		next.ServeHTTP(writer, r)

		defer func() {
			m.handler(r.Context(), r, reader.BytesRead(), writer.BytesWritten(), writer.StatusCode(), w.Header())
		}()
	})
}
