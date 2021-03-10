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

func TestNewInt(t *testing.T) {
	testCases := []struct {
		name   string
		input  []int
		expect []int
	}{
		{
			name:   "test Int New, s nothing",
			input:  nil,
			expect: nil,
		},
		{
			name:   "test Int New, inputs multiple elements",
			input:  []int{1, 2, 3},
			expect: []int{1, 2, 3},
		},
	}
	for _, tc := range testCases {
		t.Logf("running scenario: %s", tc.name)
		actual := NewInt(tc.input...)
		validateInt(t, actual, tc.expect)
	}
}

func TestInt_Add(t *testing.T) {
	testcases := []struct {
		name   string
		input  []int
		expect []int
	}{
		{
			name:   "test Int Add, inputs nothing",
			input:  []int{},
			expect: []int{},
		},
		{
			name:   "test Int Add, inputs multiple elements",
			input:  []int{1, 2, 3},
			expect: []int{1, 2, 3},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		actual := Int{}
		actual.Add(tc.input...)
		validateInt(t, actual, tc.expect)
	}
}

func TestInt_Remove(t *testing.T) {
	testcases := []struct {
		name   string
		s      Int
		input  []int
		expect []int
	}{
		{
			name:   "test Int Remove, inputs nothing",
			s:      map[int]struct{}{1: {}, 2: {}, 3: {}},
			input:  []int{},
			expect: []int{1, 2, 3},
		},
		{
			name:   "test Int Remove, inputs multiple exit elements",
			s:      map[int]struct{}{1: {}, 2: {}, 3: {}},
			input:  []int{1, 3},
			expect: []int{2},
		},
		{
			name:   "test Int Remove, inputs multiple non-exit elements",
			s:      map[int]struct{}{1: {}, 2: {}, 3: {}},
			input:  []int{0, 6, 7},
			expect: []int{1, 2, 3},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		tc.s.Remove(tc.input...)
		validateInt(t, tc.s, tc.expect)
	}
}

func TestInt_Pop(t *testing.T) {
	testcases := []struct {
		name   string
		s      Int
		expect bool
	}{
		{
			name:   "test Int Pop, s is not empty",
			s:      map[int]struct{}{1: {}, 2: {}, 3: {}},
			expect: true,
		},
		{
			name:   "test Int Pop, s is empty",
			s:      map[int]struct{}{},
			expect: false,
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		backup := make(map[int]struct{})
		for k := range tc.s {
			backup[k] = struct{}{}
		}
		key, ok := tc.s.Pop()
		if ok != tc.expect {
			t.Errorf("expect ok: %v, but got: %v", tc.expect, ok)
		}
		delete(backup, key)
		var expect []int
		for key := range backup {
			expect = append(expect, key)
		}
		validateInt(t, tc.s, expect)
	}
}

func TestInt_Size(t *testing.T) {
	testcases := []struct {
		name   string
		s      Int
		expect int
	}{
		{
			name:   "test Int Size, s is empty",
			s:      Int{},
			expect: 0,
		},
		{
			name:   "test Int Size, s is not empty",
			s:      map[int]struct{}{1: {}, 2: {}, 3: {}},
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

func TestInt_IsEmpty(t *testing.T) {
	testcases := []struct {
		name   string
		s      Int
		expect bool
	}{
		{
			name:   "test Int Empty, s is empty",
			s:      Int{},
			expect: true,
		},
		{
			name:   "test interface Empty, s is not empty",
			s:      map[int]struct{}{1: {}, 2: {}, 3: {}},
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

func TestInt_Clear(t *testing.T) {
	testcases := []struct {
		name string
		s    Int
	}{
		{
			name: "test Int Clear, s is empty",
			s:    Int{},
		},
		{
			name: "test Int Clear, s is not empty",
			s:    map[int]struct{}{1: {}, 2: {}, 3: {}},
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

func TestInt_Has(t *testing.T) {
	testcases := []struct {
		name   string
		s      Int
		input  int
		expect bool
	}{
		{
			name:   "test Int Has, s has input element",
			s:      map[int]struct{}{1: {}, 2: {}, 3: {}},
			input:  1,
			expect: true,
		},
		{
			name:   "test Int Has, s does not have element",
			s:      map[int]struct{}{1: {}, 2: {}, 3: {}},
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

func TestInt_HasAll(t *testing.T) {
	testcases := []struct {
		name   string
		s      Int
		input  []int
		expect bool
	}{
		{
			name:   "test Int HasAll, set has all input elements",
			s:      map[int]struct{}{1: {}, 2: {}, 3: {}},
			input:  []int{1, 2},
			expect: true,
		},
		{
			name:   "test Int HasAll, set does not have all input elements",
			s:      map[int]struct{}{1: {}, 2: {}, 3: {}},
			input:  []int{8, 9},
			expect: false,
		},
		{
			name:   "test Int HasAll, set does not have all input elements, but exist elements in set",
			s:      map[int]struct{}{1: {}, 2: {}, 3: {}},
			input:  []int{1, 9},
			expect: false,
		},
		{
			name:   "test Int HasAll, input empty",
			s:      map[int]struct{}{1: {}, 2: {}, 3: {}},
			input:  []int{},
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

func TestInt_HasAny(t *testing.T) {
	testcases := []struct {
		name   string
		s      Int
		input  []int
		expect bool
	}{
		{
			name:   "test Int HasAny, s has all elements",
			s:      map[int]struct{}{1: {}, 2: {}, 3: {}},
			input:  []int{1, 2},
			expect: true,
		},
		{
			name:   "test Int HasAny, s does not have all elements, but exist elements in set",
			s:      map[int]struct{}{1: {}, 2: {}, 3: {}},
			input:  []int{1, 9},
			expect: true,
		},
		{
			name:   "test Int HasAny, s does not has all elements",
			s:      map[int]struct{}{1: {}, 2: {}, 3: {}},
			input:  []int{8, 9},
			expect: false,
		},
		{
			name:   "test Int HasAny, input empty",
			s:      map[int]struct{}{1: {}, 2: {}, 3: {}},
			input:  []int{},
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

func TestInt_List(t *testing.T) {
	testcases := []struct {
		name   string
		s      Int
		input  []int
		expect bool
	}{
		{
			name: "test Int List, s is empty",
			s:    Int{},
		},
		{
			name: "test Int List, s is not empty",
			s:    map[int]struct{}{1: {}, 2: {}, 3: {}},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		actual := tc.s.List()
		validateInt(t, tc.s, actual)
	}
}

func TestInt_SortedList(t *testing.T) {
	testcases := []struct {
		name   string
		s      Int
		expect []int
	}{
		{
			name: "test Int List, s is empty",
			s:    Int{},
		},
		{
			name:   "test Int SortedList, s is not empty",
			s:      map[int]struct{}{1: {}, 2: {}, 3: {}},
			expect: []int{1, 2, 3},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		actual := tc.s.SortedList(func(i, j int) bool {
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

func TestInt_Each(t *testing.T) {
	testcases := []struct {
		name   string
		origin Int
	}{
		{
			name:   "test Int Each",
			origin: map[int]struct{}{1: {}, 2: {}, 3: {}},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		var expect []int
		tc.origin.Each(func(i int) {
			expect = append(expect, i)
		})
		validateInt(t, tc.origin, expect)
	}
}

func TestInt_EachE(t *testing.T) {
	inputErr := errors.New("s error")
	testcases := []struct {
		name      string
		origin    Int
		breakEach bool
		inputErr  error
		expectLen int
		expectErr error
	}{
		{
			name:      "test Int EachE",
			origin:    map[int]struct{}{1: {}, 2: {}, 3: {}},
			expectLen: 3,
		},
		{
			name:      "test Int EachE, return break error",
			origin:    map[int]struct{}{1: {}, 2: {}, 3: {}},
			breakEach: true,
		},
		{
			name:      "test Int Each, returns error",
			origin:    map[int]struct{}{1: {}, 2: {}, 3: {}},
			inputErr:  inputErr,
			expectErr: inputErr,
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		var expect []int
		err := tc.origin.EachE(func(i int) error {
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
			validateInt(t, tc.origin, expect)
		}
	}
}

func TestInt_Union(t *testing.T) {
	testcases := []struct {
		name   string
		s      Int
		t      Int
		expect []int
	}{
		{
			name:   "test Int Union, s and s are empty",
			s:      Int{},
			t:      Int{},
			expect: []int{},
		},
		{
			name:   "test Int Union, s is empty",
			s:      Int{},
			t:      map[int]struct{}{1: {}, 2: {}, 3: {}},
			expect: []int{1, 2, 3},
		},
		{
			name:   "test Int Union, s is empty",
			s:      map[int]struct{}{1: {}, 2: {}, 3: {}},
			t:      Int{},
			expect: []int{1, 2, 3},
		},
		{
			name:   "test Int Union, s has same element to s",
			s:      map[int]struct{}{1: {}, 2: {}, 3: {}},
			t:      map[int]struct{}{1: {}, 9: {}, 4: {}},
			expect: []int{1, 2, 3, 4, 9},
		},
		{
			name:   "test Int Union, s does not have same element to s",
			s:      map[int]struct{}{1: {}, 2: {}, 3: {}},
			t:      map[int]struct{}{6: {}, 7: {}, 8: {}},
			expect: []int{1, 2, 3, 6, 7, 8},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		actual := tc.s.Union(tc.t)
		validateInt(t, actual, tc.expect)
	}
}

func TestInt_Difference(t *testing.T) {
	testcases := []struct {
		name   string
		s      Int
		t      Int
		expect []int
	}{
		{
			name:   "test Int Difference, s is empty",
			s:      Int{},
			t:      map[int]struct{}{2: {}, 9: {}, 4: {}},
			expect: []int{},
		},
		{
			name:   "test Int Difference, s is empty",
			s:      map[int]struct{}{2: {}, 9: {}, 4: {}},
			t:      Int{},
			expect: []int{2, 9, 4},
		},
		{
			name:   "test Int Difference, s ⊂ s",
			s:      map[int]struct{}{2: {}, 9: {}},
			t:      map[int]struct{}{2: {}, 9: {}, 4: {}},
			expect: []int{},
		},
		{
			name:   "test Int Difference, s ⊃ s",
			s:      map[int]struct{}{2: {}, 9: {}, 4: {}},
			t:      map[int]struct{}{2: {}},
			expect: []int{9, 4},
		},
		{
			name:   "test Int Difference, s = s",
			s:      map[int]struct{}{2: {}, 9: {}, 4: {}},
			t:      map[int]struct{}{2: {}, 9: {}, 4: {}},
			expect: []int{},
		},
		{
			name:   "test Int Difference, s ∩ s = Ø",
			s:      map[int]struct{}{2: {}, 9: {}, 4: {}},
			t:      map[int]struct{}{1: {}, 6: {}},
			expect: []int{2, 9, 4},
		},
		{
			name:   "test Int Difference, s ∩ s ≠ Ø",
			s:      map[int]struct{}{2: {}, 9: {}, 4: {}},
			t:      map[int]struct{}{2: {}, 9: {}, 5: {}},
			expect: []int{4},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		actual := tc.s.Difference(tc.t)
		validateInt(t, actual, tc.expect)
	}
}

func TestInt_Intersection(t *testing.T) {
	testcases := []struct {
		name   string
		s      Int
		t      Int
		expect []int
	}{
		{
			name:   "test Int Intersection, s is empty",
			s:      Int{},
			t:      map[int]struct{}{2: {}, 9: {}, 4: {}},
			expect: []int{},
		},
		{
			name:   "test Int Intersection, s is empty",
			s:      map[int]struct{}{2: {}, 9: {}, 4: {}},
			t:      Int{},
			expect: []int{},
		},
		{
			name:   "test Int Intersection, s ⊂ s",
			s:      map[int]struct{}{2: {}, 9: {}},
			t:      map[int]struct{}{2: {}, 9: {}, 4: {}},
			expect: []int{2, 9},
		},
		{
			name:   "test Int Intersection, s ⊃ s",
			s:      map[int]struct{}{2: {}, 9: {}, 4: {}},
			t:      map[int]struct{}{2: {}, 9: {}},
			expect: []int{2, 9},
		},
		{
			name:   "test Int Intersection, s = s",
			s:      map[int]struct{}{2: {}, 9: {}, 4: {}},
			t:      map[int]struct{}{2: {}, 9: {}, 4: {}},
			expect: []int{2, 9, 4},
		},
		{
			name:   "test Int Intersection, s ∩ s = Ø",
			s:      map[int]struct{}{1: {}, 4: {}},
			t:      map[int]struct{}{2: {}, 6: {}},
			expect: []int{},
		},
		{
			name:   "test Int Intersection, s ∩ s ≠ Ø",
			s:      map[int]struct{}{1: {}, 4: {}},
			t:      map[int]struct{}{1: {}, 6: {}},
			expect: []int{1},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		actual := tc.s.Intersection(tc.t)
		validateInt(t, actual, tc.expect)
	}
}

func TestInt_SymmetricDifference(t *testing.T) {
	testcases := []struct {
		name   string
		s      Int
		t      Int
		expect []int
	}{
		{
			name:   "test Int SymmetricDifference, s is empty",
			s:      Int{},
			t:      map[int]struct{}{2: {}, 9: {}, 4: {}},
			expect: []int{2, 9, 4},
		},
		{
			name:   "test Int SymmetricDifference, s is empty",
			s:      map[int]struct{}{2: {}, 9: {}, 4: {}},
			t:      Int{},
			expect: []int{2, 9, 4},
		},
		{
			name:   "test Int SymmetricDifference, s ⊂ s",
			s:      map[int]struct{}{2: {}, 9: {}},
			t:      map[int]struct{}{2: {}, 9: {}, 4: {}},
			expect: []int{4},
		},
		{
			name:   "test Int SymmetricDifference, s ⊃ s",
			s:      map[int]struct{}{2: {}, 9: {}, 4: {}},
			t:      map[int]struct{}{2: {}, 9: {}},
			expect: []int{4},
		},
		{
			name:   "test Int SymmetricDifference, s = s",
			s:      map[int]struct{}{2: {}, 9: {}, 4: {}},
			t:      map[int]struct{}{2: {}, 9: {}, 4: {}},
			expect: []int{},
		},
		{
			name:   "test Int SymmetricDifference, s ∩ s = Ø",
			s:      map[int]struct{}{1: {}, 4: {}},
			t:      map[int]struct{}{2: {}, 6: {}},
			expect: []int{1, 2, 4, 6},
		},
		{
			name:   "test Int SymmetricDifference, s ∩ s ≠ Ø",
			s:      map[int]struct{}{1: {}, 4: {}},
			t:      map[int]struct{}{1: {}, 6: {}},
			expect: []int{4, 6},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		actual := tc.s.SymmetricDifference(tc.t)
		validateInt(t, actual, tc.expect)
	}
}

func TestInt_IsSubset(t *testing.T) {
	testcases := []struct {
		name   string
		s      Int
		t      Int
		expect bool
	}{
		{
			name:   "test Int IsSubset, s is empty",
			s:      Int{},
			t:      map[int]struct{}{2: {}, 9: {}, 4: {}},
			expect: true,
		},
		{
			name:   "test Int IsSubset, s is empty",
			s:      map[int]struct{}{2: {}, 9: {}, 4: {}},
			t:      Int{},
			expect: false,
		},
		{
			name:   "test Int IsSubset, s ⊂ s",
			s:      map[int]struct{}{2: {}, 9: {}},
			t:      map[int]struct{}{2: {}, 9: {}, 4: {}},
			expect: true,
		},
		{
			name:   "test Int IsSubset, s ⊃ s",
			s:      map[int]struct{}{2: {}, 9: {}, 4: {}},
			t:      map[int]struct{}{2: {}, 9: {}},
			expect: false,
		},
		{
			name:   "test Int IsSubset, s = s",
			s:      map[int]struct{}{2: {}, 9: {}, 4: {}},
			t:      map[int]struct{}{2: {}, 9: {}, 4: {}},
			expect: true,
		},
		{
			name:   "test Int IsSubset, s ∩ s = Ø",
			s:      map[int]struct{}{1: {}, 4: {}},
			t:      map[int]struct{}{2: {}, 6: {}},
			expect: false,
		},
		{
			name:   "test Int IsSubset, s ∩ s ≠ Ø && s ∩ s ≠ s",
			s:      map[int]struct{}{1: {}, 4: {}},
			t:      map[int]struct{}{1: {}, 6: {}},
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

func TestInt_IsSuperset(t *testing.T) {
	testcases := []struct {
		name   string
		s      Int
		t      Int
		expect bool
	}{
		{
			name:   "test Int IsSuperset, s is empty",
			s:      Int{},
			t:      map[int]struct{}{2: {}, 9: {}, 4: {}},
			expect: false,
		},
		{
			name:   "test Int IsSuperset, s is empty",
			s:      map[int]struct{}{2: {}, 9: {}, 4: {}},
			t:      Int{},
			expect: true,
		},
		{
			name:   "test Int IsSuperset, s ⊂ s",
			s:      map[int]struct{}{2: {}, 9: {}},
			t:      map[int]struct{}{2: {}, 9: {}, 4: {}},
			expect: false,
		},
		{
			name:   "test Int IsSuperset, s ⊃ s",
			s:      map[int]struct{}{2: {}, 9: {}, 4: {}},
			t:      map[int]struct{}{2: {}, 9: {}},
			expect: true,
		},
		{
			name:   "test Int IsSuperset, s = s",
			s:      map[int]struct{}{2: {}, 9: {}, 4: {}},
			t:      map[int]struct{}{2: {}, 9: {}, 4: {}},
			expect: true,
		},
		{
			name:   "test Int IsSuperset, s ∩ s = Ø",
			s:      map[int]struct{}{1: {}, 4: {}},
			t:      map[int]struct{}{2: {}, 6: {}},
			expect: false,
		},
		{
			name:   "test Int IsSuperset, s ∩ s ≠ Ø && s ∩ s ≠ s",
			s:      map[int]struct{}{1: {}, 4: {}},
			t:      map[int]struct{}{1: {}, 6: {}},
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

func TestInt_Equal(t *testing.T) {
	testcases := []struct {
		name   string
		s      Int
		t      Int
		expect bool
	}{
		{
			name:   "test Int Equal, s is empty",
			s:      Int{},
			t:      map[int]struct{}{2: {}, 9: {}, 4: {}},
			expect: false,
		},
		{
			name:   "test Int Equal, s is empty",
			s:      map[int]struct{}{2: {}, 9: {}, 4: {}},
			t:      Int{},
			expect: false,
		},
		{
			name:   "test Int Equal, s ⊂ s",
			s:      map[int]struct{}{2: {}, 9: {}},
			t:      map[int]struct{}{2: {}, 9: {}, 4: {}},
			expect: false,
		},
		{
			name:   "test Int Equal, s ⊃ s",
			s:      map[int]struct{}{2: {}, 9: {}, 4: {}},
			t:      map[int]struct{}{2: {}, 9: {}},
			expect: false,
		},
		{
			name:   "test Int Equal, s = s",
			s:      map[int]struct{}{2: {}, 9: {}, 4: {}},
			t:      map[int]struct{}{2: {}, 9: {}, 4: {}},
			expect: true,
		},
		{
			name:   "test Int Equal, s ∩ s = Ø",
			s:      map[int]struct{}{1: {}, 4: {}},
			t:      map[int]struct{}{2: {}, 6: {}},
			expect: false,
		},
		{
			name:   "test Int Equal, s ∩ s ≠ Ø && s ∩ s ≠ s",
			s:      map[int]struct{}{1: {}, 4: {}},
			t:      map[int]struct{}{1: {}, 6: {}},
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

func TestInt_Copy(t *testing.T) {
	testcases := []struct {
		name   string
		s      Int
		expect []int
	}{
		{
			name:   "test Int Copy, s is empty",
			s:      Int{},
			expect: []int{},
		},
		{
			name:   "test Int Copy, s is not empty",
			s:      map[int]struct{}{2: {}, 9: {}, 4: {}},
			expect: []int{2, 9, 4},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		actual := tc.s.Copy()
		validateInt(t, actual, tc.expect)
	}
}

func TestInt_String(t *testing.T) {
	testcases := []struct {
		name   string
		s      Int
		expect string
	}{
		{
			name:   "test Int String, s is empty",
			s:      Int{},
			expect: "[]",
		},
		{
			name:   "test Int String, s is not empty",
			s:      map[int]struct{}{1: {}},
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

func validateInt(t *testing.T, actual Int, expect []int) {
	if len(expect) != len(actual) {
		t.Errorf("expect set len: %d, but got: %d", len(expect), len(actual))
	}
	slice2String := func(i []int) string {
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
