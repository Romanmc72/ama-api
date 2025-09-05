package database_test

import (
	"ama/api/application"
	"ama/api/constants"
	"ama/api/database"
	"ama/api/logging"
	"ama/api/test"
	"errors"
	"fmt"
	"testing"
)

// TODO: finish these tests and get them to pass
func TestReadQuestions(t *testing.T) {
	logger := logging.GetLogger()
	makeDb := func(n int, err error) database.Database {
		qs := make([]test.MockDocumentConfig, n)
		for i := range n {
			qs[i] = test.MockDocumentConfig{
				ID: fmt.Sprintf("test-question-id-%d", i),
				Data: application.DatabaseQuestion{
					Prompt:     fmt.Sprintf("Test %d", i),
					Tags:       []string{"test"},
					SearchTags: []string{"test"},
				},
			}
		}
		return database.ManualTestConnect(
			t.Context(),
			test.NewMockDatabase(&test.MockDBConfig{
				Collections: map[string]test.MockCollectionConfig{
					constants.QuestionCollection: {
						QueryDocuments: qs,
						MockError:      err,
					},
				},
			}),
			logger,
		)
	}
	testCases := []struct {
		name    string
		limit   int
		finalId string
		tags    []string
		wantErr bool
	}{
		{
			name:    "Gets Some Questions One Tag Filter",
			limit:   10,
			finalId: "abc123",
			tags:    []string{"test"},
			wantErr: false,
		},
		{
			name:    "Gets Some Questions Multiple Tag Filters",
			limit:   10,
			finalId: "abc123",
			tags:    []string{"test", "silly"},
			wantErr: false,
		},
		{
			name:    "Gets Some Questions No Tag Filters",
			limit:   10,
			finalId: "abc123",
			tags:    []string{},
			wantErr: false,
		},
		{
			// This is just for adding coverage, results are still
			// returned if the query parameters are omitted
			name:    "Gets No Questions",
			limit:   0,
			finalId: "",
			tags:    []string{},
			wantErr: false,
		},
		{
			name:    "Iterator Error",
			limit:   10,
			finalId: "abc123",
			tags:    []string{"test"},
			wantErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var mockErr error
			if tc.wantErr {
				mockErr = errors.New("here is your error")
			}
			db := makeDb(tc.limit, mockErr)
			qs, err := db.ReadQuestions(tc.limit, tc.finalId, tc.tags)
			if (err != nil) != tc.wantErr {
				t.Errorf("ReadQuestions() wantedErr = %v; got = %v", tc.wantErr, err)
			}
			if !tc.wantErr && len(qs) != tc.limit {
				t.Errorf("ReadQuestions() wanted == %d questions, got %d", tc.limit, len(qs))
			}
		})
	}
}
