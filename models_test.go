package ipnetblocks

import (
	"encoding/json"
	"testing"
)

// TestTime tests JSON encoding/parsing functions for the time values
func TestTime(t *testing.T) {
	tests := []struct {
		name   string
		decErr string
		encErr string
	}{
		{
			name:   `"2006-01-02 15:04:05 EST"`,
			decErr: "parsing time \"2006-01-02 15:04:05 EST\" as \"2006-01-02T15:04:05Z07:00\": cannot parse \" 15:04:05 EST\" as \"T\"",
			encErr: "",
		},
		{
			name:   `"2014-03-14T00:00:00Z"`,
			decErr: "",
			encErr: "",
		},
		{
			name:   `""`,
			decErr: "",
			encErr: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var v Time

			err := json.Unmarshal([]byte(tt.name), &v)
			checkErr(t, err, tt.decErr)
			if tt.decErr != "" {
				return
			}

			bb, err := json.Marshal(v)
			checkErr(t, err, tt.encErr)
			if tt.encErr != "" {
				return
			}

			if string(bb) != tt.name {
				t.Errorf("got = %v, want %v", string(bb), tt.name)
			}
		})
	}
}

// checkErr checks for an error.
func checkErr(t *testing.T, err error, want string) {
	if (err != nil || want != "") && (err == nil || err.Error() != want) {
		t.Errorf("error = %v, wantErr %v", err, want)
	}
}
