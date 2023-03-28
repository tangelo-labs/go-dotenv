package dotenv_test

import (
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/tangelo-labs/go-dotenv"
)

func Test_Override1(t *testing.T) {
	t.Parallel()

	testOverride(t)
}

func Test_Override2(t *testing.T) {
	t.Parallel()

	testOverride(t)
}

func Test_Override3(t *testing.T) {
	t.Parallel()

	testOverride(t)
}

func Test_Override4(t *testing.T) {
	t.Parallel()

	testOverride(t)
}

func Test_Override5(t *testing.T) {
	t.Parallel()

	testOverride(t)
}

func Test_Override6(t *testing.T) {
	t.Parallel()

	testOverride(t)
}

func Test_Override7(t *testing.T) {
	t.Parallel()

	testOverride(t)
}

func Test_Override8(t *testing.T) {
	t.Parallel()

	testOverride(t)
}

func Test_Override9(t *testing.T) {
	t.Parallel()

	testOverride(t)
}

func Test_Override10(t *testing.T) {
	t.Parallel()

	testOverride(t)
}

func Test_Override11(t *testing.T) {
	t.Parallel()

	testOverride(t)
}

func Test_Override12(t *testing.T) {
	t.Parallel()

	testOverride(t)
}

func Test_Override13(t *testing.T) {
	t.Parallel()

	testOverride(t)
}

func Test_Override14(t *testing.T) {
	t.Parallel()

	testOverride(t)
}

func Test_Override15(t *testing.T) {
	t.Parallel()

	testOverride(t)
}

func Test_Override16(t *testing.T) {
	t.Parallel()

	testOverride(t)
}

func Test_Override17(t *testing.T) {
	t.Parallel()

	testOverride(t)
}

func Test_Override18(t *testing.T) {
	t.Parallel()

	testOverride(t)
}

func Test_Override19(t *testing.T) {
	t.Parallel()

	testOverride(t)
}

func Test_Override20(t *testing.T) {
	t.Parallel()

	testOverride(t)
}

func testOverride(t *testing.T) {
	gofakeit := gofakeit.New(time.Now().Unix())
	val := gofakeit.Sentence(10)

	dotenv.WithOverride(func() {
		var env tEnv

		if err := dotenv.LoadAndParse(&env); err != nil {
			t.Errorf("failed to load env: %s", err)
		}

		if env.ConcurrentValue != val {
			t.Errorf("expected %s, got %s", val, env.ConcurrentValue)
		}
	}, "CONCURRENT_STRING", val)
}

type tEnv struct {
	ConcurrentValue string `env:"CONCURRENT_STRING" default:"bar"`
}
