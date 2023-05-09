// Package random
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
package random

import (
	cRand "crypto/rand"
	"math/big"
)

var (
	Numeric             = "0123456789"                                                      // numeric: [0-9]
	AsciiAlphabetsLower = "abcdefghijklmnopqrstuvwxyz"                                      // Ascii lower alphabets: [a-z]
	AsciiAlphabetsUpper = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"                                      // Ascii upper alphabets: [A-Z]
	AsciiAlphabets      = AsciiAlphabetsLower + AsciiAlphabetsUpper                         // Ascii alphabets: [a-zA-Z]
	AsciiCharacters     = AsciiAlphabetsLower + AsciiAlphabetsUpper + Numeric               // Ascii characters: [a-zA-Z0-9]
	Hexadecimal         = "0123456789abcdefABCDEF"                                          // hexadecimal number: [0-9a-fA-F]
	Punctuation         = "!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~"                              // punctuation and special characters
	Printables          = Numeric + AsciiAlphabetsLower + AsciiAlphabetsUpper + Punctuation // printable characters
)

// Int generates a cryptographically-secure random Int.
// Provided max must be greater than 0.
func Int(max int) int {
	return int(Int64(int64(max)))
}

// Int64 generates a cryptographically-secure random Int64.
// Provided max must be greater than 0.
func Int64(max int64) int64 {
	if max <= 0 {
		panic("input must be greater than 0")
	}
	target, err := cRand.Int(cRand.Reader, big.NewInt(max))
	if err != nil {
		panic(err)
	}

	return target.Int64()
}

// String generates a cryptographically secure string.
func String(n int) string {
	return Random(n, AsciiCharacters)
}

// IntRange returns a random integer between a given range.
func IntRange(min int, max int) int {
	i := Int(max - min)
	i += min
	return i
}

// IntRange64 returns a random big integer between a given range.
func IntRange64(min int64, max int64) int64 {
	i := Int64(max - min)
	i += min
	return i
}

// Random is responsible for generating random data from a given character set.
func Random(n int, charset string) string {
	var charsetByte = []byte(charset)
	s := make([]byte, n)
	max := len(charset)
	for i := range s {
		s[i] = charsetByte[Int(max)]
	}

	return string(s)
}

// Choice makes a random choice from a slice.
func Choice[T comparable](datas []T) T {
	return datas[Int(len(datas))]
}

// ChoiceSlice select n comparable are random choice in a slice.
func ChoiceSlice[T comparable](datas []T, n int) []T {
	if n < 1 {
		n = 1
	}
	slice := make([]T, n)
	for i := 0; i < n; i++ {
		slice[i] = Choice(datas)
	}
	return slice
}
