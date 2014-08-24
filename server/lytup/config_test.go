package lytup

import (
  "testing"
)

func TestLoadConfig(t *testing.T) {
  err := loadConfig()
  if err != nil {
    t.Error(err)
  }
}
