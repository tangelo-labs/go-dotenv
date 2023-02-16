package dotenv

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestValue_AsDuration(t *testing.T) {
	tests := []struct {
		raw      string
		expected time.Duration
	}{
		{
			raw:      "30m",
			expected: 30 * time.Minute,
		},
		{
			raw:      "1h",
			expected: 1 * time.Hour,
		},
		{
			raw:      "2h",
			expected: 2 * time.Hour,
		},
		{
			raw:      "1d",
			expected: 24 * time.Hour,
		},
		{
			raw:      "3d",
			expected: 3 * 24 * time.Hour,
		},
		{
			raw:      "1w2d2h30m",
			expected: 9*24*time.Hour + 2*time.Hour + 30*time.Minute,
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("GIVEN %s raw duration WHEN parsed as duration THEN should return %s", test.raw, test.expected.String()), func(t *testing.T) {
			v := value(test.raw)

			assert.Equal(t, test.expected, v.AsDuration())
		})
	}
}

func TestValue_AsStringSlice(t *testing.T) {
	tests := []struct {
		raw       string
		delimiter string
		expected  []string
	}{
		{
			raw:       "A,B,C",
			delimiter: ",",
			expected:  []string{"A", "B", "C"},
		},
		{
			raw:       "A B C",
			delimiter: " ",
			expected:  []string{"A", "B", "C"},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("GIVEN %s raw slice WHEN parsed as string slice THEN should return %+v", test.raw, test.expected), func(t *testing.T) {
			v := value(test.raw)

			assert.Equal(t, test.expected, v.AsStringSlice(test.delimiter))
		})
	}
}
