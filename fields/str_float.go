// Package fields
/*
 * Version: 1.0.0
 * Copyright (c) 2023. Pashifika
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package fields

import (
	"strconv"
	"strings"

	"github.com/pashifika/util/conv"
)

type StrFloat float32

func (s StrFloat) Value() float32 { return float32(s) }

// MarshalJSON returns the encoded JSON string.
func (s StrFloat) MarshalJSON() ([]byte, error) {
	str := strconv.FormatFloat(float64(s), 'g', -1, 32)
	str = _jsonChar + str + _jsonChar
	return conv.StringToBytes(str), nil
}

// UnmarshalJSON sets the value that decoded JSON.
func (s *StrFloat) UnmarshalJSON(data []byte) (err error) {
	str := conv.BytesToString(data)
	str = strings.TrimPrefix(strings.TrimSuffix(str, _jsonChar), _jsonChar)
	v, err := strconv.ParseFloat(str, 32)
	*s = StrFloat(v)
	return err
}

type StrFloat64 float64

func (s StrFloat64) Value() float64 { return float64(s) }

// MarshalJSON returns the encoded JSON string.
func (s StrFloat64) MarshalJSON() ([]byte, error) {
	str := strconv.FormatFloat(float64(s), 'g', -1, 64)
	str = _jsonChar + str + _jsonChar
	return conv.StringToBytes(str), nil
}

// UnmarshalJSON sets the value that decoded JSON.
func (s *StrFloat64) UnmarshalJSON(data []byte) (err error) {
	str := conv.BytesToString(data)
	str = strings.TrimPrefix(strings.TrimSuffix(str, _jsonChar), _jsonChar)
	v, err := strconv.ParseFloat(str, 64)
	*s = StrFloat64(v)
	return err
}
