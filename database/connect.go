package database

import (
	"context"
	"os"

	"cloud.google.com/go/firestore"

	"ama/api/logging"
)

// Connect() provides a convenience method for instantiating the
// firestore client. It can be created with a single call encapsulating the
// required calls except the deferred close.
func Connect() (Database, error) {
	logger := logging.GetLogger()
	projectId := os.Getenv("PROJECT_ID")
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectId)
	if err != nil {
		logger.Error("Encountered an error connecting to the database", "error", err)
		return Database{}, err
	}
	return Database{
		client: client,
		ctx: ctx,
		logger: logger,
	}, nil
}
