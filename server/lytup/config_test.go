package lytup

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {
	if err := loadConfig(); err != nil {
		t.Error(err)
	}
}
