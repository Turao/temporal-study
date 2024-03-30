package workflows

import (
	"log"
	"time"

	"github.com/gofrs/uuid"
	"github.com/turao/temporal-study/src/api"
	projectentity "github.com/turao/temporal-study/src/entity/project"
	"github.com/turao/temporal-study/src/temporal/activities"
	"go.temporal.io/sdk/workflow"
)

type CreateProjectWorkflowRequest struct {
	Request *api.CreateProjectRequest
}

type CreateProjectWorkflowResponse struct {
	Response *api.CreateProjectResponse
}

func CreateProjectWorkflow(
	ctx workflow.Context,
	req CreateProjectWorkflowRequest,
) (*CreateProjectWorkflowResponse, error) {
	// apply activity options
	ctx = workflow.WithActivityOptions(
		ctx,
		workflow.ActivityOptions{
			StartToCloseTimeout: time.Minute,
		},
	)

	// create the entity
	ownerID, err := uuid.FromString(req.Request.OwnerID)
	if err != nil {
		return nil, err
	}
	project, err := projectentity.New(
		projectentity.WithName(req.Request.ProjectName),
		projectentity.WithOwnerID(ownerID),
	)
	if err != nil {
		return nil, err
	}

	// store in the repository
	var upsertProjectActivityResponse *activities.UpsertProjectActivityResponse
	err = workflow.ExecuteActivity(
		ctx,
		activities.UpsertProjectActivity,
		activities.UpsertProjectActivityRequest{
			Request: &api.UpsertProjectRequest{
				ProjectID: project.ID.String(),
				OwnerID:   project.OwnerID.String(),
			},
		},
	).Get(
		ctx,
		upsertProjectActivityResponse,
	)
	if err != nil {
		log.Println("unable to upsert project", err)
		return nil, err
	}

	// notify the owner
	var notifyProjectOwnerActivityResponse *activities.NotifyProjectOwnerActivityResponse
	err = workflow.ExecuteActivity(
		ctx,
		activities.NotifyProjectOwnerActivity,
		activities.NotifyProjectOwnerActivityRequest{
			Request: &api.NotifyRequest{
				EntityID: project.OwnerID.String(),
			},
		},
	).Get(
		ctx,
		notifyProjectOwnerActivityResponse,
	)
	if err != nil {
		log.Println("unable to notify project owner", err)
		return nil, err
	}

	return nil, nil
}
