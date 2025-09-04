package database_test

import (
	"ama/api/database"
	"testing"
)

// TODO: finish these tests and get them to pass
func TestReadQuestions(t *testing.T) {
	// logger := logging.GetLogger()
	// makeDb := func() database.Database {
	// 	return database.ManualTestConnect(
	// 		t.Context(),
	// 		test.NewMockDatabase(&test.MockDBConfig{
	// 			Collections: map[string]test.MockCollectionConfig{
	// 				constants.QuestionCollection: {
	// 					Documents: map[string]test.MockDocumentConfig{
	// 						"abc123": {
	// 							ID: "abc123",
	// 							Data: application.DatabaseQuestion{
	// 								Prompt:     "Test 1",
	// 								Tags:       []string{"test"},
	// 								SearchTags: []string{"test"},
	// 							},
	// 						},
	// 						"abc1234": {
	// 							ID: "abc1234",
	// 							Data: application.DatabaseQuestion{
	// 								Prompt:     "Test 2",
	// 								Tags:       []string{"test"},
	// 								SearchTags: []string{"test"},
	// 							},
	// 						},
	// 						"abc12345": {
	// 							ID: "abc12345",
	// 							Data: application.DatabaseQuestion{
	// 								Prompt:     "Test 3",
	// 								Tags:       []string{"test"},
	// 								SearchTags: []string{"test"},
	// 							},
	// 						},
	// 					},
	// 				},
	// 			},
	// 		}),
	// 		logger,
	// 	)
	// }
	testCases := []struct {
		name    string
		db      database.Database
		limit   int
		finalId string
		tags    []string
		wantErr bool
	}{
		// {
		// 	name:    "Gets Some Questions One Tag Filter",
		// 	limit:   10,
		// 	finalId: "abc123",
		// 	tags:    []string{"test"},
		// 	db:      makeDb(),
		// 	wantErr: false,
		// },
		// {
		// 	name: "Gets Some Questions Multiple Tag Filters",
		// 	wantErr: false,
		// },
		// {
		// 	name: "Gets Some Questions No Tag Filters",
		// 	wantErr: false,
		// },
		// {
		// 	name: "Gets No Questions",
		// 	wantErr: false,
		// },
		// {
		// 	name: "Iterator Error",
		// 	wantErr: true,
		// },
		// {
		// 	name: "Iterator Done",
		// 	wantErr: false,
		// },
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			qs, err := tc.db.ReadQuestions(tc.limit, tc.finalId, tc.tags)
			if (err != nil) != tc.wantErr {
				t.Errorf("ReadQuestions() wantedErr = %v; got = %v", tc.wantErr, err)
			}
			if !tc.wantErr && len(qs) == 0 {
				t.Error("ReadQuestions() wanted > 0 questions, got 0")
			}
		})
	}
}
