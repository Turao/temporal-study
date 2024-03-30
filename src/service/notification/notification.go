package notification

import (
	"context"
	"log"

	"github.com/gofrs/uuid"
	"github.com/turao/temporal-study/src/api"
	notificationentity "github.com/turao/temporal-study/src/entity/notification"
	"github.com/turao/temporal-study/src/repository"
)

type Params struct {
	NotificationRepository repository.NotificationRepository
}

type service struct {
	notificationRepository repository.NotificationRepository
}

func New(p Params) (*service, error) {
	return &service{
		notificationRepository: p.NotificationRepository,
	}, nil
}

// CreateNotification implements service.NotificationService.
func (svc *service) Notify(ctx context.Context, req *api.NotifyRequest) (*api.NotifyResponse, error) {
	log.Println("notifying entity", req.EntityID, req)

	entityID, err := uuid.FromString(req.EntityID)
	if err != nil {
		return nil, err
	}

	notification, err := notificationentity.New(
		notificationentity.WithEntityID(entityID),
	)
	if err != nil {
		return nil, err
	}

	err = svc.notificationRepository.SaveNotification(ctx, notification)
	if err != nil {
		return nil, err
	}

	log.Println("entity notified", req.EntityID, req)
	return &api.NotifyResponse{}, nil
}
