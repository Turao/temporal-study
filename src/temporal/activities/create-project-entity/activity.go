package createprojectentity

import (
	"context"

	"github.com/gofrs/uuid"
	projectentity "github.com/turao/temporal-study/src/entity/project"
)

type Activity struct{}

type Request struct {
	ProjectName string
	OwnerID     string
}

type Response struct {
	Entity *projectentity.Project
}

func (cpe *Activity) Execute(ctx context.Context, req Request) (*Response, error) {
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
	return &Response{
		Entity: project,
	}, nil
}
