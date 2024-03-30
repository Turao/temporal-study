package main

import (
	"context"
	"log"

	"github.com/gofrs/uuid"
	"github.com/turao/temporal-study/src/api"
	notificationrepository "github.com/turao/temporal-study/src/repository/notification"
	projectrepository "github.com/turao/temporal-study/src/repository/project"
	notificationservice "github.com/turao/temporal-study/src/service/notification"
	projectservice "github.com/turao/temporal-study/src/service/project"
	"github.com/turao/temporal-study/src/temporal/workers"
	temporalclient "go.temporal.io/sdk/client"
)

func main() {
	client, err := temporalclient.Dial(temporalclient.Options{})
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	notificationRepository, err := notificationrepository.New()
	if err != nil {
		log.Fatalln(err)
	}

	notificationService, err := notificationservice.New(notificationservice.Params{
		NotificationRepository: notificationRepository,
	})
	if err != nil {
		log.Fatalln(err)
	}

	projectRepository, err := projectrepository.New()
	if err != nil {
		log.Fatalln(err)
	}

	projectService, err := projectservice.New(projectservice.Params{
		Temporal:          client,
		ProjectRepository: projectRepository,
	})
	if err != nil {
		log.Fatalln(err)
	}

	numberOfWorkers := 3
	for i := 0; i < numberOfWorkers; i++ {
		defaultWorker, err := workers.New(workers.Params{
			Client:              client,
			ProjectService:      projectService,
			NotificationService: notificationService,
		})
		if err != nil {
			log.Fatalln(err)
		}

		go func() {
			err := defaultWorker.Run()
			if err != nil {
				log.Fatalln(err)
			}
		}()
	}

	res, err := projectService.CreateProject(
		context.Background(),
		&api.CreateProjectRequest{
			ProjectName: "my-project",
			OwnerID:     uuid.Must(uuid.NewV4()).String(),
		},
	)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(res)

	projects, err := projectRepository.GetProjectsByOwnerID(
		context.Background(),
		"178a7ca5-e5cd-4fa7-9d93-5732af1855c9",
	)
	if err != nil {
		log.Fatalln(err)
	}

	for _, project := range projects {
		log.Println("project:", project)
	}
}
