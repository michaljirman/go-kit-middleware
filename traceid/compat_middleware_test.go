package traceid

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"glint-backend/lib/go/httpmiddleware"
	"glint-backend/lib/go/jsonrpc2"
	"glint-backend/lib/go/jsonrpc2/metadata"
)

func TestCompatMiddleware(t *testing.T) {
	var httpHandler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if traceID, ok := FromContext(ctx); !ok {
			t.Error("trace id wasn't set in context")
		} else {
			t.Logf("traceID = %#v", traceID)
			if md, ok := metadata.FromContext(ctx); !ok {
				t.Error("can't get metadata from context")
			} else {
				t.Logf("metadata: %#v", md)
				traceID2, err := md.Get(jsonrpc2.TraceIdKey)
				if err != nil {
					t.Errorf("failed to get trace id from metadata: %#v", err)
				}
				if traceID2 != traceID {
					t.Errorf("trace id from metadata doesn't match normal one: %#v", traceID2)
				}
			}
		}
	})

	httpHandler = httpmiddleware.Apply(httpHandler, Middleware, CompatMiddleware)

	testServer := httptest.NewServer(httpHandler)
	defer testServer.Close()

	if _, err := http.Post(testServer.URL, "text/json", strings.NewReader("")); err != nil {
		t.Fatalf("failed to make request: %#v", err)
	}
}

func TestCompatMiddleware_TraceIDIsNotSet(t *testing.T) {
	var httpHandler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if _, ok := metadata.FromContext(ctx); ok {
			t.Error("it was possible get metadata from empty context")
		}
	})

	httpHandler = CompatMiddleware(httpHandler)

	testServer := httptest.NewServer(httpHandler)
	defer testServer.Close()

	if _, err := http.Post(testServer.URL, "text/json", strings.NewReader("")); err != nil {
		t.Fatalf("failed to make request: %#v", err)
	}
}
