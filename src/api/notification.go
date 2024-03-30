package api

type NotifyRequest struct {
	EntityID string `json:"entity_id"`
}

type NotifyResponse struct {
	NotificationID string `json:"notification_id"`
}
