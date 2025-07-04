package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringArray_Value(t *testing.T) {
	tests := []struct {
		name     string
		array    StringArray
		expected string
	}{
		{
			name:     "empty array",
			array:    StringArray{},
			expected: "[]",
		},
		{
			name:     "single item",
			array:    StringArray{"item1"},
			expected: `["item1"]`,
		},
		{
			name:     "multiple items",
			array:    StringArray{"item1", "item2", "item3"},
			expected: `["item1","item2","item3"]`,
		},
		{
			name:     "nil array",
			array:    nil,
			expected: "null",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, err := tt.array.Value()
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, string(value.([]byte)))
		})
	}
}

func TestStringArray_Scan(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		expected    StringArray
		expectError bool
	}{
		{
			name:     "nil value",
			input:    nil,
			expected: nil,
		},
		{
			name:     "empty array",
			input:    []byte("[]"),
			expected: StringArray{},
		},
		{
			name:     "single item",
			input:    []byte(`["item1"]`),
			expected: StringArray{"item1"},
		},
		{
			name:     "multiple items",
			input:    []byte(`["item1","item2","item3"]`),
			expected: StringArray{"item1", "item2", "item3"},
		},
		{
			name:        "invalid JSON",
			input:       []byte(`invalid json`),
			expectError: true,
		},
		{
			name:        "non-byte input",
			input:       "string",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var array StringArray
			err := array.Scan(tt.input)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, array)
			}
		})
	}
}