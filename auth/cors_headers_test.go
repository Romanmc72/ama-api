package auth_test

import (
	"ama/api/auth"
	"ama/api/test"
	"maps"
	"net/http"
	"testing"
)

func TestCORSHeaders(t *testing.T) {
	corsHeaders := auth.GetCorsHeaders()
	testCases := []struct {
		name        string
		context     *test.MockAPIContext
		wantHeaders map[string]string
		wantCode    int
	}{
		{
			name: "CORS headers set on OPTIONS call",
			context: test.NewMockAPIContext(test.MockAPIContextConfig{
				Request: http.Request{
					Method: http.MethodOptions,
				},
			}),
			wantHeaders: corsHeaders,
			wantCode:    http.StatusNoContent,
		},
		{
			name: "CORS headers set on NON-OPTIONS call",
			context: test.NewMockAPIContext(test.MockAPIContextConfig{
				Request: http.Request{
					Method: http.MethodPost,
				},
			}),
			wantHeaders: corsHeaders,
			wantCode:    0,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			auth.CORSHeaders(tc.context)
			if !maps.Equal(tc.context.GetTestHeaders(), tc.wantHeaders) {
				t.Errorf("Wanted CORS headers: '%v'; got: '%v'", tc.wantHeaders, tc.context.GetTestHeaders())
			}
			if tc.wantCode != tc.context.ResponseCode {
				t.Errorf("Wanted response code: %d; got: %d", tc.wantCode, tc.context.ResponseCode)
			}
		})
	}
}
