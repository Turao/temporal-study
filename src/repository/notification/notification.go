package notification

import (
	"context"

	notificationentity "github.com/turao/temporal-study/src/entity/notification"
)

type repository struct {
	notifications           map[string]*model
	notificationsByEntityID map[string][]*model
}

func New() (*repository, error) {
	return &repository{
		notifications:           make(map[string]*model),
		notificationsByEntityID: make(map[string][]*model),
	}, nil
}

func (r *repository) SaveNotification(
	ctx context.Context,
	notification *notificationentity.Notification,
) error {
	notificationModel, err := toModel(notification)
	if err != nil {
		return err
	}
	r.notifications[notificationModel.NotificationID] = notificationModel

	if _, ok := r.notificationsByEntityID[notificationModel.EntityID]; !ok {
		r.notificationsByEntityID[notificationModel.EntityID] = make([]*model, 0)
	}
	r.notificationsByEntityID[notificationModel.EntityID] = append(
		r.notificationsByEntityID[notificationModel.EntityID],
		notificationModel,
	)
	return nil
}

func (r *repository) GetNotificationsByEntityID(
	ctx context.Context,
	entityID string,
) ([]*notificationentity.Notification, error) {
	notificationModels, found := r.notificationsByEntityID[entityID]
	if !found {
		return []*notificationentity.Notification{}, nil
	}

	notifications := make([]*notificationentity.Notification, len(notificationModels))
	for _, notificationModel := range notificationModels {
		notification, err := toEntity(notificationModel)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, notification)
	}

	return notifications, nil
}
