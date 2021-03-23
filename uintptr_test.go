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

func TestNewUintptr(t *testing.T) {
	testCases := []struct {
		name   string
		input  []uintptr
		expect []uintptr
	}{
		{
			name:   "test Uintptr New, s nothing",
			input:  nil,
			expect: nil,
		},
		{
			name:   "test Uintptr New, inputs multiple elements",
			input:  []uintptr{1, 2, 3},
			expect: []uintptr{1, 2, 3},
		},
	}
	for _, tc := range testCases {
		t.Logf("running scenario: %s", tc.name)
		actual := NewUintptr(tc.input...)
		validateUintptr(t, actual, tc.expect)
	}
}

func TestUintptr_Add(t *testing.T) {
	testcases := []struct {
		name   string
		input  []uintptr
		expect []uintptr
	}{
		{
			name:   "test Uintptr Add, inputs nothing",
			input:  []uintptr{},
			expect: []uintptr{},
		},
		{
			name:   "test Uintptr Add, inputs multiple elements",
			input:  []uintptr{1, 2, 3},
			expect: []uintptr{1, 2, 3},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		actual := Uintptr{}
		actual.Add(tc.input...)
		validateUintptr(t, actual, tc.expect)
	}
}

func TestUintptr_Remove(t *testing.T) {
	testcases := []struct {
		name   string
		s      Uintptr
		input  []uintptr
		expect []uintptr
	}{
		{
			name:   "test Uintptr Remove, inputs nothing",
			s:      map[uintptr]struct{}{1: {}, 2: {}, 3: {}},
			input:  []uintptr{},
			expect: []uintptr{1, 2, 3},
		},
		{
			name:   "test Uintptr Remove, inputs multiple exit elements",
			s:      map[uintptr]struct{}{1: {}, 2: {}, 3: {}},
			input:  []uintptr{1, 3},
			expect: []uintptr{2},
		},
		{
			name:   "test Uintptr Remove, inputs multiple non-exit elements",
			s:      map[uintptr]struct{}{1: {}, 2: {}, 3: {}},
			input:  []uintptr{0, 6, 7},
			expect: []uintptr{1, 2, 3},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		tc.s.Remove(tc.input...)
		validateUintptr(t, tc.s, tc.expect)
	}
}

func TestUintptr_Pop(t *testing.T) {
	testcases := []struct {
		name   string
		s      Uintptr
		expect bool
	}{
		{
			name:   "test Uintptr Pop, s is not empty",
			s:      map[uintptr]struct{}{1: {}, 2: {}, 3: {}},
			expect: true,
		},
		{
			name:   "test Uintptr Pop, s is empty",
			s:      map[uintptr]struct{}{},
			expect: false,
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		backup := make(map[uintptr]struct{})
		for k := range tc.s {
			backup[k] = struct{}{}
		}
		key, ok := tc.s.Pop()
		if ok != tc.expect {
			t.Errorf("expect ok: %v, but got: %v", tc.expect, ok)
		}
		delete(backup, key)
		var expect []uintptr
		for key := range backup {
			expect = append(expect, key)
		}
		validateUintptr(t, tc.s, expect)
	}
}

func TestUintptr_Size(t *testing.T) {
	testcases := []struct {
		name   string
		s      Uintptr
		expect int
	}{
		{
			name:   "test Uintptr Size, s is empty",
			s:      Uintptr{},
			expect: 0,
		},
		{
			name:   "test Uintptr Size, s is not empty",
			s:      map[uintptr]struct{}{1: {}, 2: {}, 3: {}},
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

func TestUintptr_IsEmpty(t *testing.T) {
	testcases := []struct {
		name   string
		s      Uintptr
		expect bool
	}{
		{
			name:   "test Uintptr Empty, s is empty",
			s:      Uintptr{},
			expect: true,
		},
		{
			name:   "test interface Empty, s is not empty",
			s:      map[uintptr]struct{}{1: {}, 2: {}, 3: {}},
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

func TestUintptr_Clear(t *testing.T) {
	testcases := []struct {
		name string
		s    Uintptr
	}{
		{
			name: "test Uintptr Clear, s is empty",
			s:    Uintptr{},
		},
		{
			name: "test Uintptr Clear, s is not empty",
			s:    map[uintptr]struct{}{1: {}, 2: {}, 3: {}},
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

func TestUintptr_Has(t *testing.T) {
	testcases := []struct {
		name   string
		s      Uintptr
		input  uintptr
		expect bool
	}{
		{
			name:   "test Uintptr Has, s has input element",
			s:      map[uintptr]struct{}{1: {}, 2: {}, 3: {}},
			input:  1,
			expect: true,
		},
		{
			name:   "test Uintptr Has, s does not have element",
			s:      map[uintptr]struct{}{1: {}, 2: {}, 3: {}},
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

func TestUintptr_HasAll(t *testing.T) {
	testcases := []struct {
		name   string
		s      Uintptr
		input  []uintptr
		expect bool
	}{
		{
			name:   "test Uintptr HasAll, set has all input elements",
			s:      map[uintptr]struct{}{1: {}, 2: {}, 3: {}},
			input:  []uintptr{1, 2},
			expect: true,
		},
		{
			name:   "test Uintptr HasAll, set does not have all input elements",
			s:      map[uintptr]struct{}{1: {}, 2: {}, 3: {}},
			input:  []uintptr{8, 9},
			expect: false,
		},
		{
			name:   "test Uintptr HasAll, set does not have all input elements, but exist elements in set",
			s:      map[uintptr]struct{}{1: {}, 2: {}, 3: {}},
			input:  []uintptr{1, 9},
			expect: false,
		},
		{
			name:   "test Uintptr HasAll, input empty",
			s:      map[uintptr]struct{}{1: {}, 2: {}, 3: {}},
			input:  []uintptr{},
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

func TestUintptr_HasAny(t *testing.T) {
	testcases := []struct {
		name   string
		s      Uintptr
		input  []uintptr
		expect bool
	}{
		{
			name:   "test Uintptr HasAny, s has all elements",
			s:      map[uintptr]struct{}{1: {}, 2: {}, 3: {}},
			input:  []uintptr{1, 2},
			expect: true,
		},
		{
			name:   "test Uintptr HasAny, s does not have all elements, but exist elements in set",
			s:      map[uintptr]struct{}{1: {}, 2: {}, 3: {}},
			input:  []uintptr{1, 9},
			expect: true,
		},
		{
			name:   "test Uintptr HasAny, s does not has all elements",
			s:      map[uintptr]struct{}{1: {}, 2: {}, 3: {}},
			input:  []uintptr{8, 9},
			expect: false,
		},
		{
			name:   "test Uintptr HasAny, input empty",
			s:      map[uintptr]struct{}{1: {}, 2: {}, 3: {}},
			input:  []uintptr{},
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

func TestUintptr_List(t *testing.T) {
	testcases := []struct {
		name   string
		s      Uintptr
		input  []uintptr
		expect bool
	}{
		{
			name: "test Uintptr List, s is empty",
			s:    Uintptr{},
		},
		{
			name: "test Uintptr List, s is not empty",
			s:    map[uintptr]struct{}{1: {}, 2: {}, 3: {}},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		actual := tc.s.List()
		validateUintptr(t, tc.s, actual)
	}
}

func TestUintptr_SortedList(t *testing.T) {
	testcases := []struct {
		name   string
		s      Uintptr
		expect []uintptr
	}{
		{
			name: "test Uintptr List, s is empty",
			s:    Uintptr{},
		},
		{
			name:   "test Uintptr SortedList, s is not empty",
			s:      map[uintptr]struct{}{1: {}, 2: {}, 3: {}},
			expect: []uintptr{1, 2, 3},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		actual := tc.s.SortedList(func(i, j uintptr) bool {
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

func TestUintptr_Each(t *testing.T) {
	testcases := []struct {
		name   string
		origin Uintptr
	}{
		{
			name:   "test Uintptr Each",
			origin: map[uintptr]struct{}{1: {}, 2: {}, 3: {}},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		var expect []uintptr
		tc.origin.Each(func(i uintptr) {
			expect = append(expect, i)
		})
		validateUintptr(t, tc.origin, expect)
	}
}

func TestUintptr_EachE(t *testing.T) {
	inputErr := errors.New("s error")
	testcases := []struct {
		name      string
		origin    Uintptr
		breakEach bool
		inputErr  error
		expectLen int
		expectErr error
	}{
		{
			name:      "test Uintptr EachE",
			origin:    map[uintptr]struct{}{1: {}, 2: {}, 3: {}},
			expectLen: 3,
		},
		{
			name:      "test Uintptr EachE, return break error",
			origin:    map[uintptr]struct{}{1: {}, 2: {}, 3: {}},
			breakEach: true,
		},
		{
			name:      "test Uintptr Each, returns error",
			origin:    map[uintptr]struct{}{1: {}, 2: {}, 3: {}},
			inputErr:  inputErr,
			expectErr: inputErr,
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		var expect []uintptr
		err := tc.origin.EachE(func(i uintptr) error {
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
			validateUintptr(t, tc.origin, expect)
		}
	}
}

func TestUintptr_Union(t *testing.T) {
	testcases := []struct {
		name   string
		s      Uintptr
		t      Uintptr
		expect []uintptr
	}{
		{
			name:   "test Uintptr Union, s and s are empty",
			s:      Uintptr{},
			t:      Uintptr{},
			expect: []uintptr{},
		},
		{
			name:   "test Uintptr Union, s is empty",
			s:      Uintptr{},
			t:      map[uintptr]struct{}{1: {}, 2: {}, 3: {}},
			expect: []uintptr{1, 2, 3},
		},
		{
			name:   "test Uintptr Union, s is empty",
			s:      map[uintptr]struct{}{1: {}, 2: {}, 3: {}},
			t:      Uintptr{},
			expect: []uintptr{1, 2, 3},
		},
		{
			name:   "test Uintptr Union, s has same element to s",
			s:      map[uintptr]struct{}{1: {}, 2: {}, 3: {}},
			t:      map[uintptr]struct{}{1: {}, 9: {}, 4: {}},
			expect: []uintptr{1, 2, 3, 4, 9},
		},
		{
			name:   "test Uintptr Union, s does not have same element to s",
			s:      map[uintptr]struct{}{1: {}, 2: {}, 3: {}},
			t:      map[uintptr]struct{}{6: {}, 7: {}, 8: {}},
			expect: []uintptr{1, 2, 3, 6, 7, 8},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		actual := tc.s.Union(tc.t)
		validateUintptr(t, actual, tc.expect)
	}
}

func TestUintptr_Difference(t *testing.T) {
	testcases := []struct {
		name   string
		s      Uintptr
		t      Uintptr
		expect []uintptr
	}{
		{
			name:   "test Uintptr Difference, s is empty",
			s:      Uintptr{},
			t:      map[uintptr]struct{}{2: {}, 9: {}, 4: {}},
			expect: []uintptr{},
		},
		{
			name:   "test Uintptr Difference, s is empty",
			s:      map[uintptr]struct{}{2: {}, 9: {}, 4: {}},
			t:      Uintptr{},
			expect: []uintptr{2, 9, 4},
		},
		{
			name:   "test Uintptr Difference, s ⊂ s",
			s:      map[uintptr]struct{}{2: {}, 9: {}},
			t:      map[uintptr]struct{}{2: {}, 9: {}, 4: {}},
			expect: []uintptr{},
		},
		{
			name:   "test Uintptr Difference, s ⊃ s",
			s:      map[uintptr]struct{}{2: {}, 9: {}, 4: {}},
			t:      map[uintptr]struct{}{2: {}},
			expect: []uintptr{9, 4},
		},
		{
			name:   "test Uintptr Difference, s = s",
			s:      map[uintptr]struct{}{2: {}, 9: {}, 4: {}},
			t:      map[uintptr]struct{}{2: {}, 9: {}, 4: {}},
			expect: []uintptr{},
		},
		{
			name:   "test Uintptr Difference, s ∩ s = Ø",
			s:      map[uintptr]struct{}{2: {}, 9: {}, 4: {}},
			t:      map[uintptr]struct{}{1: {}, 6: {}},
			expect: []uintptr{2, 9, 4},
		},
		{
			name:   "test Uintptr Difference, s ∩ s ≠ Ø",
			s:      map[uintptr]struct{}{2: {}, 9: {}, 4: {}},
			t:      map[uintptr]struct{}{2: {}, 9: {}, 5: {}},
			expect: []uintptr{4},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		actual := tc.s.Difference(tc.t)
		validateUintptr(t, actual, tc.expect)
	}
}

func TestUintptr_Intersection(t *testing.T) {
	testcases := []struct {
		name   string
		s      Uintptr
		t      Uintptr
		expect []uintptr
	}{
		{
			name:   "test Uintptr Intersection, s is empty",
			s:      Uintptr{},
			t:      map[uintptr]struct{}{2: {}, 9: {}, 4: {}},
			expect: []uintptr{},
		},
		{
			name:   "test Uintptr Intersection, s is empty",
			s:      map[uintptr]struct{}{2: {}, 9: {}, 4: {}},
			t:      Uintptr{},
			expect: []uintptr{},
		},
		{
			name:   "test Uintptr Intersection, s ⊂ s",
			s:      map[uintptr]struct{}{2: {}, 9: {}},
			t:      map[uintptr]struct{}{2: {}, 9: {}, 4: {}},
			expect: []uintptr{2, 9},
		},
		{
			name:   "test Uintptr Intersection, s ⊃ s",
			s:      map[uintptr]struct{}{2: {}, 9: {}, 4: {}},
			t:      map[uintptr]struct{}{2: {}, 9: {}},
			expect: []uintptr{2, 9},
		},
		{
			name:   "test Uintptr Intersection, s = s",
			s:      map[uintptr]struct{}{2: {}, 9: {}, 4: {}},
			t:      map[uintptr]struct{}{2: {}, 9: {}, 4: {}},
			expect: []uintptr{2, 9, 4},
		},
		{
			name:   "test Uintptr Intersection, s ∩ s = Ø",
			s:      map[uintptr]struct{}{1: {}, 4: {}},
			t:      map[uintptr]struct{}{2: {}, 6: {}},
			expect: []uintptr{},
		},
		{
			name:   "test Uintptr Intersection, s ∩ s ≠ Ø",
			s:      map[uintptr]struct{}{1: {}, 4: {}},
			t:      map[uintptr]struct{}{1: {}, 6: {}},
			expect: []uintptr{1},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		actual := tc.s.Intersection(tc.t)
		validateUintptr(t, actual, tc.expect)
	}
}

func TestUintptr_SymmetricDifference(t *testing.T) {
	testcases := []struct {
		name   string
		s      Uintptr
		t      Uintptr
		expect []uintptr
	}{
		{
			name:   "test Uintptr SymmetricDifference, s is empty",
			s:      Uintptr{},
			t:      map[uintptr]struct{}{2: {}, 9: {}, 4: {}},
			expect: []uintptr{2, 9, 4},
		},
		{
			name:   "test Uintptr SymmetricDifference, s is empty",
			s:      map[uintptr]struct{}{2: {}, 9: {}, 4: {}},
			t:      Uintptr{},
			expect: []uintptr{2, 9, 4},
		},
		{
			name:   "test Uintptr SymmetricDifference, s ⊂ s",
			s:      map[uintptr]struct{}{2: {}, 9: {}},
			t:      map[uintptr]struct{}{2: {}, 9: {}, 4: {}},
			expect: []uintptr{4},
		},
		{
			name:   "test Uintptr SymmetricDifference, s ⊃ s",
			s:      map[uintptr]struct{}{2: {}, 9: {}, 4: {}},
			t:      map[uintptr]struct{}{2: {}, 9: {}},
			expect: []uintptr{4},
		},
		{
			name:   "test Uintptr SymmetricDifference, s = s",
			s:      map[uintptr]struct{}{2: {}, 9: {}, 4: {}},
			t:      map[uintptr]struct{}{2: {}, 9: {}, 4: {}},
			expect: []uintptr{},
		},
		{
			name:   "test Uintptr SymmetricDifference, s ∩ s = Ø",
			s:      map[uintptr]struct{}{1: {}, 4: {}},
			t:      map[uintptr]struct{}{2: {}, 6: {}},
			expect: []uintptr{1, 2, 4, 6},
		},
		{
			name:   "test Uintptr SymmetricDifference, s ∩ s ≠ Ø",
			s:      map[uintptr]struct{}{1: {}, 4: {}},
			t:      map[uintptr]struct{}{1: {}, 6: {}},
			expect: []uintptr{4, 6},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		actual := tc.s.SymmetricDifference(tc.t)
		validateUintptr(t, actual, tc.expect)
	}
}

func TestUintptr_IsSubset(t *testing.T) {
	testcases := []struct {
		name   string
		s      Uintptr
		t      Uintptr
		expect bool
	}{
		{
			name:   "test Uintptr IsSubset, s is empty",
			s:      Uintptr{},
			t:      map[uintptr]struct{}{2: {}, 9: {}, 4: {}},
			expect: true,
		},
		{
			name:   "test Uintptr IsSubset, s is empty",
			s:      map[uintptr]struct{}{2: {}, 9: {}, 4: {}},
			t:      Uintptr{},
			expect: false,
		},
		{
			name:   "test Uintptr IsSubset, s ⊂ s",
			s:      map[uintptr]struct{}{2: {}, 9: {}},
			t:      map[uintptr]struct{}{2: {}, 9: {}, 4: {}},
			expect: true,
		},
		{
			name:   "test Uintptr IsSubset, s ⊃ s",
			s:      map[uintptr]struct{}{2: {}, 9: {}, 4: {}},
			t:      map[uintptr]struct{}{2: {}, 9: {}},
			expect: false,
		},
		{
			name:   "test Uintptr IsSubset, s = s",
			s:      map[uintptr]struct{}{2: {}, 9: {}, 4: {}},
			t:      map[uintptr]struct{}{2: {}, 9: {}, 4: {}},
			expect: true,
		},
		{
			name:   "test Uintptr IsSubset, s ∩ s = Ø",
			s:      map[uintptr]struct{}{1: {}, 4: {}},
			t:      map[uintptr]struct{}{2: {}, 6: {}},
			expect: false,
		},
		{
			name:   "test Uintptr IsSubset, s ∩ s ≠ Ø && s ∩ s ≠ s",
			s:      map[uintptr]struct{}{1: {}, 4: {}},
			t:      map[uintptr]struct{}{1: {}, 6: {}},
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

func TestUintptr_IsSuperset(t *testing.T) {
	testcases := []struct {
		name   string
		s      Uintptr
		t      Uintptr
		expect bool
	}{
		{
			name:   "test Uintptr IsSuperset, s is empty",
			s:      Uintptr{},
			t:      map[uintptr]struct{}{2: {}, 9: {}, 4: {}},
			expect: false,
		},
		{
			name:   "test Uintptr IsSuperset, s is empty",
			s:      map[uintptr]struct{}{2: {}, 9: {}, 4: {}},
			t:      Uintptr{},
			expect: true,
		},
		{
			name:   "test Uintptr IsSuperset, s ⊂ s",
			s:      map[uintptr]struct{}{2: {}, 9: {}},
			t:      map[uintptr]struct{}{2: {}, 9: {}, 4: {}},
			expect: false,
		},
		{
			name:   "test Uintptr IsSuperset, s ⊃ s",
			s:      map[uintptr]struct{}{2: {}, 9: {}, 4: {}},
			t:      map[uintptr]struct{}{2: {}, 9: {}},
			expect: true,
		},
		{
			name:   "test Uintptr IsSuperset, s = s",
			s:      map[uintptr]struct{}{2: {}, 9: {}, 4: {}},
			t:      map[uintptr]struct{}{2: {}, 9: {}, 4: {}},
			expect: true,
		},
		{
			name:   "test Uintptr IsSuperset, s ∩ s = Ø",
			s:      map[uintptr]struct{}{1: {}, 4: {}},
			t:      map[uintptr]struct{}{2: {}, 6: {}},
			expect: false,
		},
		{
			name:   "test Uintptr IsSuperset, s ∩ s ≠ Ø && s ∩ s ≠ s",
			s:      map[uintptr]struct{}{1: {}, 4: {}},
			t:      map[uintptr]struct{}{1: {}, 6: {}},
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

func TestUintptr_Equal(t *testing.T) {
	testcases := []struct {
		name   string
		s      Uintptr
		t      Uintptr
		expect bool
	}{
		{
			name:   "test Uintptr Equal, s is empty",
			s:      Uintptr{},
			t:      map[uintptr]struct{}{2: {}, 9: {}, 4: {}},
			expect: false,
		},
		{
			name:   "test Uintptr Equal, s is empty",
			s:      map[uintptr]struct{}{2: {}, 9: {}, 4: {}},
			t:      Uintptr{},
			expect: false,
		},
		{
			name:   "test Uintptr Equal, s ⊂ s",
			s:      map[uintptr]struct{}{2: {}, 9: {}},
			t:      map[uintptr]struct{}{2: {}, 9: {}, 4: {}},
			expect: false,
		},
		{
			name:   "test Uintptr Equal, s ⊃ s",
			s:      map[uintptr]struct{}{2: {}, 9: {}, 4: {}},
			t:      map[uintptr]struct{}{2: {}, 9: {}},
			expect: false,
		},
		{
			name:   "test Uintptr Equal, s = s",
			s:      map[uintptr]struct{}{2: {}, 9: {}, 4: {}},
			t:      map[uintptr]struct{}{2: {}, 9: {}, 4: {}},
			expect: true,
		},
		{
			name:   "test Uintptr Equal, s ∩ s = Ø",
			s:      map[uintptr]struct{}{1: {}, 4: {}},
			t:      map[uintptr]struct{}{2: {}, 6: {}},
			expect: false,
		},
		{
			name:   "test Uintptr Equal, s ∩ s ≠ Ø && s ∩ s ≠ s",
			s:      map[uintptr]struct{}{1: {}, 4: {}},
			t:      map[uintptr]struct{}{1: {}, 6: {}},
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

func TestUintptr_Copy(t *testing.T) {
	testcases := []struct {
		name   string
		s      Uintptr
		expect []uintptr
	}{
		{
			name:   "test Uintptr Copy, s is empty",
			s:      Uintptr{},
			expect: []uintptr{},
		},
		{
			name:   "test Uintptr Copy, s is not empty",
			s:      map[uintptr]struct{}{2: {}, 9: {}, 4: {}},
			expect: []uintptr{2, 9, 4},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		actual := tc.s.Copy()
		validateUintptr(t, actual, tc.expect)
	}
}

func TestUintptr_String(t *testing.T) {
	testcases := []struct {
		name   string
		s      Uintptr
		expect string
	}{
		{
			name:   "test Uintptr String, s is empty",
			s:      Uintptr{},
			expect: "[]",
		},
		{
			name:   "test Uintptr String, s is not empty",
			s:      map[uintptr]struct{}{1: {}},
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

func validateUintptr(t *testing.T, actual Uintptr, expect []uintptr) {
	if len(expect) != len(actual) {
		t.Errorf("expect set len: %d, but got: %d", len(expect), len(actual))
	}
	slice2String := func(i []uintptr) string {
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
