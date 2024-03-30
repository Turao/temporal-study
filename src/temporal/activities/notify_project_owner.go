package activities

import (
	"context"

	"github.com/turao/temporal-study/src/api"
	"github.com/turao/temporal-study/src/service"
)

type NotifyProjectOwnerActivity struct {
	NotificationService service.NotificationService
}

type NotifyProjectOwnerActivityRequest struct {
	Request *api.NotifyRequest
}

type NotifyProjectOwnerActivityResponse struct {
	Response *api.NotifyResponse
}

func (npoa *NotifyProjectOwnerActivity) ExecuteNotifyProjectOwnerActivity(
	ctx context.Context,
	req NotifyProjectOwnerActivityRequest,
) (*NotifyProjectOwnerActivityResponse, error) {
	res, err := npoa.NotificationService.Notify(
		ctx,
		&api.NotifyRequest{
			EntityID: req.Request.EntityID,
		},
	)
	if err != nil {
		return nil, err
	}

	return &NotifyProjectOwnerActivityResponse{
		Response: res,
	}, nil
}
