package log

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNilRequestIsEmpty(t *testing.T) {
	ctx := context.Background()

	ctx = ContextWithRequest(ctx, nil)

	request, ok := RequestIDsFromContext(ctx)
	assert.True(t, ok)
	assert.NotNil(t, request)
	assert.Equal(t, "", request.TraceID)
	assert.Equal(t, "", request.RequestID)
	assert.Equal(t, "", request.CorrelationID)
}

func TestRequestIsValid(t *testing.T) {
	ctx := context.Background()

	// create a dummy request
	req := httptest.NewRequest(http.MethodGet, "http://example.com/foo", nil)
	req.Header.Add(TraceHeader, "trace_123_id")
	req.Header.Add(RequestHeader, "request_456_id")
	req.Header.Add(CorrelationHeader, "correlation_789_id")

	ctx = ContextWithRequest(ctx, req)

	request, ok := RequestIDsFromContext(ctx)
	assert.True(t, ok)
	assert.NotNil(t, request)
	assert.Equal(t, "trace_123_id", request.TraceID)
	assert.Equal(t, "request_456_id", request.RequestID)
	assert.Equal(t, "correlation_789_id", request.CorrelationID)
}
