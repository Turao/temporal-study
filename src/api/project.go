package api

type CreateProjectRequest struct {
	ProjectName string `json:"project_name"`
	OwnerID     string `json:"owner_id"`
}

type CreateProjectResponse struct{}

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
