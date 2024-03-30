package project

import (
	"context"
	"fmt"
	"log"

	"github.com/turao/temporal-study/src/api"
	svc "github.com/turao/temporal-study/src/service"
	"github.com/turao/temporal-study/src/temporal"
	"github.com/turao/temporal-study/src/temporal/workflows"
	temporalclient "go.temporal.io/sdk/client"
)

type Params struct {
	Temporal temporalclient.Client
}

type service struct {
	temporal temporalclient.Client
}

var _ svc.ProjectService = (*service)(nil)

func New(params Params) (*service, error) {
	return &service{
		temporal: params.Temporal,
	}, nil
}

// CreateProject implements service.ProjectService.
func (svc *service) CreateProject(ctx context.Context, req *api.CreateProjectRequest) (*api.CreateProjectResponse, error) {
	options := temporalclient.StartWorkflowOptions{
		ID:        fmt.Sprintf("%s_%s", req.OwnerID, req.ProjectName),
		TaskQueue: temporal.WorkerTaskQueue,
	}
	execution, err := svc.temporal.ExecuteWorkflow(
		ctx,
		options,
		workflows.CreateProjectWorkflow,
		workflows.CreateProjectWorkflowRequest{
			Request: req,
		},
	)
	if err != nil {
		log.Println("unable to start create project workflow", err)
		return nil, err
	}

	var createProjectWorkflowResponse workflows.CreateProjectWorkflowResponse
	err = execution.Get(ctx, &createProjectWorkflowResponse)
	if err != nil {
		log.Println("unable to get create project workflow response", err)
		return nil, err
	}

	return createProjectWorkflowResponse.Response, nil
}

// DeleteProject implements service.ProjectService.
func (svc *service) DeleteProject(ctx context.Context, req *api.DeleteProjectRequest) (*api.DeleteProjectResponse, error) {
	panic("unimplemented")
}

func (svc *service) UpsertProject(ctx context.Context, req *api.UpsertProjectRequest) (*api.UpsertProjectResponse, error) {
	panic("unimplemented")
}
