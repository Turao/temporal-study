package temporal

import (
	"log"

	"github.com/turao/temporal-study/src/service"
	"github.com/turao/temporal-study/src/temporal/activities"
	"github.com/turao/temporal-study/src/temporal/workflows"
	temporalclient "go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
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

	w.RegisterWorkflow(workflows.CreateProjectWorkflow)

	w.RegisterActivity(&activities.UpsertProjectActivity{
		ProjectService: t.projectService,
	})
	w.RegisterActivity(&activities.NotifyProjectOwnerActivity{
		NotificationService: t.notificationService,
	})

	err := w.Run(worker.InterruptCh())
	if err != nil {
		log.Println("unable to start worker", err)
		return err
	}

	return nil
}
