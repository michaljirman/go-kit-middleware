package httprequestlogger

import (
	"context"
	"net/http"
	"time"

	"glint-backend/lib/go/httpmiddleware"
)

func NewMiddleware(service Service) httpmiddleware.Middleware {
	return (&middleware{
		service: service,
	}).do
}

type middleware struct {
	service Service
}

func (m *middleware) do(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		handler := func(ctx context.Context, request *http.Request, requestContent, responseContent []byte, responseCode int, responseHeader http.Header) {
			logRequest := LogRequest{
				RequestProtocol:      r.Proto,
				RequestMethod:        r.Method,
				RequestURL:           r.URL.Path,
				RequestRemoteAddress: r.RemoteAddr,
				RequestContentLength: r.ContentLength,
				RequestHeader:        copyHeader(r.Header),
				RequestBody:          requestContent,

				ResponseHeader: copyHeader(responseHeader),
				ResponseBody:   responseContent,
				ResponseCode:   responseCode,

				Duration: time.Since(startTime),
			}

			if err := m.service.Log(ctx, logRequest); err != nil {
				// TODO: log error
			}
		}

		NewMemoMiddleware(handler)(next).ServeHTTP(w, r)
	})
}

func copyHeader(httpHeader http.Header) map[string][]string {
	headers := map[string][]string{}
	for headerName, headerValues := range httpHeader {
		valuesCopy := make([]string, 0, len(headerValues))
		for _, value := range headerValues {
			valuesCopy = append(valuesCopy, value)
		}
		headers[headerName] = valuesCopy
	}
	return headers
}
