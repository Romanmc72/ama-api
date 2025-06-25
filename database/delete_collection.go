package database

import (
	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

// Delete an entire collection of data in the firestore database.
// This could take a while and is done in batches.
func (db *Database) deleteCollection(collectionRef *firestore.CollectionRef, batchSize int) error {
	db.logger.Debug("delete collection", "collection", collectionRef.Path, "batchSize", batchSize)
	if batchSize > 500 {
		db.logger.Warn("batch size too large, setting to 500", "batchSize", batchSize)
		batchSize = 500
	}

	bulkWriter := db.client.BulkWriter(db.ctx)

	totalDeleted := 0
	for {
		// Get a batch of documents
		iter := collectionRef.Limit(batchSize).Documents(db.ctx)
		numDeleted := 0

		// Iterate through the documents, adding
		// a delete operation for each one to the BulkWriter.
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return err
			}

			_, err = bulkWriter.Delete(doc.Ref)
			if err != nil {
				db.logger.Error("error deleting document in collection deletion", "totalDeleted", totalDeleted, "error", err)
				return err
			}
			numDeleted++
		}

		// If there are no documents to delete,
		// the process is over.
		if numDeleted == 0 {
			bulkWriter.End()
			break
		}

		bulkWriter.Flush()
		totalDeleted += numDeleted
	}
	db.logger.Debug("delete succeeded", "collection", collectionRef.Path, "totalDeleted", totalDeleted)
	return nil
}
