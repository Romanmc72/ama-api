package fixtures

import (
	"ama/api/application"
)

const QuestionId = "abc123"
const prompt = "How much wood would a woodchuck chuck?"

var tags = []string{"silly", "test"}

var ValidNewQuestion = application.NewQuestion{
	Prompt: prompt,
	Tags:   tags,
}

var ValidDatabaseQuestion = application.DatabaseQuestion{
	Prompt:     prompt,
	Tags:       tags,
	SearchTags: []string{"silly", "test", "silly|test"},
}

var ValidQuestion = application.Question{
	ID:     QuestionId,
	Prompt: prompt,
	Tags:   tags,
}
