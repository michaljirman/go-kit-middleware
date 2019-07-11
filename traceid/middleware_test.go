package traceid

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMiddleware(t *testing.T) {
	var httpHandler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if traceID, ok := FromContext(ctx); !ok {
			t.Error("trace id wasn't set in context")
		} else {
			t.Logf("traceID = %#v", traceID)
		}
	})

	httpHandler = Middleware(httpHandler)

	testServer := httptest.NewServer(httpHandler)
	defer testServer.Close()

	if _, err := http.Post(testServer.URL, "text/json", strings.NewReader("")); err != nil {
		t.Fatalf("failed to make request: %#v", err)
	}
}

func TestMiddleware_TraceIDIsAlreadySet(t *testing.T) {
	traceID := Generate()

	var httpHandler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if traceID2, ok := FromContext(ctx); !ok {
			t.Error("trace id wasn't set in context")
		} else {
			t.Logf("traceID2 = %#v", traceID2)
			if traceID2 != traceID {
				t.Errorf("trace id isn't equal to original one: %#v", traceID2)
			}
		}
	})

	httpHandler = Middleware(httpHandler)

	httpHandler = func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = NewContext(ctx, traceID)
			newRequest := r.WithContext(ctx)
			next.ServeHTTP(w, newRequest)
		})
	}(httpHandler)

	testServer := httptest.NewServer(httpHandler)
	defer testServer.Close()

	if _, err := http.Post(testServer.URL, "text/json", strings.NewReader("")); err != nil {
		t.Fatalf("failed to make request: %#v", err)
	}
}
