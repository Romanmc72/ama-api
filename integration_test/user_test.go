//go:build integration
// +build integration

package integration_test

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"ama/api/application"
	"ama/api/application/list"
	"ama/api/application/requests"
	"ama/api/application/responses"
	"ama/api/application/user"
)

const (
	firebaseUrlPath   = "identitytoolkit.googleapis.com/v1/accounts"
	firebaseSignUpUrl = firebaseUrlPath + ":signUp"
	firebaseSignInUrl = firebaseUrlPath + ":signInWithPassword"
)

func getFirebaseSignUpUrl(secure bool) string {
	s := ""
	if secure {
		s = "s"
	}
	return fmt.Sprintf(
		"http%s://%s:%d/%s?key=%s",
		s,
		EmulatorHost,
		EmulatorPort,
		firebaseSignUpUrl,
		ApiKey,
	)
}

func getFirebaseSignInUrl(secure bool) string {
	s := ""
	if secure {
		s = "s"
	}
	return fmt.Sprintf(
		"http%s://%s:%d/%s?key=%s",
		s,
		EmulatorHost,
		EmulatorPort,
		firebaseSignInUrl,
		ApiKey,
	)
}

func UserSetupSuite(t *testing.T) {
	t.Log("Running user integ tests")
	client := &http.Client{}
	token, err := CreateUserAndGetToken(client, UserEmail, UserPass)
	if err != nil {
		t.Fatalf("Error getting user auth token: %v", err)
	}
	userId, err := createUser(token, client, UserEmail, "test user")
	if err != nil {
		t.Fatalf("Error creating user: %v", err)
	}
	u, err := readUser(client, token, userId)
	if err != nil {
		t.Fatalf("Error reading user: %s, err: %v", userId, err)
	}
	u.Name = "test user name change"
	err = updateUser(client, token, userId, u.BaseUser)
	if err != nil {
		t.Fatalf("Error updating user: %s, data: %v, err: %v", userId, u, err)
	}
	t.Logf("Created user with ID: %s", userId)
}

func UserTearDownSuite(t *testing.T) {
	client := &http.Client{}
	token, err := SignUserInAndGetToken(client, UserEmail, UserPass)
	if err != nil {
		t.Fatalf("unable to get sign in token for %s err: %s", UserEmail, err)
	}
	userId, err := GetIdFromToken(token)
	if err != nil {
		t.Fatalf("unable to parse id from token: %s err: %s", token, err)
	}
	err = deleteUser(token, client, userId)
	if err != nil {
		t.Fatalf("failed to delete user %s, err: %s", userId, err)
	}
}

type createUserRequest struct {
	Email             string `json:"email"`
	Password          string `json:"password"`
	ReturnSecureToken bool   `json:"returnSecureToken"`
}

func CreateUserAndGetToken(httpClient *http.Client, email string, password string) (string, error) {
	return hitFirebaseEmulatorAuth(httpClient, email, password, getFirebaseSignUpUrl(IsSecure))
}

func SignUserInAndGetToken(httpClient *http.Client, email string, password string) (string, error) {
	return hitFirebaseEmulatorAuth(httpClient, email, password, getFirebaseSignInUrl(IsSecure))
}

func GetUserBaseUrl(secure bool) string {
	s := ""
	if secure {
		s = "s"
	}
	return fmt.Sprintf(
			"http%s://%s:%d/user",
			s,
			ResourceServerHost,
			ResourceServerPort,
		)
}

func GetUserUrl(secure bool, userId string) string {
	return fmt.Sprintf(
			"%s/%s",
			GetUserBaseUrl(secure),
			userId,
		)
}

func hitFirebaseEmulatorAuth(httpClient *http.Client, email string, password string, url string) (string, error) {
	reqBody := createUserRequest{
		Email:             email,
		Password:          password,
		ReturnSecureToken: true,
	}
	var respToken map[string]any
	respToken, err := HitApi(
		httpClient,
		url,
		http.MethodPost,
		"",
		reqBody,
		respToken,
	)
	if err != nil {
		return "", err
	}
	idToken, ok := respToken["idToken"].(string)
	if !ok {
		return "", fmt.Errorf("idToken not found in response: %v", respToken)
	}
	return idToken, nil
}

func createUser(idToken string, httpClient *http.Client, email string, name string) (string, error) {
	firebaseId, err := GetIdFromToken(idToken)
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
	secure := ""
	if IsSecure {
		secure = "s"
	}
	var createdUser application.User
	createdUser, err = HitApi(
		httpClient,
		fmt.Sprintf(
			"http%s://%s:%d/user",
			secure,
			ResourceServerHost,
			ResourceServerPort,
		),
		http.MethodPost,
		idToken,
		userReqBody,
		createdUser,
	)
	if err != nil {
		return "", err
	}
	return createdUser.ID, nil
}

func readUser(client *http.Client, token string, userId string) (application.User, error) {
	var u application.User
	return HitApi(
		client,
		GetUserUrl(IsSecure, userId),
		http.MethodGet,
		token,
		nil,
		u,
	)
}

func updateUser(client *http.Client, token string, userId string, u user.BaseUser) error {
	_, err := HitApi(
		client,
		GetUserUrl(IsSecure, userId),
		http.MethodPut,
		token,
		u,
		responses.SuccessResponse{},
	)
	return err
}

func deleteUser(idToken string, client *http.Client, userId string) error {
	_, err := HitApi(
		client,
		GetUserUrl(IsSecure, userId),
		http.MethodDelete,
		idToken,
		nil,
		responses.SuccessResponse{},
	)
	return err
}

func GetIdFromToken(idToken string) (string, error) {
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
