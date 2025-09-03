package interfaces

import (
	"context"

	"firebase.google.com/go/v4/auth"
)

type JWTVerifier interface {
	VerifyIDTokenAndCheckRevoked(c context.Context, token string) (*auth.Token, error)
}
