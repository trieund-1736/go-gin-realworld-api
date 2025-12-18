package utils

import (
	"go-gin-realworld-api/internal/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateSlug(t *testing.T) {
	tests := []struct {
		name     string
		title    string
		expected string
	}{
		{
			name:     "Simple title",
			title:    "Hello World",
			expected: "hello-world",
		},
		{
			name:     "Title with special characters",
			title:    "Hello, World! @2025",
			expected: "hello-world-2025",
		},
		{
			name:     "Title with multiple spaces",
			title:    "Hello    World",
			expected: "hello-world",
		},
		{
			name:     "Title with leading and trailing spaces",
			title:    "  Hello World  ",
			expected: "hello-world",
		},
		{
			name:     "Title with hyphens",
			title:    "Hello-World",
			expected: "hello-world",
		},
		{
			name:     "Title with mixed case",
			title:    "hElLo WoRlD",
			expected: "hello-world",
		},
		{
			name:     "Empty title",
			title:    "",
			expected: "",
		},
		{
			name:     "Title with only special characters",
			title:    "!@#$%^&*",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := utils.GenerateSlug(tt.title)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
