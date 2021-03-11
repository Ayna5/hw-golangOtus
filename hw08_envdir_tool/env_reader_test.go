package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("test is valid", func(t *testing.T) {
		expectedEnv := make(Environment)
		expectedEnv["BAR"] = EnvValue{Value: "bar"}
		expectedEnv["FOO"] = EnvValue{Value: "   foo\nwith new line"}
		expectedEnv["HELLO"] = EnvValue{Value: "\"hello\""}
		expectedEnv["UNSET"] = EnvValue{Value: ""}

		env, err := ReadDir("testdata/env")

		require.NoError(t, err)
		require.Equal(t, expectedEnv, env)
	})
}
