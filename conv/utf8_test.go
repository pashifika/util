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
	"testing"
)

var (
	src = "黄昏よりも昏きもの　血の流れより紅きもの\n時の流れに埋もれし　偉大な汝の名において\n我ここに闇に誓わん　我等が前に立ち塞がりし　すべての愚かなるものに　\n我と汝が力もて　等しく滅びを与えんことを！"
)

func TestFindUnicodeString(t *testing.T) {
	type args struct {
		src  string
		find string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "血の流れより紅きもの", args: args{
			src:  src,
			find: "血の流れより紅きもの",
		}, want: true},
		{name: "偉大汝の名において", args: args{
			src:  src,
			find: "偉大汝の名において",
		}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FindUnicodeString(tt.args.src, tt.args.find); got != tt.want {
				t.Errorf("FindUnicodeString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCutUnicodeString(t *testing.T) {
	type args struct {
		str    string
		length int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "黄昏よりも昏きもの", args: args{
			str:    "黄昏よりも昏きもの　血の流れより紅きもの",
			length: 9,
		}, want: "黄昏よりも昏きもの"},
		{name: "黄昏よりも昏きもの　血の流れより紅きもの", args: args{
			str:    "黄昏よりも昏きもの　血の流れより紅きもの",
			length: 20,
		}, want: "黄昏よりも昏きもの　血の流れより紅きもの"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CutUnicodeString(tt.args.str, tt.args.length); got != tt.want {
				t.Errorf("CutUnicodeString() = %v, want %v", got, tt.want)
			}
		})
	}
}
