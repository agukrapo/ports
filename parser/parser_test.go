package parser

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	okInput = `{
  "key": {
    "name": "name",
    "city": "city",
    "province": "province",
    "country": "country",
    "alias": ["one", "two"],
    "coordinates": [
      1.111,
      2.22
    ],
    "timezone": "timezone",
    "unlocs": [
      "1", "2"
    ],
    "code": "code"
  }`
)

func TestParser(t *testing.T) {
	tests := []struct {
		name string
		r    io.Reader
		more bool
		next *Port
	}{
		{
			name: "empty object",
			r:    strings.NewReader("{}"),
		},
		{
			name: "ok",
			r:    strings.NewReader(okInput),
			more: true,
			next: &Port{
				Key:         "key",
				Timezone:    "timezone",
				Coordinates: []float64{1.111, 2.22},
				Name:        "name",
				City:        "city",
				Province:    "province",
				Country:     "country",
				Alias:       []string{"one", "two"},
				Unlocs:      []string{"1", "2"},
				Code:        "code",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := New(tt.r)
			require.NoError(t, err)

			more := p.More()
			require.Equal(t, tt.more, more)

			if more {
				var next Port
				require.NoError(t, p.Next(&next))
				require.Equal(t, tt.next, &next)
			}
		})
	}
}
