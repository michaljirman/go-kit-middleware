package httprequestlogger

import "net/http"

type MemoResponseWriter struct {
	next       http.ResponseWriter
	writer     *MemoWriter
	statusCode int
}

func NewMemoResponseWriter(next http.ResponseWriter) *MemoResponseWriter {
	return &MemoResponseWriter{
		next:       next,
		writer:     NewMemoWriter(next),
		statusCode: http.StatusOK,
	}
}

func (w *MemoResponseWriter) Header() http.Header {
	return w.next.Header()
}

func (w *MemoResponseWriter) Write(p []byte) (int, error) {
	return w.writer.Write(p)
}

func (w *MemoResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.next.WriteHeader(statusCode)
}

func (w *MemoResponseWriter) BytesWritten() []byte {
	return w.writer.BytesWritten()
}

func (w *MemoResponseWriter) StatusCode() int {
	return w.statusCode
}
