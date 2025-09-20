package auth_test

import (
	"ama/api/auth"
	"ama/api/constants"
	"ama/api/logging"
	"ama/api/test"
	"errors"
	"net/http"
	"reflect"
	"testing"

	firebaseAuth "firebase.google.com/go/v4/auth"
)

func TestVerifyToken(t *testing.T) {
	userId := "user123"
	claims := map[string]any{"role": "admin"}
	testCases := []struct {
		name        string
		context     *test.MockAPIContext
		tokenString string
		verifier    *test.MockJWTVerifier
		wantErr     bool
	}{
		{
			name:        "Valid Token",
			context:     test.NewMockAPIContext(test.MockAPIContextConfig{}),
			tokenString: "Bearer valid-token",
			verifier: test.NewMockJWTVerifier(&firebaseAuth.Token{
				UID:    userId,
				Claims: claims,
			}, nil),
			wantErr: false,
		},
		{
			name:        "Valid Token - No Bearer Prefix",
			context:     test.NewMockAPIContext(test.MockAPIContextConfig{}),
			tokenString: "valid-token",
			verifier: test.NewMockJWTVerifier(&firebaseAuth.Token{
				UID:    userId,
				Claims: claims,
			}, nil),
			wantErr: false,
		},
		{
			name:        "Token Invalid",
			context:     test.NewMockAPIContext(test.MockAPIContextConfig{}),
			tokenString: "Bearer looks-like-a-valid-token",
			verifier: test.NewMockJWTVerifier(&firebaseAuth.Token{
				UID:    userId,
				Claims: claims,
			}, errors.New("invalid token")),
			wantErr: true,
		},
		{
			name:        "Token Blank",
			context:     test.NewMockAPIContext(test.MockAPIContextConfig{}),
			tokenString: "",
			verifier:    test.NewMockJWTVerifier(nil, nil),
			wantErr:     true,
		},
		{
			name:        "Token nil",
			context:     test.NewMockAPIContext(test.MockAPIContextConfig{}),
			tokenString: "valid-token",
			verifier:    test.NewMockJWTVerifier(nil, nil),
			wantErr:     true,
		},
		{
			name:        "Blank UID",
			context:     test.NewMockAPIContext(test.MockAPIContextConfig{}),
			tokenString: "valid-token",
			verifier: test.NewMockJWTVerifier(&firebaseAuth.Token{
				UID:    "",
				Claims: claims,
			}, nil),
			wantErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.tokenString != "" {
				tc.context.Header(constants.AuthorizationHeader, tc.tokenString)
			}
			logger := logging.GetLogger()
			auth.VerifyToken(tc.context, tc.verifier, logger)
			if tc.wantErr {
				if tc.context.ResponseCode != http.StatusUnauthorized {
					t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, tc.context.ResponseCode)
				}
			} else {
				if tc.context.ResponseCode != 0 {
					t.Errorf("Did not expect an error, but got status %d", tc.context.ResponseCode)
				}
				if uid, exists := tc.context.Get(constants.AuthTokenUserIdContextKey); !exists || uid != userId {
					t.Errorf("Expected user ID '%s', got '%s'", userId, uid)
				}
				if c, exists := tc.context.Get(constants.AuthTokenClaimsContextKey); !exists || !reflect.DeepEqual(c, claims) {
					t.Errorf("Expected claims '%v', got '%v'", claims, c)
				}
			}
		})
	}
}
