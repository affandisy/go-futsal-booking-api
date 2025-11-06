package request

type CreateFieldRequest struct {
	VenueID   uint   `json:"venue_id" validate:"required"`
	Name      string `json:"name" validate:"required"`
	FieldType string `json:"field_type" validate:"required"`
}

type UpdateFieldRequest struct {
	Name      string `json:"name" validate:"required"`
	FieldType string `json:"field_type" validate:"required"`
}
