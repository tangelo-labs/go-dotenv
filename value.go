package dotenv

import (
	"strconv"
	"strings"
	"time"

	"github.com/xhit/go-str2duration"
)

// value represents a raw value which can be converted to a specific type.
type value string

// IsZero returns true if the value is empty.
func (v value) IsZero() bool {
	return strings.TrimSpace(string(v)) == ""
}

// AsInt cast this value to int type.
func (v value) AsInt() int {
	i, err := strconv.Atoi(string(v))
	if err != nil {
		return 0
	}

	return i
}

// AsInt8 cast this value to int8 type.
func (v value) AsInt8() int8 {
	i, err := strconv.ParseInt(string(v), 10, 8)
	if err != nil {
		return 0
	}

	return int8(i)
}

// AsInt16 cast this value to int16 type.
func (v value) AsInt16() int16 {
	i, err := strconv.ParseInt(string(v), 10, 16)
	if err != nil {
		return 0
	}

	return int16(i)
}

// AsInt32 cast this value to int32 type.
func (v value) AsInt32() int32 {
	i, err := strconv.ParseInt(string(v), 10, 32)
	if err != nil {
		return 0
	}

	return int32(i)
}

// AsInt64 cast this value to int64 type.
func (v value) AsInt64() int64 {
	i, err := strconv.ParseInt(string(v), 10, 64)
	if err != nil {
		return 0
	}

	return i
}

// AsUint cast this value to uint type.
func (v value) AsUint() uint {
	i, err := strconv.ParseUint(string(v), 10, 64)
	if err != nil {
		return 0
	}

	return uint(i)
}

// AsUint8 cast this value to uint8 type.
func (v value) AsUint8() uint8 {
	i, err := strconv.ParseUint(string(v), 10, 8)
	if err != nil {
		return 0
	}

	return uint8(i)
}

// AsUint16 cast this value to uint16 type.
func (v value) AsUint16() uint16 {
	i, err := strconv.ParseUint(string(v), 10, 16)
	if err != nil {
		return 0
	}

	return uint16(i)
}

// AsUint32 cast this value to uint32 type.
func (v value) AsUint32() uint32 {
	i, err := strconv.ParseUint(string(v), 10, 32)
	if err != nil {
		return 0
	}

	return uint32(i)
}

// AsUint64 cast this value to uint64 type.
func (v value) AsUint64() uint64 {
	i, err := strconv.ParseUint(string(v), 10, 64)
	if err != nil {
		return 0
	}

	return i
}

// AsFloat32 cast this value float32 type.
func (v value) AsFloat32() float32 {
	f, err := strconv.ParseFloat(string(v), 32)
	if err != nil {
		return 0
	}

	return float32(f)
}

// AsFloat64 cast this value to float64 type.
func (v value) AsFloat64() float64 {
	f, err := strconv.ParseFloat(string(v), 64)
	if err != nil {
		return 0
	}

	return f
}

// AsString cast this value to string type.
func (v value) AsString() string {
	return string(v)
}

// AsBool cast this value to bool type.
func (v value) AsBool() bool {
	b, err := strconv.ParseBool(string(v))
	if err != nil {
		return false
	}

	return b
}

// AsTime cast this value to time.Time type using the given format layout.
func (v value) AsTime(layout string) time.Time {
	t, err := time.Parse(layout, string(v))
	if err != nil {
		return time.Time{}
	}

	return t
}

// AsDuration cast this value to time.Duration type using as input values in human-readable format, such as:
// "30m", "1h30m", "2d", "1w2d12h30m5s", etc.
func (v value) AsDuration() time.Duration {
	d, err := str2duration.Str2Duration(string(v))
	if err != nil {
		return 0
	}

	return d
}

// AsStringSlice cast this value to []string type.
func (v value) AsStringSlice(delimiter string) []string {
	if v.IsZero() {
		return []string{}
	}

	return strings.Split(string(v), delimiter)
}
