package config

import (
	"testing"

	"github.com/bmizerany/assert"
)

func TestNew(t *testing.T) {
	cfg, err := New("./../bin/")
	assert.Equal(t, nil, err)
	t.Log(cfg.Database)
	assert.Equal(t, true, cfg != nil)
}
