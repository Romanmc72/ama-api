package database

import (
	"context"
	"strings"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"

	"ama/api/application"
	"ama/api/firestoreobjects"
)

// Retrieve all of the questions from the database one page at a time.
// limit = the number of items per page. 0 or less will not paginate but return
// an iterator for all documents.
// finalId = the Id field of the last document on the previous page. It will
// start the next query off after that one.
// tag = the tag to search for within the set of questions.
// tags = the set of tags to search for an inclusive match of.
func (db *Database) ReadQuestions(limit int, finalId string, tags []string) ([]application.Question, error) {
	var questions []application.Question
	collection := db.client.Collection(firestoreobjects.QuestionCollection)
	ctx := context.Background()
	query := collection.OrderBy(firestore.DocumentID, firestore.Asc)
	db.logger.Debug(
		"Read questions query params",
		"finalId", finalId,
		"limit", limit,
		"tags", tags,
	)
	if finalId != "" {
		query = query.StartAfter(finalId)
	}
	tags = application.SortDedupeAndIgnore(tags, application.SearchTagDelimiter)
	if len(tags) == 1 {
		tag := tags[0]
		db.logger.Debug(
			"Got one tag to filter with",
			"tag", tag,
		)
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
		searchTags := strings.Join(tags, application.SearchTagDelimiter)
		query = query.Where("searchTags", "array-contains", searchTags)
	}
	if limit > 0 {
		query = query.Limit(limit)
	}
	iter := query.Documents(ctx)
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
		retrievedQuestion.ID = doc.Ref.ID
		questions = append(questions, retrievedQuestion)
	}
	if len(questions) > 0 {
		return questions, nil
	} else {
		return []application.Question{}, nil
	}
}
