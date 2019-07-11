package httprequestlogger

import (
	"bytes"
	"io"
)

type MemoWriter struct {
	next io.Writer
	buf  bytes.Buffer
}

func NewMemoWriter(next io.Writer) *MemoWriter {
	return &MemoWriter{next: next}
}

func (w *MemoWriter) Write(p []byte) (int, error) {
	n, err := w.next.Write(p)
	w.buf.Write(p[:n])
	return n, err
}

func (w *MemoWriter) BytesWritten() []byte {
	return w.buf.Bytes()
}
