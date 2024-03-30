package repository

import (
	"context"

	notificationentity "github.com/turao/temporal-study/src/entity/notification"
	projectentity "github.com/turao/temporal-study/src/entity/project"
)

type ProjectRepository interface {
	SaveProject(ctx context.Context, project *projectentity.Project) error
	GetProjectsByOwnerID(ctx context.Context, ownerID string) ([]*projectentity.Project, error)
}

type NotificationRepository interface {
	SaveNotification(ctx context.Context, notification *notificationentity.Notification) error
	GetNotificationsByEntityID(ctx context.Context, entityID string) ([]*notificationentity.Notification, error)
}
