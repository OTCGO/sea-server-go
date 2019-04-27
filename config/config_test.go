package config

import (
	"github.com/hzxiao/goutil/assert"
	"testing"
)

func TestInit(t *testing.T) {
	err := Init("../seago-test.toml")
	assert.NoError(t, err)

	t.Log(*Conf)
}