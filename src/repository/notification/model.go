package notification

import (
	"github.com/gofrs/uuid"
	notificationentity "github.com/turao/temporal-study/src/entity/notification"
)

type model struct {
	NotificationID string
	EntityID       string
}

func toModel(notificationEntity *notificationentity.Notification) (*model, error) {
	return &model{
		NotificationID: notificationEntity.ID.String(),
		EntityID:       notificationEntity.EntityID.String(),
	}, nil
}

func toEntity(notificationModel *model) (*notificationentity.Notification, error) {
	notificationID, err := uuid.FromString(notificationModel.NotificationID)
	if err != nil {
		return nil, err
	}

	entityID, err := uuid.FromString(notificationModel.EntityID)
	if err != nil {
		return nil, err
	}

	return notificationentity.New(
		notificationentity.WithID(notificationID),
		notificationentity.WithEntityID(entityID),
	)
}
