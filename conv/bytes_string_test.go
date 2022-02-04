// Package conv
/*
 * Version: 1.0.0
 * Copyright (c) 2022. Pashifika
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
	"bytes"
	"testing"
)

func TestBytesToString(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "number", args: args{b: []byte("123")}, want: "123"},
		{name: "alphabet", args: args{b: []byte("abcdefg")}, want: "abcdefg"},
		{name: "utf8", args: args{b: []byte("あいうえお・")}, want: "あいうえお・"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BytesToString(tt.args.b); got != tt.want {
				t.Errorf("BytesToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringToBytes(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{name: "number", args: args{s: "123"}, want: []byte("123")},
		{name: "alphabet", args: args{s: "abcdefg"}, want: []byte("abcdefg")},
		{name: "utf8", args: args{s: "あいうえお・"}, want: []byte("あいうえお・")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringToBytes(tt.args.s); !bytes.Equal(got, tt.want) {
				t.Errorf("BytesToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkBytesToString(b *testing.B) {
	_bytes := []byte("あいうえお・あいうえお・あいうえお・あいうえお")
	for i := 0; i < b.N; i++ {
		BytesToString(_bytes)
	}
}

func BenchmarkString(b *testing.B) {
	_bytes := []byte("あいうえお・あいうえお・あいうえお・あいうえお")
	for i := 0; i < b.N; i++ {
		_ = string(_bytes)
	}
}

func BenchmarkStringToBytes(b *testing.B) {
	str := "あいうえお・あいうえお・あいうえお・あいうえお"
	for i := 0; i < b.N; i++ {
		StringToBytes(str)
	}
}

func BenchmarkBytes(b *testing.B) {
	str := "あいうえお・あいうえお・あいうえお・あいうえお"
	for i := 0; i < b.N; i++ {
		_ = []byte(str)
	}
}
