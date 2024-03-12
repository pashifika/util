// Package conv
/*
 * Version: 1.0.0
 * Copyright (c) 2024. Pashifika
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
	"strings"
	"unicode"
)

func CamelToSnake(s string) string {
	if s == "" {
		return s
	}

	delimiter := '_'
	sLen := len(s)
	snake := new(strings.Builder)
	for i, cur := range s {
		if i > 0 && i+1 < sLen {
			if cur >= 'A' && cur <= 'Z' {
				next := s[i+1]
				prev := s[i-1]
				if (next >= 'a' && next <= 'z') || (prev >= 'a' && prev <= 'z') {
					snake.WriteRune(delimiter)
				}
			}
		}
		snake.WriteRune(unicode.ToLower(cur))
	}
	return snake.String()
}

func SnakeToCamel(s string) string {
	if s == "" {
		return s
	}

	sLen := len(s)
	snake := new(strings.Builder)
	for i := 0; i < sLen; i++ {
		cur := rune(s[i])
		if i > 0 && i+1 < sLen {
			if cur == '_' {
				next := s[i+1]
				prev := s[i-1]
				if (next >= 'A' && next <= 'Z') || (prev >= 'a' && prev <= 'z') {
					cur = unicode.ToUpper(rune(next))
					i++
				}
			}
		}
		if i == 0 {
			cur = unicode.ToUpper(cur)
		}
		snake.WriteRune(cur)
	}
	return snake.String()
}
