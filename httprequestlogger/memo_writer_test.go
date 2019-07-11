package httprequestlogger

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestMemoWriter(t *testing.T) {
	const data = "abcdefghijklmnopqrstuvwxyz"

	reader := strings.NewReader(data)
	next := bytes.Buffer{}
	writer := NewMemoWriter(&next)
	io.Copy(writer, reader)

	result := string(next.Bytes())
	t.Logf("Write result = %#v", result)
	if result != data {
		t.Errorf("unexpected Read result: %#v", result)
	}

	result = string(writer.BytesWritten())
	t.Logf("BytesWritten result = %#v", result)
	if result != data {
		t.Errorf("unexpected BytesWritten result: %#v", result)
	}
}
