package createproject

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/turao/temporal-study/src/api"
	projectentity "github.com/turao/temporal-study/src/entity/project"
	"github.com/turao/temporal-study/src/service"
)

type Activity struct {
	ProjectService service.ProjectService
}

type Request struct {
	ProjectName string
	OwnerID     string
}

type Response struct {
	Entity *projectentity.Project
}

func (a *Activity) Execute(ctx context.Context, req Request) (*Response, error) {
	ownerID, err := uuid.FromString(req.OwnerID)
	if err != nil {
		return nil, err
	}
	project, err := projectentity.New(
		projectentity.WithName(req.ProjectName),
		projectentity.WithOwnerID(ownerID),
	)
	if err != nil {
		return nil, err
	}

	_, err = a.ProjectService.UpsertProject(ctx, &api.UpsertProjectRequest{
		ProjectID:   project.ID.String(),
		ProjectName: project.Name,
		OwnerID:     project.OwnerID.String(),
	})
	if err != nil {
		return nil, err
	}

	return &Response{
		Entity: project,
	}, nil
}
