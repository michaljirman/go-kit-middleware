package httpmiddleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewConditionalMiddleware(t *testing.T) {
	condition := false
	triggered := false

	var httpHandler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Log("httpHandler has triggered")
	})

	targetMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			triggered = true
			next.ServeHTTP(w, r)
		})
	}

	httpHandler = Apply(httpHandler, NewConditionMiddleware(targetMiddleware, func() bool { return condition }))

	testServer := httptest.NewServer(httpHandler)
	defer testServer.Close()

	if _, err := http.Post(testServer.URL, "text/plain ", strings.NewReader("body")); err != nil {
		t.Fatalf("failed to make request: %#v", err)
	}
	if triggered {
		t.Error("targetMiddleware was triggered on negative condition")
	}

	condition = true
	if _, err := http.Post(testServer.URL, "text/plain ", strings.NewReader("body")); err != nil {
		t.Fatalf("failed to make request: %#v", err)
	}
	if !triggered {
		t.Error("targetMiddleware wasn't triggered on positive condition")
	}
}
