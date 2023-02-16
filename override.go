package dotenv

import (
	"os"
	"reflect"
)

type empty struct{}

var withOverridePackagePath = reflect.TypeOf(empty{}).PkgPath() + ".WithOverride"

// WithOverride overrides the environment variables with the given map and
// restores them after the callback is executed.
//
// Any call to the Load, LoadAndParse or similar within the callback will be
// affected by the overridden values.
//
// This function will panic if the number of arguments is not even.
//
// Typical Usage Example:
//
//	dotenv.WithOverride(func() {
//	   functionThatCallsLoadAndParse()
//	}, "FOO", "bar")
func WithOverride(callback func(), kv ...string) {
	if len(kv)%2 != 0 {
		panic("dotenv: WithOverride requires an even number of arguments")
	}

	tuples := make(map[string]string, len(kv)/2)

	for i := 0; i < len(kv); i += 2 {
		k := kv[i]
		v := kv[i+1]

		tuples[k] = v
	}

	original := make(map[string]string)
	for k, v := range tuples {
		original[k] = os.Getenv(k)

		if err := os.Setenv(k, v); err != nil {
			panic(err)
		}
	}

	callback()

	for k, v := range original {
		if err := os.Setenv(k, v); err != nil {
			panic(err)
		}
	}
}
