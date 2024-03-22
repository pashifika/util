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

func TestFilterFunc(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}
	type args[S interface{ ~[]E }, T any, E any] struct {
		x      S
		target T
		cmp    func(E, T) bool
	}
	type testCase[S interface{ ~[]E }, T any, E any] struct {
		name  string
		args  args[S, T, E]
		want  int
		want1 S
	}
	tests := []testCase[[]Person, Person, Person]{
		{
			name: "case 1",
			args: args[[]Person, Person, Person]{
				x: []Person{
					{"Alice", 55},
					{"Gopher1", 45},
					{"Gopher2", 45},
					{"Gopher", 33},
					{"Gopher", 31},
					{"Bob", 24},
					{"Gopher3", 45},
				},
				target: Person{Name: "Gopher", Age: 0},
				cmp: func(a Person, b Person) bool {
					return a.Name == b.Name
				},
			},
			want: 2,
			want1: []Person{
				{"Gopher", 33},
				{"Gopher", 31},
			},
		},
		{
			name: "case 2",
			args: args[[]Person, Person, Person]{
				x: []Person{
					{"Alice", 55},
					{"Gopher1", 45},
					{"Gopher2", 45},
					{"Gopher", 33},
					{"Gopher", 31},
					{"Bob", 24},
					{"Gopher3", 45},
				},
				target: Person{Name: "", Age: 45},
				cmp: func(a Person, b Person) bool {
					return a.Age == b.Age
				},
			},
			want: 3,
			want1: []Person{
				{"Gopher1", 45},
				{"Gopher2", 45},
				{"Gopher3", 45},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := FilterFunc(tt.args.x, tt.args.target, tt.args.cmp)
			if got != tt.want {
				t.Errorf("FilterFunc() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("FilterFunc() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
