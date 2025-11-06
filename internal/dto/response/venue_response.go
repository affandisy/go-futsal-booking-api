package response

import (
	"go-futsal-booking-api/internal/domain"
	"time"
)

type VenueResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	City      string    `json:"city"`
	CreatedAt time.Time `json:"created_at"`
}

func ToVenueResponse(venue *domain.Venue) VenueResponse {
	return VenueResponse{
		ID:        venue.ID,
		Name:      venue.Name,
		Address:   venue.Address,
		City:      venue.City,
		CreatedAt: venue.CreatedAt,
	}
}
