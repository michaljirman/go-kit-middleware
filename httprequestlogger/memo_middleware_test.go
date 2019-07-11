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

func TestMemoMiddleware(t *testing.T) {
	const input = "input"
	const output = "output"

	var bytesRead, bytesWritten []byte
	var statusCode int
	handlerCalled := false

	handler := func(ctx context.Context, r *http.Request, requestContent, responseContent []byte, responseCode int, responseHeader http.Header) {
		t.Logf("ctx = %#v", ctx)
		t.Logf("request = %#v", r)
		t.Logf("requestContent = %#v", requestContent)
		t.Logf("responseContent = %#v", responseContent)
		t.Logf("responseCode = %#v", responseCode)
		t.Logf("responseHeader = %#v", responseHeader)

		handlerCalled = true
		bytesRead = requestContent
		bytesWritten = responseContent
		statusCode = responseCode
	}

	var httpHandler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("key", "value")
		w.WriteHeader(http.StatusNotFound)
		io.Copy(&bytes.Buffer{}, r.Body)
		fmt.Fprint(w, output)
	})

	httpHandler = NewMemoMiddleware(handler)(httpHandler)

	testServer := httptest.NewServer(httpHandler)
	defer testServer.Close()

	resp, err := http.Post(testServer.URL, "text/plain", strings.NewReader(input))
	if err != nil {
		t.Fatalf("failed to make POST request: %#v", err)
	}
	writer := &bytes.Buffer{}
	io.Copy(writer, resp.Body)
	responseBody := string(writer.Bytes())
	if responseBody != output {
		t.Errorf("invalid response body: %#v", responseBody)
	}

	if !handlerCalled {
		t.Errorf("handler wasn't called")
	}
	if statusCode != http.StatusNotFound {
		t.Errorf("unexpected statusCode: %#v", statusCode)
	}

	caughtInput := string(bytesRead)
	t.Logf("caughtInput = %#v", caughtInput)
	caughtOutput := string(bytesWritten)
	t.Logf("caughtOutput = %#v", caughtOutput)

	if caughtInput != input {
		t.Errorf("invalid caught input: %#v", caughtInput)
	}
	if caughtOutput != output {
		t.Errorf("invalid caught output: %#v", caughtOutput)
	}
}
