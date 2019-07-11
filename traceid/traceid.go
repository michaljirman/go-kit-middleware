package traceid

import (
	"context"

	"github.com/rs/xid"
)

type key struct{}

var traceIDKey = key{}

func Generate() string {
	return xid.New().String()
}

func FromContext(ctx context.Context) (string, bool) {
	traceID, ok := ctx.Value(traceIDKey).(string)
	return traceID, ok
}

func FromContextGenerateOnAbsence(ctx context.Context) string {
	if traceID, ok := FromContext(ctx); ok {
		return traceID
	}
	return Generate()
}

func NewContext(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDKey, traceID)
}
