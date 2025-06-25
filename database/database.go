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
func (db *Database) Close() error {
	return db.client.Close()
}
