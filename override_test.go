package dotenv_test

import (
	"fmt"
	"math/rand"
	"os"
	"sync"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tangelo-labs/go-dotenv"
)

func TestWithOverride(t *testing.T) {
	t.Run("GIVEN an environment struct with default values", func(t *testing.T) {
		var vars testEnv

		require.NoError(t, dotenv.LoadAndParse(&vars))
		require.EqualValues(t, "bar", vars.Foo)
		require.EqualValues(t, 42, vars.TheMeaningOfLifeTheUniverseAndEverything)
		require.EqualValues(t, []string{"foo", "bar"}, vars.FakeList)

		t.Run("WHEN overriding the environment variables with another values", func(t *testing.T) {
			overrideValue := gofakeit.UUID()
			funcWasCalled := false

			dotenv.WithOverride(func() {
				funcWasCalled = true

				require.NoError(t, dotenv.LoadAndParse(&vars))
				require.EqualValues(t, overrideValue, vars.Foo)
				require.EqualValues(t, 24, vars.TheMeaningOfLifeTheUniverseAndEverything)
				require.EqualValues(t, []string{"1", "2", "3"}, vars.FakeList)
			},
				"FOO", overrideValue,
				"DUMMY", "24",
				"FAKE_LIST", "1,2,3",
			)

			t.Run("THEN callbacks was invoked AND new values were seen AND original values are restored", func(t *testing.T) {
				require.True(t, funcWasCalled)

				require.NoError(t, dotenv.LoadAndParse(&vars))
				require.EqualValues(t, "bar", vars.Foo)
				require.EqualValues(t, 42, vars.TheMeaningOfLifeTheUniverseAndEverything)
				require.EqualValues(t, []string{"foo", "bar"}, vars.FakeList)
			})
		})
	})

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

			t.Run("THEN callbacks were invoked AND new values were seen AND original values are restored", func(t *testing.T) {
				require.True(t, funcOneWasCalled)
				require.True(t, funcTwoWasCalled)

				require.NoError(t, dotenv.LoadAndParse(&vars))
				require.EqualValues(t, "bar", vars.Foo)
				require.EqualValues(t, 42, vars.TheMeaningOfLifeTheUniverseAndEverything)
			})
		})
	})

	t.Run("GIVEN an environment variable", func(t *testing.T) {
		require.NoError(t, os.Setenv("GIT_GUT", "lol"))

		t.Run("WHEN multiple goroutines are overriding THEN each goroutine should be its value", func(t *testing.T) {
			muxAssert := sync.Mutex{}
			wg := sync.WaitGroup{}

			for i := 0; i < 100; i++ {

			}

			for i := 0; i < 100; i++ {
				wg.Add(1)
				go func(idx int) {
					val := fmt.Sprintf("%d", rand.New(rand.NewSource(int64(idx))).Int63())
					fmt.Printf("val %d: %s\n", idx, val)

					dotenv.WithOverride(func() {
						defer wg.Done()
						var env loneVarTestEnv
						err := dotenv.LoadAndParse(&env)

						muxAssert.Lock()
						assert.NoErrorf(t, err, "error on idx %d", idx)
						assert.EqualValuesf(t, val, env.YRURunning, "assert fail on idx %d", idx)
						muxAssert.Unlock()
					}, "GIT_GUT", val)
				}(i)
			}

			wg.Wait()
		})
	})
}

type testEnv struct {
	Foo                                      string   `env:"FOO" default:"bar"`
	TheMeaningOfLifeTheUniverseAndEverything int      `env:"DUMMY" default:"42"`
	FakeList                                 []string `env:"FAKE_LIST" default:"foo,bar" delimiter:","`
}

type loneVarTestEnv struct {
	YRURunning string `env:"GIT_GUT"`
}
