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

func TestNewUint(t *testing.T) {
	testCases := []struct {
		name   string
		input  []uint
		expect []uint
	}{
		{
			name:   "test Uint New, s nothing",
			input:  nil,
			expect: nil,
		},
		{
			name:   "test Uint New, inputs multiple elements",
			input:  []uint{1, 2, 3},
			expect: []uint{1, 2, 3},
		},
	}
	for _, tc := range testCases {
		t.Logf("running scenario: %s", tc.name)
		actual := NewUint(tc.input...)
		validateUint(t, actual, tc.expect)
	}
}

func TestUint_Add(t *testing.T) {
	testcases := []struct {
		name   string
		input  []uint
		expect []uint
	}{
		{
			name:   "test Uint Add, inputs nothing",
			input:  []uint{},
			expect: []uint{},
		},
		{
			name:   "test Uint Add, inputs multiple elements",
			input:  []uint{1, 2, 3},
			expect: []uint{1, 2, 3},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		actual := Uint{}
		actual.Add(tc.input...)
		validateUint(t, actual, tc.expect)
	}
}

func TestUint_Remove(t *testing.T) {
	testcases := []struct {
		name   string
		s      Uint
		input  []uint
		expect []uint
	}{
		{
			name:   "test Uint Remove, inputs nothing",
			s:      map[uint]struct{}{1: {}, 2: {}, 3: {}},
			input:  []uint{},
			expect: []uint{1, 2, 3},
		},
		{
			name:   "test Uint Remove, inputs multiple exit elements",
			s:      map[uint]struct{}{1: {}, 2: {}, 3: {}},
			input:  []uint{1, 3},
			expect: []uint{2},
		},
		{
			name:   "test Uint Remove, inputs multiple non-exit elements",
			s:      map[uint]struct{}{1: {}, 2: {}, 3: {}},
			input:  []uint{0, 6, 7},
			expect: []uint{1, 2, 3},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		tc.s.Remove(tc.input...)
		validateUint(t, tc.s, tc.expect)
	}
}

func TestUint_Pop(t *testing.T) {
	testcases := []struct {
		name   string
		s      Uint
		expect bool
	}{
		{
			name:   "test Uint Pop, s is not empty",
			s:      map[uint]struct{}{1: {}, 2: {}, 3: {}},
			expect: true,
		},
		{
			name:   "test Uint Pop, s is empty",
			s:      map[uint]struct{}{},
			expect: false,
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		backup := make(map[uint]struct{})
		for k := range tc.s {
			backup[k] = struct{}{}
		}
		key, ok := tc.s.Pop()
		if ok != tc.expect {
			t.Errorf("expect ok: %v, but got: %v", tc.expect, ok)
		}
		delete(backup, key)
		var expect []uint
		for key := range backup {
			expect = append(expect, key)
		}
		validateUint(t, tc.s, expect)
	}
}

func TestUint_Size(t *testing.T) {
	testcases := []struct {
		name   string
		s      Uint
		expect int
	}{
		{
			name:   "test Uint Size, s is empty",
			s:      Uint{},
			expect: 0,
		},
		{
			name:   "test Uint Size, s is not empty",
			s:      map[uint]struct{}{1: {}, 2: {}, 3: {}},
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

func TestUint_IsEmpty(t *testing.T) {
	testcases := []struct {
		name   string
		s      Uint
		expect bool
	}{
		{
			name:   "test Uint Empty, s is empty",
			s:      Uint{},
			expect: true,
		},
		{
			name:   "test interface Empty, s is not empty",
			s:      map[uint]struct{}{1: {}, 2: {}, 3: {}},
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

func TestUint_Clear(t *testing.T) {
	testcases := []struct {
		name string
		s    Uint
	}{
		{
			name: "test Uint Clear, s is empty",
			s:    Uint{},
		},
		{
			name: "test Uint Clear, s is not empty",
			s:    map[uint]struct{}{1: {}, 2: {}, 3: {}},
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

func TestUint_Has(t *testing.T) {
	testcases := []struct {
		name   string
		s      Uint
		input  uint
		expect bool
	}{
		{
			name:   "test Uint Has, s has input element",
			s:      map[uint]struct{}{1: {}, 2: {}, 3: {}},
			input:  1,
			expect: true,
		},
		{
			name:   "test Uint Has, s does not have element",
			s:      map[uint]struct{}{1: {}, 2: {}, 3: {}},
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

func TestUint_HasAll(t *testing.T) {
	testcases := []struct {
		name   string
		s      Uint
		input  []uint
		expect bool
	}{
		{
			name:   "test Uint HasAll, set has all input elements",
			s:      map[uint]struct{}{1: {}, 2: {}, 3: {}},
			input:  []uint{1, 2},
			expect: true,
		},
		{
			name:   "test Uint HasAll, set does not have all input elements",
			s:      map[uint]struct{}{1: {}, 2: {}, 3: {}},
			input:  []uint{8, 9},
			expect: false,
		},
		{
			name:   "test Uint HasAll, set does not have all input elements, but exist elements in set",
			s:      map[uint]struct{}{1: {}, 2: {}, 3: {}},
			input:  []uint{1, 9},
			expect: false,
		},
		{
			name:   "test Uint HasAll, input empty",
			s:      map[uint]struct{}{1: {}, 2: {}, 3: {}},
			input:  []uint{},
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

func TestUint_HasAny(t *testing.T) {
	testcases := []struct {
		name   string
		s      Uint
		input  []uint
		expect bool
	}{
		{
			name:   "test Uint HasAny, s has all elements",
			s:      map[uint]struct{}{1: {}, 2: {}, 3: {}},
			input:  []uint{1, 2},
			expect: true,
		},
		{
			name:   "test Uint HasAny, s does not have all elements, but exist elements in set",
			s:      map[uint]struct{}{1: {}, 2: {}, 3: {}},
			input:  []uint{1, 9},
			expect: true,
		},
		{
			name:   "test Uint HasAny, s does not has all elements",
			s:      map[uint]struct{}{1: {}, 2: {}, 3: {}},
			input:  []uint{8, 9},
			expect: false,
		},
		{
			name:   "test Uint HasAny, input empty",
			s:      map[uint]struct{}{1: {}, 2: {}, 3: {}},
			input:  []uint{},
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

func TestUint_List(t *testing.T) {
	testcases := []struct {
		name   string
		s      Uint
		input  []uint
		expect bool
	}{
		{
			name: "test Uint List, s is empty",
			s:    Uint{},
		},
		{
			name: "test Uint List, s is not empty",
			s:    map[uint]struct{}{1: {}, 2: {}, 3: {}},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		actual := tc.s.List()
		validateUint(t, tc.s, actual)
	}
}

func TestUint_SortedList(t *testing.T) {
	testcases := []struct {
		name   string
		s      Uint
		expect []uint
	}{
		{
			name: "test Uint List, s is empty",
			s:    Uint{},
		},
		{
			name:   "test Uint SortedList, s is not empty",
			s:      map[uint]struct{}{1: {}, 2: {}, 3: {}},
			expect: []uint{1, 2, 3},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		actual := tc.s.SortedList(func(i, j uint) bool {
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

func TestUint_Each(t *testing.T) {
	testcases := []struct {
		name   string
		origin Uint
	}{
		{
			name:   "test Uint Each",
			origin: map[uint]struct{}{1: {}, 2: {}, 3: {}},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		var expect []uint
		tc.origin.Each(func(i uint) {
			expect = append(expect, i)
		})
		validateUint(t, tc.origin, expect)
	}
}

func TestUint_EachE(t *testing.T) {
	inputErr := errors.New("s error")
	testcases := []struct {
		name      string
		origin    Uint
		breakEach bool
		inputErr  error
		expectLen int
		expectErr error
	}{
		{
			name:      "test Uint EachE",
			origin:    map[uint]struct{}{1: {}, 2: {}, 3: {}},
			expectLen: 3,
		},
		{
			name:      "test Uint EachE, return break error",
			origin:    map[uint]struct{}{1: {}, 2: {}, 3: {}},
			breakEach: true,
		},
		{
			name:      "test Uint Each, returns error",
			origin:    map[uint]struct{}{1: {}, 2: {}, 3: {}},
			inputErr:  inputErr,
			expectErr: inputErr,
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		var expect []uint
		err := tc.origin.EachE(func(i uint) error {
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
			validateUint(t, tc.origin, expect)
		}
	}
}

func TestUint_Union(t *testing.T) {
	testcases := []struct {
		name   string
		s      Uint
		t      Uint
		expect []uint
	}{
		{
			name:   "test Uint Union, s and s are empty",
			s:      Uint{},
			t:      Uint{},
			expect: []uint{},
		},
		{
			name:   "test Uint Union, s is empty",
			s:      Uint{},
			t:      map[uint]struct{}{1: {}, 2: {}, 3: {}},
			expect: []uint{1, 2, 3},
		},
		{
			name:   "test Uint Union, s is empty",
			s:      map[uint]struct{}{1: {}, 2: {}, 3: {}},
			t:      Uint{},
			expect: []uint{1, 2, 3},
		},
		{
			name:   "test Uint Union, s has same element to s",
			s:      map[uint]struct{}{1: {}, 2: {}, 3: {}},
			t:      map[uint]struct{}{1: {}, 9: {}, 4: {}},
			expect: []uint{1, 2, 3, 4, 9},
		},
		{
			name:   "test Uint Union, s does not have same element to s",
			s:      map[uint]struct{}{1: {}, 2: {}, 3: {}},
			t:      map[uint]struct{}{6: {}, 7: {}, 8: {}},
			expect: []uint{1, 2, 3, 6, 7, 8},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		actual := tc.s.Union(tc.t)
		validateUint(t, actual, tc.expect)
	}
}

func TestUint_Difference(t *testing.T) {
	testcases := []struct {
		name   string
		s      Uint
		t      Uint
		expect []uint
	}{
		{
			name:   "test Uint Difference, s is empty",
			s:      Uint{},
			t:      map[uint]struct{}{2: {}, 9: {}, 4: {}},
			expect: []uint{},
		},
		{
			name:   "test Uint Difference, s is empty",
			s:      map[uint]struct{}{2: {}, 9: {}, 4: {}},
			t:      Uint{},
			expect: []uint{2, 9, 4},
		},
		{
			name:   "test Uint Difference, s ⊂ s",
			s:      map[uint]struct{}{2: {}, 9: {}},
			t:      map[uint]struct{}{2: {}, 9: {}, 4: {}},
			expect: []uint{},
		},
		{
			name:   "test Uint Difference, s ⊃ s",
			s:      map[uint]struct{}{2: {}, 9: {}, 4: {}},
			t:      map[uint]struct{}{2: {}},
			expect: []uint{9, 4},
		},
		{
			name:   "test Uint Difference, s = s",
			s:      map[uint]struct{}{2: {}, 9: {}, 4: {}},
			t:      map[uint]struct{}{2: {}, 9: {}, 4: {}},
			expect: []uint{},
		},
		{
			name:   "test Uint Difference, s ∩ s = Ø",
			s:      map[uint]struct{}{2: {}, 9: {}, 4: {}},
			t:      map[uint]struct{}{1: {}, 6: {}},
			expect: []uint{2, 9, 4},
		},
		{
			name:   "test Uint Difference, s ∩ s ≠ Ø",
			s:      map[uint]struct{}{2: {}, 9: {}, 4: {}},
			t:      map[uint]struct{}{2: {}, 9: {}, 5: {}},
			expect: []uint{4},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		actual := tc.s.Difference(tc.t)
		validateUint(t, actual, tc.expect)
	}
}

func TestUint_Intersection(t *testing.T) {
	testcases := []struct {
		name   string
		s      Uint
		t      Uint
		expect []uint
	}{
		{
			name:   "test Uint Intersection, s is empty",
			s:      Uint{},
			t:      map[uint]struct{}{2: {}, 9: {}, 4: {}},
			expect: []uint{},
		},
		{
			name:   "test Uint Intersection, s is empty",
			s:      map[uint]struct{}{2: {}, 9: {}, 4: {}},
			t:      Uint{},
			expect: []uint{},
		},
		{
			name:   "test Uint Intersection, s ⊂ s",
			s:      map[uint]struct{}{2: {}, 9: {}},
			t:      map[uint]struct{}{2: {}, 9: {}, 4: {}},
			expect: []uint{2, 9},
		},
		{
			name:   "test Uint Intersection, s ⊃ s",
			s:      map[uint]struct{}{2: {}, 9: {}, 4: {}},
			t:      map[uint]struct{}{2: {}, 9: {}},
			expect: []uint{2, 9},
		},
		{
			name:   "test Uint Intersection, s = s",
			s:      map[uint]struct{}{2: {}, 9: {}, 4: {}},
			t:      map[uint]struct{}{2: {}, 9: {}, 4: {}},
			expect: []uint{2, 9, 4},
		},
		{
			name:   "test Uint Intersection, s ∩ s = Ø",
			s:      map[uint]struct{}{1: {}, 4: {}},
			t:      map[uint]struct{}{2: {}, 6: {}},
			expect: []uint{},
		},
		{
			name:   "test Uint Intersection, s ∩ s ≠ Ø",
			s:      map[uint]struct{}{1: {}, 4: {}},
			t:      map[uint]struct{}{1: {}, 6: {}},
			expect: []uint{1},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		actual := tc.s.Intersection(tc.t)
		validateUint(t, actual, tc.expect)
	}
}

func TestUint_SymmetricDifference(t *testing.T) {
	testcases := []struct {
		name   string
		s      Uint
		t      Uint
		expect []uint
	}{
		{
			name:   "test Uint SymmetricDifference, s is empty",
			s:      Uint{},
			t:      map[uint]struct{}{2: {}, 9: {}, 4: {}},
			expect: []uint{2, 9, 4},
		},
		{
			name:   "test Uint SymmetricDifference, s is empty",
			s:      map[uint]struct{}{2: {}, 9: {}, 4: {}},
			t:      Uint{},
			expect: []uint{2, 9, 4},
		},
		{
			name:   "test Uint SymmetricDifference, s ⊂ s",
			s:      map[uint]struct{}{2: {}, 9: {}},
			t:      map[uint]struct{}{2: {}, 9: {}, 4: {}},
			expect: []uint{4},
		},
		{
			name:   "test Uint SymmetricDifference, s ⊃ s",
			s:      map[uint]struct{}{2: {}, 9: {}, 4: {}},
			t:      map[uint]struct{}{2: {}, 9: {}},
			expect: []uint{4},
		},
		{
			name:   "test Uint SymmetricDifference, s = s",
			s:      map[uint]struct{}{2: {}, 9: {}, 4: {}},
			t:      map[uint]struct{}{2: {}, 9: {}, 4: {}},
			expect: []uint{},
		},
		{
			name:   "test Uint SymmetricDifference, s ∩ s = Ø",
			s:      map[uint]struct{}{1: {}, 4: {}},
			t:      map[uint]struct{}{2: {}, 6: {}},
			expect: []uint{1, 2, 4, 6},
		},
		{
			name:   "test Uint SymmetricDifference, s ∩ s ≠ Ø",
			s:      map[uint]struct{}{1: {}, 4: {}},
			t:      map[uint]struct{}{1: {}, 6: {}},
			expect: []uint{4, 6},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		actual := tc.s.SymmetricDifference(tc.t)
		validateUint(t, actual, tc.expect)
	}
}

func TestUint_IsSubset(t *testing.T) {
	testcases := []struct {
		name   string
		s      Uint
		t      Uint
		expect bool
	}{
		{
			name:   "test Uint IsSubset, s is empty",
			s:      Uint{},
			t:      map[uint]struct{}{2: {}, 9: {}, 4: {}},
			expect: true,
		},
		{
			name:   "test Uint IsSubset, s is empty",
			s:      map[uint]struct{}{2: {}, 9: {}, 4: {}},
			t:      Uint{},
			expect: false,
		},
		{
			name:   "test Uint IsSubset, s ⊂ s",
			s:      map[uint]struct{}{2: {}, 9: {}},
			t:      map[uint]struct{}{2: {}, 9: {}, 4: {}},
			expect: true,
		},
		{
			name:   "test Uint IsSubset, s ⊃ s",
			s:      map[uint]struct{}{2: {}, 9: {}, 4: {}},
			t:      map[uint]struct{}{2: {}, 9: {}},
			expect: false,
		},
		{
			name:   "test Uint IsSubset, s = s",
			s:      map[uint]struct{}{2: {}, 9: {}, 4: {}},
			t:      map[uint]struct{}{2: {}, 9: {}, 4: {}},
			expect: true,
		},
		{
			name:   "test Uint IsSubset, s ∩ s = Ø",
			s:      map[uint]struct{}{1: {}, 4: {}},
			t:      map[uint]struct{}{2: {}, 6: {}},
			expect: false,
		},
		{
			name:   "test Uint IsSubset, s ∩ s ≠ Ø && s ∩ s ≠ s",
			s:      map[uint]struct{}{1: {}, 4: {}},
			t:      map[uint]struct{}{1: {}, 6: {}},
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

func TestUint_IsSuperset(t *testing.T) {
	testcases := []struct {
		name   string
		s      Uint
		t      Uint
		expect bool
	}{
		{
			name:   "test Uint IsSuperset, s is empty",
			s:      Uint{},
			t:      map[uint]struct{}{2: {}, 9: {}, 4: {}},
			expect: false,
		},
		{
			name:   "test Uint IsSuperset, s is empty",
			s:      map[uint]struct{}{2: {}, 9: {}, 4: {}},
			t:      Uint{},
			expect: true,
		},
		{
			name:   "test Uint IsSuperset, s ⊂ s",
			s:      map[uint]struct{}{2: {}, 9: {}},
			t:      map[uint]struct{}{2: {}, 9: {}, 4: {}},
			expect: false,
		},
		{
			name:   "test Uint IsSuperset, s ⊃ s",
			s:      map[uint]struct{}{2: {}, 9: {}, 4: {}},
			t:      map[uint]struct{}{2: {}, 9: {}},
			expect: true,
		},
		{
			name:   "test Uint IsSuperset, s = s",
			s:      map[uint]struct{}{2: {}, 9: {}, 4: {}},
			t:      map[uint]struct{}{2: {}, 9: {}, 4: {}},
			expect: true,
		},
		{
			name:   "test Uint IsSuperset, s ∩ s = Ø",
			s:      map[uint]struct{}{1: {}, 4: {}},
			t:      map[uint]struct{}{2: {}, 6: {}},
			expect: false,
		},
		{
			name:   "test Uint IsSuperset, s ∩ s ≠ Ø && s ∩ s ≠ s",
			s:      map[uint]struct{}{1: {}, 4: {}},
			t:      map[uint]struct{}{1: {}, 6: {}},
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

func TestUint_Equal(t *testing.T) {
	testcases := []struct {
		name   string
		s      Uint
		t      Uint
		expect bool
	}{
		{
			name:   "test Uint Equal, s is empty",
			s:      Uint{},
			t:      map[uint]struct{}{2: {}, 9: {}, 4: {}},
			expect: false,
		},
		{
			name:   "test Uint Equal, s is empty",
			s:      map[uint]struct{}{2: {}, 9: {}, 4: {}},
			t:      Uint{},
			expect: false,
		},
		{
			name:   "test Uint Equal, s ⊂ s",
			s:      map[uint]struct{}{2: {}, 9: {}},
			t:      map[uint]struct{}{2: {}, 9: {}, 4: {}},
			expect: false,
		},
		{
			name:   "test Uint Equal, s ⊃ s",
			s:      map[uint]struct{}{2: {}, 9: {}, 4: {}},
			t:      map[uint]struct{}{2: {}, 9: {}},
			expect: false,
		},
		{
			name:   "test Uint Equal, s = s",
			s:      map[uint]struct{}{2: {}, 9: {}, 4: {}},
			t:      map[uint]struct{}{2: {}, 9: {}, 4: {}},
			expect: true,
		},
		{
			name:   "test Uint Equal, s ∩ s = Ø",
			s:      map[uint]struct{}{1: {}, 4: {}},
			t:      map[uint]struct{}{2: {}, 6: {}},
			expect: false,
		},
		{
			name:   "test Uint Equal, s ∩ s ≠ Ø && s ∩ s ≠ s",
			s:      map[uint]struct{}{1: {}, 4: {}},
			t:      map[uint]struct{}{1: {}, 6: {}},
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

func TestUint_Copy(t *testing.T) {
	testcases := []struct {
		name   string
		s      Uint
		expect []uint
	}{
		{
			name:   "test Uint Copy, s is empty",
			s:      Uint{},
			expect: []uint{},
		},
		{
			name:   "test Uint Copy, s is not empty",
			s:      map[uint]struct{}{2: {}, 9: {}, 4: {}},
			expect: []uint{2, 9, 4},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		actual := tc.s.Copy()
		validateUint(t, actual, tc.expect)
	}
}

func TestUint_String(t *testing.T) {
	testcases := []struct {
		name   string
		s      Uint
		expect string
	}{
		{
			name:   "test Uint String, s is empty",
			s:      Uint{},
			expect: "[]",
		},
		{
			name:   "test Uint String, s is not empty",
			s:      map[uint]struct{}{1: {}},
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

func validateUint(t *testing.T, actual Uint, expect []uint) {
	if len(expect) != len(actual) {
		t.Errorf("expect set len: %d, but got: %d", len(expect), len(actual))
	}
	slice2String := func(i []uint) string {
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
