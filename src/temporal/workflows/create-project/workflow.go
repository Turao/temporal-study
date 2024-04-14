package createproject

import (
	"log"
	"time"

	"github.com/turao/temporal-study/src/api"
	"github.com/turao/temporal-study/src/temporal/activities"
	createprojectentity "github.com/turao/temporal-study/src/temporal/activities/create-project-entity"
	notifyprojectowneractivity "github.com/turao/temporal-study/src/temporal/activities/notify-project-owner"
	upsertprojectactivity "github.com/turao/temporal-study/src/temporal/activities/upsert-project"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type Workflow struct{}

type Request struct {
	ProjectName string
	OwnerID     string
}

type Response struct{}

func (w *Workflow) Execute(ctx workflow.Context, req Request) (*Response, error) {
	log.Println("starting create project workflow", req)

	// apply activity options
	ctx = workflow.WithActivityOptions(
		ctx,
		workflow.ActivityOptions{
			StartToCloseTimeout: time.Minute,
			RetryPolicy: &temporal.RetryPolicy{
				MaximumAttempts: 3,
			},
		},
	)

	// create the entity
	log.Println("creating the project")
	var createProjectEntityResponse *createprojectentity.Response
	err := workflow.ExecuteActivity(
		ctx,
		activities.ActivityNameCreateProjectEntity,
		createprojectentity.Request{
			ProjectName: req.ProjectName,
			OwnerID:     req.OwnerID,
		},
	).Get(
		ctx,
		&createProjectEntityResponse,
	)
	if err != nil {
		log.Println("unable to create project entity", err)
		return nil, err
	}

	project := createProjectEntityResponse.Entity

	// store in the repository
	log.Println("saving the project")
	var upsertProjectActivityResponse *upsertprojectactivity.Response
	err = workflow.ExecuteActivity(
		ctx,
		activities.ActivityNameUpsertProject,
		upsertprojectactivity.Request{
			Request: &api.UpsertProjectRequest{
				ProjectID: project.ID.String(),
				OwnerID:   project.OwnerID.String(),
			},
		},
	).Get(
		ctx,
		&upsertProjectActivityResponse,
	)
	if err != nil {
		log.Println("unable to upsert project", err)
		return nil, err
	}

	// notify the owner
	log.Println("notifying the project owner")
	var notifyProjectOwnerActivityResponse *notifyprojectowneractivity.Response
	err = workflow.ExecuteActivity(
		ctx,
		activities.ActivityNameNotifyProjectOwner,
		notifyprojectowneractivity.Request{
			Request: &api.NotifyRequest{
				EntityID: project.OwnerID.String(),
			},
		},
	).Get(
		ctx,
		&notifyProjectOwnerActivityResponse,
	)
	if err != nil {
		log.Println("unable to notify project owner", err)
		return nil, err
	}

	log.Println("workflow completed")
	return &Response{}, nil
}
