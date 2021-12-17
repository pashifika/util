// Package conv
/*
 * Version: 1.0.0
 * Copyright (c) 2021. Pashifika
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
package conv

import (
	"unicode/utf8"
)

// CutUnicodeString is like RuneCount but its input is a string.
func CutUnicodeString(str string, length int) string {
	if str == "" || length <= 0 || utf8.RuneCountInString(str) < length {
		return ""
	}
	var (
		res   string
		count int
	)
	for len(str) > 0 {
		if count >= length {
			break
		}
		r, size := utf8.DecodeRuneInString(str)
		str = str[size:]
		res += string(r)
		count++
	}
	return res
}

// FindUnicodeString is use rune to find the string.
func FindUnicodeString(src, find string) bool {
	s, e := _findUnicodeString(src, find)
	if s != -1 && e != -1 {
		return true
	}
	return false
}

func _findUnicodeString(src, find string) (start, end int) {
	rf := []rune(find)
	rfl := len(rf) - 1
	var idx int
	start = -1
	end = -1
	for i, r := range src {
		if r == rf[idx] {
			if idx == 0 {
				start = i
			} else if idx == rfl {
				end = i
			}
			if start != -1 && end != -1 {
				break
			}
			idx++
		} else {
			idx = 0
		}
	}
	return
}
