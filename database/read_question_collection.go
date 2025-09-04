package database

import (
	"ama/api/application"
	"ama/api/constants"
	"ama/api/interfaces"
	"strings"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

// Generically read questions from a collection of questions using the given reference to the collection
func (db *Database) readQuestionCollection(collectionRef interfaces.CollectionRef, limit int, finalId string, tags []string) (questions []application.Question, err error) {
	query := collectionRef.OrderBy(firestore.DocumentID, firestore.Asc)
	db.logger.Debug(
		"Read questions query params",
		"finalId", finalId,
		"limit", limit,
		"tags", tags,
	)
	if finalId != "" {
		query = query.StartAfter(finalId)
	}
	tags = application.SortDedupeAndIgnore(tags, constants.SearchTagDelimiter)
	if len(tags) == 1 {
		tag := tags[0]
		db.logger.Debug(
			"Got one tag to filter with",
			"tag", tag,
		)
		query = query.Where("tags", "array-contains", tag)
	} else if len(tags) > 1 {
		db.logger.Debug(
			"Got several tags to filter with",
			"tags", tags,
		)
		searchTags := strings.Join(tags, constants.SearchTagDelimiter)
		query = query.Where("searchTags", "array-contains", searchTags)
	}
	if limit > 0 {
		query = query.Limit(limit)
	}
	iter := query.Documents(db.ctx)
	defer iter.Stop()
	for {
		doc, err := iter.Next()
		var retrievedQuestion application.Question
		if err != nil {
			if err == iterator.Done {
				break
			}
			db.logger.Error("Encountered an error iterating questions", "error", err)
		}
		doc.DataTo(&retrievedQuestion)
		retrievedQuestion.ID = doc.ID()
		questions = append(questions, retrievedQuestion)
	}
	return questions, nil
}
