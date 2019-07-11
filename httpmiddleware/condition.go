package httpmiddleware

import (
	"net/http"
)

type conditionMiddleware struct {
	targetMiddleware Middleware
	conditionFunc    func() bool
}

func NewConditionMiddleware(targetMiddleware Middleware, conditionFunc func() bool) Middleware {
	return (&conditionMiddleware{
		targetMiddleware: targetMiddleware,
		conditionFunc:    conditionFunc,
	}).do
}

func (m *conditionMiddleware) do(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		middleware := NopMiddleware
		if m.conditionFunc() {
			middleware = m.targetMiddleware
		}
		middleware(next).ServeHTTP(w, r)
	})
}

func NewBoolConditionMiddleware(targetMiddleware Middleware, flag bool) Middleware {
	return NewConditionMiddleware(targetMiddleware, func() bool { return flag })
}
