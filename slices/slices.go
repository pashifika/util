// Package slices
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
package slices

import (
	"golang.org/x/exp/constraints"
)

// MergeNotDuplicate is merge multiple slices and remove duplicate entries.
//
// see more: golang.org/x/exp/slices
func MergeNotDuplicate[E constraints.Ordered](s []E, m ...[]E) []E {
	var res []E
	check := make(map[E]struct{})
	for _, e := range s {
		if _, ok := check[e]; ok {
			continue
		}
		check[e] = struct{}{}
		res = append(res, e)
	}
	for _, rows := range m {
		for _, row := range rows {
			if _, ok := check[row]; ok {
				continue
			}
			check[row] = struct{}{}
			res = append(res, row)
		}
	}

	return res
}

// MergeNotDuplicateFunc is like MergeNotDuplicate but uses a comparison function.
func MergeNotDuplicateFunc[E comparable, K constraints.Ordered](s []E, eq func(e E) K, m ...[]E) []E {
	var res []E
	check := make(map[K]struct{})
	for _, e := range s {
		key := eq(e)
		if _, ok := check[key]; ok {
			continue
		}
		check[key] = struct{}{}
		res = append(res, e)
	}
	for _, rows := range m {
		for _, row := range rows {
			key := eq(row)
			if _, ok := check[key]; ok {
				continue
			}
			check[key] = struct{}{}
			res = append(res, row)
		}
	}

	return res
}

func FilterFunc[S ~[]E, E, T any](x S, target T, cmp func(E, T) bool) (int, S) {
	n := 0
	arr := x[:0]
	for _, e := range x {
		if cmp(e, target) {
			arr = append(arr, e)
			n++
		}
	}
	return n, arr
}
