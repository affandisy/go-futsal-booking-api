package response

import (
	"go-futsal-booking-api/internal/domain"
	"time"
)

type FieldResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Venue     string    `json:"venue"`
	CreatedAt time.Time `json:"created_at"`
}

func ToFieldResponse(field *domain.Field) FieldResponse {
	return FieldResponse{
		ID:        field.ID,
		Name:      field.Name,
		Type:      field.Type,
		Venue:     field.Venue.Name,
		CreatedAt: field.CreatedAt,
	}
}
