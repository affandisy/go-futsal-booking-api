package docs

type SuccessResponse struct {
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message" example:"Operation successful"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Success bool        `json:"success" example:"false"`
	Error   string      `json:"error" example:"ERROR_CODE"`
	Message string      `json:"message" example:"An error occurred"`
	Details interface{} `json:"details,omitempty"`
}
