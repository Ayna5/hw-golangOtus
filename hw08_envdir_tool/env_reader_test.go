package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("test is valid", func(t *testing.T) {
		expectedEnv := make(Environment)
		expectedEnv["ADDED"] = EnvValue{Value: "from original env", NeedRemove: false}
		expectedEnv["BAR"] = EnvValue{Value: "bar", NeedRemove: false}
		expectedEnv["EMPTY"] = EnvValue{Value: "", NeedRemove: true}
		expectedEnv["FOO"] = EnvValue{Value: "   foo\nwith new line", NeedRemove: false}
		expectedEnv["HELLO"] = EnvValue{Value: "\"hello\"", NeedRemove: false}
		expectedEnv["UNSET"] = EnvValue{Value: "", NeedRemove: true}

		env, err := ReadDir("testdata/env")

		fmt.Println(env)

		require.NoError(t, err)
		require.Equal(t, expectedEnv, env)
	})
}
