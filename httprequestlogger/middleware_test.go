package httprequestlogger

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMiddleware(t *testing.T) {
	const requestContent = "input"
	const responseContent = "output"

	mock := &mockService{}

	var httpHandler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("key2", "value2")
		w.WriteHeader(http.StatusNotFound)
		io.Copy(&bytes.Buffer{}, r.Body)
		fmt.Fprint(w, responseContent)
	})

	httpHandler = NewMiddleware(mock)(httpHandler)

	testServer := httptest.NewServer(httpHandler)
	defer testServer.Close()

	req, err := http.NewRequest(http.MethodPost, testServer.URL, strings.NewReader(requestContent))
	if err != nil {
		t.Fatalf("failed to prepare request %#v", err)
	}
	req.Header.Set("key1", "value1")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("failed to make POST request: %#v", err)
	}

	buf := &bytes.Buffer{}
	io.Copy(buf, resp.Body)
	responseBody := string(buf.Bytes())
	if responseBody != responseContent {
		t.Errorf("invalid response body: %#v", responseBody)
	}

	if !mock.called {
		t.Fatal("mock.Log wasn't called")
	}

	t.Logf("RequestProtocol = %#v", mock.request.RequestProtocol)
	t.Logf("RequestMethod = %#v", mock.request.RequestMethod)
	t.Logf("RequestURL = %#v", mock.request.RequestURL)
	t.Logf("RequestRemoteAddress = %#v", mock.request.RequestRemoteAddress)
	t.Logf("RequestContentLength = %#v", mock.request.RequestContentLength)
	t.Logf("RequestHeader = %#v", mock.request.RequestHeader)
	t.Logf("RequestBody = %#v", string(mock.request.RequestBody))
	t.Logf("ResponseHeader = %#v", mock.request.ResponseHeader)
	t.Logf("ResponseBody = %#v", string(mock.request.ResponseBody))
	t.Logf("ResponseCode = %#v", mock.request.ResponseCode)
	t.Logf("Duration = %v", mock.request.Duration)

	if mock.request.RequestProtocol != req.Proto {
		t.Errorf("invalid RequestProtocol = %#v", mock.request.RequestProtocol)
	}
	if mock.request.RequestMethod != req.Method {
		t.Errorf("invalid RequestMethod = %#v", mock.request.RequestMethod)
	}
	if mock.request.RequestURL != "/" {
		t.Errorf("invalid RequestURL = %#v", mock.request.RequestURL)
	}
	// TODO: test RequestRemoteAddress
	if mock.request.RequestContentLength != req.ContentLength {
		t.Errorf("invalid RequestContentLength = %#v", mock.request.RequestContentLength)
	}
	if mock.request.RequestHeader["Key1"][0] != "value1" {
		t.Errorf(`RequestHeader doesn't contain "key1"`)
	}
	if string(mock.request.RequestBody) != requestContent {
		t.Errorf("invalid RequestBody = %#v", string(mock.request.RequestBody))
	}
	if mock.request.ResponseHeader["Key2"][0] != "value2" {
		t.Errorf(`ResponseHeader doesn't contain "key2"`)
	}
	if string(mock.request.ResponseBody) != responseContent {
		t.Errorf("invalid ResponseBody = %#v", string(mock.request.ResponseBody))
	}
	if mock.request.ResponseCode != http.StatusNotFound {
		t.Errorf("invalid ResponseCode = %#v", mock.request.ResponseCode)
	}
	if mock.request.Duration <= 0 {
		t.Errorf("Duration is 0")
	}
}

type mockService struct {
	request LogRequest
	called  bool
}

func (s *mockService) Log(ctx context.Context, request LogRequest) error {
	s.called = true
	s.request = request
	return nil
}
