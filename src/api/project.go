package api

type StartNewProjectRequest struct {
	ProjectName string `json:"project_name"`
	OwnerID     string `json:"owner_id"`
}

type StartNewProjectResponse struct{}

type DeleteProjectRequest struct {
	ProjectID string `json:"project_id"`
}

type DeleteProjectResponse struct{}

type UpsertProjectRequest struct {
	ProjectID   string `json:"project_id"`
	ProjectName string `json:"project_name"`
	OwnerID     string `json:"owner_id"`
}

type UpsertProjectResponse struct{}
