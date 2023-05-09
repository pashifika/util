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

	"github.com/pashifika/util/conv"
)

type StrInt int

func (s StrInt) Value() int { return int(s) }

// MarshalJSON returns the encoded JSON string.
func (s StrInt) MarshalJSON() ([]byte, error) {
	return conv.StringToBytes(strconv.FormatInt(int64(s), 10)), nil
}

// UnmarshalJSON sets the value that decoded JSON.
func (s *StrInt) UnmarshalJSON(data []byte) (err error) {
	str := conv.BytesToString(data)
	v, err := strconv.ParseInt(str, 10, 32)
	*s = StrInt(v)
	return err
}

type StrInt64 int64

func (s StrInt64) Value() int64 { return int64(s) }

// MarshalJSON returns the encoded JSON string.
func (s StrInt64) MarshalJSON() ([]byte, error) {
	return conv.StringToBytes(strconv.FormatInt(int64(s), 10)), nil
}

// UnmarshalJSON sets the value that decoded JSON.
func (s *StrInt64) UnmarshalJSON(data []byte) (err error) {
	str := conv.BytesToString(data)
	v, err := strconv.ParseInt(str, 10, 64)
	*s = StrInt64(v)
	return err
}
