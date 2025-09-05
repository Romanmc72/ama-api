package user_test

import (
	"ama/api/application/responses"
	"ama/api/endpoints/user"
	"ama/api/test"
	"testing"
)

func TestPostUser(t *testing.T) {
	testCases := []struct {
		name string
		db test.MockUserManager
		ctx test.MockAPIContext
		wantCode int
		wantErr bool
	}{
		// {
		// 	name: "Success",
		// 	wantCode: http.StatusCreated,
		// 	wantErr: false,
		// },
		// {
		// 	name: "Failure - Bad Request",
		// 	wantCode: http.StatusBadRequest,
		// 	wantErr: true,
		// },
		// {
		// 	name: "Failure - Internal Server Error",
		// 	wantCode: http.StatusInternalServerError,
		// 	wantErr: true,
		// },
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			user.PostUser(&tc.ctx, &tc.db)
			if tc.wantCode != tc.ctx.ResponseCode {
				t.Errorf("PostUser() wanted Code = %d; got = %d", tc.wantCode, tc.ctx.ResponseCode)
			}
			if _, ok := tc.ctx.ResponseData.(responses.ErrorResponse); tc.wantErr && !ok {
				t.Errorf("PostUser() wanted error, got %v", tc.ctx.ResponseData)
			}
		})
	}
}
