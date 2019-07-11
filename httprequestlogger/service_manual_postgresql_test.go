// +build manual

// to run this test:
//  	docker run --rm -it -e POSTGRES_PASSWORD=postgres -p 15432:5432 postgres:9.6.3
// 		go test glint-backend/lib/go/httprequestlogger/... -v -tags manual

package httprequestlogger

// to prepare postgres for tests:
//  docker run --rm -it -e POSTGRES_PASSWORD=postgres -p 15432:5432 postgres:9.6.3

import (
	"context"
	"database/sql"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"glint-backend/lib/go/traceid"

	_ "github.com/lib/pq"
)

const dsn = "postgres://postgres:postgres@127.0.0.1:15432/postgres?sslmode=disable"

func TestService_Log_ManualPostgresql(t *testing.T) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("failed to initialise DB: %#v", err)
	}

	schema, err := ioutil.ReadFile("schema.postgresql.sql")
	if err != nil {
		t.Fatalf("failed to read schema file: %#v", err)
	}

	t.Log(db)
	t.Log(schema)

	if _, err := db.Exec(string(schema)); err != nil {
		t.Fatalf("failed to apply schema: %#v", err)
	}

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

	if err := service.Log(ctx, request); err != nil {
		t.Errorf("service.Log failed with %#v", err)
	}

}
