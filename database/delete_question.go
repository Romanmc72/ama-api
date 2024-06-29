package database

import (
	"context"
	"time"
)

// Deletes a question from the database and returns the time that the delete occurred
func (db *Database) DeleteQuestion(questionId string) (time.Time, error) {
	ctx := context.Background()
	collection := db.client.Collection("questions")
	deleteResult, err := collection.Doc(questionId).Delete(ctx)
	if err != nil {
		db.logger.Error("Encountered an error deleting the question", "id", questionId, "error", err)
		return time.Now(), err
	}
	return deleteResult.UpdateTime, nil
}
