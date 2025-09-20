package list_test

import (
	appErrors "ama/api/application/errors"
	appList "ama/api/application/list"
	"ama/api/application/responses"
	"ama/api/constants"
	"ama/api/endpoints/list"
	"ama/api/test"
	"ama/api/test/fixtures"
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"testing"
)

func TestPostUserList(t *testing.T) {
	validListBytes, err := json.Marshal(fixtures.ValidPostListRequest)
	if err != nil {
		t.Fatalf("Example data of %v was unable to be marshaled, err = %v", fixtures.ValidPostListRequest, err)
		return
	}
	testCases := []struct {
		name     string
		db       *test.MockListManager
		ctx      *test.MockAPIContext
		wantCode int
		wantErr  bool
	}{
		{
			name: "Success",
			db: test.NewMockListManager(test.MockListManagerConfig{
				CreateList: func(userId string, l appList.List) (appList.List, error) {
					l.ID = fixtures.ListId
					return l, nil
				},
			}),
			ctx: test.NewMockAPIContext(test.MockAPIContextConfig{
				InputJSON: validListBytes,
				Params: map[string]string{
					constants.UserIdPathIdentifier: fixtures.UserId,
				},
			}),
			wantCode: http.StatusCreated,
			wantErr:  false,
		},
		{
			name: "Failure - Blank Path Identifiers",
			db:   test.NewMockListManager(test.MockListManagerConfig{}),
			ctx: test.NewMockAPIContext(test.MockAPIContextConfig{
				InputJSON: validListBytes,
				Params:    map[string]string{},
			}),
			wantCode: http.StatusBadRequest,
			wantErr:  true,
		},
		{
			name: "Failure - Bad Request JSON Bind Error",
			db:   test.NewMockListManager(test.MockListManagerConfig{}),
			ctx: test.NewMockAPIContext(test.MockAPIContextConfig{
				InputJSON: []byte(`{"name": {"baba": "booey"}}`),
				Params: map[string]string{
					constants.UserIdPathIdentifier: fixtures.UserId,
				},
			}),
			wantCode: http.StatusBadRequest,
			wantErr:  true,
		},
		{
			name: "Failure - Bad Request Invalid Data",
			db: test.NewMockListManager(test.MockListManagerConfig{
				CreateList: func(userId string, list appList.List) (appList.List, error) {
					return list, appErrors.NewValidationError([]string{"name cannot be blank"})
				},
			}),
			ctx: test.NewMockAPIContext(test.MockAPIContextConfig{
				InputJSON: []byte(`{"name": "   "}`),
				Params: map[string]string{
					constants.UserIdPathIdentifier: fixtures.UserId,
				},
			}),
			wantCode: http.StatusBadRequest,
			wantErr:  true,
		},
		{
			name: "Failure - Internal Server Error",
			db: test.NewMockListManager(test.MockListManagerConfig{
				CreateList: func(userId string, list appList.List) (appList.List, error) {
					return list, errors.New("create error")
				},
			}),
			ctx: test.NewMockAPIContext(test.MockAPIContextConfig{
				InputJSON: validListBytes,
				Params: map[string]string{
					constants.UserIdPathIdentifier: fixtures.UserId,
				},
			}),
			wantCode: http.StatusInternalServerError,
			wantErr:  true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			list.PostUserList(tc.ctx, tc.db)
			if tc.wantCode != tc.ctx.ResponseCode {
				t.Errorf("PostUserList() wanted code = %d; got = %d", tc.wantCode, tc.ctx.ResponseCode)
			}
			if _, ok := tc.ctx.ResponseData.(responses.ErrorResponse); !ok && tc.wantErr && tc.ctx.ResponseData != nil {
				t.Errorf("PostUserList() wanted error Response, got = %v", tc.ctx.ResponseData)
			}
			if l, ok := tc.ctx.ResponseData.(appList.List); !tc.wantErr && (!ok || !reflect.DeepEqual(l, fixtures.ValidLists[0])) {
				t.Errorf("PostUserList() wanted list = %s Response, got = %v", fixtures.ValidLists[0], tc.ctx.ResponseData)
			}
		})
	}
}
