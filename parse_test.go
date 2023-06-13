package dotenv_test

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
	"github.com/tangelo-labs/go-dotenv"
)

func TestParse(t *testing.T) {
	t.Run("GIVEN a env struct with string, int8, time and duration fields AND env variables declared for such fields", func(t *testing.T) {
		tv := "2021-12-24T17:04:05"
		tl := "2006-01-02T15:04:05"
		pt, err := time.Parse(tl, tv)
		require.NoError(t, err)

		expectedEnv := dummyStruct{
			String: gofakeit.LoremIpsumWord(),
			StringList: []string{
				gofakeit.UUID(),
				gofakeit.UUID(),
				gofakeit.UUID(),
			},
			Int8:     gofakeit.Int8(),
			Time:     pt,
			Duration: time.Second,
		}

		require.NoError(t, os.Setenv("TEST_STRING", expectedEnv.String))
		require.NoError(t, os.Setenv("TEST_STRING_LIST", strings.Join(expectedEnv.StringList, ";")))
		require.NoError(t, os.Setenv("TEST_INT8", fmt.Sprintf("%d", expectedEnv.Int8)))
		require.NoError(t, os.Setenv("TEST_TIME", tv))
		require.NoError(t, os.Setenv("TEST_DURATION", "1s"))

		t.Run("WHEN parsing the env variables into the struct THEN struct is filled with expected values", func(t *testing.T) {
			env := dummyStruct{}
			require.NoError(t, dotenv.Parse(&env))
			require.EqualValues(t, expectedEnv, env)
		})
	})

	t.Run("GIVEN a struct with required variables", func(t *testing.T) {
		env := dummyStructWithRequire{}

		t.Run("WHEN parsing passing by value THEN an error is raised", func(t *testing.T) {
			err := dotenv.Parse(env)
			require.ErrorIs(t, err, dotenv.ErrNotAPointer)
		})
	})

	t.Run("GIVEN a struct with time.Time variable but no template defined", func(t *testing.T) {
		env := dummyStructWithNoTimeLayout{}

		t.Run("WHEN parsing passing THEN an error is raised", func(t *testing.T) {
			err := dotenv.Parse(&env)
			require.ErrorIs(t, err, dotenv.ErrTimeLayoutRequired)
		})
	})

	t.Run("GIVEN a struct with required variables", func(t *testing.T) {
		env := dummyStructWithRequire{}

		t.Run("WHEN parsing and variable is not defined THEN an error is raised", func(t *testing.T) {
			err := dotenv.Parse(&env)
			require.ErrorIs(t, err, dotenv.ErrRequiredField)
		})
	})

	t.Run("GIVEN a struct with notEmpty variable and one variable defined but with no value", func(t *testing.T) {
		env := dummyStructNotEmpty{}
		require.NoError(t, os.Setenv("TEST_NOT_EMPTY", ""))

		t.Run("WHEN parsing THEN an error is raised", func(t *testing.T) {
			err := dotenv.Parse(&env)
			require.ErrorIs(t, err, dotenv.ErrEmptyField)
		})
	})

	t.Run("GIVEN a struct with a string slice and one variable defined with three elements and another with default 3 elements", func(t *testing.T) {
		env := dummyStringSlice{}
		require.NoError(t, os.Setenv("TEST_STRING_SLICE", "A;B;C"))

		t.Run("WHEN parsing THEN slices are populated", func(t *testing.T) {
			require.NoError(t, dotenv.Parse(&env))

			require.EqualValues(t, []string{"A", "B", "C"}, env.StringSlice)
			require.EqualValues(t, []string{"X", "Y", "Z"}, env.StringSliceWithDefaults)
		})
	})
}

func TestMustParse(t *testing.T) {
	t.Run("GIVEN a struct with notEmpty variable and one variable defined but with no value", func(t *testing.T) {
		env := dummyStructNotEmpty{}
		require.NoError(t, os.Setenv("TEST_NOT_EMPTY", ""))

		t.Run("WHEN using must-parseTHEN it panics", func(t *testing.T) {
			require.Panics(t, func() {
				dotenv.MustParse(&env)
			})
		})

		t.Run("WHEN using must-load-and-parse THEN it panics", func(t *testing.T) {
			require.Panics(t, func() {
				dotenv.MustLoadAndParse(&env)
			})
		})
	})
}

type dummyStruct struct {
	String     string        `env:"TEST_STRING"`
	StringList []string      `env:"TEST_STRING_LIST" delimiter:";"`
	Int8       int8          `env:"TEST_INT8"`
	Time       time.Time     `env:"TEST_TIME" timeLayout:"2006-01-02T15:04:05"`
	Duration   time.Duration `env:"TEST_DURATION"`
}

type dummyStructWithRequire struct {
	String string `env:"TEST_STRING_REQUIRED,required"`
}

type dummyStructWithNoTimeLayout struct {
	Time time.Time `env:"TEST_TIME"`
}

type dummyStructNotEmpty struct {
	String string `env:"TEST_NOT_EMPTY,notEmpty"`
}

type dummyStringSlice struct {
	StringSlice             []string `env:"TEST_STRING_SLICE" delimiter:";"`
	StringSliceWithDefaults []string `env:"TEST_STRING_SLICE_WITH_DEFAULTS" delimiter:";" default:"X;Y;Z"`
}
