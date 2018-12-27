package resources

import (
	"context"
	"time"

	"google.golang.org/appengine/datastore"

	ds "github.com/sh0e1/translation-konjac/pkg/datastore"
)

// NewAudio ...
func NewAudio(id, uid, path string) *Audio {
	return &Audio{
		ID:     id,
		UserID: uid,
		Path:   path,
	}
}

// Audio ...
type Audio struct {
	ID             string
	UserID         string
	Path           string
	SourceLanguage string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	key            *datastore.Key `datastore:"-"`
}

// KindName ...
func (a *Audio) KindName() string {
	return "Audios"
}

// Key ..
func (a *Audio) Key(ctx context.Context) *datastore.Key {
	if a.key == nil {
		a.key = datastore.NewKey(ctx, a.KindName(), a.ID, 0, nil)
	}
	return a.key
}

// Save ...
func (a *Audio) Save(ctx context.Context) error {
	now := time.Now()
	if a.CreatedAt.IsZero() {
		a.CreatedAt = now
	}
	a.UpdatedAt = now
	return ds.Save(ctx, a)
}

// Load...
func (a *Audio) Load(ctx context.Context) error {
	return ds.Load(ctx, a)
}

// Delete ...
func (a *Audio) Delete(ctx context.Context) error {
	return ds.Delete(ctx, a.key)
}
