package activities

import (
	"context"

	"github.com/turao/temporal-study/src/api"
	"github.com/turao/temporal-study/src/service"
)

type UpsertProjectActivity struct {
	ProjectService service.ProjectService
}

type UpsertProjectActivityRequest struct {
	Request *api.UpsertProjectRequest
}

type UpsertProjectActivityResponse struct {
	Response *api.UpsertProjectResponse
}

func (upa *UpsertProjectActivity) ExecuteUpsertProjectActivity(
	ctx context.Context,
	req UpsertProjectActivityRequest,
) (*UpsertProjectActivityResponse, error) {
	res, err := upa.ProjectService.UpsertProject(ctx, req.Request)
	if err != nil {
		return nil, err
	}

	return &UpsertProjectActivityResponse{
		Response: res,
	}, nil
}
