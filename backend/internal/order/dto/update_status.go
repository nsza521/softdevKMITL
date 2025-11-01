package dto

type UpdateStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

type UpdateStatusResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}
