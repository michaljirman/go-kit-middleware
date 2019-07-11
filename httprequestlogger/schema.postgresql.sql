CREATE TABLE IF NOT EXISTS request_log
(
  id        BIGSERIAL NOT NULL,
  trace_id  TEXT NULL,
  timestamp TIMESTAMPTZ(6) NULL,
  duration  BIGINT NULL,

  request_protocol       TEXT NULL,
  request_method         TEXT NULL,
  request_url            TEXT NULL,
  request_remote_address TEXT NULL,
  request_content_length BIGINT NULL,
  request_header         JSONB NULL,
  request_body           TEXT NULL,

  response_header  JSONB NULL,
  response_body    TEXT NULL,
  response_code    INTEGER NULL,

  CONSTRAINT pk_request_log PRIMARY KEY (id)
);