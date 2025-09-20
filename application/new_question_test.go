package application

import (
	"reflect"
	"testing"
)

// Ensure it prints right
func TestNewQuestionToString(t *testing.T) {
	nq := NewQuestion{
		Prompt: "prompt",
		Tags:   []string{"1", "2", "3"},
	}
	actual := nq.String()
	want := "NewQuestion(Prompt=prompt, Tags=[1, 2, 3])"
	if actual != want {
		t.Fatalf("Wanted %s but received %s for new questions string method", want, actual)
	}
}

// Ensure that it creates a new question correctly
func TestNewQuestionToQuestion(t *testing.T) {
	nq := NewQuestion{
		Prompt: "prompt",
		Tags:   []string{"1", "2", "3"},
	}
	q := nq.Question("id")
	if q.ID != "id" || !reflect.DeepEqual(nq.Tags, q.Tags) || nq.Prompt != q.Prompt {
		t.Fatalf("new question %s did not convert exactly to question %s", nq, q)
	}
}
