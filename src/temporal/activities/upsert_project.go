package activities

import (
	"context"

	"github.com/turao/temporal-study/src/api"
)

type UpsertProjectActivityRequest struct {
	Request *api.UpsertProjectRequest
}

type UpsertProjectActivityResponse struct {
	Response *api.UpsertProjectResponse
}

func UpsertProjectActivity(
	ctx context.Context,
	req *UpsertProjectActivityRequest,
) (*UpsertProjectActivityResponse, error) {
	return nil, nil
}
