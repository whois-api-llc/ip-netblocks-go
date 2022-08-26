package ipnetblocks

import (
	"net/url"
	"reflect"
	"testing"
)

// TestOptions tests the Options functions.
func TestOptions(t *testing.T) {
	from := "8.2.17.0-8.2.17.255"
	tests := []struct {
		name   string
		values url.Values
		option Option
		want   string
	}{
		{
			name:   "output format",
			values: url.Values{},
			option: OptionOutputFormat("JSON"),
			want:   "outputFormat=JSON",
		},
		{
			name:   "limit",
			values: url.Values{},
			option: OptionLimit(150),
			want:   "limit=150",
		},
		{
			name:   "from",
			values: url.Values{},
			option: OptionFrom(nil),
			want:   "",
		},
		{
			name:   "from2",
			values: url.Values{},
			option: OptionFrom(&from),
			want:   "from=8.2.17.0-8.2.17.255",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.option(tt.values)
			if got := tt.values.Encode(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Option() = %v, want %v", got, tt.want)
			}
		})
	}
}
