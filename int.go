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

// Int is a int collection that contains no duplicate elements, without any particular order.
// It supports typical set operations: Core set-theoretical operations, Static sets, Dynamic
// sets, Additional operations.
type Int map[int]struct{}

// NewInt initializes a new Int.
func NewInt(elements ...int) Int {
	s := Int{}
	s.Add(elements...)
	return s
}

// NewIntWithSize initializes a new Int with the specified size.
func NewIntWithSize(size int) Int {
	return make(map[int]struct{}, size)
}

// Add adds the elements to Int, if it is not present already.
func (s Int) Add(elements ...int) {
	for _, element := range elements {
		s[element] = struct{}{}
	}
}

// Remove removes the element from Int, if it is present.
func (s Int) Remove(elements ...int) {
	for _, element := range elements {
		delete(s, element)
	}
}

// Pop returns an arbitrary element of Int, deleting it from Int.
// The second value is a bool that is true if the elements existed in
// the Int, and false if not.
func (s Int) Pop() (int, bool) {
	for k := range s {
		delete(s, k)
		return k, true
	}
	return 0, false
}

// Size returns the number of elements in Int.
func (s Int) Size() int {
	return len(s)
}

// IsEmpty returns whether the Int is Empty.
func (s Int) IsEmpty() bool {
	return len(s) == 0
}

// Clear removes all items from the Int.
func (s *Int) Clear() {
	*s = make(map[int]struct{})
}

// Has judges the specified element whether exists in the Int.
// it returns true if existed, and false if not.
func (s Int) Has(element int) bool {
	_, ok := s[element]
	return ok
}

// HasAll looks for the specified elements to judge
// whether all exist in the Int.
// it returns true if existed, and false if not.
func (s Int) HasAll(elements ...int) bool {
	for _, element := range elements {
		if _, ok := s[element]; !ok {
			return false
		}
	}
	return true
}

// HasAny looks for the specified elements to judge
// whether at least one of the element exists in the Int.
// it returns true if existed, and false if not.
func (s Int) HasAny(elements ...int) bool {
	for _, element := range elements {
		if _, ok := s[element]; ok {
			return true
		}
	}
	return false
}

// List returns the all elements as a slice.
func (s Int) List() []int {
	var dest []int
	for k := range s {
		dest = append(dest, k)
	}
	return dest
}

// SortedList returns the all elements as a slice sorted by less func.
func (s Int) SortedList(less func(i, j int) bool) []int {
	dest := s.List()
	sort.Slice(dest, func(i, j int) bool {
		return less(dest[i], dest[j])
	})
	return dest
}

// EachE traverses the elements in the Int, calling do func for each
// Int member. the cycle will be stopped when the do func returns error.
// if err is ErrBreakEach, break the cycle and return nil,
// else, break the cycle and return error.
func (s Int) EachE(do func(i int) error) error {
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

// Each traverses the elements in the Int, calling do func for each
// Int member.
func (s Int) Each(do func(i int)) {
	for k := range s {
		do(k)
	}
}

// Union returns the union of Int s and t.
// For example:
// s = {a, b, c}
// s = {a, c, d, e, f}
// s.Union(s) = {a, b, c, d, e, f}
// s.Union(s) = {a, b, c, d, e, f}
// s.Union(s) = s.Union(s)
func (s Int) Union(t Int) Int {
	// in order to reduce the number of growing map, copy the largest map here
	var max, min Int
	if s.Size() > t.Size() {
		max = s
		min = t
	} else {
		max = t
		min = s
	}
	if max.Size() == 0 {
		return NewInt()
	}
	u := max.Copy()
	for k := range min {
		u[k] = struct{}{}
	}
	return u
}

// Difference returns the difference of Int s and t.
// For example:
// s = {a, b, c}
// s = {a, c, d, e, f}
// s.Difference(s) = {b}
// s.Difference(s) = {d, e, f}
func (s Int) Difference(t Int) Int {
	u := NewInt()
	for k := range s {
		if !t.Has(k) {
			u.Add(k)
		}
	}
	return u
}

// Intersection returns the intersection of Int s and t.
// For example:
// s = {a, b, c}
// s = {a, c, d, e, f}
// s.Intersection(s) = {a, c}
// s.Intersection(s) = {a, c}
// s.Intersection(s) = s.Intersection(s)
func (s Int) Intersection(t Int) Int {
	var max, min Int
	if s.Size() > t.Size() {
		max = s
		min = t
	} else {
		max = t
		min = s
	}
	u := NewInt()
	if min.Size() > 0 {
		for k := range min {
			if max.Has(k) {
				u[k] = struct{}{}
			}
		}
	}
	return u
}

// SymmetricDifference returns a new Int with the elements that are either in this Int
// or in the given Int, but not in both.
// For example:
// s = {a, c}
// s = {a, b, d}
// s.SymmetricDifference(s) = {c, b, d}
// s.SymmetricDifference(s) = {c, b, d}
// s.SymmetricDifference(s) = s.SymmetricDifference(s)
func (s Int) SymmetricDifference(t Int) Int {
	return s.Difference(t).Union(t.Difference(s))
}

// IsSubset predicates that tests whether the Int s is a subset of Int t.
// For example:
// s is subset of s
// s = {a, b, c}
// s = {a, b, c, d}
// s is not subset of s
// s = {a, f}
// s = {a, b, c, d}
func (s Int) IsSubset(t Int) bool {
	for k := range s {
		if !t.Has(k) {
			return false
		}
	}
	return true
}

// IsSuperset predicates that tests whether the Int s is a super of Int t.
// For example:
// s is super of s
// s = {a, b, c, d}
// s = {a, b, c}
// s is not super of s
// s = {a, f}
// s = {a, b, c, d}
func (s Int) IsSuperset(t Int) bool {
	return t.IsSubset(s)
}

// Equal predicates that tests whether the Int s equals of Int t.
// For example:
// s equals of s
// s = {a, b, c}
// s = {a, b, c}
// s does not equal of s
// s = {a, f}
// s = {a, b, c, d}
func (s Int) Equal(t Int) bool {
	return len(s) == len(t) && s.IsSubset(t)
}

// Copy returns new Int that clones from Int.
func (s Int) Copy() Int {
	t := NewIntWithSize(len(s))
	for k := range s {
		t[k] = struct{}{}
	}
	return t
}

// String returns a string representation of Int
func (s Int) String() string {
	v := make([]string, 0, s.Size())
	for element := range s {
		v = append(v, fmt.Sprintf("%v", element))
	}
	return fmt.Sprintf("[%s]", strings.Join(v, ", "))
}
