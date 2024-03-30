package notification

import (
	"context"

	"github.com/turao/temporal-study/src/api"
	svc "github.com/turao/temporal-study/src/service"
)

type service struct{}

var _ svc.NotificationService = (*service)(nil)

func New() (*service, error) {
	return &service{}, nil
}

// CreateNotification implements service.NotificationService.
func (svc *service) Notify(ctx context.Context, req *api.NotifyRequest) (*api.NotifyResponse, error) {
	panic("unimplemented")
}
