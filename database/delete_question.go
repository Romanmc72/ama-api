package database

import "time"

// Deletes a question from the database and returns the time that the delete occurred
func (db *Database) DeleteQuestion(questionId string) (time.Time, error) {
	collection := db.client.Collection("questions")
	deleteResult, err := collection.Doc(questionId).Delete(db.ctx)
	if err != nil {
		db.logger.Error("Encountered an error deleting the question", "id", questionId, "error", err)
		return time.Now(), err
	}
	return deleteResult.UpdateTime, nil
}
