package auth_test

import (
	"ama/api/auth"
	"ama/api/constants"
	"ama/api/logging"
	"ama/api/test"
	"net/http"
	"testing"
)

func TestVerifyRequiredScope(t *testing.T) {
	testCases := []struct {
		name           string
		context        *test.MockAPIContext
		claims         any
		requiredScopes map[string]any
		wantErr        bool
	}{
		{
			name:    "All Required Scopes Present",
			context: test.NewMockAPIContext(),
			claims: map[string]any{
				"role": "admin",
				"tier": "premium",
			},
			requiredScopes: map[string]any{
				"role": "admin",
				"tier": "premium",
			},
			wantErr: false,
		},
		{
			name:    "Missing One Scope",
			context: test.NewMockAPIContext(),
			claims: map[string]any{
				"tier": "premium",
			},
			requiredScopes: map[string]any{
				"role": "admin",
				"tier": "premium",
			},
			wantErr: true,
		},
		{
			name:    "Missing All Scopes",
			context: test.NewMockAPIContext(),
			claims:  nil,
			requiredScopes: map[string]any{
				"role": "admin",
				"tier": "premium",
			},
			wantErr: true,
		},
		{
			name:           "No Required Scopes",
			context:        test.NewMockAPIContext(),
			claims:         nil,
			requiredScopes: map[string]any{},
			wantErr:        false,
		},
		{
			name:    "Malformatted claims",
			context: test.NewMockAPIContext(),
			claims:  21,
			requiredScopes: map[string]any{
				"role": "admin",
			},
			wantErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.claims != nil {
				tc.context.Set(constants.AuthTokenClaimsContextKey, tc.claims)
			}

			// TODO set up test logging
			logger := logging.GetLogger()
			auth.VerifyRequiredScope(tc.context, logger, tc.requiredScopes)

			if tc.wantErr {
				if tc.context.ResponseCode != http.StatusForbidden {
					t.Errorf("Expected %d, got %d", http.StatusForbidden, tc.context.ResponseCode)
				}
			} else {
				if tc.context.ResponseCode != 0 {
					t.Errorf("Expected no response code, got %d", tc.context.ResponseCode)
				}
			}
		})
	}
}
