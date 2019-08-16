package xpgtypes

import (
	"encoding/json"
	"testing"
)

func TestNullTime_UnmarshalJSON(t *testing.T) {
	var data = struct {
		Time NullTime `json:"time"`
	}{}

	if err := json.Unmarshal([]byte(`{"time":"2019-08-16 09:45:36"}`), &data); err != nil {
		t.Fatal(err)
	}

	if data.Time.Valid == false {
		t.Logf("%#v", data.Time)
		t.Error("Time not valid")
	}
}
