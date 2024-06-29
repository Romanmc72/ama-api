package database

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"ama/api/application"
	"ama/api/firestoreobjects"
)

// Retrieve a particular question from the database
func (db *Database) ReadQuestion(id string) (application.Question, error) {
	collection := db.client.Collection(firestoreobjects.QuestionCollection)
	var questionFetched application.NewQuestion
	document, err := collection.Doc(id).Get(db.ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			db.logger.Warn(
				"Unable to find document with id. Received not found error.",
				"id",
				id,
				"error",
				err,
			)
			return application.Question{}, err
		}
		db.logger.Error("Encountered an error fetching the question", "id", id, "error", err)
		return application.Question{}, err
	}
	document.DataTo(&questionFetched)
	return questionFetched.Question(id), nil
}
