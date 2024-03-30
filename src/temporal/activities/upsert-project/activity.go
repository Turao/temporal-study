package upsertproject

import (
	"context"

	"github.com/turao/temporal-study/src/api"
	"github.com/turao/temporal-study/src/service"
)

type Activity struct {
	ProjectService service.ProjectService
}

type Request struct {
	Request *api.UpsertProjectRequest
}

type Response struct {
	Response *api.UpsertProjectResponse
}

func (upa *Activity) Execute(ctx context.Context, req Request) (*Response, error) {
	res, err := upa.ProjectService.UpsertProject(ctx, req.Request)
	if err != nil {
		return nil, err
	}

	return &Response{
		Response: res,
	}, nil
}
