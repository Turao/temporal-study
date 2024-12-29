package project

import (
	"context"
	"fmt"
	"log"

	"github.com/gofrs/uuid"
	"github.com/turao/temporal-study/src/api"
	projectentity "github.com/turao/temporal-study/src/entity/project"
	"github.com/turao/temporal-study/src/repository"
	"github.com/turao/temporal-study/src/temporal/workers"

	"github.com/turao/temporal-study/src/temporal/workflows"
	startnewprojectworkflow "github.com/turao/temporal-study/src/temporal/workflows/start-new-project"
	"go.temporal.io/api/enums/v1"
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

// StartNewProject implements service.ProjectService.
func (svc *service) StartNewProject(ctx context.Context, req *api.StartNewProjectRequest) (*api.StartNewProjectResponse, error) {
	options := temporalclient.StartWorkflowOptions{
		ID:                    fmt.Sprintf("%s_%s_%s", workflows.WorkflowNameStartNewProject, req.OwnerID, req.ProjectName),
		TaskQueue:             workers.TaskQueueDefault,
		WorkflowIDReusePolicy: enums.WORKFLOW_ID_REUSE_POLICY_REJECT_DUPLICATE,
	}
	execution, err := svc.temporal.ExecuteWorkflow(
		ctx,
		options,
		workflows.WorkflowNameStartNewProject,
		startnewprojectworkflow.Request{
			ProjectName: req.ProjectName,
			OwnerID:     req.OwnerID,
		},
	)
	if err != nil {
		log.Println("unable to start workflow", err)
		return nil, err
	}

	var startNewProjectWorkflowResponse startnewprojectworkflow.Response
	err = execution.Get(ctx, &startNewProjectWorkflowResponse)
	if err != nil {
		log.Println("unable to get workflow response", err)
		return nil, err
	}

	return &api.StartNewProjectResponse{}, nil
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
