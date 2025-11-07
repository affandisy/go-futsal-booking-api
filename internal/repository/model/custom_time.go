package model

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// TimeOfDay represents a time without date (HH:MM:SS)
type TimeOfDay struct {
	time.Time
}

// Scan implements the sql.Scanner interface
func (t *TimeOfDay) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case []byte:
		return t.parseTime(string(v))
	case string:
		return t.parseTime(v)
	case time.Time:
		t.Time = v
		return nil
	default:
		return fmt.Errorf("cannot scan type %T into TimeOfDay", value)
	}
}

// Value implements the driver.Valuer interface
func (t TimeOfDay) Value() (driver.Value, error) {
	if t.Time.IsZero() {
		return nil, nil
	}
	return t.Time.Format("15:04:05"), nil
}

// parseTime parses various time formats
func (t *TimeOfDay) parseTime(s string) error {
	formats := []string{
		"15:04:05",
		"15:04",
		time.TimeOnly,
	}

	var err error
	for _, format := range formats {
		t.Time, err = time.Parse(format, s)
		if err == nil {
			return nil
		}
	}

	return fmt.Errorf("unable to parse time: %s", s)
}

// ToTime converts TimeOfDay to time.Time (using today's date)
func (t TimeOfDay) ToTime() time.Time {
	return t.Time
}

// NewTimeOfDay creates a new TimeOfDay from time.Time
func NewTimeOfDay(t time.Time) TimeOfDay {
	return TimeOfDay{Time: t}
}
