package httprequestlogger

import (
	"bytes"
	"io"
	"io/ioutil"
	"strings"
	"testing"
)

func TestMemoReadCloser(t *testing.T) {
	const data = "abcdefghijklmnopqrstuvwxyz"

	next := ioutil.NopCloser(strings.NewReader(data))
	reader := NewMemoReadCloser(next)
	resultWriter := &bytes.Buffer{}
	io.Copy(resultWriter, reader)

	result := string(resultWriter.Bytes())
	t.Logf("Read result = %#v", result)
	if result != data {
		t.Errorf("unexpected Read result: %#v", result)
	}

	result = string(reader.BytesRead())
	t.Logf("BytesRead result = %#v", result)
	if result != data {
		t.Errorf("unexpected BytesRead result: %#v", result)
	}
}
