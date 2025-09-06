package list_test

import (
	appList "ama/api/application/list"
	"ama/api/application/responses"
	"ama/api/constants"
	"ama/api/endpoints/list"
	"ama/api/test"
	"ama/api/test/fixtures"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
)

func TestPutUserListById(t *testing.T) {
	validListBytes, _ := json.Marshal(fixtures.ValidListRequest)
	testCases := []struct {
		name     string
		db       *test.MockListManager
		ctx      *test.MockAPIContext
		wantCode int
		wantErr  bool
	}{
		{
			name: "Success",
			db:   test.NewMockListManager(test.MockListManagerConfig{}),
			ctx: test.NewMockAPIContext(test.MockAPIContextConfig{
				InputJSON: validListBytes,
				Params: map[string]string{
					constants.UserIdPathIdentifier: fixtures.UserId,
					constants.ListIdPathIdentifier: fixtures.ListId,
				},
			}),
			wantCode: http.StatusOK,
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
			name: "Failure - Bad Request",
			db:   test.NewMockListManager(test.MockListManagerConfig{}),
			ctx: test.NewMockAPIContext(test.MockAPIContextConfig{
				InputJSON: []byte(`{"name": {"baba": "booey"}}`),
				Params: map[string]string{
					constants.UserIdPathIdentifier: fixtures.UserId,
					constants.ListIdPathIdentifier: fixtures.ListId,
				},
			}),
			wantCode: http.StatusBadRequest,
			wantErr:  true,
		},
		{
			name: "Failure - Internal Server Error",
			db: test.NewMockListManager(test.MockListManagerConfig{
				UpdateList: func(userId string, updatedList appList.List) error {
					return errors.New("update error")
				},
			}),
			ctx: test.NewMockAPIContext(test.MockAPIContextConfig{
				InputJSON: validListBytes,
				Params: map[string]string{
					constants.UserIdPathIdentifier: fixtures.UserId,
					constants.ListIdPathIdentifier: fixtures.ListId,
				},
			}),
			wantCode: http.StatusInternalServerError,
			wantErr:  true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			list.PutUserListById(tc.ctx, tc.db)
			if tc.wantCode != tc.ctx.ResponseCode {
				t.Errorf("PutUserListById() wanted code = %d; got = %d", tc.wantCode, tc.ctx.ResponseCode)
			}
			if _, ok := tc.ctx.ResponseData.(responses.ErrorResponse); !ok && tc.wantErr && tc.ctx.ResponseData != nil {
				t.Errorf("PutUserListById() wanted error Response, got = %v", tc.ctx.ResponseData)
			}
		})
	}
}
