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

func TestNewUint64(t *testing.T) {
	testCases := []struct {
		name   string
		input  []uint64
		expect []uint64
	}{
		{
			name:   "test Uint64 New, s nothing",
			input:  nil,
			expect: nil,
		},
		{
			name:   "test Uint64 New, inputs multiple elements",
			input:  []uint64{1, 2, 3},
			expect: []uint64{1, 2, 3},
		},
	}
	for _, tc := range testCases {
		t.Logf("running scenario: %s", tc.name)
		actual := NewUint64(tc.input...)
		validateUint64(t, actual, tc.expect)
	}
}

func TestUint64_Add(t *testing.T) {
	testcases := []struct {
		name   string
		input  []uint64
		expect []uint64
	}{
		{
			name:   "test Uint64 Add, inputs nothing",
			input:  []uint64{},
			expect: []uint64{},
		},
		{
			name:   "test Uint64 Add, inputs multiple elements",
			input:  []uint64{1, 2, 3},
			expect: []uint64{1, 2, 3},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		actual := Uint64{}
		actual.Add(tc.input...)
		validateUint64(t, actual, tc.expect)
	}
}

func TestUint64_Remove(t *testing.T) {
	testcases := []struct {
		name   string
		s      Uint64
		input  []uint64
		expect []uint64
	}{
		{
			name:   "test Uint64 Remove, inputs nothing",
			s:      map[uint64]struct{}{1: {}, 2: {}, 3: {}},
			input:  []uint64{},
			expect: []uint64{1, 2, 3},
		},
		{
			name:   "test Uint64 Remove, inputs multiple exit elements",
			s:      map[uint64]struct{}{1: {}, 2: {}, 3: {}},
			input:  []uint64{1, 3},
			expect: []uint64{2},
		},
		{
			name:   "test Uint64 Remove, inputs multiple non-exit elements",
			s:      map[uint64]struct{}{1: {}, 2: {}, 3: {}},
			input:  []uint64{0, 6, 7},
			expect: []uint64{1, 2, 3},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		tc.s.Remove(tc.input...)
		validateUint64(t, tc.s, tc.expect)
	}
}

func TestUint64_Pop(t *testing.T) {
	testcases := []struct {
		name   string
		s      Uint64
		expect bool
	}{
		{
			name:   "test Uint64 Pop, s is not empty",
			s:      map[uint64]struct{}{1: {}, 2: {}, 3: {}},
			expect: true,
		},
		{
			name:   "test Uint64 Pop, s is empty",
			s:      map[uint64]struct{}{},
			expect: false,
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		backup := make(map[uint64]struct{})
		for k := range tc.s {
			backup[k] = struct{}{}
		}
		key, ok := tc.s.Pop()
		if ok != tc.expect {
			t.Errorf("expect ok: %v, but got: %v", tc.expect, ok)
		}
		delete(backup, key)
		var expect []uint64
		for key := range backup {
			expect = append(expect, key)
		}
		validateUint64(t, tc.s, expect)
	}
}

func TestUint64_Size(t *testing.T) {
	testcases := []struct {
		name   string
		s      Uint64
		expect int
	}{
		{
			name:   "test Uint64 Size, s is empty",
			s:      Uint64{},
			expect: 0,
		},
		{
			name:   "test Uint64 Size, s is not empty",
			s:      map[uint64]struct{}{1: {}, 2: {}, 3: {}},
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

func TestUint64_IsEmpty(t *testing.T) {
	testcases := []struct {
		name   string
		s      Uint64
		expect bool
	}{
		{
			name:   "test Uint64 Empty, s is empty",
			s:      Uint64{},
			expect: true,
		},
		{
			name:   "test interface Empty, s is not empty",
			s:      map[uint64]struct{}{1: {}, 2: {}, 3: {}},
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

func TestUint64_Clear(t *testing.T) {
	testcases := []struct {
		name string
		s    Uint64
	}{
		{
			name: "test Uint64 Clear, s is empty",
			s:    Uint64{},
		},
		{
			name: "test Uint64 Clear, s is not empty",
			s:    map[uint64]struct{}{1: {}, 2: {}, 3: {}},
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

func TestUint64_Has(t *testing.T) {
	testcases := []struct {
		name   string
		s      Uint64
		input  uint64
		expect bool
	}{
		{
			name:   "test Uint64 Has, s has input element",
			s:      map[uint64]struct{}{1: {}, 2: {}, 3: {}},
			input:  1,
			expect: true,
		},
		{
			name:   "test Uint64 Has, s does not have element",
			s:      map[uint64]struct{}{1: {}, 2: {}, 3: {}},
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

func TestUint64_HasAll(t *testing.T) {
	testcases := []struct {
		name   string
		s      Uint64
		input  []uint64
		expect bool
	}{
		{
			name:   "test Uint64 HasAll, set has all input elements",
			s:      map[uint64]struct{}{1: {}, 2: {}, 3: {}},
			input:  []uint64{1, 2},
			expect: true,
		},
		{
			name:   "test Uint64 HasAll, set does not have all input elements",
			s:      map[uint64]struct{}{1: {}, 2: {}, 3: {}},
			input:  []uint64{8, 9},
			expect: false,
		},
		{
			name:   "test Uint64 HasAll, set does not have all input elements, but exist elements in set",
			s:      map[uint64]struct{}{1: {}, 2: {}, 3: {}},
			input:  []uint64{1, 9},
			expect: false,
		},
		{
			name:   "test Uint64 HasAll, input empty",
			s:      map[uint64]struct{}{1: {}, 2: {}, 3: {}},
			input:  []uint64{},
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

func TestUint64_HasAny(t *testing.T) {
	testcases := []struct {
		name   string
		s      Uint64
		input  []uint64
		expect bool
	}{
		{
			name:   "test Uint64 HasAny, s has all elements",
			s:      map[uint64]struct{}{1: {}, 2: {}, 3: {}},
			input:  []uint64{1, 2},
			expect: true,
		},
		{
			name:   "test Uint64 HasAny, s does not have all elements, but exist elements in set",
			s:      map[uint64]struct{}{1: {}, 2: {}, 3: {}},
			input:  []uint64{1, 9},
			expect: true,
		},
		{
			name:   "test Uint64 HasAny, s does not has all elements",
			s:      map[uint64]struct{}{1: {}, 2: {}, 3: {}},
			input:  []uint64{8, 9},
			expect: false,
		},
		{
			name:   "test Uint64 HasAny, input empty",
			s:      map[uint64]struct{}{1: {}, 2: {}, 3: {}},
			input:  []uint64{},
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

func TestUint64_List(t *testing.T) {
	testcases := []struct {
		name   string
		s      Uint64
		input  []uint64
		expect bool
	}{
		{
			name: "test Uint64 List, s is empty",
			s:    Uint64{},
		},
		{
			name: "test Uint64 List, s is not empty",
			s:    map[uint64]struct{}{1: {}, 2: {}, 3: {}},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		actual := tc.s.List()
		validateUint64(t, tc.s, actual)
	}
}

func TestUint64_SortedList(t *testing.T) {
	testcases := []struct {
		name   string
		s      Uint64
		expect []uint64
	}{
		{
			name: "test Uint64 List, s is empty",
			s:    Uint64{},
		},
		{
			name:   "test Uint64 SortedList, s is not empty",
			s:      map[uint64]struct{}{1: {}, 2: {}, 3: {}},
			expect: []uint64{1, 2, 3},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		actual := tc.s.SortedList(func(i, j uint64) bool {
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

func TestUint64_Each(t *testing.T) {
	testcases := []struct {
		name   string
		origin Uint64
	}{
		{
			name:   "test Uint64 Each",
			origin: map[uint64]struct{}{1: {}, 2: {}, 3: {}},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		var expect []uint64
		tc.origin.Each(func(i uint64) {
			expect = append(expect, i)
		})
		validateUint64(t, tc.origin, expect)
	}
}

func TestUint64_EachE(t *testing.T) {
	inputErr := errors.New("s error")
	testcases := []struct {
		name      string
		origin    Uint64
		breakEach bool
		inputErr  error
		expectLen int
		expectErr error
	}{
		{
			name:      "test Uint64 EachE",
			origin:    map[uint64]struct{}{1: {}, 2: {}, 3: {}},
			expectLen: 3,
		},
		{
			name:      "test Uint64 EachE, return break error",
			origin:    map[uint64]struct{}{1: {}, 2: {}, 3: {}},
			breakEach: true,
		},
		{
			name:      "test Uint64 Each, returns error",
			origin:    map[uint64]struct{}{1: {}, 2: {}, 3: {}},
			inputErr:  inputErr,
			expectErr: inputErr,
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		var expect []uint64
		err := tc.origin.EachE(func(i uint64) error {
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
			validateUint64(t, tc.origin, expect)
		}
	}
}

func TestUint64_Union(t *testing.T) {
	testcases := []struct {
		name   string
		s      Uint64
		t      Uint64
		expect []uint64
	}{
		{
			name:   "test Uint64 Union, s and s are empty",
			s:      Uint64{},
			t:      Uint64{},
			expect: []uint64{},
		},
		{
			name:   "test Uint64 Union, s is empty",
			s:      Uint64{},
			t:      map[uint64]struct{}{1: {}, 2: {}, 3: {}},
			expect: []uint64{1, 2, 3},
		},
		{
			name:   "test Uint64 Union, s is empty",
			s:      map[uint64]struct{}{1: {}, 2: {}, 3: {}},
			t:      Uint64{},
			expect: []uint64{1, 2, 3},
		},
		{
			name:   "test Uint64 Union, s has same element to s",
			s:      map[uint64]struct{}{1: {}, 2: {}, 3: {}},
			t:      map[uint64]struct{}{1: {}, 9: {}, 4: {}},
			expect: []uint64{1, 2, 3, 4, 9},
		},
		{
			name:   "test Uint64 Union, s does not have same element to s",
			s:      map[uint64]struct{}{1: {}, 2: {}, 3: {}},
			t:      map[uint64]struct{}{6: {}, 7: {}, 8: {}},
			expect: []uint64{1, 2, 3, 6, 7, 8},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		actual := tc.s.Union(tc.t)
		validateUint64(t, actual, tc.expect)
	}
}

func TestUint64_Difference(t *testing.T) {
	testcases := []struct {
		name   string
		s      Uint64
		t      Uint64
		expect []uint64
	}{
		{
			name:   "test Uint64 Difference, s is empty",
			s:      Uint64{},
			t:      map[uint64]struct{}{2: {}, 9: {}, 4: {}},
			expect: []uint64{},
		},
		{
			name:   "test Uint64 Difference, s is empty",
			s:      map[uint64]struct{}{2: {}, 9: {}, 4: {}},
			t:      Uint64{},
			expect: []uint64{2, 9, 4},
		},
		{
			name:   "test Uint64 Difference, s ⊂ s",
			s:      map[uint64]struct{}{2: {}, 9: {}},
			t:      map[uint64]struct{}{2: {}, 9: {}, 4: {}},
			expect: []uint64{},
		},
		{
			name:   "test Uint64 Difference, s ⊃ s",
			s:      map[uint64]struct{}{2: {}, 9: {}, 4: {}},
			t:      map[uint64]struct{}{2: {}},
			expect: []uint64{9, 4},
		},
		{
			name:   "test Uint64 Difference, s = s",
			s:      map[uint64]struct{}{2: {}, 9: {}, 4: {}},
			t:      map[uint64]struct{}{2: {}, 9: {}, 4: {}},
			expect: []uint64{},
		},
		{
			name:   "test Uint64 Difference, s ∩ s = Ø",
			s:      map[uint64]struct{}{2: {}, 9: {}, 4: {}},
			t:      map[uint64]struct{}{1: {}, 6: {}},
			expect: []uint64{2, 9, 4},
		},
		{
			name:   "test Uint64 Difference, s ∩ s ≠ Ø",
			s:      map[uint64]struct{}{2: {}, 9: {}, 4: {}},
			t:      map[uint64]struct{}{2: {}, 9: {}, 5: {}},
			expect: []uint64{4},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		actual := tc.s.Difference(tc.t)
		validateUint64(t, actual, tc.expect)
	}
}

func TestUint64_Intersection(t *testing.T) {
	testcases := []struct {
		name   string
		s      Uint64
		t      Uint64
		expect []uint64
	}{
		{
			name:   "test Uint64 Intersection, s is empty",
			s:      Uint64{},
			t:      map[uint64]struct{}{2: {}, 9: {}, 4: {}},
			expect: []uint64{},
		},
		{
			name:   "test Uint64 Intersection, s is empty",
			s:      map[uint64]struct{}{2: {}, 9: {}, 4: {}},
			t:      Uint64{},
			expect: []uint64{},
		},
		{
			name:   "test Uint64 Intersection, s ⊂ s",
			s:      map[uint64]struct{}{2: {}, 9: {}},
			t:      map[uint64]struct{}{2: {}, 9: {}, 4: {}},
			expect: []uint64{2, 9},
		},
		{
			name:   "test Uint64 Intersection, s ⊃ s",
			s:      map[uint64]struct{}{2: {}, 9: {}, 4: {}},
			t:      map[uint64]struct{}{2: {}, 9: {}},
			expect: []uint64{2, 9},
		},
		{
			name:   "test Uint64 Intersection, s = s",
			s:      map[uint64]struct{}{2: {}, 9: {}, 4: {}},
			t:      map[uint64]struct{}{2: {}, 9: {}, 4: {}},
			expect: []uint64{2, 9, 4},
		},
		{
			name:   "test Uint64 Intersection, s ∩ s = Ø",
			s:      map[uint64]struct{}{1: {}, 4: {}},
			t:      map[uint64]struct{}{2: {}, 6: {}},
			expect: []uint64{},
		},
		{
			name:   "test Uint64 Intersection, s ∩ s ≠ Ø",
			s:      map[uint64]struct{}{1: {}, 4: {}},
			t:      map[uint64]struct{}{1: {}, 6: {}},
			expect: []uint64{1},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		actual := tc.s.Intersection(tc.t)
		validateUint64(t, actual, tc.expect)
	}
}

func TestUint64_SymmetricDifference(t *testing.T) {
	testcases := []struct {
		name   string
		s      Uint64
		t      Uint64
		expect []uint64
	}{
		{
			name:   "test Uint64 SymmetricDifference, s is empty",
			s:      Uint64{},
			t:      map[uint64]struct{}{2: {}, 9: {}, 4: {}},
			expect: []uint64{2, 9, 4},
		},
		{
			name:   "test Uint64 SymmetricDifference, s is empty",
			s:      map[uint64]struct{}{2: {}, 9: {}, 4: {}},
			t:      Uint64{},
			expect: []uint64{2, 9, 4},
		},
		{
			name:   "test Uint64 SymmetricDifference, s ⊂ s",
			s:      map[uint64]struct{}{2: {}, 9: {}},
			t:      map[uint64]struct{}{2: {}, 9: {}, 4: {}},
			expect: []uint64{4},
		},
		{
			name:   "test Uint64 SymmetricDifference, s ⊃ s",
			s:      map[uint64]struct{}{2: {}, 9: {}, 4: {}},
			t:      map[uint64]struct{}{2: {}, 9: {}},
			expect: []uint64{4},
		},
		{
			name:   "test Uint64 SymmetricDifference, s = s",
			s:      map[uint64]struct{}{2: {}, 9: {}, 4: {}},
			t:      map[uint64]struct{}{2: {}, 9: {}, 4: {}},
			expect: []uint64{},
		},
		{
			name:   "test Uint64 SymmetricDifference, s ∩ s = Ø",
			s:      map[uint64]struct{}{1: {}, 4: {}},
			t:      map[uint64]struct{}{2: {}, 6: {}},
			expect: []uint64{1, 2, 4, 6},
		},
		{
			name:   "test Uint64 SymmetricDifference, s ∩ s ≠ Ø",
			s:      map[uint64]struct{}{1: {}, 4: {}},
			t:      map[uint64]struct{}{1: {}, 6: {}},
			expect: []uint64{4, 6},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		actual := tc.s.SymmetricDifference(tc.t)
		validateUint64(t, actual, tc.expect)
	}
}

func TestUint64_IsSubset(t *testing.T) {
	testcases := []struct {
		name   string
		s      Uint64
		t      Uint64
		expect bool
	}{
		{
			name:   "test Uint64 IsSubset, s is empty",
			s:      Uint64{},
			t:      map[uint64]struct{}{2: {}, 9: {}, 4: {}},
			expect: true,
		},
		{
			name:   "test Uint64 IsSubset, s is empty",
			s:      map[uint64]struct{}{2: {}, 9: {}, 4: {}},
			t:      Uint64{},
			expect: false,
		},
		{
			name:   "test Uint64 IsSubset, s ⊂ s",
			s:      map[uint64]struct{}{2: {}, 9: {}},
			t:      map[uint64]struct{}{2: {}, 9: {}, 4: {}},
			expect: true,
		},
		{
			name:   "test Uint64 IsSubset, s ⊃ s",
			s:      map[uint64]struct{}{2: {}, 9: {}, 4: {}},
			t:      map[uint64]struct{}{2: {}, 9: {}},
			expect: false,
		},
		{
			name:   "test Uint64 IsSubset, s = s",
			s:      map[uint64]struct{}{2: {}, 9: {}, 4: {}},
			t:      map[uint64]struct{}{2: {}, 9: {}, 4: {}},
			expect: true,
		},
		{
			name:   "test Uint64 IsSubset, s ∩ s = Ø",
			s:      map[uint64]struct{}{1: {}, 4: {}},
			t:      map[uint64]struct{}{2: {}, 6: {}},
			expect: false,
		},
		{
			name:   "test Uint64 IsSubset, s ∩ s ≠ Ø && s ∩ s ≠ s",
			s:      map[uint64]struct{}{1: {}, 4: {}},
			t:      map[uint64]struct{}{1: {}, 6: {}},
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

func TestUint64_IsSuperset(t *testing.T) {
	testcases := []struct {
		name   string
		s      Uint64
		t      Uint64
		expect bool
	}{
		{
			name:   "test Uint64 IsSuperset, s is empty",
			s:      Uint64{},
			t:      map[uint64]struct{}{2: {}, 9: {}, 4: {}},
			expect: false,
		},
		{
			name:   "test Uint64 IsSuperset, s is empty",
			s:      map[uint64]struct{}{2: {}, 9: {}, 4: {}},
			t:      Uint64{},
			expect: true,
		},
		{
			name:   "test Uint64 IsSuperset, s ⊂ s",
			s:      map[uint64]struct{}{2: {}, 9: {}},
			t:      map[uint64]struct{}{2: {}, 9: {}, 4: {}},
			expect: false,
		},
		{
			name:   "test Uint64 IsSuperset, s ⊃ s",
			s:      map[uint64]struct{}{2: {}, 9: {}, 4: {}},
			t:      map[uint64]struct{}{2: {}, 9: {}},
			expect: true,
		},
		{
			name:   "test Uint64 IsSuperset, s = s",
			s:      map[uint64]struct{}{2: {}, 9: {}, 4: {}},
			t:      map[uint64]struct{}{2: {}, 9: {}, 4: {}},
			expect: true,
		},
		{
			name:   "test Uint64 IsSuperset, s ∩ s = Ø",
			s:      map[uint64]struct{}{1: {}, 4: {}},
			t:      map[uint64]struct{}{2: {}, 6: {}},
			expect: false,
		},
		{
			name:   "test Uint64 IsSuperset, s ∩ s ≠ Ø && s ∩ s ≠ s",
			s:      map[uint64]struct{}{1: {}, 4: {}},
			t:      map[uint64]struct{}{1: {}, 6: {}},
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

func TestUint64_Equal(t *testing.T) {
	testcases := []struct {
		name   string
		s      Uint64
		t      Uint64
		expect bool
	}{
		{
			name:   "test Uint64 Equal, s is empty",
			s:      Uint64{},
			t:      map[uint64]struct{}{2: {}, 9: {}, 4: {}},
			expect: false,
		},
		{
			name:   "test Uint64 Equal, s is empty",
			s:      map[uint64]struct{}{2: {}, 9: {}, 4: {}},
			t:      Uint64{},
			expect: false,
		},
		{
			name:   "test Uint64 Equal, s ⊂ s",
			s:      map[uint64]struct{}{2: {}, 9: {}},
			t:      map[uint64]struct{}{2: {}, 9: {}, 4: {}},
			expect: false,
		},
		{
			name:   "test Uint64 Equal, s ⊃ s",
			s:      map[uint64]struct{}{2: {}, 9: {}, 4: {}},
			t:      map[uint64]struct{}{2: {}, 9: {}},
			expect: false,
		},
		{
			name:   "test Uint64 Equal, s = s",
			s:      map[uint64]struct{}{2: {}, 9: {}, 4: {}},
			t:      map[uint64]struct{}{2: {}, 9: {}, 4: {}},
			expect: true,
		},
		{
			name:   "test Uint64 Equal, s ∩ s = Ø",
			s:      map[uint64]struct{}{1: {}, 4: {}},
			t:      map[uint64]struct{}{2: {}, 6: {}},
			expect: false,
		},
		{
			name:   "test Uint64 Equal, s ∩ s ≠ Ø && s ∩ s ≠ s",
			s:      map[uint64]struct{}{1: {}, 4: {}},
			t:      map[uint64]struct{}{1: {}, 6: {}},
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

func TestUint64_Copy(t *testing.T) {
	testcases := []struct {
		name   string
		s      Uint64
		expect []uint64
	}{
		{
			name:   "test Uint64 Copy, s is empty",
			s:      Uint64{},
			expect: []uint64{},
		},
		{
			name:   "test Uint64 Copy, s is not empty",
			s:      map[uint64]struct{}{2: {}, 9: {}, 4: {}},
			expect: []uint64{2, 9, 4},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		actual := tc.s.Copy()
		validateUint64(t, actual, tc.expect)
	}
}

func TestUint64_String(t *testing.T) {
	testcases := []struct {
		name   string
		s      Uint64
		expect string
	}{
		{
			name:   "test Uint64 String, s is empty",
			s:      Uint64{},
			expect: "[]",
		},
		{
			name:   "test Uint64 String, s is not empty",
			s:      map[uint64]struct{}{1: {}},
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

func validateUint64(t *testing.T, actual Uint64, expect []uint64) {
	if len(expect) != len(actual) {
		t.Errorf("expect set len: %d, but got: %d", len(expect), len(actual))
	}
	slice2String := func(i []uint64) string {
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
