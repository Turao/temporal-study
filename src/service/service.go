package service

import (
	"context"

	"github.com/turao/temporal-study/src/api"
)

type ProjectService interface {
	StartNewProject(ctx context.Context, req *api.StartNewProjectRequest) (*api.StartNewProjectResponse, error)
	DeleteProject(ctx context.Context, req *api.DeleteProjectRequest) (*api.DeleteProjectResponse, error)
	UpsertProject(ctx context.Context, req *api.UpsertProjectRequest) (*api.UpsertProjectResponse, error)
}

type NotificationService interface {
	Notify(ctx context.Context, req *api.NotifyRequest) (*api.NotifyResponse, error)
}
