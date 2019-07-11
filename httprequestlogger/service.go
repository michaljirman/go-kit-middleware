package httprequestlogger

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"glint-backend/lib/go/traceid"

	"github.com/lib/pq"
	"github.com/pkg/errors"
)

type Service interface {
	Log(ctx context.Context, request LogRequest) error
}

type LogRequest struct {
	RequestProtocol      string
	RequestMethod        string
	RequestURL           string
	RequestRemoteAddress string
	RequestContentLength int64
	RequestHeader        map[string][]string
	RequestBody          []byte

	ResponseHeader map[string][]string
	ResponseBody   []byte
	ResponseCode   int

	Duration time.Duration
}

func New(db *sql.DB) Service {
	service := service{
		db: db,
	}
	return &service
}

type RequestLogDTO struct {
	ID        int64
	TraceID   string
	Timestamp time.Time
	Duration  time.Duration

	RequestProtocol      string
	RequestMethod        string
	RequestURL           string
	RequestRemoteAddress string
	RequestContentLength int64
	RequestHeader        map[string][]string
	RequestBody          []byte

	ResponseHeader map[string][]string
	ResponseBody   []byte
	ResponseCode   int
}

type service struct {
	db *sql.DB
}

func (s *service) Log(ctx context.Context, request LogRequest) error {
	traceID, _ := traceid.FromContext(ctx)

	// ID and Timestamp are ignored
	dto := RequestLogDTO{
		TraceID:  traceID,
		Duration: request.Duration,

		RequestProtocol:      request.RequestProtocol,
		RequestMethod:        request.RequestMethod,
		RequestURL:           request.RequestURL,
		RequestRemoteAddress: request.RequestRemoteAddress,
		RequestContentLength: request.RequestContentLength,
		RequestHeader:        request.RequestHeader,
		RequestBody:          request.RequestBody,

		ResponseHeader: request.ResponseHeader,
		ResponseBody:   request.ResponseBody,
		ResponseCode:   request.ResponseCode,
	}

	err := s.insertRequest(ctx, dto)
	return errors.Wrap(err, "failed to save request")
}

func (s *service) insertRequest(ctx context.Context, request RequestLogDTO) error {
	row, err := newRequestLogRowFromDTO(request)
	if err != nil {
		return errors.Wrap(err, "failed to create request_log row object from DTO")
	}

	err = row.insert(ctx, s.db)
	return errors.Wrap(err, "failed to insert to request_log table")
}

type requestLogRow struct {
	ID        int64
	TraceID   sql.NullString
	Timestamp pq.NullTime
	Duration  sql.NullInt64

	RequestProtocol      sql.NullString
	RequestMethod        sql.NullString
	RequestURL           sql.NullString
	RequestRemoteAddress sql.NullString
	RequestContentLength sql.NullInt64
	RequestHeader        sql.NullString
	RequestBody          sql.NullString

	ResponseHeader sql.NullString
	ResponseBody   sql.NullString
	ResponseCode   sql.NullInt64
}

func newRequestLogRowFromDTO(dto RequestLogDTO) (*requestLogRow, error) {
	requestHeader, err := json.Marshal(dto.RequestHeader)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal request header")
	}
	responseHeader, err := json.Marshal(dto.ResponseHeader)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal response header")
	}

	row := requestLogRow{
		ID:        dto.ID,
		TraceID:   sql.NullString{String: dto.TraceID, Valid: true},
		Timestamp: pq.NullTime{Time: dto.Timestamp, Valid: true},
		Duration:  sql.NullInt64{Int64: int64(dto.Duration), Valid: true},

		RequestProtocol:      sql.NullString{String: dto.RequestProtocol, Valid: true},
		RequestMethod:        sql.NullString{String: dto.RequestMethod, Valid: true},
		RequestURL:           sql.NullString{String: dto.RequestURL, Valid: true},
		RequestRemoteAddress: sql.NullString{String: dto.RequestRemoteAddress, Valid: true},
		RequestContentLength: sql.NullInt64{Int64: dto.RequestContentLength, Valid: true},
		RequestHeader:        sql.NullString{String: string(requestHeader), Valid: true},
		RequestBody:          sql.NullString{String: string(dto.RequestBody), Valid: true},

		ResponseHeader: sql.NullString{String: string(responseHeader), Valid: true},
		ResponseBody:   sql.NullString{String: string(dto.ResponseBody), Valid: true},
		ResponseCode:   sql.NullInt64{Int64: int64(dto.ResponseCode), Valid: true},
	}

	return &row, nil
}

func (r *requestLogRow) insert(ctx context.Context, db *sql.DB) error {
	const query = `
INSERT INTO request_log (
	trace_id, 
	timestamp, 
	duration, 
	request_protocol, 
	request_method,
	request_url,
	request_remote_address,
	request_content_length,
	request_header,
	request_body,
	response_header,
	response_body,
	response_code
) VALUES (
    -- id should be autoincrement field
	$1, -- trace_id
	current_timestamp(6), -- timestamp
	$2, -- duration
	$3, -- request_protocol
	$4, -- request_method
	$5, -- request_url
	$6, -- request_remote_address
	$7, -- request_content_length
	$8, -- request_header
	$9, -- request_body
	$10, -- response_header
	$11, -- response_body
	$12 -- response_code	
)`
	_, err := db.ExecContext(ctx, query,
		r.TraceID,
		r.Duration,
		r.RequestProtocol,
		r.RequestMethod,
		r.RequestURL,
		r.RequestRemoteAddress,
		r.RequestContentLength,
		r.RequestHeader,
		r.RequestBody,
		r.ResponseHeader,
		r.ResponseBody,
		r.ResponseCode,
	)

	return errors.Wrap(err, "failed to execute insert statement")
}
