package notification

import (
	"errors"

	"github.com/gofrs/uuid"
)

type Notification struct {
	ID       uuid.UUID
	EntityID uuid.UUID
}

func New(opts ...NotificationOption) (*Notification, error) {
	uuid, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	notification := &Notification{
		ID: uuid,
	}

	errs := []error{}
	for _, opt := range opts {
		errs = append(errs, opt(notification))
	}
	err = errors.Join(errs...)
	if err != nil {
		return nil, err
	}

	return notification, nil
}

type NotificationOption func(n *Notification) error

func WithID(id uuid.UUID) NotificationOption {
	return func(n *Notification) error {
		n.ID = id
		return nil
	}
}

func WithEntityID(entityID uuid.UUID) NotificationOption {
	return func(n *Notification) error {
		n.EntityID = entityID
		return nil
	}
}
