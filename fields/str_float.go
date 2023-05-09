// Package fields
package fields

import (
	"strconv"

	"github.com/pashifika/util/conv"
)

type StrFloat float32

func (s StrFloat) Value() float32 { return float32(s) }

// MarshalJSON returns the encoded JSON string.
func (s StrFloat) MarshalJSON() ([]byte, error) {
	return conv.StringToBytes(strconv.FormatFloat(float64(s), 'g', -1, 32)), nil
}

// UnmarshalJSON sets the value that decoded JSON.
func (s *StrFloat) UnmarshalJSON(data []byte) (err error) {
	str := conv.BytesToString(data)
	v, err := strconv.ParseFloat(str, 32)
	*s = StrFloat(v)
	return err
}

type StrFloat64 float64

func (s StrFloat64) Value() float64 { return float64(s) }

// MarshalJSON returns the encoded JSON string.
func (s StrFloat64) MarshalJSON() ([]byte, error) {
	return conv.StringToBytes(strconv.FormatFloat(float64(s), 'g', -1, 64)), nil
}

// UnmarshalJSON sets the value that decoded JSON.
func (s *StrFloat64) UnmarshalJSON(data []byte) (err error) {
	str := conv.BytesToString(data)
	v, err := strconv.ParseFloat(str, 64)
	*s = StrFloat64(v)
	return err
}
