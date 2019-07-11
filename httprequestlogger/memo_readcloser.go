package httprequestlogger

import (
	"bytes"
	"io"
)

type MemoReadCloser struct {
	next io.ReadCloser
	buf  bytes.Buffer
}

func NewMemoReadCloser(next io.ReadCloser) *MemoReadCloser {
	return &MemoReadCloser{next: next}
}

func (r *MemoReadCloser) Read(p []byte) (int, error) {
	n, err := r.next.Read(p)
	r.buf.Write(p[:n])
	return n, err
}

func (r *MemoReadCloser) Close() error {
	return r.next.Close()
}

func (r *MemoReadCloser) BytesRead() []byte {
	return r.buf.Bytes()
}
