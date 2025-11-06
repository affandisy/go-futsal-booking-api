package request

type CreateVenueRequest struct {
	Name    string `json:"name" validate:"required"`
	Address string `json:"address" validate:"required"`
	City    string `json:"city" validate:"required"`
}

type UpdateVenueRequest struct {
	Name    string `json:"name" validate:"required"`
	Address string `json:"address" validate:"required"`
	City    string `json:"city" validate:"required"`
}
