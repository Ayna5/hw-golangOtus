package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("should ok", func(t *testing.T) {
		cmd := []string{"testdata/echo.sh", "arg1", "arg2"}
		env := make(Environment)
		env["FOO"] = EnvValue{Value: "foo"}
		env["BAR"] = EnvValue{Value: "bar"}
		env["ADDED"] = EnvValue{Value: "added"}

		code := RunCmd(cmd, env)

		require.Equal(t, 0, code)
	})

	t.Run("cmd is nil", func(t *testing.T) {
		code := RunCmd(nil, nil)

		require.Equal(t, 1, code)
	})
}
