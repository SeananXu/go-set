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

func TestNewInt8(t *testing.T) {
	testCases := []struct {
		name   string
		input  []int8
		expect []int8
	}{
		{
			name:   "test Int8 New, s nothing",
			input:  nil,
			expect: nil,
		},
		{
			name:   "test Int8 New, inputs multiple elements",
			input:  []int8{1, 2, 3},
			expect: []int8{1, 2, 3},
		},
	}
	for _, tc := range testCases {
		t.Logf("running scenario: %s", tc.name)
		actual := NewInt8(tc.input...)
		validateInt8(t, actual, tc.expect)
	}
}

func TestInt8_Add(t *testing.T) {
	testcases := []struct {
		name   string
		input  []int8
		expect []int8
	}{
		{
			name:   "test Int8 Add, inputs nothing",
			input:  []int8{},
			expect: []int8{},
		},
		{
			name:   "test Int8 Add, inputs multiple elements",
			input:  []int8{1, 2, 3},
			expect: []int8{1, 2, 3},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		actual := Int8{}
		actual.Add(tc.input...)
		validateInt8(t, actual, tc.expect)
	}
}

func TestInt8_Remove(t *testing.T) {
	testcases := []struct {
		name   string
		s      Int8
		input  []int8
		expect []int8
	}{
		{
			name:   "test Int8 Remove, inputs nothing",
			s:      map[int8]struct{}{1: {}, 2: {}, 3: {}},
			input:  []int8{},
			expect: []int8{1, 2, 3},
		},
		{
			name:   "test Int8 Remove, inputs multiple exit elements",
			s:      map[int8]struct{}{1: {}, 2: {}, 3: {}},
			input:  []int8{1, 3},
			expect: []int8{2},
		},
		{
			name:   "test Int8 Remove, inputs multiple non-exit elements",
			s:      map[int8]struct{}{1: {}, 2: {}, 3: {}},
			input:  []int8{0, 6, 7},
			expect: []int8{1, 2, 3},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		tc.s.Remove(tc.input...)
		validateInt8(t, tc.s, tc.expect)
	}
}

func TestInt8_Pop(t *testing.T) {
	testcases := []struct {
		name   string
		s      Int8
		expect bool
	}{
		{
			name:   "test Int8 Pop, s is not empty",
			s:      map[int8]struct{}{1: {}, 2: {}, 3: {}},
			expect: true,
		},
		{
			name:   "test Int8 Pop, s is empty",
			s:      map[int8]struct{}{},
			expect: false,
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		backup := make(map[int8]struct{})
		for k := range tc.s {
			backup[k] = struct{}{}
		}
		key, ok := tc.s.Pop()
		if ok != tc.expect {
			t.Errorf("expect ok: %v, but got: %v", tc.expect, ok)
		}
		delete(backup, key)
		var expect []int8
		for key := range backup {
			expect = append(expect, key)
		}
		validateInt8(t, tc.s, expect)
	}
}

func TestInt8_Size(t *testing.T) {
	testcases := []struct {
		name   string
		s      Int8
		expect int
	}{
		{
			name:   "test Int8 Size, s is empty",
			s:      Int8{},
			expect: 0,
		},
		{
			name:   "test Int8 Size, s is not empty",
			s:      map[int8]struct{}{1: {}, 2: {}, 3: {}},
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

func TestInt8_IsEmpty(t *testing.T) {
	testcases := []struct {
		name   string
		s      Int8
		expect bool
	}{
		{
			name:   "test Int8 Empty, s is empty",
			s:      Int8{},
			expect: true,
		},
		{
			name:   "test interface Empty, s is not empty",
			s:      map[int8]struct{}{1: {}, 2: {}, 3: {}},
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

func TestInt8_Clear(t *testing.T) {
	testcases := []struct {
		name string
		s    Int8
	}{
		{
			name: "test Int8 Clear, s is empty",
			s:    Int8{},
		},
		{
			name: "test Int8 Clear, s is not empty",
			s:    map[int8]struct{}{1: {}, 2: {}, 3: {}},
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

func TestInt8_Has(t *testing.T) {
	testcases := []struct {
		name   string
		s      Int8
		input  int8
		expect bool
	}{
		{
			name:   "test Int8 Has, s has input element",
			s:      map[int8]struct{}{1: {}, 2: {}, 3: {}},
			input:  1,
			expect: true,
		},
		{
			name:   "test Int8 Has, s does not have element",
			s:      map[int8]struct{}{1: {}, 2: {}, 3: {}},
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

func TestInt8_HasAll(t *testing.T) {
	testcases := []struct {
		name   string
		s      Int8
		input  []int8
		expect bool
	}{
		{
			name:   "test Int8 HasAll, set has all input elements",
			s:      map[int8]struct{}{1: {}, 2: {}, 3: {}},
			input:  []int8{1, 2},
			expect: true,
		},
		{
			name:   "test Int8 HasAll, set does not have all input elements",
			s:      map[int8]struct{}{1: {}, 2: {}, 3: {}},
			input:  []int8{8, 9},
			expect: false,
		},
		{
			name:   "test Int8 HasAll, set does not have all input elements, but exist elements in set",
			s:      map[int8]struct{}{1: {}, 2: {}, 3: {}},
			input:  []int8{1, 9},
			expect: false,
		},
		{
			name:   "test Int8 HasAll, input empty",
			s:      map[int8]struct{}{1: {}, 2: {}, 3: {}},
			input:  []int8{},
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

func TestInt8_HasAny(t *testing.T) {
	testcases := []struct {
		name   string
		s      Int8
		input  []int8
		expect bool
	}{
		{
			name:   "test Int8 HasAny, s has all elements",
			s:      map[int8]struct{}{1: {}, 2: {}, 3: {}},
			input:  []int8{1, 2},
			expect: true,
		},
		{
			name:   "test Int8 HasAny, s does not have all elements, but exist elements in set",
			s:      map[int8]struct{}{1: {}, 2: {}, 3: {}},
			input:  []int8{1, 9},
			expect: true,
		},
		{
			name:   "test Int8 HasAny, s does not has all elements",
			s:      map[int8]struct{}{1: {}, 2: {}, 3: {}},
			input:  []int8{8, 9},
			expect: false,
		},
		{
			name:   "test Int8 HasAny, input empty",
			s:      map[int8]struct{}{1: {}, 2: {}, 3: {}},
			input:  []int8{},
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

func TestInt8_List(t *testing.T) {
	testcases := []struct {
		name   string
		s      Int8
		input  []int8
		expect bool
	}{
		{
			name: "test Int8 List, s is empty",
			s:    Int8{},
		},
		{
			name: "test Int8 List, s is not empty",
			s:    map[int8]struct{}{1: {}, 2: {}, 3: {}},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		actual := tc.s.List()
		validateInt8(t, tc.s, actual)
	}
}

func TestInt8_SortedList(t *testing.T) {
	testcases := []struct {
		name   string
		s      Int8
		expect []int8
	}{
		{
			name: "test Int8 List, s is empty",
			s:    Int8{},
		},
		{
			name:   "test Int8 SortedList, s is not empty",
			s:      map[int8]struct{}{1: {}, 2: {}, 3: {}},
			expect: []int8{1, 2, 3},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		actual := tc.s.SortedList(func(i, j int8) bool {
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

func TestInt8_Each(t *testing.T) {
	testcases := []struct {
		name   string
		origin Int8
	}{
		{
			name:   "test Int8 Each",
			origin: map[int8]struct{}{1: {}, 2: {}, 3: {}},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		var expect []int8
		tc.origin.Each(func(i int8) {
			expect = append(expect, i)
		})
		validateInt8(t, tc.origin, expect)
	}
}

func TestInt8_EachE(t *testing.T) {
	inputErr := errors.New("s error")
	testcases := []struct {
		name      string
		origin    Int8
		breakEach bool
		inputErr  error
		expectLen int
		expectErr error
	}{
		{
			name:      "test Int8 EachE",
			origin:    map[int8]struct{}{1: {}, 2: {}, 3: {}},
			expectLen: 3,
		},
		{
			name:      "test Int8 EachE, return break error",
			origin:    map[int8]struct{}{1: {}, 2: {}, 3: {}},
			breakEach: true,
		},
		{
			name:      "test Int8 Each, returns error",
			origin:    map[int8]struct{}{1: {}, 2: {}, 3: {}},
			inputErr:  inputErr,
			expectErr: inputErr,
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		var expect []int8
		err := tc.origin.EachE(func(i int8) error {
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
			validateInt8(t, tc.origin, expect)
		}
	}
}

func TestInt8_Union(t *testing.T) {
	testcases := []struct {
		name   string
		s      Int8
		t      Int8
		expect []int8
	}{
		{
			name:   "test Int8 Union, s and s are empty",
			s:      Int8{},
			t:      Int8{},
			expect: []int8{},
		},
		{
			name:   "test Int8 Union, s is empty",
			s:      Int8{},
			t:      map[int8]struct{}{1: {}, 2: {}, 3: {}},
			expect: []int8{1, 2, 3},
		},
		{
			name:   "test Int8 Union, s is empty",
			s:      map[int8]struct{}{1: {}, 2: {}, 3: {}},
			t:      Int8{},
			expect: []int8{1, 2, 3},
		},
		{
			name:   "test Int8 Union, s has same element to s",
			s:      map[int8]struct{}{1: {}, 2: {}, 3: {}},
			t:      map[int8]struct{}{1: {}, 9: {}, 4: {}},
			expect: []int8{1, 2, 3, 4, 9},
		},
		{
			name:   "test Int8 Union, s does not have same element to s",
			s:      map[int8]struct{}{1: {}, 2: {}, 3: {}},
			t:      map[int8]struct{}{6: {}, 7: {}, 8: {}},
			expect: []int8{1, 2, 3, 6, 7, 8},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		actual := tc.s.Union(tc.t)
		validateInt8(t, actual, tc.expect)
	}
}

func TestInt8_Difference(t *testing.T) {
	testcases := []struct {
		name   string
		s      Int8
		t      Int8
		expect []int8
	}{
		{
			name:   "test Int8 Difference, s is empty",
			s:      Int8{},
			t:      map[int8]struct{}{2: {}, 9: {}, 4: {}},
			expect: []int8{},
		},
		{
			name:   "test Int8 Difference, s is empty",
			s:      map[int8]struct{}{2: {}, 9: {}, 4: {}},
			t:      Int8{},
			expect: []int8{2, 9, 4},
		},
		{
			name:   "test Int8 Difference, s ⊂ s",
			s:      map[int8]struct{}{2: {}, 9: {}},
			t:      map[int8]struct{}{2: {}, 9: {}, 4: {}},
			expect: []int8{},
		},
		{
			name:   "test Int8 Difference, s ⊃ s",
			s:      map[int8]struct{}{2: {}, 9: {}, 4: {}},
			t:      map[int8]struct{}{2: {}},
			expect: []int8{9, 4},
		},
		{
			name:   "test Int8 Difference, s = s",
			s:      map[int8]struct{}{2: {}, 9: {}, 4: {}},
			t:      map[int8]struct{}{2: {}, 9: {}, 4: {}},
			expect: []int8{},
		},
		{
			name:   "test Int8 Difference, s ∩ s = Ø",
			s:      map[int8]struct{}{2: {}, 9: {}, 4: {}},
			t:      map[int8]struct{}{1: {}, 6: {}},
			expect: []int8{2, 9, 4},
		},
		{
			name:   "test Int8 Difference, s ∩ s ≠ Ø",
			s:      map[int8]struct{}{2: {}, 9: {}, 4: {}},
			t:      map[int8]struct{}{2: {}, 9: {}, 5: {}},
			expect: []int8{4},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		actual := tc.s.Difference(tc.t)
		validateInt8(t, actual, tc.expect)
	}
}

func TestInt8_Intersection(t *testing.T) {
	testcases := []struct {
		name   string
		s      Int8
		t      Int8
		expect []int8
	}{
		{
			name:   "test Int8 Intersection, s is empty",
			s:      Int8{},
			t:      map[int8]struct{}{2: {}, 9: {}, 4: {}},
			expect: []int8{},
		},
		{
			name:   "test Int8 Intersection, s is empty",
			s:      map[int8]struct{}{2: {}, 9: {}, 4: {}},
			t:      Int8{},
			expect: []int8{},
		},
		{
			name:   "test Int8 Intersection, s ⊂ s",
			s:      map[int8]struct{}{2: {}, 9: {}},
			t:      map[int8]struct{}{2: {}, 9: {}, 4: {}},
			expect: []int8{2, 9},
		},
		{
			name:   "test Int8 Intersection, s ⊃ s",
			s:      map[int8]struct{}{2: {}, 9: {}, 4: {}},
			t:      map[int8]struct{}{2: {}, 9: {}},
			expect: []int8{2, 9},
		},
		{
			name:   "test Int8 Intersection, s = s",
			s:      map[int8]struct{}{2: {}, 9: {}, 4: {}},
			t:      map[int8]struct{}{2: {}, 9: {}, 4: {}},
			expect: []int8{2, 9, 4},
		},
		{
			name:   "test Int8 Intersection, s ∩ s = Ø",
			s:      map[int8]struct{}{1: {}, 4: {}},
			t:      map[int8]struct{}{2: {}, 6: {}},
			expect: []int8{},
		},
		{
			name:   "test Int8 Intersection, s ∩ s ≠ Ø",
			s:      map[int8]struct{}{1: {}, 4: {}},
			t:      map[int8]struct{}{1: {}, 6: {}},
			expect: []int8{1},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		actual := tc.s.Intersection(tc.t)
		validateInt8(t, actual, tc.expect)
	}
}

func TestInt8_SymmetricDifference(t *testing.T) {
	testcases := []struct {
		name   string
		s      Int8
		t      Int8
		expect []int8
	}{
		{
			name:   "test Int8 SymmetricDifference, s is empty",
			s:      Int8{},
			t:      map[int8]struct{}{2: {}, 9: {}, 4: {}},
			expect: []int8{2, 9, 4},
		},
		{
			name:   "test Int8 SymmetricDifference, s is empty",
			s:      map[int8]struct{}{2: {}, 9: {}, 4: {}},
			t:      Int8{},
			expect: []int8{2, 9, 4},
		},
		{
			name:   "test Int8 SymmetricDifference, s ⊂ s",
			s:      map[int8]struct{}{2: {}, 9: {}},
			t:      map[int8]struct{}{2: {}, 9: {}, 4: {}},
			expect: []int8{4},
		},
		{
			name:   "test Int8 SymmetricDifference, s ⊃ s",
			s:      map[int8]struct{}{2: {}, 9: {}, 4: {}},
			t:      map[int8]struct{}{2: {}, 9: {}},
			expect: []int8{4},
		},
		{
			name:   "test Int8 SymmetricDifference, s = s",
			s:      map[int8]struct{}{2: {}, 9: {}, 4: {}},
			t:      map[int8]struct{}{2: {}, 9: {}, 4: {}},
			expect: []int8{},
		},
		{
			name:   "test Int8 SymmetricDifference, s ∩ s = Ø",
			s:      map[int8]struct{}{1: {}, 4: {}},
			t:      map[int8]struct{}{2: {}, 6: {}},
			expect: []int8{1, 2, 4, 6},
		},
		{
			name:   "test Int8 SymmetricDifference, s ∩ s ≠ Ø",
			s:      map[int8]struct{}{1: {}, 4: {}},
			t:      map[int8]struct{}{1: {}, 6: {}},
			expect: []int8{4, 6},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		actual := tc.s.SymmetricDifference(tc.t)
		validateInt8(t, actual, tc.expect)
	}
}

func TestInt8_IsSubset(t *testing.T) {
	testcases := []struct {
		name   string
		s      Int8
		t      Int8
		expect bool
	}{
		{
			name:   "test Int8 IsSubset, s is empty",
			s:      Int8{},
			t:      map[int8]struct{}{2: {}, 9: {}, 4: {}},
			expect: true,
		},
		{
			name:   "test Int8 IsSubset, s is empty",
			s:      map[int8]struct{}{2: {}, 9: {}, 4: {}},
			t:      Int8{},
			expect: false,
		},
		{
			name:   "test Int8 IsSubset, s ⊂ s",
			s:      map[int8]struct{}{2: {}, 9: {}},
			t:      map[int8]struct{}{2: {}, 9: {}, 4: {}},
			expect: true,
		},
		{
			name:   "test Int8 IsSubset, s ⊃ s",
			s:      map[int8]struct{}{2: {}, 9: {}, 4: {}},
			t:      map[int8]struct{}{2: {}, 9: {}},
			expect: false,
		},
		{
			name:   "test Int8 IsSubset, s = s",
			s:      map[int8]struct{}{2: {}, 9: {}, 4: {}},
			t:      map[int8]struct{}{2: {}, 9: {}, 4: {}},
			expect: true,
		},
		{
			name:   "test Int8 IsSubset, s ∩ s = Ø",
			s:      map[int8]struct{}{1: {}, 4: {}},
			t:      map[int8]struct{}{2: {}, 6: {}},
			expect: false,
		},
		{
			name:   "test Int8 IsSubset, s ∩ s ≠ Ø && s ∩ s ≠ s",
			s:      map[int8]struct{}{1: {}, 4: {}},
			t:      map[int8]struct{}{1: {}, 6: {}},
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

func TestInt8_IsSuperset(t *testing.T) {
	testcases := []struct {
		name   string
		s      Int8
		t      Int8
		expect bool
	}{
		{
			name:   "test Int8 IsSuperset, s is empty",
			s:      Int8{},
			t:      map[int8]struct{}{2: {}, 9: {}, 4: {}},
			expect: false,
		},
		{
			name:   "test Int8 IsSuperset, s is empty",
			s:      map[int8]struct{}{2: {}, 9: {}, 4: {}},
			t:      Int8{},
			expect: true,
		},
		{
			name:   "test Int8 IsSuperset, s ⊂ s",
			s:      map[int8]struct{}{2: {}, 9: {}},
			t:      map[int8]struct{}{2: {}, 9: {}, 4: {}},
			expect: false,
		},
		{
			name:   "test Int8 IsSuperset, s ⊃ s",
			s:      map[int8]struct{}{2: {}, 9: {}, 4: {}},
			t:      map[int8]struct{}{2: {}, 9: {}},
			expect: true,
		},
		{
			name:   "test Int8 IsSuperset, s = s",
			s:      map[int8]struct{}{2: {}, 9: {}, 4: {}},
			t:      map[int8]struct{}{2: {}, 9: {}, 4: {}},
			expect: true,
		},
		{
			name:   "test Int8 IsSuperset, s ∩ s = Ø",
			s:      map[int8]struct{}{1: {}, 4: {}},
			t:      map[int8]struct{}{2: {}, 6: {}},
			expect: false,
		},
		{
			name:   "test Int8 IsSuperset, s ∩ s ≠ Ø && s ∩ s ≠ s",
			s:      map[int8]struct{}{1: {}, 4: {}},
			t:      map[int8]struct{}{1: {}, 6: {}},
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

func TestInt8_Equal(t *testing.T) {
	testcases := []struct {
		name   string
		s      Int8
		t      Int8
		expect bool
	}{
		{
			name:   "test Int8 Equal, s is empty",
			s:      Int8{},
			t:      map[int8]struct{}{2: {}, 9: {}, 4: {}},
			expect: false,
		},
		{
			name:   "test Int8 Equal, s is empty",
			s:      map[int8]struct{}{2: {}, 9: {}, 4: {}},
			t:      Int8{},
			expect: false,
		},
		{
			name:   "test Int8 Equal, s ⊂ s",
			s:      map[int8]struct{}{2: {}, 9: {}},
			t:      map[int8]struct{}{2: {}, 9: {}, 4: {}},
			expect: false,
		},
		{
			name:   "test Int8 Equal, s ⊃ s",
			s:      map[int8]struct{}{2: {}, 9: {}, 4: {}},
			t:      map[int8]struct{}{2: {}, 9: {}},
			expect: false,
		},
		{
			name:   "test Int8 Equal, s = s",
			s:      map[int8]struct{}{2: {}, 9: {}, 4: {}},
			t:      map[int8]struct{}{2: {}, 9: {}, 4: {}},
			expect: true,
		},
		{
			name:   "test Int8 Equal, s ∩ s = Ø",
			s:      map[int8]struct{}{1: {}, 4: {}},
			t:      map[int8]struct{}{2: {}, 6: {}},
			expect: false,
		},
		{
			name:   "test Int8 Equal, s ∩ s ≠ Ø && s ∩ s ≠ s",
			s:      map[int8]struct{}{1: {}, 4: {}},
			t:      map[int8]struct{}{1: {}, 6: {}},
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

func TestInt8_Copy(t *testing.T) {
	testcases := []struct {
		name   string
		s      Int8
		expect []int8
	}{
		{
			name:   "test Int8 Copy, s is empty",
			s:      Int8{},
			expect: []int8{},
		},
		{
			name:   "test Int8 Copy, s is not empty",
			s:      map[int8]struct{}{2: {}, 9: {}, 4: {}},
			expect: []int8{2, 9, 4},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		actual := tc.s.Copy()
		validateInt8(t, actual, tc.expect)
	}
}

func TestInt8_String(t *testing.T) {
	testcases := []struct {
		name   string
		s      Int8
		expect string
	}{
		{
			name:   "test Int8 String, s is empty",
			s:      Int8{},
			expect: "[]",
		},
		{
			name:   "test Int8 String, s is not empty",
			s:      map[int8]struct{}{1: {}},
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

func validateInt8(t *testing.T, actual Int8, expect []int8) {
	if len(expect) != len(actual) {
		t.Errorf("expect set len: %d, but got: %d", len(expect), len(actual))
	}
	slice2String := func(i []int8) string {
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
