package dotenv_test

import (
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
	"github.com/tangelolabs/dotenv"
)

func TestWithOverride(t *testing.T) {
	t.Run("GIVEN an environment struct with default values", func(t *testing.T) {
		var vars testEnv

		require.NoError(t, dotenv.LoadAndParse(&vars))
		require.EqualValues(t, "bar", vars.Foo)
		require.EqualValues(t, 42, vars.TheMeaningOfLifeTheUniverseAndEverything)

		t.Run("WHEN overriding the environment variables with another values in a nested fashion", func(t *testing.T) {
			overrideValue := gofakeit.UUID()
			funcOneWasCalled := false
			funcTwoWasCalled := false

			dotenv.WithOverride(func() {
				funcOneWasCalled = true

				require.NoError(t, dotenv.LoadAndParse(&vars))
				require.EqualValues(t, "testing", vars.Foo)
				require.EqualValues(t, 42, vars.TheMeaningOfLifeTheUniverseAndEverything)

				dotenv.WithOverride(func() {
					funcTwoWasCalled = true

					require.NoError(t, dotenv.LoadAndParse(&vars))
					require.EqualValues(t, overrideValue, vars.Foo)
					require.EqualValues(t, 24, vars.TheMeaningOfLifeTheUniverseAndEverything)
				},
					"FOO", overrideValue,
					"DUMMY", "24",
				)
			},
				"FOO", "testing",
			)

			t.Run("THEN callbacks were invoked and new values were seen AND original values are restored", func(t *testing.T) {
				require.True(t, funcOneWasCalled)
				require.True(t, funcTwoWasCalled)

				require.NoError(t, dotenv.LoadAndParse(&vars))
				require.EqualValues(t, "bar", vars.Foo)
				require.EqualValues(t, 42, vars.TheMeaningOfLifeTheUniverseAndEverything)
			})
		})
	})
}

type testEnv struct {
	Foo                                      string `env:"FOO" default:"bar"`
	TheMeaningOfLifeTheUniverseAndEverything int    `env:"DUMMY" default:"42"`
}
