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

import "testing"

func TestCamelToSnake(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "Case 1",
			input: "AsBsCs",
			want:  "as_bs_cs",
		},
		{
			name:  "Case 2",
			input: "AsBsCs_",
			want:  "as_bs_cs_",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CamelToSnake(tt.input); got != tt.want {
				t.Errorf("CamelToSnake() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSnakeToCamel(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "Case 1",
			input: "as_bs_cs",
			want:  "AsBsCs",
		},
		{
			name:  "Case 2",
			input: "as_bs_cs_",
			want:  "AsBsCs_",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SnakeToCamel(tt.input); got != tt.want {
				t.Errorf("SnakeToCamel() = %v, want %v", got, tt.want)
			}
		})
	}
}
