package httprequestlogger

import (
	"context"
	"net/http"
	"testing"
	"time"

	"glint-backend/lib/go/traceid"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
)

func TestService_Log(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to initialise sqlmock: %v", err)
	}
	defer db.Close()

	service := New(db)

	request := LogRequest{
		RequestProtocol:      "HTTP/1.1",
		RequestMethod:        http.MethodPost,
		RequestURL:           "/",
		RequestRemoteAddress: "127.0.0.1:54982",
		RequestContentLength: 42,
		RequestHeader:        map[string][]string{"Key1": {"value1"}},
		RequestBody:          []byte("request"),

		ResponseHeader: map[string][]string{"Key2": {"value2"}},
		ResponseBody:   []byte("response"),
		ResponseCode:   http.StatusOK,

		Duration: time.Microsecond * 42,
	}

	ctx := context.Background()
	traceID := traceid.Generate()
	ctx = traceid.NewContext(ctx, traceID)

	mock.ExpectExec("INSERT INTO request_log").WithArgs(
		traceID,
		request.Duration,
		request.RequestProtocol,
		request.RequestMethod,
		request.RequestURL,
		request.RequestRemoteAddress,
		request.RequestContentLength,
		`{"Key1":["value1"]}`,
		"request",
		`{"Key2":["value2"]}`,
		"response",
		request.ResponseCode,
	).WillReturnResult(sqlmock.NewResult(1, 1))

	if err := service.Log(ctx, request); err != nil {
		t.Errorf("service.Log failed with %#v", err)
	}
}
