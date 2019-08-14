package xpgtypes

import (
	"database/sql/driver"
	"time"
)

type NullTime struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}

// Scan implements the Scanner interface.
func (nt *NullTime) Scan(value interface{}) error {
	nt.Time, nt.Valid = value.(time.Time)
	return nil
}

// Value implements the driver Valuer interface.
func (nt NullTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}

func (nt NullTime) MarshalJSON() ([]byte, error) {
	if nt.Valid {
		st := nt.Time.Format("2006-01-02 15:04:05")
		if st == "0001-01-01 00:00:00" {
			return []byte("null"), nil
		}
		return []byte(`"` + st + `"`), nil
	}
	return []byte("null"), nil
}

func (nt NullTime) UnmarshalJSON(b []byte) error {
	t, err := time.Parse("2006-01-02 15:04:05", string(b))
	nt.Time = t
	nt.Valid = err == nil
	return nil
}
