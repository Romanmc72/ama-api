//go:build integration
// +build integration

package integration_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"ama/api/application"
	"ama/api/auth"
	"ama/api/constants"
)

type returnedToken struct {
	Kind         string `json:"kind"`
	IsNewUser    bool   `json:"isNewUser"`
	IdToken      string `json:"idToken"`
	RefreshToken string `json:"refreshToken"`
	ExpirationIn string `json:"expiresIn"`
}

func TestQuestion(t *testing.T) {
	questionsToCreate := 100
	client := &http.Client{}
	authToken, err := signIn(client)
	if err != nil {
		t.Errorf("failed to sign in: %v", err)
		return
	}
	questionIds := make([]string, questionsToCreate)
	for i := range questionsToCreate {
		questionId, err := createQuestion(i, authToken, *client)
		if err != nil {
			t.Errorf("failed to create question %d: %v", i, err)
			return
		}
		t.Logf("Created question %d with ID: %s", i, questionId)
		questionIds[i] = questionId
	}
}

func signIn(httpClient *http.Client) (string, error) {
	authClient, err := auth.NewAuthClient()
	if err != nil {
		return "", err
	}
	token, err := authClient.CustomTokenWithClaims(
		context.Background(),
		"integration-test-client",
		constants.GetAdminScopes(),
	)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest(
		http.MethodPost,
		// TODO: Get this to be configurable via env var
		"http://localhost:9099/identitytoolkit.googleapis.com/v1/accounts:signInWithCustomToken?key=fake-api-key",
		strings.NewReader(fmt.Sprintf(`{"token": "%s", "returnSecureToken": true}`, token)))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var validToken returnedToken
	err = json.Unmarshal(respBody, &validToken)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %s; err: %w", respBody, err)
	}
	return validToken.IdToken, nil
}

func createQuestion(i int, t string, client http.Client) (string, error) {
	newQuestion := application.NewQuestion{
		Prompt: fmt.Sprintf("Test Question %d", i),
		Tags:   []string{"test"},
	}
	body, err := json.Marshal(newQuestion)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest(
		http.MethodPost,
		// TODO: Get this to be configurable via env var
		"http://localhost:8088"+constants.QuestionBasePath,
		strings.NewReader(string(body)))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+t)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 201 {
		return "", fmt.Errorf("expected status code 201, got %d", resp.StatusCode)
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var question application.Question
	err = json.Unmarshal(respBody, &question)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %s; err: %w", respBody, err)
	}
	return question.ID, nil
}
