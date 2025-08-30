package auth

import (
	"context"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
)

func NewAuthClient() (*auth.Client, error) {
	// The Firebase Admin SDK automatically detects the FIREBASE_AUTH_EMULATOR_HOST
	// environment variable and connects to the emulator.
	var config *firebase.Config
	projectId := os.Getenv("PROJECT_ID")
	if projectId != "" {
		config = &firebase.Config{
			ProjectID: projectId,
		}
	}
	app, err := firebase.NewApp(context.Background(), config)
	if err != nil {
		return nil, err
	}
	jwtVerifier, err := app.Auth(context.Background())
	if err != nil {
		return nil, err
	}
	return jwtVerifier, nil
}
