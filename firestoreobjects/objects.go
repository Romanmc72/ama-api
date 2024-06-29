package firestoreobjects

const (
	// collection Names
	QuestionCollection = "questions"
	IdentifierCollection = "ids"
	// This is a nested sub-collection within the ids collection that will exist
	// on each identifier doc.
	OrphanIdCollection = "orphan-id"

	// special document names
	QuestionIdDoc = "question-id"

	// paths within fields
	Identifier = "id"
)

