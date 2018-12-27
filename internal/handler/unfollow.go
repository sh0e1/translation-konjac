package handler

import (
	"context"

	"github.com/sh0e1/translation-konjac/pkg/datastore/resources"
)

type UnfollowHandler struct {
	*BaseEventHandler
}

// Handle ...
func (h *UnfollowHandler) Handle(ctx context.Context) error {
	user := &resources.User{ID: h.Event.Source.UserID}
	return user.Delete(ctx)
}
