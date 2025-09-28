package database

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"ama/api/application"
	"ama/api/interfaces"
)

// Retrieve a particular question from the database
func (db *Database) readQuestionFromCollection(collection interfaces.CollectionRef, id string) (q application.Question, err error) {
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
			return q, err
		}
		db.logger.Error("Encountered an error fetching the question", "id", id, "error", err)
		return q, err
	}
	err = document.DataTo(&questionFetched)
	if err != nil {
		return q, err
	}
	return questionFetched.Question(id), nil
}
