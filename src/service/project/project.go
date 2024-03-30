package project

import (
	"context"
	"fmt"
	"log"

	"github.com/gofrs/uuid"
	"github.com/turao/temporal-study/src/api"
	projectentity "github.com/turao/temporal-study/src/entity/project"
	"github.com/turao/temporal-study/src/repository"
	"github.com/turao/temporal-study/src/temporal"

	"github.com/turao/temporal-study/src/temporal/workflows"
	createprojectworkflow "github.com/turao/temporal-study/src/temporal/workflows/create-project"
	temporalclient "go.temporal.io/sdk/client"
)

type Params struct {
	Temporal          temporalclient.Client
	ProjectRepository repository.ProjectRepository
}

type service struct {
	temporal          temporalclient.Client
	projectRepository repository.ProjectRepository
}

func New(params Params) (*service, error) {
	return &service{
		temporal:          params.Temporal,
		projectRepository: params.ProjectRepository,
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
		workflows.WorkflowNameCreateProject,
		createprojectworkflow.Request{
			ProjectName: req.ProjectName,
			OwnerID:     req.OwnerID,
		},
	)
	if err != nil {
		log.Println("unable to start create project workflow", err)
		return nil, err
	}

	var createProjectWorkflowResponse createprojectworkflow.Response
	err = execution.Get(ctx, &createProjectWorkflowResponse)
	if err != nil {
		log.Println("unable to get create project workflow response", err)
		return nil, err
	}

	return &api.CreateProjectResponse{}, nil
}

// DeleteProject implements service.ProjectService.
func (svc *service) DeleteProject(ctx context.Context, req *api.DeleteProjectRequest) (*api.DeleteProjectResponse, error) {
	panic("unimplemented")
}

func (svc *service) UpsertProject(ctx context.Context, req *api.UpsertProjectRequest) (*api.UpsertProjectResponse, error) {
	log.Println("upserting project", req.ProjectID, req)

	projectID, err := uuid.FromString(req.ProjectID)
	if err != nil {
		return nil, err
	}

	ownerID, err := uuid.FromString(req.OwnerID)
	if err != nil {
		return nil, err
	}

	project, err := projectentity.New(
		projectentity.WithID(projectID),
		projectentity.WithName(req.ProjectName),
		projectentity.WithOwnerID(ownerID),
	)
	if err != nil {
		return nil, err
	}

	err = svc.projectRepository.SaveProject(ctx, project)
	if err != nil {
		return nil, err
	}

	log.Println("project upserted", req.ProjectID, req)
	return &api.UpsertProjectResponse{}, nil
}
