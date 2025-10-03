package app

import (
	"testing"

	"github.com/vitalykrupin/url-shortener/cmd/shortener/config"
)

func TestNewApp(t *testing.T) {
	// Verify that the NewApp function exists
	// We can't fully test it without proper mocks for the storage and delete service
	// but we can at least verify the function exists
	_ = NewApp

	// Create a config
	conf := config.NewConfig()

	// Verify the config was created correctly
	if conf == nil {
		t.Error("Expected config to be created, but got nil")
	}
}
