package auth

import (
	"context"
	"errors"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
)

func NewAuthClient() (*auth.Client, error) {
	// The Firebase Admin SDK automatically detects the FIREBASE_AUTH_EMULATOR_HOST
	// environment variable and connects to the emulator.
	projectId := os.Getenv("PROJECT_ID")
	if projectId == "" {
		return nil, errors.New("PROJECT_ID is not set")
	}
	app, err := firebase.NewApp(context.Background(), &firebase.Config{
		ProjectID: projectId,
	})
	if err != nil {
		return nil, err
	}
	jwtVerifier, err := app.Auth(context.Background())
	if err != nil {
		return nil, err
	}
	return jwtVerifier, nil
}