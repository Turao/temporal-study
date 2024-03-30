package temporal

import (
	"log"

	"github.com/turao/temporal-study/src/temporal/activities"
	"github.com/turao/temporal-study/src/temporal/workflows"
	temporalclient "go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

const (
	WorkerTaskQueue = "WorkerTaskQueue"
)

type Params struct {
	Temporal temporalclient.Client
}

type service struct {
	temporal temporalclient.Client
}

func New(params Params) (*service, error) {
	return &service{
		temporal: params.Temporal,
	}, nil
}

func (svc *service) Run() error {
	w := worker.New(
		svc.temporal,
		WorkerTaskQueue,
		worker.Options{},
	)

	w.RegisterWorkflow(workflows.CreateProjectWorkflow)
	w.RegisterActivity(activities.UpsertProjectActivity)
	w.RegisterActivity(activities.NotifyProjectOwnerActivity)

	err := w.Run(worker.InterruptCh())
	if err != nil {
		log.Println("unable to start worker", err)
		return err
	}

	return nil
}
