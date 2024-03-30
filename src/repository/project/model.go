package project

import (
	"github.com/gofrs/uuid"
	projectentity "github.com/turao/temporal-study/src/entity/project"
)

type model struct {
	ProjectID   string
	ProjectName string
	OwnerID     string
}

func toModel(projectEntity *projectentity.Project) (*model, error) {
	return &model{
		ProjectID:   projectEntity.ID.String(),
		ProjectName: projectEntity.Name,
		OwnerID:     projectEntity.OwnerID.String(),
	}, nil
}

func toEntity(projectModel *model) (*projectentity.Project, error) {
	projectID, err := uuid.FromString(projectModel.ProjectID)
	if err != nil {
		return nil, err
	}

	ownerID, err := uuid.FromString(projectModel.OwnerID)
	if err != nil {
		return nil, err
	}

	return projectentity.New(
		projectentity.WithID(projectID),
		projectentity.WithName(projectModel.ProjectName),
		projectentity.WithOwnerID(ownerID),
	)
}
