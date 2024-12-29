package workers

import (
	"github.com/turao/temporal-study/src/service"

	"github.com/turao/temporal-study/src/temporal/activities"
	"github.com/turao/temporal-study/src/temporal/workflows"

	createproject "github.com/turao/temporal-study/src/temporal/activities/create-project"
	notifyprojectowneractivity "github.com/turao/temporal-study/src/temporal/activities/notify-project-owner"
	startnewprojectworkflow "github.com/turao/temporal-study/src/temporal/workflows/start-new-project"

	"go.temporal.io/sdk/activity"
	temporalclient "go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

const (
	TaskQueueDefault = "default"
)

type Params struct {
	Client              temporalclient.Client
	ProjectService      service.ProjectService
	NotificationService service.NotificationService
}

type defaultWorker struct {
	delegate worker.Worker
}

func New(params Params) (*defaultWorker, error) {
	delegate := worker.New(
		params.Client,
		TaskQueueDefault,
		worker.Options{},
	)

	startNewProjectWorkflow := &startnewprojectworkflow.Workflow{}
	delegate.RegisterWorkflowWithOptions(
		startNewProjectWorkflow.Execute,
		workflow.RegisterOptions{
			Name: workflows.WorkflowNameStartNewProject,
		},
	)

	createProjectActivity := &createproject.Activity{
		ProjectService: params.ProjectService,
	}
	delegate.RegisterActivityWithOptions(
		createProjectActivity.Execute,
		activity.RegisterOptions{
			Name: activities.ActivityNameCreateProject,
		},
	)

	notifyProjectOwneractivity := &notifyprojectowneractivity.Activity{
		NotificationService: params.NotificationService,
	}
	delegate.RegisterActivityWithOptions(
		notifyProjectOwneractivity.Execute,
		activity.RegisterOptions{
			Name: activities.ActivityNameNotifyProjectOwner,
		},
	)

	return &defaultWorker{
		delegate: delegate,
	}, nil
}

func (dw *defaultWorker) Run() error {
	return dw.delegate.Run(worker.InterruptCh())
}
