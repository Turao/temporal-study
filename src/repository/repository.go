package repository

import (
	"context"

	projectentity "github.com/turao/temporal-study/src/entity/project"
)

type ProjectRepository interface {
	SaveProject(ctx context.Context, project *projectentity.Project) error
	GetProjectsByOwnerID(ctx context.Context, ownerID string) ([]*projectentity.Project, error)
}
