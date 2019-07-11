package traceid

import (
	"context"
	"testing"
)

func TestFromContext(t *testing.T) {
	ctx := context.Background()

	if _, ok := FromContext(ctx); ok {
		t.Error("FromContext incorrectly indicates presence of trace id in empty context")
	}

	traceID := Generate()
	t.Logf("traceID = %#v", traceID)
	ctx = NewContext(context.Background(), traceID)

	traceID2, ok := FromContext(ctx)
	if !ok {
		t.Error("FromContext incorrectly indicates absence of trace id")
	}
	if traceID2 != traceID {
		t.Errorf("unexpected value of trace id: %#v", traceID2)
	}
}

func TestFromContextGenerateOnAbsence(t *testing.T) {
	ctx := context.Background()
	traceID := FromContextGenerateOnAbsence(ctx)
	t.Logf("traceID = %#v", traceID)
	if traceID == "" {
		t.Error("FromContextGenerateOnAbsence returns empty trace id on empty context")
	}

	traceID = Generate()
	ctx = NewContext(context.Background(), traceID)
	traceID2 := FromContextGenerateOnAbsence(ctx)
	if traceID2 != traceID {
		t.Errorf("unexpected value of trace id: %#v", traceID2)
	}
}
