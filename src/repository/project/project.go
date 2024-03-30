package project

import (
	"context"

	projectentity "github.com/turao/temporal-study/src/entity/project"
)

type repository struct {
	projects          map[string]*model
	projectsByOwnerID map[string][]*model
}

func New() (*repository, error) {
	return &repository{
		projects:          make(map[string]*model),
		projectsByOwnerID: make(map[string][]*model),
	}, nil
}

func (r *repository) SaveProject(
	ctx context.Context,
	project *projectentity.Project,
) error {
	projectModel, err := toModel(project)
	if err != nil {
		return err
	}
	r.projects[projectModel.ProjectID] = projectModel

	if _, ok := r.projectsByOwnerID[projectModel.OwnerID]; !ok {
		r.projectsByOwnerID[projectModel.OwnerID] = make([]*model, 0)
	}
	r.projectsByOwnerID[projectModel.OwnerID] = append(
		r.projectsByOwnerID[projectModel.OwnerID],
		projectModel,
	)
	return nil
}

func (r *repository) GetProjectsByOwnerID(
	ctx context.Context,
	ownerID string,
) ([]*projectentity.Project, error) {
	projectModels, found := r.projectsByOwnerID[ownerID]
	if !found {
		return []*projectentity.Project{}, nil
	}

	projects := make([]*projectentity.Project, len(projectModels))
	for _, projectModel := range projectModels {
		project, err := toEntity(projectModel)
		if err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}

	return projects, nil
}
