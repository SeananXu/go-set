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
	"fmt"
	"sort"
	"strings"
)

// String is a string collection that contains no duplicate elements, without any particular order.
// It supports typical set operations: Core set-theoretical operations, Static sets, Dynamic
// sets, Additional operations.
type String map[string]struct{}

// NewString initializes a new String.
func NewString(elements ...string) String {
	s := String{}
	s.Add(elements...)
	return s
}

// NewStringWithSize initializes a new String with the specified size.
func NewStringWithSize(size int) String {
	return make(map[string]struct{}, size)
}

// Add adds the elements to String, if it is not present already.
func (s String) Add(elements ...string) {
	for _, element := range elements {
		s[element] = struct{}{}
	}
}

// Remove removes the element from String, if it is present.
func (s String) Remove(elements ...string) {
	for _, element := range elements {
		delete(s, element)
	}
}

// Pop returns an arbitrary element of String, deleting it from String.
// The second value is a bool that is true if the elements existed in
// the String, and false if not.
func (s String) Pop() (string, bool) {
	for k := range s {
		delete(s, k)
		return k, true
	}
	return "", false
}

// Size returns the number of elements in String.
func (s String) Size() int {
	return len(s)
}

// IsEmpty returns whether the String is Empty.
func (s String) IsEmpty() bool {
	return len(s) == 0
}

// Clear removes all items from the String.
func (s *String) Clear() {
	*s = make(map[string]struct{})
}

// Has judges the specified element whether exists in the String.
// it returns true if existed, and false if not.
func (s String) Has(element string) bool {
	_, ok := s[element]
	return ok
}

// HasAll looks for the specified elements to judge
// whether all exist in the String.
// it returns true if existed, and false if not.
func (s String) HasAll(elements ...string) bool {
	for _, element := range elements {
		if _, ok := s[element]; !ok {
			return false
		}
	}
	return true
}

// HasAny looks for the specified elements to judge
// whether at least one of the element exists in the String.
// it returns true if existed, and false if not.
func (s String) HasAny(elements ...string) bool {
	for _, element := range elements {
		if _, ok := s[element]; ok {
			return true
		}
	}
	return false
}

// List returns the all elements as a slice.
func (s String) List() []string {
	var dest []string
	for k := range s {
		dest = append(dest, k)
	}
	return dest
}

// SortedList returns the all elements as a slice sorted by less func.
func (s String) SortedList(less func(i, j string) bool) []string {
	dest := s.List()
	sort.Slice(dest, func(i, j int) bool {
		return less(dest[i], dest[j])
	})
	return dest
}

// EachE traverses the elements in the String, calling do func for each
// String member. the cycle will be stopped when the do func returns error.
// if err is ErrBreakEach, break the cycle and return nil,
// else, break the cycle and return error.
func (s String) EachE(do func(i string) error) error {
	for k := range s {
		if err := do(k); err != nil {
			if err == ErrBreakEach {
				return nil
			}
			return err
		}
	}
	return nil
}

// Each traverses the elements in the String, calling do func for each
// String member.
func (s String) Each(do func(i string)) {
	for k := range s {
		do(k)
	}
}

// Union returns the union of String s and t.
// For example:
// s = {a, b, c}
// s = {a, c, d, e, f}
// s.Union(s) = {a, b, c, d, e, f}
// s.Union(s) = {a, b, c, d, e, f}
// s.Union(s) = s.Union(s)
func (s String) Union(t String) String {
	// in order to reduce the number of growing map, copy the largest map here
	var max, min String
	if s.Size() > t.Size() {
		max = s
		min = t
	} else {
		max = t
		min = s
	}
	if max.Size() == 0 {
		return NewString()
	}
	u := max.Copy()
	for k := range min {
		u[k] = struct{}{}
	}
	return u
}

// Difference returns the difference of String s and t.
// For example:
// s = {a, b, c}
// s = {a, c, d, e, f}
// s.Difference(s) = {b}
// s.Difference(s) = {d, e, f}
func (s String) Difference(t String) String {
	u := NewString()
	for k := range s {
		if !t.Has(k) {
			u.Add(k)
		}
	}
	return u
}

// Intersection returns the intersection of String s and t.
// For example:
// s = {a, b, c}
// s = {a, c, d, e, f}
// s.Intersection(s) = {a, c}
// s.Intersection(s) = {a, c}
// s.Intersection(s) = s.Intersection(s)
func (s String) Intersection(t String) String {
	var max, min String
	if s.Size() > t.Size() {
		max = s
		min = t
	} else {
		max = t
		min = s
	}
	u := NewString()
	if min.Size() > 0 {
		for k := range min {
			if max.Has(k) {
				u[k] = struct{}{}
			}
		}
	}
	return u
}

// SymmetricDifference returns a new String with the elements that are either in this String
// or in the given String, but not in both.
// For example:
// s = {a, c}
// s = {a, b, d}
// s.SymmetricDifference(s) = {c, b, d}
// s.SymmetricDifference(s) = {c, b, d}
// s.SymmetricDifference(s) = s.SymmetricDifference(s)
func (s String) SymmetricDifference(t String) String {
	return s.Difference(t).Union(t.Difference(s))
}

// IsSubset predicates that tests whether the String s is a subset of String t.
// For example:
// s is subset of s
// s = {a, b, c}
// s = {a, b, c, d}
// s is not subset of s
// s = {a, f}
// s = {a, b, c, d}
func (s String) IsSubset(t String) bool {
	for k := range s {
		if !t.Has(k) {
			return false
		}
	}
	return true
}

// IsSuperset predicates that tests whether the String s is a super of String t.
// For example:
// s is super of s
// s = {a, b, c, d}
// s = {a, b, c}
// s is not super of s
// s = {a, f}
// s = {a, b, c, d}
func (s String) IsSuperset(t String) bool {
	return t.IsSubset(s)
}

// Equal predicates that tests whether the String s equals of String t.
// For example:
// s equals of s
// s = {a, b, c}
// s = {a, b, c}
// s does not equal of s
// s = {a, f}
// s = {a, b, c, d}
func (s String) Equal(t String) bool {
	return len(s) == len(t) && s.IsSubset(t)
}

// Copy returns new String that clones from String.
func (s String) Copy() String {
	t := NewStringWithSize(len(s))
	for k := range s {
		t[k] = struct{}{}
	}
	return t
}

// String returns a string representation of String
func (s String) String() string {
	v := make([]string, 0, s.Size())
	for element := range s {
		v = append(v, fmt.Sprintf("%v", element))
	}
	return fmt.Sprintf("[%s]", strings.Join(v, ", "))
}
