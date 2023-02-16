package dotenv_test

import (
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
	"github.com/tangelolabs/dotenv"
)

func TestWithOverride(t *testing.T) {
	var vars testEnv

	require.NoError(t, dotenv.LoadAndParse(&vars))
	require.EqualValues(t, "bar", vars.Foo)
	require.EqualValues(t, 42, vars.TheMeaningOfLifeTheUniverseAndEverything)

	overrideValue := gofakeit.UUID()

	dotenv.WithOverride(func() {
		require.NoError(t, dotenv.LoadAndParse(&vars))
		require.EqualValues(t, overrideValue, vars.Foo)
		require.EqualValues(t, 24, vars.TheMeaningOfLifeTheUniverseAndEverything)
	},
		"FOO", overrideValue,
		"DUMMY", "24",
	)
}

type testEnv struct {
	Foo                                      string `env:"FOO" default:"bar"`
	TheMeaningOfLifeTheUniverseAndEverything int    `env:"DUMMY" default:"42"`
}
