package log

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFieldsWithRequestTracing(t *testing.T) {
	fields := Fields{
		"resource":    "resource_id",
		"test-number": 1,
	}

	// First test nil Request
	f := fields.WithRequestTracing(nil)
	assert.Nil(t, f["system"])

	// Next with Request but no headers
	req := httptest.NewRequest(http.MethodGet, "http://example.com/foo", nil)

	f = fields.WithRequestTracing(req)
	assert.NotNil(t, f["system"])
	system, ok := f["system"].(Fields)
	assert.True(t, ok)
	assert.NotNil(t, system["request_id"])
	assert.Equal(t, "", system["request_id"])

	// Finally with headers set
	req = httptest.NewRequest(http.MethodGet, "http://example.com/foo", nil)
	req.Header.Add(TraceIDHeader, "trace_123_id")
	req.Header.Add(RequestIDHeader, "request_456_id")
	req.Header.Add(CorrelationIDHeader, "correlation_789_id")

	f = fields.WithRequestTracing(req)
	assert.NotNil(t, f["system"])
	system, ok = f["system"].(Fields)
	assert.True(t, ok)
	assert.NotNil(t, system["request_id"])
	assert.Equal(t, "request_456_id", system["request_id"])
}

func TestFieldsWithAuthenticatedUserTracing(t *testing.T) {
	fields := Fields{
		"resource":    "resource_id",
		"test-number": 1,
	}

	// First test nil AuthPayload
	f := fields.WithAuthenticatedUserTracing(nil)
	assert.Nil(t, f["authentication"])

	// Next with empty Auth Payload
	auth := &AuthPayload{}

	f = fields.WithAuthenticatedUserTracing(auth)
	assert.NotNil(t, f["authentication"])
	subdoc, ok := f["authentication"].(Fields)
	assert.True(t, ok)
	assert.NotNil(t, subdoc["user_id"])
	assert.Equal(t, "", subdoc["user_id"])

	// Finally with Auth Payload set
	auth = &AuthPayload{
		CustomerAccountID: "account_123_id",
		RealUserID:        "real_456_id",
		UserID:            "user_789_id",
	}

	f = fields.WithAuthenticatedUserTracing(auth)
	assert.NotNil(t, f["authentication"])
	subdoc, ok = f["authentication"].(Fields)
	assert.True(t, ok)
	assert.NotNil(t, subdoc["user_id"])
	assert.Equal(t, "user_789_id", subdoc["user_id"])
}
