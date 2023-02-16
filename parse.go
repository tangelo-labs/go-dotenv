package dotenv

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/fatih/structtag"
)

var (
	timeType        = reflect.TypeOf(time.Time{})
	durationType    = reflect.TypeOf(time.Duration(0))
	stringSliceType = reflect.TypeOf([]string{})
)

// List of injection errors.
var (
	ErrNotAPointer        = errors.New("not a pointer")
	ErrTimeLayoutRequired = errors.New("missing timeLayout tag")
	ErrRequiredField      = errors.New("required field")
	ErrEmptyField         = errors.New("empty field")
)

var valueMapper = map[reflect.Kind]func(v value) interface{}{
	reflect.Int: func(v value) interface{} {
		return v.AsInt()
	},
	reflect.Int8: func(v value) interface{} {
		return v.AsInt8()
	},
	reflect.Int16: func(v value) interface{} {
		return v.AsInt16()
	},
	reflect.Int32: func(v value) interface{} {
		return v.AsInt32()
	},
	reflect.Int64: func(v value) interface{} {
		return v.AsInt64()
	},
	reflect.Uint: func(v value) interface{} {
		return v.AsUint()
	},
	reflect.Uint8: func(v value) interface{} {
		return v.AsUint8()
	},
	reflect.Uint16: func(v value) interface{} {
		return v.AsUint16()
	},
	reflect.Uint32: func(v value) interface{} {
		return v.AsUint32()
	},
	reflect.Uint64: func(v value) interface{} {
		return v.AsUint64()
	},
	reflect.Float32: func(v value) interface{} {
		return v.AsFloat32()
	},
	reflect.Float64: func(v value) interface{} {
		return v.AsFloat64()
	},
	reflect.String: func(v value) interface{} {
		return v.AsString()
	},
	reflect.Bool: func(v value) interface{} {
		return v.AsBool()
	},
}

// Parse injects environment variables into the given struct using tag
// annotations.
//
// The given argument must be a pointer to a struct where values will be
// injected.
//
// Struct must tag its fields with `env:"VAR_NAME"` to specify the environment
// variable value to be injected.
//
// Fields may be marked as required using the `required` env option. If a field
// is required and no environment variable is found, an error will be returned.
//
// However, if you want to make sure that a field is not empty, you can use the
// `notEmpty` option. In which case an error will be returned if not value is
// found.
//
// Time fields:
//
// Optionally, the tag `default` may be used to specify a default value for the
// field. In the case of `time.Time` fields, the `timeLayout` tag must be used
// to specify the format of the time string.
//
// String slices:
//
// Optionally, the tag `delimiter` may be used to specify a separator for string
// slices, by default `,` will be used.
//
// For example:
//
//	type Config struct {
//		Foo  string		`env:"ENV_FOO,required" default:"fooValue"`
//		Bar  int		`env:"ENV_BAR,notEmpty"`
//		IPs  []string	`env:"ENV_IPS" delimiter:";"`
//		When time.Time	`env:"ENV_WHEN" default:"2021-12-24T17:04:05Z07:00" timeLayout:"2006-01-02T15:04:05Z07:00"`
//	}
//
// Fields without an `env` tag will not be injected.
func Parse(st interface{}) error {
	if err := Load(); err != nil {
		return err
	}

	val := reflect.ValueOf(st)
	if val.Kind() != reflect.Ptr {
		return fmt.Errorf("%w: given `%s` is not a pointer", ErrNotAPointer, val.Kind())
	}

	val = val.Elem()
	typ := val.Type()

	if val.Kind() != reflect.Struct {
		return fmt.Errorf("%w: given `%s` is not a pointer", ErrNotAPointer, val.Kind())
	}

	for idx := 0; idx < val.NumField(); idx++ {
		field := val.Field(idx)
		tag := typ.Field(idx).Tag

		tags, err := structtag.Parse(string(tag))
		if err != nil {
			return err
		}

		envTag, err := tags.Get("env")
		if err != nil {
			// skip not tagged fields
			continue
		}

		defaultValue := ""
		defaultTag, err := tags.Get("default")

		if err == nil {
			defaultValue = defaultTag.Name
		}

		isRequired := false
		notEmpty := false

		for i := range envTag.Options {
			if envTag.Options[i] == "required" {
				isRequired = true
			}

			if envTag.Options[i] == "notEmpty" {
				notEmpty = true
			}
		}

		v, defined := lookup(envTag.Name, defaultValue)

		if isRequired && !defined {
			return fmt.Errorf("%w: environment variable `%s` must be defined", ErrRequiredField, envTag.Name)
		}

		if notEmpty && v.IsZero() {
			return fmt.Errorf("%w: environment variable `%s` cannot be empty", ErrEmptyField, envTag.Name)
		}

		writeValue, err := valueForField(field, v, tags, envTag.Name)
		if err != nil {
			return err
		}

		field.Set(reflect.ValueOf(writeValue))
	}

	return nil
}

// LoadAndParse convenience function which first loads environment variables
// and then injects them into the given struct.
//
// This simplifies the following use case:
//
//	if err := dotenv.Load(); err != nil {
//		panic(err)
//	}
//
//	env := Config{}
//
//	if err := dotenv.Parse(&env); err != nil {
//		panic(err)
//	}
//
// The above code is equivalent to:
//
//	env := Config{}
//
//	if err := dotenv.LoadAndParse(); err != nil {
//		panic(err)
//	}
//
// See Parse function for more information.
func LoadAndParse(st interface{}) error {
	if err := Load(); err != nil {
		return err
	}

	return Parse(st)
}

func valueForField(field reflect.Value, value value, tags *structtag.Tags, varName string) (interface{}, error) {
	fieldType := field.Type()

	if fieldType.AssignableTo(timeType) {
		timeLayoutTag, gErr := tags.Get("timeLayout")
		if gErr != nil {
			return nil, fmt.Errorf("%w: expecting tag `timeLayout` for environment variables of type `time.Time`", ErrTimeLayoutRequired)
		}

		return value.AsTime(timeLayoutTag.Name), nil
	}

	if fieldType.AssignableTo(stringSliceType) {
		delimiter := ","
		delimiterTag, gErr := tags.Get("delimiter")

		if gErr == nil {
			delimiter = delimiterTag.Name
		}

		return value.AsStringSlice(delimiter), nil
	}

	if fieldType.AssignableTo(durationType) {
		return value.AsDuration(), nil
	}

	t, ok := valueMapper[field.Kind()]
	if !ok {
		return nil, fmt.Errorf("unsupported environment data type `%s` for variable `%s`", fieldType.Kind(), varName)
	}

	return t(value), nil
}

// lookup similar to Get but returns whether the variable is present or not.
func lookup(name string, def ...string) (value, bool) {
	d := ""
	if len(def) > 0 {
		d = def[0]
	}

	v := value(d)
	val, defined := os.LookupEnv(name)

	if defined {
		v = value(val)
	}

	return v, defined
}
