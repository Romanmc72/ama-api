package database

import (
	"context"
	"log/slog"

	"cloud.google.com/go/firestore"
)

// Abstracts the interactions of the database implementation details away from
// the user to make a simple, easy to use database client which takes care of
// all of the details behind interacting with the database.
type Database struct {
	// the firestore client that interacts directly with google cloud firestore
	client *firestore.Client
	// The context object from the application
	ctx context.Context
	// The instantiated logging client
	logger *slog.Logger
}

// The close function to be deferred
func (db *Database) Close() (error) {
	return db.client.Close()
}

// Takes the name of the object that needs incrementing and generates the
// next identifier number for it using that document in the identifier collection.
// func (db *Database) incrementIdentifier(objectName string) (string, error) {
// 	var newId int64
// 	documentReference := db.client.
// 		Collection(firestoreobjects.IdentifierCollection).
// 		Doc(objectName)
// 	err := db.client.RunTransaction(db.ctx, func(ctx context.Context, tx *firestore.Transaction) error {
// 		document, err := tx.Get(documentReference)
// 		if err != nil {
// 			if status.Code(err) == codes.NotFound {
// 				newId = initialIdentifier
// 				return tx.Set(
// 					documentReference,
// 					map[string]interface{}{
// 						firestoreobjects.Identifier: newId,
// 					},
// 					firestore.MergeAll,
// 				)
// 			}
// 			return err
// 		}
// 		currentId, err := document.DataAt(firestoreobjects.Identifier)
// 		if err != nil {
// 			return err
// 		}
// 		newId = currentId.(int64) + 1
// 		return tx.Set(
// 			documentReference,
// 			map[string]interface{}{
// 				firestoreobjects.Identifier: newId,
// 			},
// 			firestore.MergeAll,
// 		)
// 	})
// 	if err != nil {
// 		return "", err
// 	}
// 	return strconv.FormatInt(newId, 10), err
// }

// // if an identifier has been used and is no longer associated to a question
// // then toss it into the orphan list.
// func (db *Database) orphanId(id string, objectName string) error {
// 	collection := db.client.Collection(firestoreobjects.IdentifierCollection).
// 		Doc(objectName).
// 		Collection(firestoreobjects.OrphanIdCollection)
// 	documentRef, _, err := collection.Add(db.ctx, map[string]string {
// 		firestoreobjects.Identifier: id,
// 	})
// 	if err != nil {
// 		db.logger.Error(
// 			"Encountered an issue writing the id to the orphan collection",
// 			"error", err,
// 			"orphan", id,
// 		)
// 		return err
// 	}
// 	db.logger.Debug(
// 		"Successfully orphaned an id",
// 		"orphan", id,
// 		"documentRefId", documentRef.ID,
// 	)
// 	return nil
// }

// Apply server side validation to the question data prior to writing it to the database.
// func (db *Database) validateQuestion(preQ interfaces.QuestionConverter) (error) {
// 	q := preQ.Question("")
// 	var errorBuilder strings.Builder
// 	if strings.Join(q.Tags, "") == "" {
// 		errorBuilder.WriteString(`question "tags" field is required to have at least 1 tag`)
// 	}
// 	slices.Sort(q.Tags)
// 	dedupedTags := slices.Compact(q.Tags)
// 	if len(q.Tags) != len(dedupedTags) {
// 		errorBuilder.WriteString(`question "tags" cannot have duplicate tag values`)
// 	}
// 	if strings.TrimSpace(q.Prompt) == "" {
// 		errorBuilder.WriteString(`question "prompt" field cannot be blank`)
// 	}
// 	errorMessage := errorBuilder.String()
// 	if errorMessage != "" {
// 		return errors.New(errorMessage)
// 	}
// 	return nil
// }
