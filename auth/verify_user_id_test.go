package auth_test

import (
	"ama/api/auth"
	"ama/api/constants"
	"ama/api/logging"
	"ama/api/test"
	"net/http"
	"testing"
)

func TestVerifyUserID(t *testing.T) {
	testCases := []struct {
		name    string
		context *test.MockAPIContext
		tokenId string
		userId  string
		wantErr bool
	}{
		{
			name:    "Matching User IDs",
			context: test.NewMockAPIContext(test.MockAPIContextConfig{}),
			tokenId: "user123",
			userId:  "user123",
			wantErr: false,
		},
		{
			name:    "Non-Matching User IDs",
			context: test.NewMockAPIContext(test.MockAPIContextConfig{}),
			tokenId: "user123",
			userId:  "user456",
			wantErr: true,
		},
		{
			name:    "Missing Token User ID",
			context: test.NewMockAPIContext(test.MockAPIContextConfig{}),
			tokenId: "",
			userId:  "user456",
			wantErr: true,
		},
		{
			name:    "Missing Path User ID",
			context: test.NewMockAPIContext(test.MockAPIContextConfig{}),
			tokenId: "user123",
			userId:  "",
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.tokenId != "" {
				tc.context.Set(constants.AuthTokenUserIdContextKey, tc.tokenId)
			}
			if tc.userId != "" {
				tc.context.SetParam(constants.UserIdPathIdentifier, tc.userId)
			}

			logger := logging.GetLogger()
			auth.VerifyUserID(tc.context, logger)
			if tc.wantErr {
				if tc.context.ResponseCode != http.StatusForbidden {
					t.Errorf("Expected status %d, got %d", http.StatusForbidden, tc.context.ResponseCode)
				}
			} else {
				if tc.context.ResponseCode != 0 {
					t.Errorf("Did not expect an error, but got status %d", tc.context.ResponseCode)
				}
			}
		})
	}
}
