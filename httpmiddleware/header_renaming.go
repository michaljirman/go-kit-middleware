package httpmiddleware

import "net/http"

func NewHeaderRenamingMiddleware(renameTable map[string]string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for from, to := range renameTable {
				if _, ok := r.Header[from]; ok {
					r.Header[to] = r.Header[from]
					delete(r.Header, from)
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

func NewSingleHeaderRenamingMiddleware(from, to string) Middleware {
	return NewHeaderRenamingMiddleware(map[string]string{from: to})
}
