package startnewproject

import (
	"log"
	"time"

	"github.com/turao/temporal-study/src/api"
	"github.com/turao/temporal-study/src/temporal/activities"
	createproject "github.com/turao/temporal-study/src/temporal/activities/create-project"
	notifyprojectowneractivity "github.com/turao/temporal-study/src/temporal/activities/notify-project-owner"
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
	log.Println("starting workflow", req)

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

	// create the project
	log.Println("creating the project")
	var createProjectResponse *createproject.Response
	err := workflow.ExecuteActivity(
		ctx,
		activities.ActivityNameCreateProject,
		createproject.Request{
			ProjectName: req.ProjectName,
			OwnerID:     req.OwnerID,
		},
	).Get(
		ctx,
		&createProjectResponse,
	)
	if err != nil {
		log.Println("unable to create project", err)
		return nil, err
	}

	project := createProjectResponse.Entity

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
