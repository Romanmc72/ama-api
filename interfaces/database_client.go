package interfaces

import (
	"context"

	"cloud.google.com/go/firestore"
)

type DatabaseClient interface {
	NewID() string
	Collection(name string) CollectionRef
	RunTransaction(
		ctx context.Context,
		f func(context.Context, *firestore.Transaction) error,
		t ...firestore.TransactionOption,
	) error
	Close() error
	BulkWriter(context.Context) BulkWriter
}

type Transaction interface {
	Get(doc *firestore.DocumentRef) (DocumentSnapshot, error)
	Set(doc *firestore.DocumentRef, data interface{}, opts ...firestore.SetOption) error
}

type BulkWriter interface {
	Create(doc *firestore.DocumentRef, data interface{}) (*firestore.BulkWriterJob, error)
	Set(doc *firestore.DocumentRef, data interface{}, opts ...firestore.SetOption) (*firestore.BulkWriterJob, error)
	Update(doc *firestore.DocumentRef, updates []firestore.Update) (*firestore.BulkWriterJob, error)
	Delete(doc *firestore.DocumentRef) (*firestore.BulkWriterJob, error)
	End()
	Flush()
}

type CollectionRef interface {
	Path() string
	Doc(id string) DocumentRef
	Add(ctx context.Context, data any) (DocumentRef, *firestore.WriteResult, error)
	Query
}

type DocumentRef interface {
	Ref() *firestore.DocumentRef
	ID() string
	Get(ctx context.Context) (DocumentSnapshot, error)
	Set(ctx context.Context, data interface{}, opts ...firestore.SetOption) (*firestore.WriteResult, error)
	Collection(name string) CollectionRef
	Delete(ctx context.Context) (*firestore.WriteResult, error)
}

type DocumentSnapshot interface {
	DataTo(v any) error
	ID() string
	Ref() *firestore.DocumentRef
	Exists() bool
}

type DocumentIterator interface {
	Next() (DocumentSnapshot, error)
	Stop()
}

type Query interface {
	Documents(ctx context.Context) DocumentIterator
	StartAfter(docID string) Query
	Where(path, op string, value any) Query
	Limit(n int) Query
	OrderBy(path string, direction firestore.Direction) Query
}
