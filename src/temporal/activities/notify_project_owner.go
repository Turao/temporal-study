package activities

import (
	"context"

	"github.com/turao/temporal-study/src/api"
)

type NotifyProjectOwnerActivityRequest struct {
	Request *api.NotifyRequest
}

type NotifyProjectOwnerActivityResponse struct {
	Response *api.NotifyResponse
}

func NotifyProjectOwnerActivity(
	ctx context.Context,
	req NotifyProjectOwnerActivityRequest,
) (*NotifyProjectOwnerActivityResponse, error) {
	return nil, nil
}
