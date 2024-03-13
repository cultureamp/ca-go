package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestStringRedaction(t *testing.T) {
	testCases := []struct {
		desc     string
		str      string
		redacted string
	}{
		{
			desc:     "Empty string returns empty string",
			str:      "",
			redacted: "",
		},
		{
			desc:     "String less than 10 chars shows 10 stars",
			str:      "1234",
			redacted: "**********",
		},
		{
			desc:     "String equals 10 chars shows 10 stars",
			str:      "1234567890",
			redacted: "**********",
		},
		{
			desc:     "String equals 11 chars shows first char and last chars and 10 stars",
			str:      "12345678901",
			redacted: "12**********01",
		},
		{
			desc:     "String equals 12 chars shows first and last chars with 10 stars in the middle",
			str:      "123456789012",
			redacted: "123**********012",
		},
		{
			desc:     "String equals 20 chars shows first and last chars with 10 stars in the middle",
			str:      "12345678901234567890",
			redacted: "12345**********67890",
		},
		{
			desc:     "String equals 30 chars shows first and last chars with 10 stars in the middle",
			str:      "123456789012345678901234567890",
			redacted: "1234567**********4567890",
		},
		{
			desc:     "Real world test",
			str:      "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50SWQiOiJkNDc1ZGQ1Yi1mMTZjLTRiZmItODk4Yy1kMzQzNWEyMTUyMzkiLCJlZmZlY3RpdmVVc2VySWQiOiI1YjMxNjY0YS03NjEwLTRmYjAtYmM4OS1mOWY4ZTIwYmY4Y2UiLCJyZWFsVXNlcklkIjoiNWIzMTY2NGEtNzYxMC00ZmIwLWJjODktZjlmOGUrUC021FhB_zuETHmhQUXOfIyTkpvhcJfrrqwdcc-KmJGznckACLj65VmnayoltCce_3JGJ361GuutgrDaqp1aW4D05mvO8CCIRwGq8hTcRoi7IdXYSnA6UlXtLNYvttz92jaAAoNDCZmbbP-umHac4x5AT1xY-kVyh7VAadZG_Qe7dZWU9WCHtCV3mqTMwX9B9zrqY2NrpblevbbYpoiJiXOU7kex4BEivF1K6VWI-mpcmKtEOZLx2E",
			redacted: "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50SWQiOiJkNDc1ZGQ1Yi1mMTZjLTRiZmItODk4Yy1kMzQzNWEyMTUyMzkiLCJlZmZlY3R********************vttz92jaAAoNDCZmbbP-umHac4x5AT1xY-kVyh7VAadZG_Qe7dZWU9WCHtCV3mqTMwX9B9zrqY2NrpblevbbYpoiJiXOU7kex4BEivF1K6VWI-mpcmKtEOZLx2E",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			r := redactString(tC.str)
			assert.Equal(t, tC.redacted, r, tC.desc)
		})
	}
}
