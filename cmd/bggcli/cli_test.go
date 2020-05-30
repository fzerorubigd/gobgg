package main

import (
	"context"
	"flag"
	"io"
	"testing"

	"github.com/fzerorubigd/gobgg"
	"github.com/stretchr/testify/require"
)

type devNull struct{}

func (devNull) Write(in []byte) (int, error) {
	return len(in), nil
}

func TestDispatch(t *testing.T) {
	bgg := gobgg.NewBGGClient()
	var (
		err  error
		args []string
	)

	commands = nil
	addCommand("test", "test command", func(ctx context.Context, funcBGG *gobgg.BGG, funcArgs ...string) error {
		require.True(t, len(funcArgs) > 1)
		require.Equal(t, "test", args[0])
		require.Equal(t, bgg, funcBGG)
		require.Equal(t, args, funcArgs)
		return err
	})

	ctx := context.Background()
	flag.CommandLine.SetOutput(&devNull{})
	args = []string{"test", "arg1", "arg2"}
	require.NoError(t, dispatch(ctx, bgg, args...))
	require.Error(t, dispatch(ctx, bgg, "invalid", "arg"))
	require.Error(t, dispatch(ctx, bgg))

	err = io.EOF
	require.Error(t, dispatch(ctx, bgg, args...))
}
