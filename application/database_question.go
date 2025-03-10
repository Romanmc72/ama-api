package application

// The database representation of a question
type DatabaseQuestion struct {
	// The text for the prompt
	Prompt string `firestore:"prompt"`
	// The tags stored in a way that is searchable by a firestore query array
	// operation for tags that look for an array subset match. Until there is a
	// firestore `array-contains-all` query operation available this will be the
	// solution. If we need to have more than a small number of tags on a
	// question (let's say less than 5 at all times) borrowing from.
	// https://code.build/p/firestore-many-to-many-array-contains-all-dyFZgf
	SearchTags []string `firestore:"searchTags"`
	Tags []string `firestore:"tags"`
}

func (q *DatabaseQuestion) Question(id string) Question {
	return Question{
		ID: id,
		Prompt: q.Prompt,
		Tags: q.Tags,
	}
}
