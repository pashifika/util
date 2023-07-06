// Package slices
package slices

import (
	"reflect"
	"testing"

	"golang.org/x/exp/constraints"
)

func TestMergeNotDuplicate(t *testing.T) {
	type args[E constraints.Ordered] struct {
		s []E
		m [][]E
	}
	type testCase[E constraints.Ordered] struct {
		name string
		args args[E]
		want []E
	}
	tests := []testCase[int]{
		{
			name: "case 1",
			args: args[int]{
				s: []int{1, 2, 3, 5, 6},
				m: [][]int{{1, 2, 4}, {7, 8, 9}},
			},
			want: []int{1, 2, 3, 5, 6, 4, 7, 8, 9},
		},
		{
			name: "case 2",
			args: args[int]{
				s: []int{1, 2},
				m: [][]int{{3}, {1, 2, 3}},
			},
			want: []int{1, 2, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MergeNotDuplicate(tt.args.s, tt.args.m...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MergeNotDuplicate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMergeNotDuplicateFunc(t *testing.T) {
	type args[E comparable, K constraints.Ordered] struct {
		s  []E
		eq func(e E) K
		m  [][]E
	}
	type testCase[E comparable, K constraints.Ordered] struct {
		name string
		args args[E, K]
		want []E
	}
	type entry struct {
		data int
	}
	tests := []testCase[entry, int]{
		{
			name: "case 1",
			args: args[entry, int]{
				s: []entry{
					{data: 1},
					{data: 2},
					{data: 3},
				},
				eq: func(e entry) int { return e.data },
				m: [][]entry{
					{{data: 1}, {data: 4}},
					{{data: 3}, {data: 5}},
				},
			},
			want: []entry{
				{data: 1},
				{data: 2},
				{data: 3},
				{data: 4},
				{data: 5},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MergeNotDuplicateFunc(tt.args.s, tt.args.eq, tt.args.m...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MergeNotDuplicateFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}
