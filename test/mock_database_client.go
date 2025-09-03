/** This file was brought to you by Gemini */
package test

import (
	"ama/api/interfaces"
	"context"
	"errors"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

// MockDatabaseClient is a mock for DatabaseClient.
type MockDatabaseClient struct {
	MockCollection     func(name string) interfaces.CollectionRef
	MockRunTransaction func(
		ctx context.Context,
		f func(context.Context, *firestore.Transaction) error,
		t ...firestore.TransactionOption,
	) error
	MockClose      func() error
	MockBulkWriter func(ctx context.Context) interfaces.BulkWriter
}

func (m *MockDatabaseClient) Collection(name string) interfaces.CollectionRef {
	if m.MockCollection != nil {
		return m.MockCollection(name)
	}
	return &MockCollectionRef{}
}

func (m *MockDatabaseClient) RunTransaction(
	ctx context.Context,
	f func(context.Context, *firestore.Transaction) error,
	t ...firestore.TransactionOption,
) error {
	if m.MockRunTransaction != nil {
		return m.MockRunTransaction(ctx, f, t...)
	}
	return nil
}

func (m *MockDatabaseClient) Close() error {
	if m.MockClose != nil {
		return m.MockClose()
	}
	return nil
}

func (m *MockDatabaseClient) BulkWriter(ctx context.Context) interfaces.BulkWriter {
	if m.MockBulkWriter != nil {
		return m.MockBulkWriter(ctx)
	}
	return &MockBulkWriter{}
}

// MockTransaction is a mock for Transaction.
type MockTransaction struct {
	MockGet func(doc *firestore.DocumentRef) (interfaces.DocumentSnapshot, error)
	MockSet func(doc *firestore.DocumentRef, data interface{}, opts ...firestore.SetOption) error
}

func (m *MockTransaction) Get(doc *firestore.DocumentRef) (interfaces.DocumentSnapshot, error) {
	if m.MockGet != nil {
		return m.MockGet(doc)
	}
	return &MockDocumentSnapshot{}, nil
}

func (m *MockTransaction) Set(doc *firestore.DocumentRef, data interface{}, opts ...firestore.SetOption) error {
	if m.MockSet != nil {
		return m.MockSet(doc, data, opts...)
	}
	return nil
}

// MockBulkWriter is a mock for BulkWriter.
type MockBulkWriter struct {
	MockCreate func(doc *firestore.DocumentRef, data interface{}) (*firestore.BulkWriterJob, error)
	MockSet    func(doc *firestore.DocumentRef, data interface{}, opts ...firestore.SetOption) (*firestore.BulkWriterJob, error)
	MockUpdate func(doc *firestore.DocumentRef, updates []firestore.Update) (*firestore.BulkWriterJob, error)
	MockDelete func(doc *firestore.DocumentRef) (*firestore.BulkWriterJob, error)
	MockEnd    func()
	MockFlush  func()
}

func (m *MockBulkWriter) Create(doc *firestore.DocumentRef, data interface{}) (*firestore.BulkWriterJob, error) {
	if m.MockCreate != nil {
		return m.MockCreate(doc, data)
	}
	return &firestore.BulkWriterJob{}, nil
}

func (m *MockBulkWriter) Set(doc *firestore.DocumentRef, data interface{}, opts ...firestore.SetOption) (*firestore.BulkWriterJob, error) {
	if m.MockSet != nil {
		return m.MockSet(doc, data, opts...)
	}
	return &firestore.BulkWriterJob{}, nil
}

func (m *MockBulkWriter) Update(doc *firestore.DocumentRef, updates []firestore.Update) (*firestore.BulkWriterJob, error) {
	if m.MockUpdate != nil {
		return m.MockUpdate(doc, updates)
	}
	return &firestore.BulkWriterJob{}, nil
}

func (m *MockBulkWriter) Delete(doc *firestore.DocumentRef) (*firestore.BulkWriterJob, error) {
	if m.MockDelete != nil {
		return m.MockDelete(doc)
	}
	return &firestore.BulkWriterJob{}, nil
}

func (m *MockBulkWriter) End() {
	if m.MockEnd != nil {
		m.MockEnd()
	}
}

func (m *MockBulkWriter) Flush() {
	if m.MockFlush != nil {
		m.MockFlush()
	}
}

// MockCollectionRef is a mock for CollectionRef.
type MockCollectionRef struct {
	MockPath       func() string
	MockDoc        func(id string) interfaces.DocumentRef
	MockAdd        func(ctx context.Context, data any) (interfaces.DocumentRef, *firestore.WriteResult, error)
	MockDocuments  func(ctx context.Context) interfaces.DocumentIterator
	MockStartAfter func(docID string) interfaces.Query
	MockWhere      func(path, op string, value any) interfaces.Query
	MockLimit      func(n int) interfaces.Query
	MockOrderBy    func(path string, direction firestore.Direction) interfaces.Query
}

func (m *MockCollectionRef) Path() string {
	if m.MockPath != nil {
		return m.MockPath()
	}
	return "mock/path"
}

func (m *MockCollectionRef) Doc(id string) interfaces.DocumentRef {
	if m.MockDoc != nil {
		return m.MockDoc(id)
	}
	return &MockDocumentRef{}
}

func (m *MockCollectionRef) Add(ctx context.Context, data any) (interfaces.DocumentRef, *firestore.WriteResult, error) {
	if m.MockAdd != nil {
		return m.MockAdd(ctx, data)
	}
	return &MockDocumentRef{}, &firestore.WriteResult{UpdateTime: time.Now()}, nil
}

func (m *MockCollectionRef) Documents(ctx context.Context) interfaces.DocumentIterator {
	if m.MockDocuments != nil {
		return m.MockDocuments(ctx)
	}
	return &MockDocumentIterator{}
}

func (m *MockCollectionRef) StartAfter(docID string) interfaces.Query {
	if m.MockStartAfter != nil {
		return m.MockStartAfter(docID)
	}
	return &MockCollectionRef{}
}

func (m *MockCollectionRef) Where(path, op string, value any) interfaces.Query {
	if m.MockWhere != nil {
		return m.MockWhere(path, op, value)
	}
	return &MockCollectionRef{}
}

func (m *MockCollectionRef) Limit(n int) interfaces.Query {
	if m.MockLimit != nil {
		return m.MockLimit(n)
	}
	return &MockCollectionRef{}
}

func (m *MockCollectionRef) OrderBy(path string, direction firestore.Direction) interfaces.Query {
	if m.MockOrderBy != nil {
		return m.MockOrderBy(path, direction)
	}
	return &MockCollectionRef{}
}

// MockDocumentRef is a mock for DocumentRef.
type MockDocumentRef struct {
	MockRef        func() *firestore.DocumentRef
	MockID         func() string
	MockGet        func(ctx context.Context) (interfaces.DocumentSnapshot, error)
	MockSet        func(ctx context.Context, data interface{}, opts ...firestore.SetOption) (*firestore.WriteResult, error)
	MockCollection func(name string) interfaces.CollectionRef
	MockDelete     func(ctx context.Context) (*firestore.WriteResult, error)
	// You can add public fields here to allow direct access from tests.
	RefID string
}

func (m *MockDocumentRef) Ref() *firestore.DocumentRef {
	if m.MockRef != nil {
		return m.MockRef()
	}
	return &firestore.DocumentRef{}
}

func (m *MockDocumentRef) ID() string {
	if m.MockID != nil {
		return m.MockID()
	}
	return m.RefID
}

func (m *MockDocumentRef) Get(ctx context.Context) (interfaces.DocumentSnapshot, error) {
	if m.MockGet != nil {
		return m.MockGet(ctx)
	}
	return &MockDocumentSnapshot{}, nil
}

func (m *MockDocumentRef) Set(ctx context.Context, data interface{}, opts ...firestore.SetOption) (*firestore.WriteResult, error) {
	if m.MockSet != nil {
		return m.MockSet(ctx, data, opts...)
	}
	return &firestore.WriteResult{UpdateTime: time.Now()}, nil
}

func (m *MockDocumentRef) Collection(name string) interfaces.CollectionRef {
	if m.MockCollection != nil {
		return m.MockCollection(name)
	}
	return &MockCollectionRef{}
}

func (m *MockDocumentRef) Delete(ctx context.Context) (*firestore.WriteResult, error) {
	if m.MockDelete != nil {
		return m.MockDelete(ctx)
	}
	return &firestore.WriteResult{UpdateTime: time.Now()}, nil
}

// MockDocumentSnapshot is a mock for DocumentSnapshot.
type MockDocumentSnapshot struct {
	MockDataTo func(v any) error
	MockID     func() string
	MockRef    func() *firestore.DocumentRef
	// You can add a public field for the ID to make it easy to set.
	RefID string
}

func (m *MockDocumentSnapshot) DataTo(v any) error {
	if m.MockDataTo != nil {
		return m.MockDataTo(v)
	}
	return nil
}

func (m *MockDocumentSnapshot) ID() string {
	if m.MockID != nil {
		return m.MockID()
	}
	return m.RefID
}

func (m *MockDocumentSnapshot) Ref() *firestore.DocumentRef {
	if m.MockRef != nil {
		return m.MockRef()
	}
	return &firestore.DocumentRef{}
}

// MockDocumentIterator is a mock for DocumentIterator.
type MockDocumentIterator struct {
	Snapshots []interfaces.DocumentSnapshot
	Index     int
}

func (m *MockDocumentIterator) Next() (interfaces.DocumentSnapshot, error) {
	if m.Index >= len(m.Snapshots) {
		return nil, iterator.Done
	}
	doc := m.Snapshots[m.Index]
	m.Index++
	return doc, nil
}

func (m *MockDocumentIterator) Stop() {}

// --- Test Setup Helpers ---

// MockDBConfig allows for easy configuration of the mock database.
// This is the "clever" part that simplifies your tests.
type MockDBConfig struct {
	Collections map[string]MockCollectionConfig
}

type MockCollectionConfig struct {
	Documents map[string]MockDocumentConfig
	// For queries, you can define a slice of documents to be returned.
	QueryDocuments []MockDocumentConfig
	MockError      error
}

type MockDocumentConfig struct {
	Data interface{}
	ID   string
	Err  error
}

// NewMockDatabase sets up a new mock database client based on a config.
func NewMockDatabase(cfg *MockDBConfig) *MockDatabaseClient {
	mockClient := &MockDatabaseClient{}
	mockClient.MockCollection = func(name string) interfaces.CollectionRef {
		collectionConfig, ok := cfg.Collections[name]
		if !ok {
			return &MockCollectionRef{
				MockDoc: func(id string) interfaces.DocumentRef {
					return &MockDocumentRef{
						RefID: id,
						MockGet: func(ctx context.Context) (interfaces.DocumentSnapshot, error) {
							return nil, errors.New("document not found")
						},
					}
				},
			}
		}

		// A mock to handle query-like calls
		mockIter := &MockDocumentIterator{}
		if len(collectionConfig.QueryDocuments) > 0 {
			for _, docCfg := range collectionConfig.QueryDocuments {
				mockSnap := &MockDocumentSnapshot{RefID: docCfg.ID}
				mockSnap.MockDataTo = func(v any) error {
					if docCfg.Data != nil {
						// Simple way to copy data for the mock
						val := docCfg.Data.(interface {
							String() string
						})
						v.(interface {
							Set(s string)
						}).Set(val.String())
					}
					return nil
				}
				mockIter.Snapshots = append(mockIter.Snapshots, mockSnap)
			}
		} else {
			mockIter.Snapshots = []interfaces.DocumentSnapshot{}
		}

		mockCollection := &MockCollectionRef{
			MockDocuments: func(ctx context.Context) interfaces.DocumentIterator {
				return mockIter
			},
			MockAdd: func(ctx context.Context, data any) (interfaces.DocumentRef, *firestore.WriteResult, error) {
				return &MockDocumentRef{}, &firestore.WriteResult{UpdateTime: time.Now()}, collectionConfig.MockError
			},
		}

		mockCollection.MockDoc = func(id string) interfaces.DocumentRef {
			docConfig, ok := collectionConfig.Documents[id]
			if !ok {
				return &MockDocumentRef{
					RefID: id,
					MockGet: func(ctx context.Context) (interfaces.DocumentSnapshot, error) {
						return nil, errors.New("document not found")
					},
				}
			}
			return &MockDocumentRef{
				RefID: id,
				MockGet: func(ctx context.Context) (interfaces.DocumentSnapshot, error) {
					return &MockDocumentSnapshot{
						RefID: docConfig.ID,
						MockDataTo: func(v any) error {
							return nil
						},
					}, docConfig.Err
				},
				MockSet: func(ctx context.Context, data interface{}, opts ...firestore.SetOption) (*firestore.WriteResult, error) {
					return &firestore.WriteResult{UpdateTime: time.Now()}, docConfig.Err
				},
				MockDelete: func(ctx context.Context) (*firestore.WriteResult, error) {
					return &firestore.WriteResult{UpdateTime: time.Now()}, docConfig.Err
				},
			}
		}

		return mockCollection
	}
	return mockClient
}
