package request

type CreateFieldRequest struct {
	VenueID   uint   `json:"venue_id"`
	Name      string `json:"name"`
	FieldType string `json:"field_type"`
}

type UpdateFieldRequest struct {
	Name      string `json:"name"`
	FieldType string `json:"field_type"`
}
