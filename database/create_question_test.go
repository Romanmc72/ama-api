package database

import (
	"ama/api/application"
	"ama/api/logging"
	"testing"
)

// Questions need prompts
func TestInvalidCreateQuestionNoPrompt(t *testing.T) {
	db := Database{
		client: nil,
		logger: logging.GetLogger(),
		ctx:    nil,
	}
	testQ := application.NewQuestion{
		Prompt: "    ",
		Tags:   []string{"1", "2"},
	}
	_, err := db.CreateQuestion(&testQ)
	if err == nil {
		t.Fatalf("Test question should have failed validation on no prompt %s", &testQ)
	}
}

// Questions also need some kind of tag
func TestInvalidCreateQuestionNoTags(t *testing.T) {
	db := Database{
		client: nil,
		logger: logging.GetLogger(),
		ctx:    nil,
	}
	testQ := application.NewQuestion{
		Prompt: "A question!",
		Tags:   []string{},
	}
	_, err := db.CreateQuestion(&testQ)
	if err == nil {
		t.Fatalf("Test question should have failed validation on no tags %s", &testQ)
	}
}

// Question tags cannot be duplicated
func TestInvalidCreateQuestionDuplicateTags(t *testing.T) {
	db := Database{
		client: nil,
		logger: logging.GetLogger(),
		ctx:    nil,
	}
	testQ := application.NewQuestion{
		Prompt: "A question!",
		Tags:   []string{"1", "2", "1"},
	}
	_, err := db.CreateQuestion(&testQ)
	if err == nil {
		t.Fatalf("Test question should have failed validation on duplicate tags %s", &testQ)
	}
}
