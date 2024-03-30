package service

import (
	"context"

	"github.com/turao/temporal-study/src/api"
)

type ProjectService interface {
	CreateProject(ctx context.Context, req *api.CreateProjectRequest) (*api.CreateProjectResponse, error)
	DeleteProject(ctx context.Context, req *api.DeleteProjectRequest) (*api.DeleteProjectResponse, error)
	UpsertProject(ctx context.Context, req *api.UpsertProjectRequest) (*api.UpsertProjectResponse, error)
}

type NotificationService interface {
	Notify(ctx context.Context, req *api.NotifyRequest) (*api.NotifyResponse, error)
}
