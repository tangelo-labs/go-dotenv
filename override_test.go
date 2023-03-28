package dotenv_test

import (
	"os"
	"sync"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
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

		t.Run("WHEN multiple goroutines are overriding the same variable THEN each goroutine should see its own overridden value", func(t *testing.T) {
			wg := sync.WaitGroup{}
			ng := 100

			faker := gofakeit.New(time.Now().Unix())
			ready := make(chan struct{})

			for i := 0; i < ng; i++ {
				wg.Add(1)

				go func(idx int) {
					defer wg.Done()

					value := faker.LoremIpsumSentence(5)

					dotenv.WithOverride(func() {
						var env testEnv

						<-ready

						if err := dotenv.LoadAndParse(&env); err != nil {
							t.Errorf("failed to load env: %s", err)
						}

						if env.ConcurrentString != value {
							t.Errorf("expected %s, got %s", value, env.ConcurrentString)
						}
					}, "CONCURRENT_STRING", value)
				}(i)
			}

			close(ready)
			wg.Wait()
		})
	})
}

type testEnv struct {
	Foo                                      string   `env:"FOO" default:"bar"`
	TheMeaningOfLifeTheUniverseAndEverything int      `env:"DUMMY" default:"42"`
	FakeList                                 []string `env:"FAKE_LIST" default:"foo,bar" delimiter:","`
	ConcurrentString                         string   `env:"CONCURRENT_STRING" default:"foo"`
}
