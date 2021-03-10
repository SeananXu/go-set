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

func TestNewInterface(t *testing.T) {
	testCases := []struct {
		name   string
		input  []interface{}
		expect []interface{}
	}{
		{
			name:   "test Interface New, s nothing",
			input:  nil,
			expect: nil,
		},
		{
			name:   "test Interface New, s nil element",
			input:  []interface{}{nil},
			expect: []interface{}{nil},
		},
		{
			name:   "test Interface New, inputs multiple elements",
			input:  []interface{}{1, "test", 1.2},
			expect: []interface{}{1, "test", 1.2},
		},
	}
	for _, tc := range testCases {
		t.Logf("running scenario: %s", tc.name)
		actual := NewInterface(tc.input...)
		validateInterface(t, actual, tc.expect)
	}
}

func TestInterface_Add(t *testing.T) {
	testcases := []struct {
		name   string
		input  []interface{}
		expect []interface{}
	}{
		{
			name:   "test Interface Add, inputs nothing",
			input:  []interface{}{},
			expect: []interface{}{},
		},
		{
			name:   "test Interface Add, inputs multiple elements",
			input:  []interface{}{1, 2, "3"},
			expect: []interface{}{1, 2, "3"},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		actual := Interface{}
		actual.Add(tc.input...)
		validateInterface(t, actual, tc.expect)
	}
}

func TestInterface_Remove(t *testing.T) {
	testcases := []struct {
		name   string
		s      Interface
		input  []interface{}
		expect []interface{}
	}{
		{
			name:   "test Interface Remove, inputs nothing",
			s:      map[interface{}]struct{}{1: {}, "test": {}, 1.2: {}},
			input:  []interface{}{},
			expect: []interface{}{1, "test", 1.2},
		},
		{
			name:   "test Interface Remove, inputs multiple exit elements",
			s:      map[interface{}]struct{}{1: {}, "test": {}, 1.2: {}},
			input:  []interface{}{1, 1.2},
			expect: []interface{}{"test"},
		},
		{
			name:   "test Interface Remove, inputs multiple non-exit elements",
			s:      map[interface{}]struct{}{1: {}, "test": {}, 1.2: {}},
			input:  []interface{}{0, "not exist", 1.3},
			expect: []interface{}{1, "test", 1.2},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		tc.s.Remove(tc.input...)
		validateInterface(t, tc.s, tc.expect)
	}
}

func TestInterface_Pop(t *testing.T) {
	testcases := []struct {
		name   string
		s      Interface
		expect bool
	}{
		{
			name:   "test Interface Pop, s is not empty",
			s:      map[interface{}]struct{}{1: {}, "test": {}, 1.2: {}},
			expect: true,
		},
		{
			name:   "test Interface Pop, s is empty",
			s:      map[interface{}]struct{}{},
			expect: false,
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		backup := make(map[interface{}]struct{})
		for k := range tc.s {
			backup[k] = struct{}{}
		}
		key, ok := tc.s.Pop()
		if ok != tc.expect {
			t.Errorf("expect ok: %v, but got: %v", tc.expect, ok)
		}
		delete(backup, key)
		var expect []interface{}
		for key := range backup {
			expect = append(expect, key)
		}
		validateInterface(t, tc.s, expect)
	}
}

func TestInterface_Size(t *testing.T) {
	testcases := []struct {
		name   string
		s      Interface
		expect int
	}{
		{
			name:   "test Interface Size, s is empty",
			s:      Interface{},
			expect: 0,
		},
		{
			name:   "test Interface Size, s is not empty",
			s:      map[interface{}]struct{}{1: {}, "test": {}, 1.2: {}},
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

func TestInterface_IsEmpty(t *testing.T) {
	testcases := []struct {
		name   string
		s      Interface
		expect bool
	}{
		{
			name:   "test Interface Empty, s is empty",
			s:      Interface{},
			expect: true,
		},
		{
			name:   "test interface Empty, s is not empty",
			s:      map[interface{}]struct{}{1: {}, "test": {}, 1.2: {}},
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

func TestInterface_Clear(t *testing.T) {
	testcases := []struct {
		name string
		s    Interface
	}{
		{
			name: "test Interface Clear, s is empty",
			s:    Interface{},
		},
		{
			name: "test Interface Clear, s is not empty",
			s:    map[interface{}]struct{}{1: {}, "test": {}, 1.2: {}},
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

func TestInterface_Has(t *testing.T) {
	testcases := []struct {
		name   string
		s      Interface
		input  interface{}
		expect bool
	}{
		{
			name:   "test Interface Has, s has input element",
			s:      map[interface{}]struct{}{1: {}, "test": {}, 1.2: {}},
			input:  1,
			expect: true,
		},
		{
			name:   "test Interface Has, s does not have element",
			s:      map[interface{}]struct{}{1: {}, "test": {}, 1.2: {}},
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

func TestInterface_HasAll(t *testing.T) {
	testcases := []struct {
		name   string
		s      Interface
		input  []interface{}
		expect bool
	}{
		{
			name:   "test Interface HasAll, set has all input elements",
			s:      map[interface{}]struct{}{1: {}, "test": {}, 1.2: {}},
			input:  []interface{}{1, "test"},
			expect: true,
		},
		{
			name:   "test Interface HasAll, set does not have all input elements",
			s:      map[interface{}]struct{}{1: {}, "test": {}, 1.2: {}},
			input:  []interface{}{"1", "3"},
			expect: false,
		},
		{
			name:   "test Interface HasAll, set does not have all input elements, but exist elements in set",
			s:      map[interface{}]struct{}{1: {}, "test": {}, 1.2: {}},
			input:  []interface{}{1, "3"},
			expect: false,
		},
		{
			name:   "test Interface HasAll, input empty",
			s:      map[interface{}]struct{}{1: {}, "test": {}, 1.2: {}},
			input:  []interface{}{},
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

func TestInterface_HasAny(t *testing.T) {
	testcases := []struct {
		name   string
		s      Interface
		input  []interface{}
		expect bool
	}{
		{
			name:   "test Interface HasAny, s has all elements",
			s:      map[interface{}]struct{}{1: {}, "test": {}, 1.2: {}},
			input:  []interface{}{1, "test"},
			expect: true,
		},
		{
			name:   "test Interface HasAny, s does not have all elements, but exist elements in set",
			s:      map[interface{}]struct{}{1: {}, "test": {}, 1.2: {}},
			input:  []interface{}{1, "3"},
			expect: true,
		},
		{
			name:   "test Interface HasAny, s does not has all elements",
			s:      map[interface{}]struct{}{1: {}, "test": {}, 1.2: {}},
			input:  []interface{}{2, "3"},
			expect: false,
		},
		{
			name:   "test Interface HasAll, input empty",
			s:      map[interface{}]struct{}{1: {}, "test": {}, 1.2: {}},
			input:  []interface{}{},
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

func TestInterface_List(t *testing.T) {
	testcases := []struct {
		name   string
		s      Interface
		input  []interface{}
		expect bool
	}{
		{
			name: "test Interface List, s is empty",
			s:    Interface{},
		},
		{
			name: "test Interface List, s is not empty",
			s:    map[interface{}]struct{}{1: {}, "test": {}, 1.2: {}},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		actual := tc.s.List()
		validateInterface(t, tc.s, actual)
	}
}

func TestInterface_SortedList(t *testing.T) {
	testcases := []struct {
		name   string
		s      Interface
		expect []interface{}
	}{
		{
			name: "test Interface List, s is empty",
			s:    Interface{},
		},
		{
			name:   "test Interface SortedList, s is not empty",
			s:      map[interface{}]struct{}{1: {}, 2: {}, 3: {}},
			expect: []interface{}{1, 2, 3},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		actual := tc.s.SortedList(func(i, j interface{}) bool {
			return i.(int) < j.(int)
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

func TestInterface_Each(t *testing.T) {
	testcases := []struct {
		name   string
		origin Interface
	}{
		{
			name:   "test Interface Each",
			origin: map[interface{}]struct{}{1: {}, 2: {}, 3: {}},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		var expect []interface{}
		tc.origin.Each(func(i interface{}) {
			expect = append(expect, i)
		})
		validateInterface(t, tc.origin, expect)
	}
}

func TestInterface_EachE(t *testing.T) {
	inputErr := errors.New("s error")
	testcases := []struct {
		name      string
		origin    Interface
		breakEach bool
		inputErr  error
		expectLen int
		expectErr error
	}{
		{
			name:      "test Interface EachE",
			origin:    map[interface{}]struct{}{1: {}, 2: {}, 3: {}},
			expectLen: 3,
		},
		{
			name:      "test Interface EachE, return break error",
			origin:    map[interface{}]struct{}{1: {}, 2: {}, 3: {}},
			breakEach: true,
		},
		{
			name:      "test Interface Each, returns error",
			origin:    map[interface{}]struct{}{1: {}, 2: {}, 3: {}},
			inputErr:  inputErr,
			expectErr: inputErr,
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		var expect []interface{}
		err := tc.origin.EachE(func(i interface{}) error {
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
			validateInterface(t, tc.origin, expect)
		}
	}
}

func TestInterface_Union(t *testing.T) {
	testcases := []struct {
		name   string
		s      Interface
		t      Interface
		expect []interface{}
	}{
		{
			name:   "test Interface Union, s and s are empty",
			s:      Interface{},
			t:      Interface{},
			expect: []interface{}{},
		},
		{
			name:   "test Interface Union, s is empty",
			s:      Interface{},
			t:      map[interface{}]struct{}{1: {}, "2": {}, 3.1: {}},
			expect: []interface{}{1, "2", 3.1},
		},
		{
			name:   "test Interface Union, s is empty",
			s:      map[interface{}]struct{}{1: {}, "2": {}, 3.1: {}},
			t:      Interface{},
			expect: []interface{}{1, "2", 3.1},
		},
		{
			name:   "test Interface Union, s has same element to s",
			s:      map[interface{}]struct{}{1: {}, "2": {}, 3.1: {}},
			t:      map[interface{}]struct{}{1: {}, "3": {}, 4.1: {}},
			expect: []interface{}{1, "2", 3.1, 4.1, "3"},
		},
		{
			name:   "test Interface Union, s does not have same element to s",
			s:      map[interface{}]struct{}{1: {}, "2": {}, 3.1: {}},
			t:      map[interface{}]struct{}{2: {}, "3": {}, 4.1: {}},
			expect: []interface{}{1, 2, "2", 3.1, 4.1, "3"},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		actual := tc.s.Union(tc.t)
		validateInterface(t, actual, tc.expect)
	}
}

func TestInterface_Difference(t *testing.T) {
	testcases := []struct {
		name   string
		s      Interface
		t      Interface
		expect []interface{}
	}{
		{
			name:   "test Interface Difference, s is empty",
			s:      Interface{},
			t:      map[interface{}]struct{}{2: {}, "3": {}, 4.1: {}},
			expect: []interface{}{},
		},
		{
			name:   "test Interface Difference, s is empty",
			s:      map[interface{}]struct{}{2: {}, "3": {}, 4.1: {}},
			t:      Interface{},
			expect: []interface{}{2, "3", 4.1},
		},
		{
			name:   "test Interface Difference, s ⊂ s",
			s:      map[interface{}]struct{}{2: {}, "3": {}},
			t:      map[interface{}]struct{}{2: {}, "3": {}, 4.1: {}},
			expect: []interface{}{},
		},
		{
			name:   "test Interface Difference, s ⊃ s",
			s:      map[interface{}]struct{}{2: {}, "3": {}, 4.1: {}},
			t:      map[interface{}]struct{}{2: {}},
			expect: []interface{}{"3", 4.1},
		},
		{
			name:   "test Interface Difference, s = s",
			s:      map[interface{}]struct{}{2: {}, "3": {}, 4.1: {}},
			t:      map[interface{}]struct{}{2: {}, "3": {}, 4.1: {}},
			expect: []interface{}{},
		},
		{
			name:   "test Interface Difference, s ∩ s = Ø",
			s:      map[interface{}]struct{}{2: {}, "3": {}, 4.1: {}},
			t:      map[interface{}]struct{}{1: {}, "6": {}},
			expect: []interface{}{2, "3", 4.1},
		},
		{
			name:   "test Interface Difference, s ∩ s ≠ Ø",
			s:      map[interface{}]struct{}{2: {}, "3": {}, 4.1: {}},
			t:      map[interface{}]struct{}{2: {}, "3": {}, 4.2: {}},
			expect: []interface{}{4.1},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		actual := tc.s.Difference(tc.t)
		validateInterface(t, actual, tc.expect)
	}
}

func TestInterface_Intersection(t *testing.T) {
	testcases := []struct {
		name   string
		s      Interface
		t      Interface
		expect []interface{}
	}{
		{
			name:   "test Interface Intersection, s is empty",
			s:      Interface{},
			t:      map[interface{}]struct{}{2: {}, "3": {}, 4.1: {}},
			expect: []interface{}{},
		},
		{
			name:   "test Interface Intersection, s is empty",
			s:      map[interface{}]struct{}{2: {}, "3": {}, 4.1: {}},
			t:      Interface{},
			expect: []interface{}{},
		},
		{
			name:   "test Interface Intersection, s ⊂ s",
			s:      map[interface{}]struct{}{2: {}, "3": {}},
			t:      map[interface{}]struct{}{2: {}, "3": {}, 4.1: {}},
			expect: []interface{}{2, "3"},
		},
		{
			name:   "test Interface Intersection, s ⊃ s",
			s:      map[interface{}]struct{}{2: {}, "3": {}, 4.1: {}},
			t:      map[interface{}]struct{}{2: {}, "3": {}},
			expect: []interface{}{2, "3"},
		},
		{
			name:   "test Interface Intersection, s = s",
			s:      map[interface{}]struct{}{2: {}, "3": {}, 4.1: {}},
			t:      map[interface{}]struct{}{2: {}, "3": {}, 4.1: {}},
			expect: []interface{}{2, "3", 4.1},
		},
		{
			name:   "test Interface Intersection, s ∩ s = Ø",
			s:      map[interface{}]struct{}{1: {}, "4": {}},
			t:      map[interface{}]struct{}{2: {}, "6": {}},
			expect: []interface{}{},
		},
		{
			name:   "test Interface Intersection, s ∩ s ≠ Ø",
			s:      map[interface{}]struct{}{1: {}, "4": {}},
			t:      map[interface{}]struct{}{1: {}, "6": {}},
			expect: []interface{}{1},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		actual := tc.s.Intersection(tc.t)
		validateInterface(t, actual, tc.expect)
	}
}

func TestInterface_SymmetricDifference(t *testing.T) {
	testcases := []struct {
		name   string
		s      Interface
		t      Interface
		expect []interface{}
	}{
		{
			name:   "test Interface SymmetricDifference, s is empty",
			s:      Interface{},
			t:      map[interface{}]struct{}{2: {}, "3": {}, 4.1: {}},
			expect: []interface{}{2, "3", 4.1},
		},
		{
			name:   "test Interface SymmetricDifference, s is empty",
			s:      map[interface{}]struct{}{2: {}, "3": {}, 4.1: {}},
			t:      Interface{},
			expect: []interface{}{2, "3", 4.1},
		},
		{
			name:   "test Interface SymmetricDifference, s ⊂ s",
			s:      map[interface{}]struct{}{2: {}, "3": {}},
			t:      map[interface{}]struct{}{2: {}, "3": {}, 4.1: {}},
			expect: []interface{}{4.1},
		},
		{
			name:   "test Interface SymmetricDifference, s ⊃ s",
			s:      map[interface{}]struct{}{2: {}, "3": {}, 4.1: {}},
			t:      map[interface{}]struct{}{2: {}, "3": {}},
			expect: []interface{}{4.1},
		},
		{
			name:   "test Interface SymmetricDifference, s = s",
			s:      map[interface{}]struct{}{2: {}, "3": {}, 4.1: {}},
			t:      map[interface{}]struct{}{2: {}, "3": {}, 4.1: {}},
			expect: []interface{}{},
		},
		{
			name:   "test Interface SymmetricDifference, s ∩ s = Ø",
			s:      map[interface{}]struct{}{1: {}, "4": {}},
			t:      map[interface{}]struct{}{2: {}, "6": {}},
			expect: []interface{}{1, 2, "4", "6"},
		},
		{
			name:   "test Interface SymmetricDifference, s ∩ s ≠ Ø",
			s:      map[interface{}]struct{}{1: {}, "4": {}},
			t:      map[interface{}]struct{}{1: {}, "6": {}},
			expect: []interface{}{"4", "6"},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		actual := tc.s.SymmetricDifference(tc.t)
		validateInterface(t, actual, tc.expect)
	}
}

func TestInterface_IsSubset(t *testing.T) {
	testcases := []struct {
		name   string
		s      Interface
		t      Interface
		expect bool
	}{
		{
			name:   "test Interface IsSubset, s is empty",
			s:      Interface{},
			t:      map[interface{}]struct{}{2: {}, "3": {}, 4.1: {}},
			expect: true,
		},
		{
			name:   "test Interface IsSubset, s is empty",
			s:      map[interface{}]struct{}{2: {}, "3": {}, 4.1: {}},
			t:      Interface{},
			expect: false,
		},
		{
			name:   "test Interface IsSubset, s ⊂ s",
			s:      map[interface{}]struct{}{2: {}, "3": {}},
			t:      map[interface{}]struct{}{2: {}, "3": {}, 4.1: {}},
			expect: true,
		},
		{
			name:   "test Interface IsSubset, s ⊃ s",
			s:      map[interface{}]struct{}{2: {}, "3": {}, 4.1: {}},
			t:      map[interface{}]struct{}{2: {}, "3": {}},
			expect: false,
		},
		{
			name:   "test Interface IsSubset, s = s",
			s:      map[interface{}]struct{}{2: {}, "3": {}, 4.1: {}},
			t:      map[interface{}]struct{}{2: {}, "3": {}, 4.1: {}},
			expect: true,
		},
		{
			name:   "test Interface IsSubset, s ∩ s = Ø",
			s:      map[interface{}]struct{}{1: {}, "4": {}},
			t:      map[interface{}]struct{}{2: {}, "6": {}},
			expect: false,
		},
		{
			name:   "test Interface IsSubset, s ∩ s ≠ Ø && s ∩ s ≠ s",
			s:      map[interface{}]struct{}{1: {}, "4": {}},
			t:      map[interface{}]struct{}{1: {}, "6": {}},
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

func TestInterface_IsSuperset(t *testing.T) {
	testcases := []struct {
		name   string
		s      Interface
		t      Interface
		expect bool
	}{
		{
			name:   "test Interface IsSuperset, s is empty",
			s:      Interface{},
			t:      map[interface{}]struct{}{2: {}, "3": {}, 4.1: {}},
			expect: false,
		},
		{
			name:   "test Interface IsSuperset, s is empty",
			s:      map[interface{}]struct{}{2: {}, "3": {}, 4.1: {}},
			t:      Interface{},
			expect: true,
		},
		{
			name:   "test Interface IsSuperset, s ⊂ s",
			s:      map[interface{}]struct{}{2: {}, "3": {}},
			t:      map[interface{}]struct{}{2: {}, "3": {}, 4.1: {}},
			expect: false,
		},
		{
			name:   "test Interface IsSuperset, s ⊃ s",
			s:      map[interface{}]struct{}{2: {}, "3": {}, 4.1: {}},
			t:      map[interface{}]struct{}{2: {}, "3": {}},
			expect: true,
		},
		{
			name:   "test Interface IsSuperset, s = s",
			s:      map[interface{}]struct{}{2: {}, "3": {}, 4.1: {}},
			t:      map[interface{}]struct{}{2: {}, "3": {}, 4.1: {}},
			expect: true,
		},
		{
			name:   "test Interface IsSuperset, s ∩ s = Ø",
			s:      map[interface{}]struct{}{1: {}, "4": {}},
			t:      map[interface{}]struct{}{2: {}, "6": {}},
			expect: false,
		},
		{
			name:   "test Interface IsSuperset, s ∩ s ≠ Ø && s ∩ s ≠ s",
			s:      map[interface{}]struct{}{1: {}, "4": {}},
			t:      map[interface{}]struct{}{1: {}, "6": {}},
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

func TestInterface_Equal(t *testing.T) {
	testcases := []struct {
		name   string
		s      Interface
		t      Interface
		expect bool
	}{
		{
			name:   "test Interface Equal, s is empty",
			s:      Interface{},
			t:      map[interface{}]struct{}{2: {}, "3": {}, 4.1: {}},
			expect: false,
		},
		{
			name:   "test Interface Equal, s is empty",
			s:      map[interface{}]struct{}{2: {}, "3": {}, 4.1: {}},
			t:      Interface{},
			expect: false,
		},
		{
			name:   "test Interface Equal, s ⊂ s",
			s:      map[interface{}]struct{}{2: {}, "3": {}},
			t:      map[interface{}]struct{}{2: {}, "3": {}, 4.1: {}},
			expect: false,
		},
		{
			name:   "test Interface Equal, s ⊃ s",
			s:      map[interface{}]struct{}{2: {}, "3": {}, 4.1: {}},
			t:      map[interface{}]struct{}{2: {}, "3": {}},
			expect: false,
		},
		{
			name:   "test Interface Equal, s = s",
			s:      map[interface{}]struct{}{2: {}, "3": {}, 4.1: {}},
			t:      map[interface{}]struct{}{2: {}, "3": {}, 4.1: {}},
			expect: true,
		},
		{
			name:   "test Interface Equal, s ∩ s = Ø",
			s:      map[interface{}]struct{}{1: {}, "4": {}},
			t:      map[interface{}]struct{}{2: {}, "6": {}},
			expect: false,
		},
		{
			name:   "test Interface Equal, s ∩ s ≠ Ø && s ∩ s ≠ s",
			s:      map[interface{}]struct{}{1: {}, "4": {}},
			t:      map[interface{}]struct{}{1: {}, "6": {}},
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

func TestInterface_Copy(t *testing.T) {
	testcases := []struct {
		name   string
		s      Interface
		expect []interface{}
	}{
		{
			name:   "test Interface Copy, s is empty",
			s:      Interface{},
			expect: []interface{}{},
		},
		{
			name:   "test Interface Copy, s is not empty",
			s:      map[interface{}]struct{}{2: {}, "3": {}, 4.1: {}},
			expect: []interface{}{2, "3", 4.1},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		actual := tc.s.Copy()
		validateInterface(t, actual, tc.expect)
	}
}

func TestInterface_String(t *testing.T) {
	testcases := []struct {
		name   string
		s      Interface
		expect string
	}{
		{
			name:   "test Interface String, s is empty",
			s:      Interface{},
			expect: "[]",
		},
		{
			name:   "test Interface String, s is not empty",
			s:      map[interface{}]struct{}{1: {}},
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

func validateInterface(t *testing.T, actual Interface, expect []interface{}) {
	if len(expect) != len(actual) {
		t.Errorf("expect set len: %d, but got: %d", len(expect), len(actual))
	}
	slice2String := func(i []interface{}) string {
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
