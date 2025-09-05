package fixtures

import (
	"ama/api/application"
)

const prompt = "How much wood would a woodchuck chuck?"
var Tags = []string{"silly", "test"}

var ValidNewQuestion = application.NewQuestion{
	Prompt: prompt,
	Tags:   Tags,
}

var ValidDatabaseQuestion = application.DatabaseQuestion{
	Prompt:     prompt,
	Tags:       Tags,
	SearchTags: []string{"silly", "test", "silly|test"},
}

var ValidQuestion = application.Question{
	ID:     QuestionId,
	Prompt: prompt,
	Tags:   Tags,
}
