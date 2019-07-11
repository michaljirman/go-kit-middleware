package httpmiddleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewHeaderRenamingMiddleware(t *testing.T) {
	var httpHandler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Logf("request header: %#v", r.Header)
		val, ok := r.Header["User-ID"]
		t.Logf("User-ID value: %#v", val)
		if !ok {
			t.Error("User-ID header is absent")
		}
		if val[0] != "value" {
			t.Errorf("User-ID value is incorrect: %#v", val)
		}
		if _, ok := r.Header["User-Id"]; ok {
			t.Error("User-Id header is still present")
		}
	})

	httpHandler = Apply(httpHandler, NewHeaderRenamingMiddleware(map[string]string{"User-Id": "User-ID", "Some-Header": "SOME-HEADER"}))

	testServer := httptest.NewServer(httpHandler)
	defer testServer.Close()

	request, err := http.NewRequest(http.MethodPost, testServer.URL, strings.NewReader("request"))
	if err != nil {
		t.Fatalf("failed to prepare request %#v", err)
	}
	request.Header.Set("User-Id", "value")

	if _, err := (&http.Client{}).Do(request); err != nil {
		t.Fatalf("failed to make POST request: %#v", err)
	}
}
