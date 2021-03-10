/*
MIT License

Copyright (c) 2021 Seanan Xu

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

// Package set implements a set backed by a hash map.
//
// Reference: https://en.wikipedia.org/wiki/Set_%28abstract_data_type%29
package set

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestNewInt64(t *testing.T) {
	testCases := []struct {
		name   string
		input  []int64
		expect []int64
	}{
		{
			name:   "test Int64 New, s nothing",
			input:  nil,
			expect: nil,
		},
		{
			name:   "test Int64 New, inputs multiple elements",
			input:  []int64{1, 2, 3},
			expect: []int64{1, 2, 3},
		},
	}
	for _, tc := range testCases {
		t.Logf("running scenario: %s", tc.name)
		actual := NewInt64(tc.input...)
		validateInt64(t, actual, tc.expect)
	}
}

func TestInt64_Add(t *testing.T) {
	testcases := []struct {
		name   string
		input  []int64
		expect []int64
	}{
		{
			name:   "test Int64 Add, inputs nothing",
			input:  []int64{},
			expect: []int64{},
		},
		{
			name:   "test Int64 Add, inputs multiple elements",
			input:  []int64{1, 2, 3},
			expect: []int64{1, 2, 3},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		actual := Int64{}
		actual.Add(tc.input...)
		validateInt64(t, actual, tc.expect)
	}
}

func TestInt64_Remove(t *testing.T) {
	testcases := []struct {
		name   string
		s      Int64
		input  []int64
		expect []int64
	}{
		{
			name:   "test Int64 Remove, inputs nothing",
			s:      map[int64]struct{}{1: {}, 2: {}, 3: {}},
			input:  []int64{},
			expect: []int64{1, 2, 3},
		},
		{
			name:   "test Int64 Remove, inputs multiple exit elements",
			s:      map[int64]struct{}{1: {}, 2: {}, 3: {}},
			input:  []int64{1, 3},
			expect: []int64{2},
		},
		{
			name:   "test Int64 Remove, inputs multiple non-exit elements",
			s:      map[int64]struct{}{1: {}, 2: {}, 3: {}},
			input:  []int64{0, 6, 7},
			expect: []int64{1, 2, 3},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		tc.s.Remove(tc.input...)
		validateInt64(t, tc.s, tc.expect)
	}
}

func TestInt64_Pop(t *testing.T) {
	testcases := []struct {
		name   string
		s      Int64
		expect bool
	}{
		{
			name:   "test Int64 Pop, s is not empty",
			s:      map[int64]struct{}{1: {}, 2: {}, 3: {}},
			expect: true,
		},
		{
			name:   "test Int64 Pop, s is empty",
			s:      map[int64]struct{}{},
			expect: false,
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		backup := make(map[int64]struct{})
		for k := range tc.s {
			backup[k] = struct{}{}
		}
		key, ok := tc.s.Pop()
		if ok != tc.expect {
			t.Errorf("expect ok: %v, but got: %v", tc.expect, ok)
		}
		delete(backup, key)
		var expect []int64
		for key := range backup {
			expect = append(expect, key)
		}
		validateInt64(t, tc.s, expect)
	}
}

func TestInt64_Size(t *testing.T) {
	testcases := []struct {
		name   string
		s      Int64
		expect int
	}{
		{
			name:   "test Int64 Size, s is empty",
			s:      Int64{},
			expect: 0,
		},
		{
			name:   "test Int64 Size, s is not empty",
			s:      map[int64]struct{}{1: {}, 2: {}, 3: {}},
			expect: 3,
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		if tc.s.Size() != tc.expect {
			t.Errorf("expect size: %d, but got: %d", tc.expect, tc.s.Size())
		}
	}
}

func TestInt64_IsEmpty(t *testing.T) {
	testcases := []struct {
		name   string
		s      Int64
		expect bool
	}{
		{
			name:   "test Int64 Empty, s is empty",
			s:      Int64{},
			expect: true,
		},
		{
			name:   "test interface Empty, s is not empty",
			s:      map[int64]struct{}{1: {}, 2: {}, 3: {}},
			expect: false,
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		actual := tc.s.IsEmpty()
		if actual != tc.expect {
			t.Errorf("expect return: %v, but got: %v", tc.expect, actual)
		}
	}
}

func TestInt64_Clear(t *testing.T) {
	testcases := []struct {
		name string
		s    Int64
	}{
		{
			name: "test Int64 Clear, s is empty",
			s:    Int64{},
		},
		{
			name: "test Int64 Clear, s is not empty",
			s:    map[int64]struct{}{1: {}, 2: {}, 3: {}},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		tc.s.Clear()
		if len(tc.s) != 0 {
			t.Errorf("expect empty, but got: %s", tc.s)
		}
	}
}

func TestInt64_Has(t *testing.T) {
	testcases := []struct {
		name   string
		s      Int64
		input  int64
		expect bool
	}{
		{
			name:   "test Int64 Has, s has input element",
			s:      map[int64]struct{}{1: {}, 2: {}, 3: {}},
			input:  1,
			expect: true,
		},
		{
			name:   "test Int64 Has, s does not have element",
			s:      map[int64]struct{}{1: {}, 2: {}, 3: {}},
			input:  4,
			expect: false,
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		actual := tc.s.Has(tc.input)
		if actual != tc.expect {
			t.Errorf("expect return: %v, but got: %v", tc.expect, actual)
		}
	}
}

func TestInt64_HasAll(t *testing.T) {
	testcases := []struct {
		name   string
		s      Int64
		input  []int64
		expect bool
	}{
		{
			name:   "test Int64 HasAll, set has all input elements",
			s:      map[int64]struct{}{1: {}, 2: {}, 3: {}},
			input:  []int64{1, 2},
			expect: true,
		},
		{
			name:   "test Int64 HasAll, set does not have all input elements",
			s:      map[int64]struct{}{1: {}, 2: {}, 3: {}},
			input:  []int64{8, 9},
			expect: false,
		},
		{
			name:   "test Int64 HasAll, set does not have all input elements, but exist elements in set",
			s:      map[int64]struct{}{1: {}, 2: {}, 3: {}},
			input:  []int64{1, 9},
			expect: false,
		},
		{
			name:   "test Int64 HasAll, input empty",
			s:      map[int64]struct{}{1: {}, 2: {}, 3: {}},
			input:  []int64{},
			expect: true,
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		actual := tc.s.HasAll(tc.input...)
		if actual != tc.expect {
			t.Errorf("expect return: %v, but got: %v", tc.expect, actual)
		}
	}
}

func TestInt64_HasAny(t *testing.T) {
	testcases := []struct {
		name   string
		s      Int64
		input  []int64
		expect bool
	}{
		{
			name:   "test Int64 HasAny, s has all elements",
			s:      map[int64]struct{}{1: {}, 2: {}, 3: {}},
			input:  []int64{1, 2},
			expect: true,
		},
		{
			name:   "test Int64 HasAny, s does not have all elements, but exist elements in set",
			s:      map[int64]struct{}{1: {}, 2: {}, 3: {}},
			input:  []int64{1, 9},
			expect: true,
		},
		{
			name:   "test Int64 HasAny, s does not has all elements",
			s:      map[int64]struct{}{1: {}, 2: {}, 3: {}},
			input:  []int64{8, 9},
			expect: false,
		},
		{
			name:   "test Int64 HasAny, input empty",
			s:      map[int64]struct{}{1: {}, 2: {}, 3: {}},
			input:  []int64{},
			expect: false,
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		actual := tc.s.HasAny(tc.input...)
		if actual != tc.expect {
			t.Errorf("expect return: %v, but got: %v", tc.expect, actual)
		}
	}
}

func TestInt64_List(t *testing.T) {
	testcases := []struct {
		name   string
		s      Int64
		input  []int64
		expect bool
	}{
		{
			name: "test Int64 List, s is empty",
			s:    Int64{},
		},
		{
			name: "test Int64 List, s is not empty",
			s:    map[int64]struct{}{1: {}, 2: {}, 3: {}},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		actual := tc.s.List()
		validateInt64(t, tc.s, actual)
	}
}

func TestInt64_SortedList(t *testing.T) {
	testcases := []struct {
		name   string
		s      Int64
		expect []int64
	}{
		{
			name: "test Int64 List, s is empty",
			s:    Int64{},
		},
		{
			name:   "test Int64 SortedList, s is not empty",
			s:      map[int64]struct{}{1: {}, 2: {}, 3: {}},
			expect: []int64{1, 2, 3},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		actual := tc.s.SortedList(func(i, j int64) bool {
			return i < j
		})
		if len(actual) != len(tc.expect) {
			t.Errorf("expect set len: %d, but got: %d", len(tc.expect), len(actual))
		}
		for i := 0; i < len(tc.expect); i++ {
			if actual[i] != tc.expect[i] {
				t.Errorf("expect slice: %v, but got: %v", tc.expect, actual)
			}
		}
	}
}

func TestInt64_Each(t *testing.T) {
	testcases := []struct {
		name   string
		origin Int64
	}{
		{
			name:   "test Int64 Each",
			origin: map[int64]struct{}{1: {}, 2: {}, 3: {}},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		var expect []int64
		tc.origin.Each(func(i int64) {
			expect = append(expect, i)
		})
		validateInt64(t, tc.origin, expect)
	}
}

func TestInt64_EachE(t *testing.T) {
	inputErr := errors.New("s error")
	testcases := []struct {
		name      string
		origin    Int64
		breakEach bool
		inputErr  error
		expectLen int
		expectErr error
	}{
		{
			name:      "test Int64 EachE",
			origin:    map[int64]struct{}{1: {}, 2: {}, 3: {}},
			expectLen: 3,
		},
		{
			name:      "test Int64 EachE, return break error",
			origin:    map[int64]struct{}{1: {}, 2: {}, 3: {}},
			breakEach: true,
		},
		{
			name:      "test Int64 Each, returns error",
			origin:    map[int64]struct{}{1: {}, 2: {}, 3: {}},
			inputErr:  inputErr,
			expectErr: inputErr,
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		var expect []int64
		err := tc.origin.EachE(func(i int64) error {
			if tc.breakEach {
				return ErrBreakEach
			}
			if tc.inputErr != nil {
				return tc.inputErr
			}
			expect = append(expect, i)
			return nil
		})
		if err != nil {
			if tc.expectErr != err {
				t.Errorf("expect error: %v, but got: %v", tc.expectErr, err)
			}
			return
		}
		if tc.expectLen > 0 {
			validateInt64(t, tc.origin, expect)
		}
	}
}

func TestInt64_Union(t *testing.T) {
	testcases := []struct {
		name   string
		s      Int64
		t      Int64
		expect []int64
	}{
		{
			name:   "test Int64 Union, s and s are empty",
			s:      Int64{},
			t:      Int64{},
			expect: []int64{},
		},
		{
			name:   "test Int64 Union, s is empty",
			s:      Int64{},
			t:      map[int64]struct{}{1: {}, 2: {}, 3: {}},
			expect: []int64{1, 2, 3},
		},
		{
			name:   "test Int64 Union, s is empty",
			s:      map[int64]struct{}{1: {}, 2: {}, 3: {}},
			t:      Int64{},
			expect: []int64{1, 2, 3},
		},
		{
			name:   "test Int64 Union, s has same element to s",
			s:      map[int64]struct{}{1: {}, 2: {}, 3: {}},
			t:      map[int64]struct{}{1: {}, 9: {}, 4: {}},
			expect: []int64{1, 2, 3, 4, 9},
		},
		{
			name:   "test Int64 Union, s does not have same element to s",
			s:      map[int64]struct{}{1: {}, 2: {}, 3: {}},
			t:      map[int64]struct{}{6: {}, 7: {}, 8: {}},
			expect: []int64{1, 2, 3, 6, 7, 8},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		actual := tc.s.Union(tc.t)
		validateInt64(t, actual, tc.expect)
	}
}

func TestInt64_Difference(t *testing.T) {
	testcases := []struct {
		name   string
		s      Int64
		t      Int64
		expect []int64
	}{
		{
			name:   "test Int64 Difference, s is empty",
			s:      Int64{},
			t:      map[int64]struct{}{2: {}, 9: {}, 4: {}},
			expect: []int64{},
		},
		{
			name:   "test Int64 Difference, s is empty",
			s:      map[int64]struct{}{2: {}, 9: {}, 4: {}},
			t:      Int64{},
			expect: []int64{2, 9, 4},
		},
		{
			name:   "test Int64 Difference, s ⊂ s",
			s:      map[int64]struct{}{2: {}, 9: {}},
			t:      map[int64]struct{}{2: {}, 9: {}, 4: {}},
			expect: []int64{},
		},
		{
			name:   "test Int64 Difference, s ⊃ s",
			s:      map[int64]struct{}{2: {}, 9: {}, 4: {}},
			t:      map[int64]struct{}{2: {}},
			expect: []int64{9, 4},
		},
		{
			name:   "test Int64 Difference, s = s",
			s:      map[int64]struct{}{2: {}, 9: {}, 4: {}},
			t:      map[int64]struct{}{2: {}, 9: {}, 4: {}},
			expect: []int64{},
		},
		{
			name:   "test Int64 Difference, s ∩ s = Ø",
			s:      map[int64]struct{}{2: {}, 9: {}, 4: {}},
			t:      map[int64]struct{}{1: {}, 6: {}},
			expect: []int64{2, 9, 4},
		},
		{
			name:   "test Int64 Difference, s ∩ s ≠ Ø",
			s:      map[int64]struct{}{2: {}, 9: {}, 4: {}},
			t:      map[int64]struct{}{2: {}, 9: {}, 5: {}},
			expect: []int64{4},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		actual := tc.s.Difference(tc.t)
		validateInt64(t, actual, tc.expect)
	}
}

func TestInt64_Intersection(t *testing.T) {
	testcases := []struct {
		name   string
		s      Int64
		t      Int64
		expect []int64
	}{
		{
			name:   "test Int64 Intersection, s is empty",
			s:      Int64{},
			t:      map[int64]struct{}{2: {}, 9: {}, 4: {}},
			expect: []int64{},
		},
		{
			name:   "test Int64 Intersection, s is empty",
			s:      map[int64]struct{}{2: {}, 9: {}, 4: {}},
			t:      Int64{},
			expect: []int64{},
		},
		{
			name:   "test Int64 Intersection, s ⊂ s",
			s:      map[int64]struct{}{2: {}, 9: {}},
			t:      map[int64]struct{}{2: {}, 9: {}, 4: {}},
			expect: []int64{2, 9},
		},
		{
			name:   "test Int64 Intersection, s ⊃ s",
			s:      map[int64]struct{}{2: {}, 9: {}, 4: {}},
			t:      map[int64]struct{}{2: {}, 9: {}},
			expect: []int64{2, 9},
		},
		{
			name:   "test Int64 Intersection, s = s",
			s:      map[int64]struct{}{2: {}, 9: {}, 4: {}},
			t:      map[int64]struct{}{2: {}, 9: {}, 4: {}},
			expect: []int64{2, 9, 4},
		},
		{
			name:   "test Int64 Intersection, s ∩ s = Ø",
			s:      map[int64]struct{}{1: {}, 4: {}},
			t:      map[int64]struct{}{2: {}, 6: {}},
			expect: []int64{},
		},
		{
			name:   "test Int64 Intersection, s ∩ s ≠ Ø",
			s:      map[int64]struct{}{1: {}, 4: {}},
			t:      map[int64]struct{}{1: {}, 6: {}},
			expect: []int64{1},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		actual := tc.s.Intersection(tc.t)
		validateInt64(t, actual, tc.expect)
	}
}

func TestInt64_SymmetricDifference(t *testing.T) {
	testcases := []struct {
		name   string
		s      Int64
		t      Int64
		expect []int64
	}{
		{
			name:   "test Int64 SymmetricDifference, s is empty",
			s:      Int64{},
			t:      map[int64]struct{}{2: {}, 9: {}, 4: {}},
			expect: []int64{2, 9, 4},
		},
		{
			name:   "test Int64 SymmetricDifference, s is empty",
			s:      map[int64]struct{}{2: {}, 9: {}, 4: {}},
			t:      Int64{},
			expect: []int64{2, 9, 4},
		},
		{
			name:   "test Int64 SymmetricDifference, s ⊂ s",
			s:      map[int64]struct{}{2: {}, 9: {}},
			t:      map[int64]struct{}{2: {}, 9: {}, 4: {}},
			expect: []int64{4},
		},
		{
			name:   "test Int64 SymmetricDifference, s ⊃ s",
			s:      map[int64]struct{}{2: {}, 9: {}, 4: {}},
			t:      map[int64]struct{}{2: {}, 9: {}},
			expect: []int64{4},
		},
		{
			name:   "test Int64 SymmetricDifference, s = s",
			s:      map[int64]struct{}{2: {}, 9: {}, 4: {}},
			t:      map[int64]struct{}{2: {}, 9: {}, 4: {}},
			expect: []int64{},
		},
		{
			name:   "test Int64 SymmetricDifference, s ∩ s = Ø",
			s:      map[int64]struct{}{1: {}, 4: {}},
			t:      map[int64]struct{}{2: {}, 6: {}},
			expect: []int64{1, 2, 4, 6},
		},
		{
			name:   "test Int64 SymmetricDifference, s ∩ s ≠ Ø",
			s:      map[int64]struct{}{1: {}, 4: {}},
			t:      map[int64]struct{}{1: {}, 6: {}},
			expect: []int64{4, 6},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		actual := tc.s.SymmetricDifference(tc.t)
		validateInt64(t, actual, tc.expect)
	}
}

func TestInt64_IsSubset(t *testing.T) {
	testcases := []struct {
		name   string
		s      Int64
		t      Int64
		expect bool
	}{
		{
			name:   "test Int64 IsSubset, s is empty",
			s:      Int64{},
			t:      map[int64]struct{}{2: {}, 9: {}, 4: {}},
			expect: true,
		},
		{
			name:   "test Int64 IsSubset, s is empty",
			s:      map[int64]struct{}{2: {}, 9: {}, 4: {}},
			t:      Int64{},
			expect: false,
		},
		{
			name:   "test Int64 IsSubset, s ⊂ s",
			s:      map[int64]struct{}{2: {}, 9: {}},
			t:      map[int64]struct{}{2: {}, 9: {}, 4: {}},
			expect: true,
		},
		{
			name:   "test Int64 IsSubset, s ⊃ s",
			s:      map[int64]struct{}{2: {}, 9: {}, 4: {}},
			t:      map[int64]struct{}{2: {}, 9: {}},
			expect: false,
		},
		{
			name:   "test Int64 IsSubset, s = s",
			s:      map[int64]struct{}{2: {}, 9: {}, 4: {}},
			t:      map[int64]struct{}{2: {}, 9: {}, 4: {}},
			expect: true,
		},
		{
			name:   "test Int64 IsSubset, s ∩ s = Ø",
			s:      map[int64]struct{}{1: {}, 4: {}},
			t:      map[int64]struct{}{2: {}, 6: {}},
			expect: false,
		},
		{
			name:   "test Int64 IsSubset, s ∩ s ≠ Ø && s ∩ s ≠ s",
			s:      map[int64]struct{}{1: {}, 4: {}},
			t:      map[int64]struct{}{1: {}, 6: {}},
			expect: false,
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		actual := tc.s.IsSubset(tc.t)
		if actual != tc.expect {
			t.Errorf("expect retrun: %v, but got: %v", tc.expect, actual)
		}
	}
}

func TestInt64_IsSuperset(t *testing.T) {
	testcases := []struct {
		name   string
		s      Int64
		t      Int64
		expect bool
	}{
		{
			name:   "test Int64 IsSuperset, s is empty",
			s:      Int64{},
			t:      map[int64]struct{}{2: {}, 9: {}, 4: {}},
			expect: false,
		},
		{
			name:   "test Int64 IsSuperset, s is empty",
			s:      map[int64]struct{}{2: {}, 9: {}, 4: {}},
			t:      Int64{},
			expect: true,
		},
		{
			name:   "test Int64 IsSuperset, s ⊂ s",
			s:      map[int64]struct{}{2: {}, 9: {}},
			t:      map[int64]struct{}{2: {}, 9: {}, 4: {}},
			expect: false,
		},
		{
			name:   "test Int64 IsSuperset, s ⊃ s",
			s:      map[int64]struct{}{2: {}, 9: {}, 4: {}},
			t:      map[int64]struct{}{2: {}, 9: {}},
			expect: true,
		},
		{
			name:   "test Int64 IsSuperset, s = s",
			s:      map[int64]struct{}{2: {}, 9: {}, 4: {}},
			t:      map[int64]struct{}{2: {}, 9: {}, 4: {}},
			expect: true,
		},
		{
			name:   "test Int64 IsSuperset, s ∩ s = Ø",
			s:      map[int64]struct{}{1: {}, 4: {}},
			t:      map[int64]struct{}{2: {}, 6: {}},
			expect: false,
		},
		{
			name:   "test Int64 IsSuperset, s ∩ s ≠ Ø && s ∩ s ≠ s",
			s:      map[int64]struct{}{1: {}, 4: {}},
			t:      map[int64]struct{}{1: {}, 6: {}},
			expect: false,
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		actual := tc.s.IsSuperset(tc.t)
		if actual != tc.expect {
			t.Errorf("expect retrun: %v, but got: %v", tc.expect, actual)
		}
	}
}

func TestInt64_Equal(t *testing.T) {
	testcases := []struct {
		name   string
		s      Int64
		t      Int64
		expect bool
	}{
		{
			name:   "test Int64 Equal, s is empty",
			s:      Int64{},
			t:      map[int64]struct{}{2: {}, 9: {}, 4: {}},
			expect: false,
		},
		{
			name:   "test Int64 Equal, s is empty",
			s:      map[int64]struct{}{2: {}, 9: {}, 4: {}},
			t:      Int64{},
			expect: false,
		},
		{
			name:   "test Int64 Equal, s ⊂ s",
			s:      map[int64]struct{}{2: {}, 9: {}},
			t:      map[int64]struct{}{2: {}, 9: {}, 4: {}},
			expect: false,
		},
		{
			name:   "test Int64 Equal, s ⊃ s",
			s:      map[int64]struct{}{2: {}, 9: {}, 4: {}},
			t:      map[int64]struct{}{2: {}, 9: {}},
			expect: false,
		},
		{
			name:   "test Int64 Equal, s = s",
			s:      map[int64]struct{}{2: {}, 9: {}, 4: {}},
			t:      map[int64]struct{}{2: {}, 9: {}, 4: {}},
			expect: true,
		},
		{
			name:   "test Int64 Equal, s ∩ s = Ø",
			s:      map[int64]struct{}{1: {}, 4: {}},
			t:      map[int64]struct{}{2: {}, 6: {}},
			expect: false,
		},
		{
			name:   "test Int64 Equal, s ∩ s ≠ Ø && s ∩ s ≠ s",
			s:      map[int64]struct{}{1: {}, 4: {}},
			t:      map[int64]struct{}{1: {}, 6: {}},
			expect: false,
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		actual := tc.s.Equal(tc.t)
		if actual != tc.expect {
			t.Errorf("expect retrun: %v, but got: %v", tc.expect, actual)
		}
	}
}

func TestInt64_Copy(t *testing.T) {
	testcases := []struct {
		name   string
		s      Int64
		expect []int64
	}{
		{
			name:   "test Int64 Copy, s is empty",
			s:      Int64{},
			expect: []int64{},
		},
		{
			name:   "test Int64 Copy, s is not empty",
			s:      map[int64]struct{}{2: {}, 9: {}, 4: {}},
			expect: []int64{2, 9, 4},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		actual := tc.s.Copy()
		validateInt64(t, actual, tc.expect)
	}
}

func TestInt64_String(t *testing.T) {
	testcases := []struct {
		name   string
		s      Int64
		expect string
	}{
		{
			name:   "test Int64 String, s is empty",
			s:      Int64{},
			expect: "[]",
		},
		{
			name:   "test Int64 String, s is not empty",
			s:      map[int64]struct{}{1: {}},
			expect: "[1]",
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		actual := tc.s.String()
		if actual != tc.expect {
			t.Errorf("expect string: %s, but got: %s", tc.expect, actual)
		}
	}
}

func validateInt64(t *testing.T, actual Int64, expect []int64) {
	if len(expect) != len(actual) {
		t.Errorf("expect set len: %d, but got: %d", len(expect), len(actual))
	}
	slice2String := func(i []int64) string {
		v := make([]string, 0, len(i))
		for _, element := range i {
			v = append(v, fmt.Sprintf("%v", element))
		}
		return fmt.Sprintf("[%s]", strings.Join(v, ", "))
	}
	for _, k := range expect {
		if _, ok := actual[k]; !ok {
			t.Errorf("expect set: %s, but got: %s", slice2String(expect), actual)
		}
	}
}
