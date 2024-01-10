package requestdata

import (
	"context"
)

type contextKey struct{}

var ctxKey = contextKey{}

// NewContext is for storing data into given context.
func NewContext(ctx context.Context, data RequestHeaders) context.Context {
	return context.WithValue(ctx, ctxKey, &data)
}

// FromContext is for getting data from given context.
func FromContext(ctx context.Context) (RequestHeaders, bool) {
	d := ctx.Value(ctxKey)
	if d == nil {
		return RequestHeaders{}, false
	}
	if data, ok := d.(*RequestHeaders); ok {
		return *data, true
	}
	return RequestHeaders{}, false
}
