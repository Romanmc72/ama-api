/**
 * The ENTIRE reason that this file exists is because the document snapshots
 * when using a document iteratore CANNOT access the underlying document ID
 * values without relying directly on chained calls to public fields. While
 * this is "convenient" for the people writing the code to access the ID
 * (no need to call a method, bye bye parentheses), it is a massive pain
 * in the ass when you're trying to get your unit test coverage up. I
 * need to wrap everything all the way down to the individual document snapshot
 * level to ensure that I can get the ID with my mocks and my production data
 * the same way.
 */
package database

import (
	"ama/api/interfaces"
	"context"

	"cloud.google.com/go/firestore"
	"github.com/google/uuid"
)

type FirestoreClient struct {
	client *firestore.Client
}

func (fc *FirestoreClient) NewID() string {
	return uuid.NewString()
}

func (fc *FirestoreClient) Collection(name string) interfaces.CollectionRef {
	return &FirestoreCollectionRef{
		collectionRef: *fc.client.Collection(name),
	}
}

type FirestoreCollectionRef struct {
	collectionRef firestore.CollectionRef
}

func (fcr *FirestoreCollectionRef) Path() string {
	return fcr.collectionRef.Path
}

func (fcr *FirestoreCollectionRef) Doc(id string) interfaces.DocumentRef {
	return &FirestoreDocumentRef{
		docRef: *fcr.collectionRef.Doc(id),
	}
}

func (fcr *FirestoreCollectionRef) Add(ctx context.Context, data any) (interfaces.DocumentRef, *firestore.WriteResult, error) {
	docRef, wr, err := fcr.collectionRef.Add(ctx, data)
	if err != nil {
		return &FirestoreDocumentRef{}, wr, err
	}
	return &FirestoreDocumentRef{docRef: *docRef}, wr, err
}

func (fcr *FirestoreCollectionRef) Documents(ctx context.Context) interfaces.DocumentIterator {
	return &FirestoreDocumentIterator{iterator: *fcr.collectionRef.Documents(ctx)}
}

func (fcr *FirestoreCollectionRef) StartAfter(docID string) interfaces.Query {
	return &FirestoreQuery{query: fcr.collectionRef.StartAfter(docID)}
}

func (fcr *FirestoreCollectionRef) Where(path string, op string, value any) interfaces.Query {
	return &FirestoreQuery{query: fcr.collectionRef.Where(path, op, value)}
}

func (fcr *FirestoreCollectionRef) Limit(n int) interfaces.Query {
	return &FirestoreQuery{query: fcr.collectionRef.Limit(n)}
}

func (fcr *FirestoreCollectionRef) OrderBy(path string, direction firestore.Direction) interfaces.Query {
	return &FirestoreQuery{query: fcr.collectionRef.OrderBy(path, direction)}
}

func (fc *FirestoreClient) RunTransaction(
	ctx context.Context,
	f func(context.Context, *firestore.Transaction) error,
	t ...firestore.TransactionOption,
) error {
	return fc.client.RunTransaction(ctx, f, t...)
}
func (fc *FirestoreClient) Close() error {
	return fc.client.Close()
}
func (fc *FirestoreClient) BulkWriter(ctx context.Context) interfaces.BulkWriter {
	return &FirestoreBulkWriter{bw: *fc.client.BulkWriter(ctx)}
}

type FirestoreTransaction struct {
	tx firestore.Transaction
}

func (ft *FirestoreTransaction) Get(doc *firestore.DocumentRef) (interfaces.DocumentSnapshot, error) {
	ds, err := ft.tx.Get(doc)
	if err != nil {
		return &FirestoreDocumentSnapshot{}, err
	}
	return &FirestoreDocumentSnapshot{ds: *ds}, err
}

func (ft *FirestoreTransaction) Set(doc *firestore.DocumentRef, data interface{}, opts ...firestore.SetOption) error {
	return ft.tx.Set(doc, data, opts...)
}

type FirestoreBulkWriter struct {
	bw firestore.BulkWriter
}

func (fbw *FirestoreBulkWriter) Create(doc *firestore.DocumentRef, data interface{}) (*firestore.BulkWriterJob, error) {
	return fbw.bw.Create(doc, data)
}

func (fbw *FirestoreBulkWriter) Set(doc *firestore.DocumentRef, data interface{}, opts ...firestore.SetOption) (*firestore.BulkWriterJob, error) {
	return fbw.bw.Set(doc, data, opts...)
}

func (fbw *FirestoreBulkWriter) Update(doc *firestore.DocumentRef, updates []firestore.Update) (*firestore.BulkWriterJob, error) {
	return fbw.bw.Update(doc, updates)
}

func (fbw *FirestoreBulkWriter) Delete(doc *firestore.DocumentRef) (*firestore.BulkWriterJob, error) {
	return fbw.bw.Delete(doc)
}

func (fbw *FirestoreBulkWriter) End() {
	fbw.bw.End()
}

func (fbw *FirestoreBulkWriter) Flush() {
	fbw.bw.Flush()
}

type FirestoreDocumentRef struct {
	docRef firestore.DocumentRef
}

func (fdr *FirestoreDocumentRef) Ref() *firestore.DocumentRef {
	// This is also the same problem as the document snapshot
	return &fdr.docRef
}

func (fdr *FirestoreDocumentRef) ID() string {
	// This is also the same problem as the document snapshot
	return fdr.docRef.ID
}

func (fdr *FirestoreDocumentRef) Get(ctx context.Context) (interfaces.DocumentSnapshot, error) {
	ds, err := fdr.docRef.Get(ctx)
	if err != nil {
		return &FirestoreDocumentSnapshot{}, err
	}
	return &FirestoreDocumentSnapshot{ds: *ds}, err
}

func (fdr *FirestoreDocumentRef) Set(ctx context.Context, data interface{}, opts ...firestore.SetOption) (*firestore.WriteResult, error) {
	return fdr.docRef.Set(ctx, data, opts...)
}

func (fdr *FirestoreDocumentRef) Collection(name string) interfaces.CollectionRef {
	return &FirestoreCollectionRef{collectionRef: *fdr.docRef.Collection(name)}
}

func (fdr *FirestoreDocumentRef) Delete(ctx context.Context) (*firestore.WriteResult, error) {
	return fdr.docRef.Delete(ctx)
}

type FirestoreDocumentSnapshot struct {
	ds firestore.DocumentSnapshot
}

func (fds *FirestoreDocumentSnapshot) DataTo(v any) error {
	return fds.ds.DataTo(v)
}

func (fds *FirestoreDocumentSnapshot) Ref() *firestore.DocumentRef {
	// another problem child here
	return fds.ds.Ref
}

func (fds *FirestoreDocumentSnapshot) ID() string {
	// THIS right here is the whole reason for this whole file, otherwise interfaces would have sufficed
	// <insert emoji="middle finger" />
	return fds.ds.Ref.ID
}

func (fds *FirestoreDocumentSnapshot) Exists() bool {
	// THIS right here is the whole reason for this whole file, otherwise interfaces would have sufficed
	// <insert emoji="middle finger" />
	return fds.ds.Exists()
}

type FirestoreDocumentIterator struct {
	iterator firestore.DocumentIterator
}

func (fdi *FirestoreDocumentIterator) Next() (interfaces.DocumentSnapshot, error) {
	ds, err := fdi.iterator.Next()
	if err != nil {
		return &FirestoreDocumentSnapshot{}, err
	}
	return &FirestoreDocumentSnapshot{ds: *ds}, err
}

func (fdi *FirestoreDocumentIterator) Stop() {
	fdi.iterator.Stop()
}

type FirestoreQuery struct {
	query firestore.Query
}

func (fq *FirestoreQuery) Documents(ctx context.Context) interfaces.DocumentIterator {
	return &FirestoreDocumentIterator{
		iterator: *fq.query.Documents(ctx),
	}
}

func (fq *FirestoreQuery) StartAfter(docID string) interfaces.Query {
	return &FirestoreQuery{
		query: fq.query.StartAfter(docID),
	}
}

func (fq *FirestoreQuery) Where(path, op string, value interface{}) interfaces.Query {
	return &FirestoreQuery{
		query: fq.query.Where(path, op, value),
	}
}

func (fq *FirestoreQuery) Limit(n int) interfaces.Query {
	return &FirestoreQuery{
		query: fq.query.Limit(n),
	}
}

func (fq *FirestoreQuery) OrderBy(path string, direction firestore.Direction) interfaces.Query {
	return &FirestoreQuery{
		query: fq.query.OrderBy(path, direction),
	}
}
