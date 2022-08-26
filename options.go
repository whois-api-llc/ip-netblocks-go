package ipnetblocks

import (
	"net/url"
	"strconv"
	"strings"
)

// Option adds parameters to the query.
type Option func(v url.Values)

var _ = []Option{
	OptionOutputFormat("JSON"),
	OptionLimit(100),
	OptionFrom(nil),
}

// OptionOutputFormat sets Response output format JSON | XML. Default: JSON.
func OptionOutputFormat(outputFormat string) Option {
	return func(v url.Values) {
		v.Set("outputFormat", strings.ToUpper(outputFormat))
	}
}

// OptionLimit sets max count of returned records. Acceptable values: 1 - 1000. Default: 100.
func OptionLimit(value int) Option {
	return func(v url.Values) {
		v.Set("limit", strconv.Itoa(value))
	}
}

// OptionFrom sets the IP netblock range that is used as an offset for the returned results.
func OptionFrom(value *string) Option {
	return func(v url.Values) {
		if value != nil {
			v.Set("from", strings.ToUpper(*value))
		}
	}
}
