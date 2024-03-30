package notifyprojectowner

import (
	"context"

	"github.com/turao/temporal-study/src/api"
	"github.com/turao/temporal-study/src/service"
)

type Activity struct {
	NotificationService service.NotificationService
}

type Request struct {
	Request *api.NotifyRequest
}

type Response struct {
	Response *api.NotifyResponse
}

func (npoa *Activity) Execute(ctx context.Context, req Request) (*Response, error) {
	res, err := npoa.NotificationService.Notify(
		ctx,
		&api.NotifyRequest{
			EntityID: req.Request.EntityID,
		},
	)
	if err != nil {
		return nil, err
	}

	return &Response{
		Response: res,
	}, nil
}
