package xpgtypes

import (
	"database/sql/driver"
	"time"
)

var localZone, localOffset = time.Now().Local().Zone()

// NullTime Аналог sql.NullTime
type NullTime struct {
	Time  time.Time
	Error error
	Valid bool // Valid is true if Time is not NULL
}

// Scan implements the Scanner interface.
func (nt *NullTime) Scan(value interface{}) error {
	if t, ok := value.(time.Time); ok {
		nt.Valid = true
		if zone, offset := t.Zone(); localZone != zone {
			t = t.Local().Add(time.Duration(offset - localOffset) * time.Second)
		}
		nt.Time = t
	}
	return nil
}

// Value implements the driver Valuer interface.
func (nt NullTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nt.Error
	}
	return nt.Time, nil
}

// MarshalJSON Правила для упаковки в Json
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

// UnmarshalJSON Правила для распаковки из Json
func (nt *NullTime) UnmarshalJSON(b []byte) error {
	nt.Time, nt.Error = time.Parse(`"2006-01-02 15:04:05"`, string(b))
	nt.Valid = nt.Error == nil
	return nil
}
