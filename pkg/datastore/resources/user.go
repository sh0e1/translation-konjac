package resources

import (
	"context"
	"time"

	"google.golang.org/appengine/datastore"

	ds "github.com/sh0e1/translation-konjac/pkg/datastore"
)

// NewUser ...
func NewUser(id string, lang string) *User {
	return &User{
		ID:             id,
		SelectLanguage: lang,
	}
}

// User ...
type User struct {
	ID             string
	SelectLanguage string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	key            *datastore.Key `datastore:"-"`
}

// KindName ...
func (u *User) KindName() string {
	return "Users"
}

// Key ..
func (u *User) Key(ctx context.Context) *datastore.Key {
	if u.key == nil {
		u.key = datastore.NewKey(ctx, u.KindName(), u.ID, 0, nil)
	}
	return u.key
}

// Save ...
func (u *User) Save(ctx context.Context) error {
	now := time.Now()
	if u.CreatedAt.IsZero() {
		u.CreatedAt = now
	}
	u.UpdatedAt = now
	return ds.Save(ctx, u)
}

// Load...
func (u *User) Load(ctx context.Context) error {
	return ds.Load(ctx, u)
}

// Delete ...
func (u *User) Delete(ctx context.Context) error {
	return ds.Delete(ctx, u.key)
}
