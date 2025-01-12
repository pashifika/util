// Package datetimes
/*
 * Version: 1.0.0
 * Copyright (c) 2025. Pashifika
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
package datetimes

import (
	"reflect"
	"testing"
	"time"
)

func TestUnixtimeToTime(t *testing.T) {
	type args struct {
		v float64
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{
			name: "normal",
			args: args{v: 1736640000},
			want: time.Date(2025, 1, 12, 0, 0, 0, 0, time.UTC),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UnixtimeToTime(tt.args.v); !reflect.DeepEqual(got.UTC(), tt.want) {
				t.Errorf("UnixtimeToTime() = %v, want %v", got, tt.want)
			}
		})
	}
}
