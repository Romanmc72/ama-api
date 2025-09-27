//go:build integration
// +build integration

package integration_test

import (
	"ama/api/application/list"
	"ama/api/application/requests"
	"ama/api/application/responses"
	"fmt"
	"net/http"
	"testing"
)

func GetListBaseUrl(secure bool, userId string) string {
	s := ""
	if secure {
		s = "s"
	}
	return fmt.Sprintf(
		"http%s://%s:%d/user/%s/list",
		s,
		ResourceServerHost,
		ResourceServerPort,
		userId,
	)
}

func GetListUrlPath(secure bool, userId string, listId string) string {
	return fmt.Sprintf(
		"%s/%s",
		GetListBaseUrl(secure, userId),
		listId,
	)
}

func GetListQuestionUrlPath(secure bool, userId string, listId string, questionId string) string {
	return fmt.Sprintf(
		"%s/question/%s",
		GetListUrlPath(secure, userId, listId),
		questionId,
	)
}

func ListSuite(t *testing.T) {
	client := &http.Client{}
	token, err := SignUserInAndGetToken(client, UserEmail, UserPass)
	if err != nil {
		t.Fatalf("Error getting user auth token: %v", err)
	}
	userId, err := GetIdFromToken(token)
	if err != nil {
		t.Fatalf("Error getting user id from token: %s err: %v", token, err)
	}
	listClient := NewListTestClient(client, token)
	listName := "test list"
	newList, err := listClient.CreateList(userId, listName)
	if err != nil {
		t.Fatalf("Error creating list %s, err: %s", listName, err)
	}
	listId := newList.ID
	newListName := "new + " + listName
	newList.Name = newListName
	err = listClient.UpdateList(userId, newList)
	if err != nil {
		t.Fatalf("Error updating list %s, err: %s", newList.ID, err)
	}
	limit := 10
	qs, err := ReadQuestions(client, token, TestTags, limit, "", false)
	if err != nil {
		t.Fatalf("failed to read questions %s", err)
	}
	for i, q := range qs {
		err = listClient.AddToList(userId, listId, q.ID)
		if err != nil {
			t.Fatalf("iteration %d of %d; failed to add question %s to list %s, err: %s", i, len(qs), q.ID, listId, err)
		}
	}
	r, err := listClient.ReadList(userId, listId)
	if err != nil {
		t.Fatalf("failed to read list %s, err: %s", listId, err)
	}
	if r.List.Name != newListName {
		t.Fatalf("List name failed to update want: %s, got: %s", newListName, r.List.Name)
	}
	if r.List.ID != listId {
		t.Fatalf("wanted list: %s, got list: %s", listId, r.List.ID)
	}
	if len(qs) != len(r.Questions) {
		t.Fatalf("wanted to see %d questions in the list, got %d instead", len(qs), len(r.Questions))
	}
}

type ListTestClient struct {
	client  *http.Client
	idToken string
}

func NewListTestClient(client *http.Client, token string) ListTestClient {
	return ListTestClient{client: client, idToken: token}
}

func (ltc *ListTestClient) CreateList(userId string, listName string) (list.List, error) {
	data := requests.PostUserListRequest{Name: listName}
	var l list.List
	l, err := HitApi(
		ltc.client,
		GetListBaseUrl(IsSecure, userId),
		http.MethodPost,
		ltc.idToken,
		data,
		l,
	)
	if err != nil {
		return list.List{}, err
	}
	return l, nil
}

func (ltc *ListTestClient) ReadList(userId string, listId string) (r responses.GetUserListByIdResponse, err error) {
	return HitApi(
		ltc.client,
		GetListUrlPath(IsSecure, userId, listId),
		http.MethodGet,
		ltc.idToken,
		nil,
		r,
	)
}

func (ltc *ListTestClient) UpdateList(userId string, l list.List) error {
	_, err := HitApi(
		ltc.client,
		GetListUrlPath(IsSecure, userId, l.ID),
		http.MethodPut,
		ltc.idToken,
		l,
		responses.SuccessResponse{},
	)
	return err
}

func (ltc *ListTestClient) DeleteList(userId string, listId string) error {
	_, err := HitApi(
		ltc.client,
		GetListUrlPath(IsSecure, userId, listId),
		http.MethodDelete,
		ltc.idToken,
		nil,
		responses.SuccessResponse{},
	)
	return err
}

func (ltc *ListTestClient) AddToList(userId string, listId string, questionId string) error {
	_, err := HitApi(
		ltc.client,
		GetListQuestionUrlPath(IsSecure, userId, listId, questionId),
		http.MethodPost,
		ltc.idToken,
		nil,
		responses.SuccessResponse{},
	)
	return err
}

func (ltc *ListTestClient) RemoveFromList(userId string, listId string, questionId string) error {
	_, err := HitApi(
		ltc.client,
		GetListQuestionUrlPath(IsSecure, userId, listId, questionId),
		http.MethodDelete,
		ltc.idToken,
		nil,
		responses.SuccessResponse{},
	)
	return err
}
