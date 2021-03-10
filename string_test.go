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

func TestNewString(t *testing.T) {
	testCases := []struct {
		name   string
		input  []string
		expect []string
	}{
		{
			name:   "test String New, s nothing",
			input:  nil,
			expect: nil,
		},
		{
			name:   "test String New, s nil element",
			input:  []string{""},
			expect: []string{""},
		},
		{
			name:   "test String New, inputs multiple elements",
			input:  []string{"1", "test", "test1.2"},
			expect: []string{"1", "test", "test1.2"},
		},
	}
	for _, tc := range testCases {
		t.Logf("running scenario: %s", tc.name)
		actual := NewString(tc.input...)
		validateString(t, actual, tc.expect)
	}
}

func TestString_Add(t *testing.T) {
	testcases := []struct {
		name   string
		input  []string
		expect []string
	}{
		{
			name:   "test String Add, inputs nothing",
			input:  []string{},
			expect: []string{},
		},
		{
			name:   "test String Add, inputs multiple elements",
			input:  []string{"1", "test2", "3"},
			expect: []string{"1", "test2", "3"},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		actual := String{}
		actual.Add(tc.input...)
		validateString(t, actual, tc.expect)
	}
}

func TestString_Remove(t *testing.T) {
	testcases := []struct {
		name   string
		s      String
		input  []string
		expect []string
	}{
		{
			name:   "test String Remove, inputs nothing",
			s:      map[string]struct{}{"test1": {}, "test": {}, "test1.2": {}},
			input:  []string{},
			expect: []string{"test1", "test", "test1.2"},
		},
		{
			name:   "test String Remove, inputs multiple exit elements",
			s:      map[string]struct{}{"test1": {}, "test": {}, "test1.2": {}},
			input:  []string{"test1", "test1.2"},
			expect: []string{"test"},
		},
		{
			name:   "test String Remove, inputs multiple non-exit elements",
			s:      map[string]struct{}{"test1": {}, "test": {}, "test1.2": {}},
			input:  []string{"test0", "not exist", "test1.3"},
			expect: []string{"test1", "test", "test1.2"},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		tc.s.Remove(tc.input...)
		validateString(t, tc.s, tc.expect)
	}
}

func TestString_Pop(t *testing.T) {
	testcases := []struct {
		name   string
		s      String
		expect bool
	}{
		{
			name:   "test String Pop, s is not empty",
			s:      map[string]struct{}{"test1": {}, "test": {}, "test1.2": {}},
			expect: true,
		},
		{
			name:   "test String Pop, s is empty",
			s:      map[string]struct{}{},
			expect: false,
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		backup := make(map[string]struct{})
		for k := range tc.s {
			backup[k] = struct{}{}
		}
		key, ok := tc.s.Pop()
		if ok != tc.expect {
			t.Errorf("expect ok: %v, but got: %v", tc.expect, ok)
		}
		delete(backup, key)
		var expect []string
		for key := range backup {
			expect = append(expect, key)
		}
		validateString(t, tc.s, expect)
	}
}

func TestString_Size(t *testing.T) {
	testcases := []struct {
		name   string
		s      String
		expect int
	}{
		{
			name:   "test String Size, s is empty",
			s:      String{},
			expect: 0,
		},
		{
			name:   "test String Size, s is not empty",
			s:      map[string]struct{}{"test1": {}, "test": {}, "test1.2": {}},
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

func TestString_IsEmpty(t *testing.T) {
	testcases := []struct {
		name   string
		s      String
		expect bool
	}{
		{
			name:   "test String Empty, s is empty",
			s:      String{},
			expect: true,
		},
		{
			name:   "test interface Empty, s is not empty",
			s:      map[string]struct{}{"test1": {}, "test": {}, "test1.2": {}},
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

func TestString_Clear(t *testing.T) {
	testcases := []struct {
		name string
		s    String
	}{
		{
			name: "test String Clear, s is empty",
			s:    String{},
		},
		{
			name: "test String Clear, s is not empty",
			s:    map[string]struct{}{"test1": {}, "test": {}, "test1.2": {}},
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

func TestString_Has(t *testing.T) {
	testcases := []struct {
		name   string
		s      String
		input  string
		expect bool
	}{
		{
			name:   "test String Has, s has input element",
			s:      map[string]struct{}{"test1": {}, "test": {}, "test1.2": {}},
			input:  "test1",
			expect: true,
		},
		{
			name:   "test String Has, s does not have element",
			s:      map[string]struct{}{"test1": {}, "test": {}, "test1.2": {}},
			input:  "test2",
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

func TestString_HasAll(t *testing.T) {
	testcases := []struct {
		name   string
		s      String
		input  []string
		expect bool
	}{
		{
			name:   "test String HasAll, set has all input elements",
			s:      map[string]struct{}{"test1": {}, "test": {}, "test1.2": {}},
			input:  []string{"test1", "test"},
			expect: true,
		},
		{
			name:   "test String HasAll, set does not have all input elements",
			s:      map[string]struct{}{"test1": {}, "test": {}, "test1.2": {}},
			input:  []string{"test1", "3"},
			expect: false,
		},
		{
			name:   "test String HasAll, set does not have all input elements, but exist elements in set",
			s:      map[string]struct{}{"test1": {}, "test": {}, "test1.2": {}},
			input:  []string{"test1", "3"},
			expect: false,
		},
		{
			name:   "test String HasAll, input empty",
			s:      map[string]struct{}{"test1": {}, "test": {}, "test1.2": {}},
			input:  []string{},
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

func TestString_HasAny(t *testing.T) {
	testcases := []struct {
		name   string
		s      String
		input  []string
		expect bool
	}{
		{
			name:   "test String HasAny, s has all elements",
			s:      map[string]struct{}{"test1": {}, "test": {}, "test1.2": {}},
			input:  []string{"test1", "test"},
			expect: true,
		},
		{
			name:   "test String HasAny, s does not have all elements, but exist elements in set",
			s:      map[string]struct{}{"test1": {}, "test": {}, "test1.2": {}},
			input:  []string{"test1", "3"},
			expect: true,
		},
		{
			name:   "test String HasAny, s does not has all elements",
			s:      map[string]struct{}{"test1": {}, "test": {}, "test1.2": {}},
			input:  []string{"test2", "3"},
			expect: false,
		},
		{
			name:   "test String HasAll, input empty",
			s:      map[string]struct{}{"test1": {}, "test": {}, "test1.2": {}},
			input:  []string{},
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

func TestString_List(t *testing.T) {
	testcases := []struct {
		name   string
		s      String
		input  []string
		expect bool
	}{
		{
			name: "test String List, s is empty",
			s:    String{},
		},
		{
			name: "test String List, s is not empty",
			s:    map[string]struct{}{"test1": {}, "test": {}, "test1.2": {}},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		actual := tc.s.List()
		validateString(t, tc.s, actual)
	}
}

func TestString_SortedList(t *testing.T) {
	testcases := []struct {
		name   string
		s      String
		expect []string
	}{
		{
			name: "test String List, s is empty",
			s:    String{},
		},
		{
			name:   "test String SortedList, s is not empty",
			s:      map[string]struct{}{"test1": {}, "test2": {}, "test3": {}},
			expect: []string{"test1", "test2", "test3"},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		actual := tc.s.SortedList(func(i, j string) bool {
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

func TestString_Each(t *testing.T) {
	testcases := []struct {
		name   string
		origin String
	}{
		{
			name:   "test String Each",
			origin: map[string]struct{}{"test1": {}, "test2": {}, "test3": {}},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		var expect []string
		tc.origin.Each(func(i string) {
			expect = append(expect, i)
		})
		validateString(t, tc.origin, expect)
	}
}

func TestString_EachE(t *testing.T) {
	inputErr := errors.New("s error")
	testcases := []struct {
		name      string
		origin    String
		breakEach bool
		inputErr  error
		expectLen int
		expectErr error
	}{
		{
			name:      "test String EachE",
			origin:    map[string]struct{}{"test1": {}, "test2": {}, "test3": {}},
			expectLen: 3,
		},
		{
			name:      "test String EachE, return break error",
			origin:    map[string]struct{}{"test1": {}, "test2": {}, "test3": {}},
			breakEach: true,
		},
		{
			name:      "test String Each, returns error",
			origin:    map[string]struct{}{"test1": {}, "test2": {}, "test3": {}},
			inputErr:  inputErr,
			expectErr: inputErr,
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		var expect []string
		err := tc.origin.EachE(func(i string) error {
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
			validateString(t, tc.origin, expect)
		}
	}
}

func TestString_Union(t *testing.T) {
	testcases := []struct {
		name   string
		s      String
		t      String
		expect []string
	}{
		{
			name:   "test String Union, s and s are empty",
			s:      String{},
			t:      String{},
			expect: []string{},
		},
		{
			name:   "test String Union, s is empty",
			s:      String{},
			t:      map[string]struct{}{"test1": {}, "2": {}, "test3.1": {}},
			expect: []string{"test1", "2", "test3.1"},
		},
		{
			name:   "test String Union, s is empty",
			s:      map[string]struct{}{"test1": {}, "2": {}, "test3.1": {}},
			t:      String{},
			expect: []string{"test1", "2", "test3.1"},
		},
		{
			name:   "test String Union, s has same element to s",
			s:      map[string]struct{}{"test1": {}, "2": {}, "test3.1": {}},
			t:      map[string]struct{}{"test1": {}, "3": {}, "test4.1": {}},
			expect: []string{"test1", "2", "test3.1", "test4.1", "3"},
		},
		{
			name:   "test String Union, s does not have same element to s",
			s:      map[string]struct{}{"test1": {}, "2": {}, "test3.1": {}},
			t:      map[string]struct{}{"test2": {}, "3": {}, "test4.1": {}},
			expect: []string{"test1", "test2", "2", "test3.1", "test4.1", "3"},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		actual := tc.s.Union(tc.t)
		validateString(t, actual, tc.expect)
	}
}

func TestString_Difference(t *testing.T) {
	testcases := []struct {
		name   string
		s      String
		t      String
		expect []string
	}{
		{
			name:   "test String Difference, s is empty",
			s:      String{},
			t:      map[string]struct{}{"test2": {}, "3": {}, "test4.1": {}},
			expect: []string{},
		},
		{
			name:   "test String Difference, s is empty",
			s:      map[string]struct{}{"test2": {}, "3": {}, "test4.1": {}},
			t:      String{},
			expect: []string{"test2", "3", "test4.1"},
		},
		{
			name:   "test String Difference, s ⊂ s",
			s:      map[string]struct{}{"test2": {}, "3": {}},
			t:      map[string]struct{}{"test2": {}, "3": {}, "test4.1": {}},
			expect: []string{},
		},
		{
			name:   "test String Difference, s ⊃ s",
			s:      map[string]struct{}{"test2": {}, "3": {}, "test4.1": {}},
			t:      map[string]struct{}{"test2": {}},
			expect: []string{"3", "test4.1"},
		},
		{
			name:   "test String Difference, s = s",
			s:      map[string]struct{}{"test2": {}, "3": {}, "test4.1": {}},
			t:      map[string]struct{}{"test2": {}, "3": {}, "test4.1": {}},
			expect: []string{},
		},
		{
			name:   "test String Difference, s ∩ s = Ø",
			s:      map[string]struct{}{"test2": {}, "3": {}, "test4.1": {}},
			t:      map[string]struct{}{"test1": {}, "6": {}},
			expect: []string{"test2", "3", "test4.1"},
		},
		{
			name:   "test String Difference, s ∩ s ≠ Ø",
			s:      map[string]struct{}{"test2": {}, "3": {}, "test4.1": {}},
			t:      map[string]struct{}{"test2": {}, "3": {}, "test4.2": {}},
			expect: []string{"test4.1"},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		actual := tc.s.Difference(tc.t)
		validateString(t, actual, tc.expect)
	}
}

func TestString_Intersection(t *testing.T) {
	testcases := []struct {
		name   string
		s      String
		t      String
		expect []string
	}{
		{
			name:   "test String Intersection, s is empty",
			s:      String{},
			t:      map[string]struct{}{"test2": {}, "3": {}, "test4.1": {}},
			expect: []string{},
		},
		{
			name:   "test String Intersection, s is empty",
			s:      map[string]struct{}{"test2": {}, "3": {}, "test4.1": {}},
			t:      String{},
			expect: []string{},
		},
		{
			name:   "test String Intersection, s ⊂ s",
			s:      map[string]struct{}{"test2": {}, "3": {}},
			t:      map[string]struct{}{"test2": {}, "3": {}, "test4.1": {}},
			expect: []string{"test2", "3"},
		},
		{
			name:   "test String Intersection, s ⊃ s",
			s:      map[string]struct{}{"test2": {}, "3": {}, "test4.1": {}},
			t:      map[string]struct{}{"test2": {}, "3": {}},
			expect: []string{"test2", "3"},
		},
		{
			name:   "test String Intersection, s = s",
			s:      map[string]struct{}{"test2": {}, "3": {}, "test4.1": {}},
			t:      map[string]struct{}{"test2": {}, "3": {}, "test4.1": {}},
			expect: []string{"test2", "3", "test4.1"},
		},
		{
			name:   "test String Intersection, s ∩ s = Ø",
			s:      map[string]struct{}{"test1": {}, "4": {}},
			t:      map[string]struct{}{"test2": {}, "6": {}},
			expect: []string{},
		},
		{
			name:   "test String Intersection, s ∩ s ≠ Ø",
			s:      map[string]struct{}{"test1": {}, "4": {}},
			t:      map[string]struct{}{"test1": {}, "6": {}},
			expect: []string{"test1"},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		actual := tc.s.Intersection(tc.t)
		validateString(t, actual, tc.expect)
	}
}

func TestString_SymmetricDifference(t *testing.T) {
	testcases := []struct {
		name   string
		s      String
		t      String
		expect []string
	}{
		{
			name:   "test String SymmetricDifference, s is empty",
			s:      String{},
			t:      map[string]struct{}{"test2": {}, "3": {}, "test4.1": {}},
			expect: []string{"test2", "3", "test4.1"},
		},
		{
			name:   "test String SymmetricDifference, s is empty",
			s:      map[string]struct{}{"test2": {}, "3": {}, "test4.1": {}},
			t:      String{},
			expect: []string{"test2", "3", "test4.1"},
		},
		{
			name:   "test String SymmetricDifference, s ⊂ s",
			s:      map[string]struct{}{"test2": {}, "3": {}},
			t:      map[string]struct{}{"test2": {}, "3": {}, "test4.1": {}},
			expect: []string{"test4.1"},
		},
		{
			name:   "test String SymmetricDifference, s ⊃ s",
			s:      map[string]struct{}{"test2": {}, "3": {}, "test4.1": {}},
			t:      map[string]struct{}{"test2": {}, "3": {}},
			expect: []string{"test4.1"},
		},
		{
			name:   "test String SymmetricDifference, s = s",
			s:      map[string]struct{}{"test2": {}, "3": {}, "test4.1": {}},
			t:      map[string]struct{}{"test2": {}, "3": {}, "test4.1": {}},
			expect: []string{},
		},
		{
			name:   "test String SymmetricDifference, s ∩ s = Ø",
			s:      map[string]struct{}{"test1": {}, "4": {}},
			t:      map[string]struct{}{"test2": {}, "6": {}},
			expect: []string{"test1", "test2", "4", "6"},
		},
		{
			name:   "test String SymmetricDifference, s ∩ s ≠ Ø",
			s:      map[string]struct{}{"test1": {}, "4": {}},
			t:      map[string]struct{}{"test1": {}, "6": {}},
			expect: []string{"4", "6"},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		actual := tc.s.SymmetricDifference(tc.t)
		validateString(t, actual, tc.expect)
	}
}

func TestString_IsSubset(t *testing.T) {
	testcases := []struct {
		name   string
		s      String
		t      String
		expect bool
	}{
		{
			name:   "test String IsSubset, s is empty",
			s:      String{},
			t:      map[string]struct{}{"test2": {}, "3": {}, "test4.1": {}},
			expect: true,
		},
		{
			name:   "test String IsSubset, s is empty",
			s:      map[string]struct{}{"test2": {}, "3": {}, "test4.1": {}},
			t:      String{},
			expect: false,
		},
		{
			name:   "test String IsSubset, s ⊂ s",
			s:      map[string]struct{}{"test2": {}, "3": {}},
			t:      map[string]struct{}{"test2": {}, "3": {}, "test4.1": {}},
			expect: true,
		},
		{
			name:   "test String IsSubset, s ⊃ s",
			s:      map[string]struct{}{"test2": {}, "3": {}, "test4.1": {}},
			t:      map[string]struct{}{"test2": {}, "3": {}},
			expect: false,
		},
		{
			name:   "test String IsSubset, s = s",
			s:      map[string]struct{}{"test2": {}, "3": {}, "test4.1": {}},
			t:      map[string]struct{}{"test2": {}, "3": {}, "test4.1": {}},
			expect: true,
		},
		{
			name:   "test String IsSubset, s ∩ s = Ø",
			s:      map[string]struct{}{"test1": {}, "4": {}},
			t:      map[string]struct{}{"test2": {}, "6": {}},
			expect: false,
		},
		{
			name:   "test String IsSubset, s ∩ s ≠ Ø && s ∩ s ≠ s",
			s:      map[string]struct{}{"test1": {}, "4": {}},
			t:      map[string]struct{}{"test1": {}, "6": {}},
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

func TestString_IsSuperset(t *testing.T) {
	testcases := []struct {
		name   string
		s      String
		t      String
		expect bool
	}{
		{
			name:   "test String IsSuperset, s is empty",
			s:      String{},
			t:      map[string]struct{}{"test2": {}, "3": {}, "test4.1": {}},
			expect: false,
		},
		{
			name:   "test String IsSuperset, s is empty",
			s:      map[string]struct{}{"test2": {}, "3": {}, "test4.1": {}},
			t:      String{},
			expect: true,
		},
		{
			name:   "test String IsSuperset, s ⊂ s",
			s:      map[string]struct{}{"test2": {}, "3": {}},
			t:      map[string]struct{}{"test2": {}, "3": {}, "test4.1": {}},
			expect: false,
		},
		{
			name:   "test String IsSuperset, s ⊃ s",
			s:      map[string]struct{}{"test2": {}, "3": {}, "test4.1": {}},
			t:      map[string]struct{}{"test2": {}, "3": {}},
			expect: true,
		},
		{
			name:   "test String IsSuperset, s = s",
			s:      map[string]struct{}{"test2": {}, "3": {}, "test4.1": {}},
			t:      map[string]struct{}{"test2": {}, "3": {}, "test4.1": {}},
			expect: true,
		},
		{
			name:   "test String IsSuperset, s ∩ s = Ø",
			s:      map[string]struct{}{"test1": {}, "4": {}},
			t:      map[string]struct{}{"test2": {}, "6": {}},
			expect: false,
		},
		{
			name:   "test String IsSuperset, s ∩ s ≠ Ø && s ∩ s ≠ s",
			s:      map[string]struct{}{"test1": {}, "4": {}},
			t:      map[string]struct{}{"test1": {}, "6": {}},
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

func TestString_Equal(t *testing.T) {
	testcases := []struct {
		name   string
		s      String
		t      String
		expect bool
	}{
		{
			name:   "test String Equal, s is empty",
			s:      String{},
			t:      map[string]struct{}{"test2": {}, "3": {}, "test4.1": {}},
			expect: false,
		},
		{
			name:   "test String Equal, s is empty",
			s:      map[string]struct{}{"test2": {}, "3": {}, "test4.1": {}},
			t:      String{},
			expect: false,
		},
		{
			name:   "test String Equal, s ⊂ s",
			s:      map[string]struct{}{"test2": {}, "3": {}},
			t:      map[string]struct{}{"test2": {}, "3": {}, "test4.1": {}},
			expect: false,
		},
		{
			name:   "test String Equal, s ⊃ s",
			s:      map[string]struct{}{"test2": {}, "3": {}, "test4.1": {}},
			t:      map[string]struct{}{"test2": {}, "3": {}},
			expect: false,
		},
		{
			name:   "test String Equal, s = s",
			s:      map[string]struct{}{"test2": {}, "3": {}, "test4.1": {}},
			t:      map[string]struct{}{"test2": {}, "3": {}, "test4.1": {}},
			expect: true,
		},
		{
			name:   "test String Equal, s ∩ s = Ø",
			s:      map[string]struct{}{"test1": {}, "4": {}},
			t:      map[string]struct{}{"test2": {}, "6": {}},
			expect: false,
		},
		{
			name:   "test String Equal, s ∩ s ≠ Ø && s ∩ s ≠ s",
			s:      map[string]struct{}{"test1": {}, "4": {}},
			t:      map[string]struct{}{"test1": {}, "6": {}},
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

func TestString_Copy(t *testing.T) {
	testcases := []struct {
		name   string
		s      String
		expect []string
	}{
		{
			name:   "test String Copy, s is empty",
			s:      String{},
			expect: []string{},
		},
		{
			name:   "test String Copy, s is not empty",
			s:      map[string]struct{}{"test2": {}, "3": {}, "test4.1": {}},
			expect: []string{"test2", "3", "test4.1"},
		},
	}
	for _, tc := range testcases {
		t.Logf("running scenario: %s", tc.name)
		actual := tc.s.Copy()
		validateString(t, actual, tc.expect)
	}
}

func TestString_String(t *testing.T) {
	testcases := []struct {
		name   string
		s      String
		expect string
	}{
		{
			name:   "test String String, s is empty",
			s:      String{},
			expect: "[]",
		},
		{
			name:   "test String String, s is not empty",
			s:      map[string]struct{}{"test1": {}},
			expect: "[test1]",
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

func validateString(t *testing.T, actual String, expect []string) {
	if len(expect) != len(actual) {
		t.Errorf("expect set len: %d, but got: %d", len(expect), len(actual))
	}
	slice2String := func(i []string) string {
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
