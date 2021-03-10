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

func TestNewFloat32(t *testing.T) {
	testCases := []struct {
		name   string
		input  []float32
		expect []float32
	}{
		{
			name:   "test Float32 New, s nothing",
			input:  nil,
			expect: nil,
		},
		{
			name:   "test Float32 New, inputs multiple elements",
			input:  []float32{1, 1.5, 1.2},
			expect: []float32{1, 1.5, 1.2},
		},
	}
	for _, tc := range testCases {
		t.Logf("running scenario: %s", tc.name)
		actual := NewFloat32(tc.input...)
		validateFloat32(t, actual, tc.expect)
	}
}

func TestFloat32_Add(t *testing.T) {
	testcases := []struct {
		name   string
		input  []float32
		expect []float32
	}{
		{
			name:   "test Float32 Add, inputs nothing",
			input:  []float32{},
			expect: []float32{},
		},
		{
			name:   "test Float32 Add, inputs multiple elements",
			input:  []float32{1, 2, 1.3},
			expect: []float32{1, 2, 1.3},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		actual := Float32{}
		actual.Add(tc.input...)
		validateFloat32(t, actual, tc.expect)
	}
}

func TestFloat32_Remove(t *testing.T) {
	testcases := []struct {
		name   string
		s      Float32
		input  []float32
		expect []float32
	}{
		{
			name:   "test Float32 Remove, inputs nothing",
			s:      map[float32]struct{}{1: {}, 1.5: {}, 1.2: {}},
			input:  []float32{},
			expect: []float32{1, 1.5, 1.2},
		},
		{
			name:   "test Float32 Remove, inputs multiple exit elements",
			s:      map[float32]struct{}{1: {}, 1.5: {}, 1.2: {}},
			input:  []float32{1, 1.2},
			expect: []float32{1.5},
		},
		{
			name:   "test Float32 Remove, inputs multiple non-exit elements",
			s:      map[float32]struct{}{1: {}, 1.5: {}, 1.2: {}},
			input:  []float32{0, 1.9, 1.3},
			expect: []float32{1, 1.5, 1.2},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		tc.s.Remove(tc.input...)
		validateFloat32(t, tc.s, tc.expect)
	}
}

func TestFloat32_Pop(t *testing.T) {
	testcases := []struct {
		name   string
		s      Float32
		expect bool
	}{
		{
			name:   "test Float32 Pop, s is not empty",
			s:      map[float32]struct{}{1: {}, 1.5: {}, 1.2: {}},
			expect: true,
		},
		{
			name:   "test Float32 Pop, s is empty",
			s:      map[float32]struct{}{},
			expect: false,
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		backup := make(map[float32]struct{})
		for k := range tc.s {
			backup[k] = struct{}{}
		}
		key, ok := tc.s.Pop()
		if ok != tc.expect {
			t.Errorf("expect ok: %v, but got: %v", tc.expect, ok)
		}
		delete(backup, key)
		var expect []float32
		for key := range backup {
			expect = append(expect, key)
		}
		validateFloat32(t, tc.s, expect)
	}
}

func TestFloat32_Size(t *testing.T) {
	testcases := []struct {
		name   string
		s      Float32
		expect int
	}{
		{
			name:   "test Float32 Size, s is empty",
			s:      Float32{},
			expect: 0,
		},
		{
			name:   "test Float32 Size, s is not empty",
			s:      map[float32]struct{}{1: {}, 1.5: {}, 1.2: {}},
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

func TestFloat32_IsEmpty(t *testing.T) {
	testcases := []struct {
		name   string
		s      Float32
		expect bool
	}{
		{
			name:   "test Float32 Empty, s is empty",
			s:      Float32{},
			expect: true,
		},
		{
			name:   "test interface Empty, s is not empty",
			s:      map[float32]struct{}{1: {}, 1.5: {}, 1.2: {}},
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

func TestFloat32_Clear(t *testing.T) {
	testcases := []struct {
		name string
		s    Float32
	}{
		{
			name: "test Float32 Clear, s is empty",
			s:    Float32{},
		},
		{
			name: "test Float32 Clear, s is not empty",
			s:    map[float32]struct{}{1: {}, 1.5: {}, 1.2: {}},
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

func TestFloat32_Has(t *testing.T) {
	testcases := []struct {
		name   string
		s      Float32
		input  float32
		expect bool
	}{
		{
			name:   "test Float32 Has, s has input element",
			s:      map[float32]struct{}{1: {}, 1.5: {}, 1.2: {}},
			input:  1,
			expect: true,
		},
		{
			name:   "test Float32 Has, s does not have element",
			s:      map[float32]struct{}{1: {}, 1.5: {}, 1.2: {}},
			input:  2,
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

func TestFloat32_HasAll(t *testing.T) {
	testcases := []struct {
		name   string
		s      Float32
		input  []float32
		expect bool
	}{
		{
			name:   "test Float32 HasAll, set has all input elements",
			s:      map[float32]struct{}{1: {}, 1.5: {}, 1.2: {}},
			input:  []float32{1, 1.5},
			expect: true,
		},
		{
			name:   "test Float32 HasAll, set does not have all input elements",
			s:      map[float32]struct{}{1: {}, 1.5: {}, 1.2: {}},
			input:  []float32{1.1, 1.3},
			expect: false,
		},
		{
			name:   "test Float32 HasAll, set does not have all input elements, but exist elements in set",
			s:      map[float32]struct{}{1: {}, 1.5: {}, 1.2: {}},
			input:  []float32{1, 1.3},
			expect: false,
		},
		{
			name:   "test Float32 HasAll, input empty",
			s:      map[float32]struct{}{1: {}, 1.5: {}, 1.2: {}},
			input:  []float32{},
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

func TestFloat32_HasAny(t *testing.T) {
	testcases := []struct {
		name   string
		s      Float32
		input  []float32
		expect bool
	}{
		{
			name:   "test Float32 HasAny, s has all elements",
			s:      map[float32]struct{}{1: {}, 1.5: {}, 1.2: {}},
			input:  []float32{1, 1.5},
			expect: true,
		},
		{
			name:   "test Float32 HasAny, s does not have all elements, but exist elements in set",
			s:      map[float32]struct{}{1: {}, 1.5: {}, 1.2: {}},
			input:  []float32{1, 1.3},
			expect: true,
		},
		{
			name:   "test Float32 HasAny, s does not has all elements",
			s:      map[float32]struct{}{1: {}, 1.5: {}, 1.2: {}},
			input:  []float32{2, 1.3},
			expect: false,
		},
		{
			name:   "test Float32 HasAll, input empty",
			s:      map[float32]struct{}{1: {}, 1.5: {}, 1.2: {}},
			input:  []float32{},
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

func TestFloat32_List(t *testing.T) {
	testcases := []struct {
		name   string
		s      Float32
		input  []float32
		expect bool
	}{
		{
			name: "test Float32 List, s is empty",
			s:    Float32{},
		},
		{
			name: "test Float32 List, s is not empty",
			s:    map[float32]struct{}{1: {}, 1.5: {}, 1.2: {}},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		actual := tc.s.List()
		validateFloat32(t, tc.s, actual)
	}
}

func TestFloat32_SortedList(t *testing.T) {
	testcases := []struct {
		name   string
		s      Float32
		expect []float32
	}{
		{
			name: "test Float32 List, s is empty",
			s:    Float32{},
		},
		{
			name:   "test Float32 SortedList, s is not empty",
			s:      map[float32]struct{}{1: {}, 2: {}, 3: {}},
			expect: []float32{1, 2, 3},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		actual := tc.s.SortedList(func(i, j float32) bool {
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

func TestFloat32_Each(t *testing.T) {
	testcases := []struct {
		name   string
		origin Float32
	}{
		{
			name:   "test Float32 Each",
			origin: map[float32]struct{}{1: {}, 2: {}, 3: {}},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		var expect []float32
		tc.origin.Each(func(i float32) {
			expect = append(expect, i)
		})
		validateFloat32(t, tc.origin, expect)
	}
}

func TestFloat32_EachE(t *testing.T) {
	inputErr := errors.New("s error")
	testcases := []struct {
		name      string
		origin    Float32
		breakEach bool
		inputErr  error
		expectLen int
		expectErr error
	}{
		{
			name:      "test Float32 EachE",
			origin:    map[float32]struct{}{1: {}, 2: {}, 3: {}},
			expectLen: 3,
		},
		{
			name:      "test Float32 EachE, return break error",
			origin:    map[float32]struct{}{1: {}, 2: {}, 3: {}},
			breakEach: true,
		},
		{
			name:      "test Float32 Each, returns error",
			origin:    map[float32]struct{}{1: {}, 2: {}, 3: {}},
			inputErr:  inputErr,
			expectErr: inputErr,
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		var expect []float32
		err := tc.origin.EachE(func(i float32) error {
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
			validateFloat32(t, tc.origin, expect)
		}
	}
}

func TestFloat32_Union(t *testing.T) {
	testcases := []struct {
		name   string
		s      Float32
		t      Float32
		expect []float32
	}{
		{
			name:   "test Float32 Union, s and s are empty",
			s:      Float32{},
			t:      Float32{},
			expect: []float32{},
		},
		{
			name:   "test Float32 Union, s is empty",
			s:      Float32{},
			t:      map[float32]struct{}{1: {}, 2: {}, 3.1: {}},
			expect: []float32{1, 2, 3.1},
		},
		{
			name:   "test Float32 Union, s is empty",
			s:      map[float32]struct{}{1: {}, 2: {}, 3.1: {}},
			t:      Float32{},
			expect: []float32{1, 2, 3.1},
		},
		{
			name:   "test Float32 Union, s has same element to s",
			s:      map[float32]struct{}{1: {}, 2: {}, 3.1: {}},
			t:      map[float32]struct{}{1: {}, 1.3: {}, 4.1: {}},
			expect: []float32{1, 2, 3.1, 4.1, 1.3},
		},
		{
			name:   "test Float32 Union, s does not have same element to s",
			s:      map[float32]struct{}{1: {}, 2: {}, 3.1: {}},
			t:      map[float32]struct{}{2.1: {}, 1.3: {}, 4.1: {}},
			expect: []float32{1, 2, 2.1, 3.1, 4.1, 1.3},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		actual := tc.s.Union(tc.t)
		validateFloat32(t, actual, tc.expect)
	}
}

func TestFloat32_Difference(t *testing.T) {
	testcases := []struct {
		name   string
		s      Float32
		t      Float32
		expect []float32
	}{
		{
			name:   "test Float32 Difference, s is empty",
			s:      Float32{},
			t:      map[float32]struct{}{2: {}, 1.3: {}, 4.1: {}},
			expect: []float32{},
		},
		{
			name:   "test Float32 Difference, s is empty",
			s:      map[float32]struct{}{2: {}, 1.3: {}, 4.1: {}},
			t:      Float32{},
			expect: []float32{2, 1.3, 4.1},
		},
		{
			name:   "test Float32 Difference, s ⊂ s",
			s:      map[float32]struct{}{2: {}, 1.3: {}},
			t:      map[float32]struct{}{2: {}, 1.3: {}, 4.1: {}},
			expect: []float32{},
		},
		{
			name:   "test Float32 Difference, s ⊃ s",
			s:      map[float32]struct{}{2: {}, 1.3: {}, 4.1: {}},
			t:      map[float32]struct{}{2: {}},
			expect: []float32{1.3, 4.1},
		},
		{
			name:   "test Float32 Difference, s = s",
			s:      map[float32]struct{}{2: {}, 1.3: {}, 4.1: {}},
			t:      map[float32]struct{}{2: {}, 1.3: {}, 4.1: {}},
			expect: []float32{},
		},
		{
			name:   "test Float32 Difference, s ∩ s = Ø",
			s:      map[float32]struct{}{2: {}, 1.3: {}, 4.1: {}},
			t:      map[float32]struct{}{1: {}, 6.0: {}},
			expect: []float32{2, 1.3, 4.1},
		},
		{
			name:   "test Float32 Difference, s ∩ s ≠ Ø",
			s:      map[float32]struct{}{2: {}, 1.3: {}, 4.1: {}},
			t:      map[float32]struct{}{2: {}, 1.3: {}, 4.2: {}},
			expect: []float32{4.1},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		actual := tc.s.Difference(tc.t)
		validateFloat32(t, actual, tc.expect)
	}
}

func TestFloat32_Intersection(t *testing.T) {
	testcases := []struct {
		name   string
		s      Float32
		t      Float32
		expect []float32
	}{
		{
			name:   "test Float32 Intersection, s is empty",
			s:      Float32{},
			t:      map[float32]struct{}{2: {}, 1.3: {}, 4.1: {}},
			expect: []float32{},
		},
		{
			name:   "test Float32 Intersection, s is empty",
			s:      map[float32]struct{}{2: {}, 1.3: {}, 4.1: {}},
			t:      Float32{},
			expect: []float32{},
		},
		{
			name:   "test Float32 Intersection, s ⊂ s",
			s:      map[float32]struct{}{2: {}, 1.3: {}},
			t:      map[float32]struct{}{2: {}, 1.3: {}, 4.1: {}},
			expect: []float32{2, 1.3},
		},
		{
			name:   "test Float32 Intersection, s ⊃ s",
			s:      map[float32]struct{}{2: {}, 1.3: {}, 4.1: {}},
			t:      map[float32]struct{}{2: {}, 1.3: {}},
			expect: []float32{2, 1.3},
		},
		{
			name:   "test Float32 Intersection, s = s",
			s:      map[float32]struct{}{2: {}, 1.3: {}, 4.1: {}},
			t:      map[float32]struct{}{2: {}, 1.3: {}, 4.1: {}},
			expect: []float32{2, 1.3, 4.1},
		},
		{
			name:   "test Float32 Intersection, s ∩ s = Ø",
			s:      map[float32]struct{}{1: {}, 4.1: {}},
			t:      map[float32]struct{}{2: {}, 6.0: {}},
			expect: []float32{},
		},
		{
			name:   "test Float32 Intersection, s ∩ s ≠ Ø",
			s:      map[float32]struct{}{1: {}, 4.1: {}},
			t:      map[float32]struct{}{1: {}, 6.0: {}},
			expect: []float32{1},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		actual := tc.s.Intersection(tc.t)
		validateFloat32(t, actual, tc.expect)
	}
}

func TestFloat32_SymmetricDifference(t *testing.T) {
	testcases := []struct {
		name   string
		s      Float32
		t      Float32
		expect []float32
	}{
		{
			name:   "test Float32 SymmetricDifference, s is empty",
			s:      Float32{},
			t:      map[float32]struct{}{2: {}, 1.3: {}, 4.1: {}},
			expect: []float32{2, 1.3, 4.1},
		},
		{
			name:   "test Float32 SymmetricDifference, s is empty",
			s:      map[float32]struct{}{2: {}, 1.3: {}, 4.1: {}},
			t:      Float32{},
			expect: []float32{2, 1.3, 4.1},
		},
		{
			name:   "test Float32 SymmetricDifference, s ⊂ s",
			s:      map[float32]struct{}{2: {}, 1.3: {}},
			t:      map[float32]struct{}{2: {}, 1.3: {}, 4.1: {}},
			expect: []float32{4.1},
		},
		{
			name:   "test Float32 SymmetricDifference, s ⊃ s",
			s:      map[float32]struct{}{2: {}, 1.3: {}, 4.1: {}},
			t:      map[float32]struct{}{2: {}, 1.3: {}},
			expect: []float32{4.1},
		},
		{
			name:   "test Float32 SymmetricDifference, s = s",
			s:      map[float32]struct{}{2: {}, 1.3: {}, 4.1: {}},
			t:      map[float32]struct{}{2: {}, 1.3: {}, 4.1: {}},
			expect: []float32{},
		},
		{
			name:   "test Float32 SymmetricDifference, s ∩ s = Ø",
			s:      map[float32]struct{}{1: {}, 4.1: {}},
			t:      map[float32]struct{}{2: {}, 6.0: {}},
			expect: []float32{1, 2, 4.1, 6.0},
		},
		{
			name:   "test Float32 SymmetricDifference, s ∩ s ≠ Ø",
			s:      map[float32]struct{}{1: {}, 4.1: {}},
			t:      map[float32]struct{}{1: {}, 6.0: {}},
			expect: []float32{4.1, 6.0},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		actual := tc.s.SymmetricDifference(tc.t)
		validateFloat32(t, actual, tc.expect)
	}
}

func TestFloat32_IsSubset(t *testing.T) {
	testcases := []struct {
		name   string
		s      Float32
		t      Float32
		expect bool
	}{
		{
			name:   "test Float32 IsSubset, s is empty",
			s:      Float32{},
			t:      map[float32]struct{}{2: {}, 1.3: {}, 4.1: {}},
			expect: true,
		},
		{
			name:   "test Float32 IsSubset, s is empty",
			s:      map[float32]struct{}{2: {}, 1.3: {}, 4.1: {}},
			t:      Float32{},
			expect: false,
		},
		{
			name:   "test Float32 IsSubset, s ⊂ s",
			s:      map[float32]struct{}{2: {}, 1.3: {}},
			t:      map[float32]struct{}{2: {}, 1.3: {}, 4.1: {}},
			expect: true,
		},
		{
			name:   "test Float32 IsSubset, s ⊃ s",
			s:      map[float32]struct{}{2: {}, 1.3: {}, 4.1: {}},
			t:      map[float32]struct{}{2: {}, 1.3: {}},
			expect: false,
		},
		{
			name:   "test Float32 IsSubset, s = s",
			s:      map[float32]struct{}{2: {}, 1.3: {}, 4.1: {}},
			t:      map[float32]struct{}{2: {}, 1.3: {}, 4.1: {}},
			expect: true,
		},
		{
			name:   "test Float32 IsSubset, s ∩ s = Ø",
			s:      map[float32]struct{}{1: {}, 4.1: {}},
			t:      map[float32]struct{}{2: {}, 6.0: {}},
			expect: false,
		},
		{
			name:   "test Float32 IsSubset, s ∩ s ≠ Ø && s ∩ s ≠ s",
			s:      map[float32]struct{}{1: {}, 4.1: {}},
			t:      map[float32]struct{}{1: {}, 6.0: {}},
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

func TestFloat32_IsSuperset(t *testing.T) {
	testcases := []struct {
		name   string
		s      Float32
		t      Float32
		expect bool
	}{
		{
			name:   "test Float32 IsSuperset, s is empty",
			s:      Float32{},
			t:      map[float32]struct{}{2: {}, 1.3: {}, 4.1: {}},
			expect: false,
		},
		{
			name:   "test Float32 IsSuperset, s is empty",
			s:      map[float32]struct{}{2: {}, 1.3: {}, 4.1: {}},
			t:      Float32{},
			expect: true,
		},
		{
			name:   "test Float32 IsSuperset, s ⊂ s",
			s:      map[float32]struct{}{2: {}, 1.3: {}},
			t:      map[float32]struct{}{2: {}, 1.3: {}, 4.1: {}},
			expect: false,
		},
		{
			name:   "test Float32 IsSuperset, s ⊃ s",
			s:      map[float32]struct{}{2: {}, 1.3: {}, 4.1: {}},
			t:      map[float32]struct{}{2: {}, 1.3: {}},
			expect: true,
		},
		{
			name:   "test Float32 IsSuperset, s = s",
			s:      map[float32]struct{}{2: {}, 1.3: {}, 4.1: {}},
			t:      map[float32]struct{}{2: {}, 1.3: {}, 4.1: {}},
			expect: true,
		},
		{
			name:   "test Float32 IsSuperset, s ∩ s = Ø",
			s:      map[float32]struct{}{1: {}, 4.1: {}},
			t:      map[float32]struct{}{2: {}, 6.0: {}},
			expect: false,
		},
		{
			name:   "test Float32 IsSuperset, s ∩ s ≠ Ø && s ∩ s ≠ s",
			s:      map[float32]struct{}{1: {}, 4.1: {}},
			t:      map[float32]struct{}{1: {}, 6.0: {}},
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

func TestFloat32_Equal(t *testing.T) {
	testcases := []struct {
		name   string
		s      Float32
		t      Float32
		expect bool
	}{
		{
			name:   "test Float32 Equal, s is empty",
			s:      Float32{},
			t:      map[float32]struct{}{2: {}, 1.3: {}, 4.1: {}},
			expect: false,
		},
		{
			name:   "test Float32 Equal, s is empty",
			s:      map[float32]struct{}{2: {}, 1.3: {}, 4.1: {}},
			t:      Float32{},
			expect: false,
		},
		{
			name:   "test Float32 Equal, s ⊂ s",
			s:      map[float32]struct{}{2: {}, 1.3: {}},
			t:      map[float32]struct{}{2: {}, 1.3: {}, 4.1: {}},
			expect: false,
		},
		{
			name:   "test Float32 Equal, s ⊃ s",
			s:      map[float32]struct{}{2: {}, 1.3: {}, 4.1: {}},
			t:      map[float32]struct{}{2: {}, 1.3: {}},
			expect: false,
		},
		{
			name:   "test Float32 Equal, s = s",
			s:      map[float32]struct{}{2: {}, 1.3: {}, 4.1: {}},
			t:      map[float32]struct{}{2: {}, 1.3: {}, 4.1: {}},
			expect: true,
		},
		{
			name:   "test Float32 Equal, s ∩ s = Ø",
			s:      map[float32]struct{}{1: {}, 4.1: {}},
			t:      map[float32]struct{}{2: {}, 6.0: {}},
			expect: false,
		},
		{
			name:   "test Float32 Equal, s ∩ s ≠ Ø && s ∩ s ≠ s",
			s:      map[float32]struct{}{1: {}, 4.1: {}},
			t:      map[float32]struct{}{1: {}, 6.0: {}},
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

func TestFloat32_Copy(t *testing.T) {
	testcases := []struct {
		name   string
		s      Float32
		expect []float32
	}{
		{
			name:   "test Float32 Copy, s is empty",
			s:      Float32{},
			expect: []float32{},
		},
		{
			name:   "test Float32 Copy, s is not empty",
			s:      map[float32]struct{}{2: {}, 1.3: {}, 4.1: {}},
			expect: []float32{2, 1.3, 4.1},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		actual := tc.s.Copy()
		validateFloat32(t, actual, tc.expect)
	}
}

func TestFloat32_String(t *testing.T) {
	testcases := []struct {
		name   string
		s      Float32
		expect string
	}{
		{
			name:   "test Float32 String, s is empty",
			s:      Float32{},
			expect: "[]",
		},
		{
			name:   "test Float32 String, s is not empty",
			s:      map[float32]struct{}{1: {}},
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

func validateFloat32(t *testing.T, actual Float32, expect []float32) {
	if len(expect) != len(actual) {
		t.Errorf("expect set len: %d, but got: %d", len(expect), len(actual))
	}
	slice2String := func(i []float32) string {
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
