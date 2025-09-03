package test

import (
	"context"

	"firebase.google.com/go/v4/auth"
)

type MockJWTVerifier struct {
	token *auth.Token
	err   error
}

func NewMockJWTVerifier(token *auth.Token, err error) *MockJWTVerifier {
	return &MockJWTVerifier{
		token: token,
		err:   err,
	}
}

func (m *MockJWTVerifier) VerifyIDTokenAndCheckRevoked(c context.Context, tokenString string) (*auth.Token, error) {
	return m.token, m.err
}
