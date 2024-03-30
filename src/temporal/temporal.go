package temporal

import (
	"log"

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
	WorkerTaskQueue = "WorkerTaskQueue"
)

type Params struct {
	Client              temporalclient.Client
	ProjectService      service.ProjectService
	NotificationService service.NotificationService
}

type temporal struct {
	client              temporalclient.Client
	projectService      service.ProjectService
	notificationService service.NotificationService
}

func New(params Params) (*temporal, error) {
	return &temporal{
		client:              params.Client,
		projectService:      params.ProjectService,
		notificationService: params.NotificationService,
	}, nil
}

func (t *temporal) Run() error {
	w := worker.New(
		t.client,
		WorkerTaskQueue,
		worker.Options{},
	)

	createProjectWorkflow := &createprojectworkflow.Workflow{}
	w.RegisterWorkflowWithOptions(
		createProjectWorkflow.Execute,
		workflow.RegisterOptions{
			Name: workflows.WorkflowNameCreateProject,
		},
	)

	upsertProjectActivity := &upsertprojectactivity.Activity{
		ProjectService: t.projectService,
	}
	w.RegisterActivityWithOptions(
		upsertProjectActivity.Execute,
		activity.RegisterOptions{
			Name: activities.ActivityNameUpsertProject,
		},
	)

	notifyProjectOwneractivity := &notifyprojectowneractivity.Activity{
		NotificationService: t.notificationService,
	}
	w.RegisterActivityWithOptions(
		notifyProjectOwneractivity.Execute,
		activity.RegisterOptions{
			Name: activities.ActivityNameNotifyProjectOwner,
		},
	)

	err := w.Run(worker.InterruptCh())
	if err != nil {
		log.Println("unable to start worker", err)
		return err
	}

	return nil
}
