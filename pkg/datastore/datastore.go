package datastore

import (
	"context"

	"google.golang.org/appengine/datastore"
)

// Save ...
func Save(ctx context.Context, src Resource) (err error) {
	_, err = datastore.Put(ctx, src.Key(ctx), src)
	return
}

// Load ...
func Load(ctx context.Context, src Resource) (err error) {
	return datastore.Get(ctx, src.Key(ctx), src)
}

// Delete ...
func Delete(ctx context.Context, key *datastore.Key) error {
	return datastore.Delete(ctx, key)
}

// Resource ...
type Resource interface {
	KindName() string
	Key(ctx context.Context) *datastore.Key
	Save(ctx context.Context) error
	Load(ctx context.Context) error
	Delete(ctx context.Context) error
}
