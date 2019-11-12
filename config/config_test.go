package config

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestConfig(t *testing.T) {
	config, e := LoadConfig()
	require.Nil(t, e)
	require.NotNil(t, config)
}
