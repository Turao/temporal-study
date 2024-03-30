package workers

import (
	"github.com/turao/temporal-study/src/service"

	"github.com/turao/temporal-study/src/temporal/activities"
	"github.com/turao/temporal-study/src/temporal/workflows"

	notifyprojectowneractivity "github.com/turao/temporal-study/src/temporal/activities/notify-project-owner"
	upsertprojectactivity "github.com/turao/temporal-study/src/temporal/activities/upsert-project"
	createprojectworkflow "github.com/turao/temporal-study/src/temporal/workflows/create-project"

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

	createProjectWorkflow := &createprojectworkflow.Workflow{}
	delegate.RegisterWorkflowWithOptions(
		createProjectWorkflow.Execute,
		workflow.RegisterOptions{
			Name: workflows.WorkflowNameCreateProject,
		},
	)

	upsertProjectActivity := &upsertprojectactivity.Activity{
		ProjectService: params.ProjectService,
	}
	delegate.RegisterActivityWithOptions(
		upsertProjectActivity.Execute,
		activity.RegisterOptions{
			Name: activities.ActivityNameUpsertProject,
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
