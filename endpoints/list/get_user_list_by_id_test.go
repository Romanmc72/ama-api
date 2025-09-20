package list_test

import (
	"ama/api/application"
	appList "ama/api/application/list"
	"ama/api/application/responses"
	"ama/api/constants"
	"ama/api/endpoints/list"
	"ama/api/test"
	"ama/api/test/fixtures"
	"errors"
	"net/http"
	"reflect"
	"testing"
)

func TestGetUserListById(t *testing.T) {
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
				ReadList: func(userId, listId string, limit int, finalId string, tags []string) (appList.List, []application.Question, error) {
					return fixtures.ValidLists[0], []application.Question{fixtures.ValidQuestion}, nil
				},
			}),
			ctx: test.NewMockAPIContext(test.MockAPIContextConfig{
				Params: map[string]string{
					constants.UserIdPathIdentifier: fixtures.UserId,
					constants.ListIdPathIdentifier: fixtures.ListId,
				},
				QueryValues: map[string][]string{
					constants.FinalIdParam: {"abc123"},
					constants.LimitParam:   {"10"},
					constants.TagParam:     {"test"},
				},
			}),
			wantCode: http.StatusOK,
			wantErr:  false,
		},
		{
			name: "Failure - Blank Path Identifiers",
			db:   test.NewMockListManager(test.MockListManagerConfig{}),
			ctx: test.NewMockAPIContext(test.MockAPIContextConfig{
				Params: map[string]string{
					constants.UserIdPathIdentifier: "     ",
					constants.ListIdPathIdentifier: "     ",
				},
			}),
			wantCode: http.StatusBadRequest,
			wantErr:  true,
		},
		{
			name: "Failure - Read Error",
			db: test.NewMockListManager(test.MockListManagerConfig{
				ReadList: func(userId, listId string, limit int, finalId string, tags []string) (appList.List, []application.Question, error) {
					return appList.List{}, []application.Question{}, errors.New("read error")
				},
			}),
			ctx: test.NewMockAPIContext(test.MockAPIContextConfig{
				Params: map[string]string{
					constants.UserIdPathIdentifier: fixtures.UserId,
					constants.ListIdPathIdentifier: fixtures.ListId,
				},
			}),
			wantCode: http.StatusNotFound,
			wantErr:  true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			list.GetUserListById(tc.ctx, tc.db)
			if tc.wantCode != tc.ctx.ResponseCode {
				t.Errorf("GetUserLists() wanted code = %d; got = %d", tc.wantCode, tc.ctx.ResponseCode)
			}
			if _, ok := tc.ctx.ResponseData.(responses.ErrorResponse); !ok && tc.wantErr && tc.ctx.ResponseData != nil {
				t.Errorf("GetUserLists() wanted error Response, got = %v", tc.ctx.ResponseData)
			}
			if l, ok := tc.ctx.ResponseData.(responses.GetUserListByIdResponse); !tc.wantErr && (!ok || !reflect.DeepEqual(l, fixtures.ValidGetUserListByIdResponse)) {
				t.Errorf("GetUserLists() wanted list = %s Response, got = %v", fixtures.ValidLists[0], tc.ctx.ResponseData)
			}
		})
	}
}
