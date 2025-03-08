package models

import (
	"snippetbox.alexedwards.net/internal/assert"
	"testing"
)

func TestUserModelExists(t *testing.T) {

	// Skip the test if the -short flag is provided when running the test suite.
	if testing.Short() {
		t.Skip("models: skipping integration test")
	}

	tests := []struct {
		name   string
		userID int
		want   bool
	}{
		{"Valid ID", 1, true},
		{"Zero ID", 0, false},
		{"Non-existent ID", 2, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := newTestDB(t)
			m := UserModel{db}
			exists, err := m.Exists(tt.userID)
			assert.NilError(t, err)
			assert.Equal(t, exists, tt.want)
		})
	}

}
