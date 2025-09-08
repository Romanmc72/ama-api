//go:build integration
// +build integration

package integration_test

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"ama/api/application"
	"ama/api/application/list"
	"ama/api/application/requests"
	"ama/api/application/user"
)

func TestUser(t *testing.T) {
	t.Log("Running user integ tests")
	client := &http.Client{}
	token, err := getUserAuthToken(client, UserEmail, UserPass)
	if err != nil {
		t.Fatalf("Error getting user auth token: %v", err)
	}
	userId, err := createUser(token, client, UserEmail, "test user")
	if err != nil {
		t.Fatalf("Error creating user: %v", err)
	}
	t.Logf("Created user with ID: %s", userId)
}

type createUserRequest struct {
	Email             string `json:"email"`
	Password          string `json:"password"`
	ReturnSecureToken bool   `json:"returnSecureToken"`
}

func getUserAuthToken(httpClient *http.Client, email string, password string) (string, error) {
	reqBody := createUserRequest{
		Email:             email,
		Password:          password,
		ReturnSecureToken: true,
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}
	secure := ""
	if IsSecure {
		secure = "s"
	}
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf(
			"http%s://%s:%d/identitytoolkit.googleapis.com/v1/accounts:signUp?key=%s",
			secure,
			EmulatorHost,
			EmulatorPort,
			ApiKey,
		),
		strings.NewReader(string(body)))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var respToken map[string]any
	err = json.Unmarshal(respBody, &respToken)
	idToken, ok := respToken["idToken"].(string)
	if !ok {
		return "", fmt.Errorf("idToken not found in response: %s", respBody)
	}
	return idToken, nil
}

func createUser(idToken string, httpClient *http.Client, email string, name string) (string, error) {
	firebaseId, err := getIdFromToken(idToken)
	if err != nil {
		return "", err
	}
	userReqBody := requests.PostUserRequest{
		Name:       name,
		Email:      email,
		Tier:       "free",
		FirebaseID: firebaseId,
		Subscription: user.UserSubscription{
			PayCadence:  "monthly",
			RenewalDate: time.Now().AddDate(0, 0, 30),
		},
		Settings: user.UserSettings{
			ColorScheme: user.UserColorScheme{
				Background:            "default",
				Foreground:            "default",
				HighlightedBackground: "default",
				HighlightedForeground: "default",
			},
		},
		Lists: []list.List{},
	}
	body, err := json.Marshal(userReqBody)
	if err != nil {
		return "", err
	}
	secure := ""
	if IsSecure {
		secure = "s"
	}
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf(
			"http%s://%s:%d/user",
			secure,
			ResourceServerHost,
			ResourceServerPort,
		),
		strings.NewReader(string(body)))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", idToken))
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var createdUser application.User
	err = json.Unmarshal(respBody, &createdUser)
	if err != nil {
		return "", err
	}
	return createdUser.ID, nil
}

func getIdFromToken(idToken string) (string, error) {
	splitToken := strings.Split(idToken, ".")
	if len(splitToken) != 3 {
		return "", fmt.Errorf("invalid token format")
	}
	payload := splitToken[1]
	decodedBytes, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		return "", err
	}
	var payloadData map[string]any
	err = json.Unmarshal(decodedBytes, &payloadData)
	if err != nil {
		return "", err
	}
	sub, ok := payloadData["sub"].(string)
	if !ok {
		return "", fmt.Errorf("sub field not found in token payload")
	}
	return sub, nil
}
