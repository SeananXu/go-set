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

// Int64 is a int64 collection that contains no duplicate elements, without any particular order.
// It supports typical set operations: Core set-theoretical operations, Static sets, Dynamic
// sets, Additional operations.
type Int64 map[int64]struct{}

// NewInt64 initializes a new Int64.
func NewInt64(elements ...int64) Int64 {
	s := Int64{}
	s.Add(elements...)
	return s
}

// NewInt64WithSize initializes a new Int64 with the specified size.
func NewInt64WithSize(size int) Int64 {
	return make(map[int64]struct{}, size)
}

// Add adds the elements to Int64, if it is not present already.
func (s Int64) Add(elements ...int64) {
	for _, element := range elements {
		s[element] = struct{}{}
	}
}

// Remove removes the element from Int64, if it is present.
func (s Int64) Remove(elements ...int64) {
	for _, element := range elements {
		delete(s, element)
	}
}

// Pop returns an arbitrary element of Int64, deleting it from Int64.
// The second value is a bool that is true if the elements existed in
// the Int64, and false if not.
func (s Int64) Pop() (int64, bool) {
	for k := range s {
		delete(s, k)
		return k, true
	}
	return 0, false
}

// Size returns the number of elements in Int64.
func (s Int64) Size() int {
	return len(s)
}

// IsEmpty returns whether the Int64 is Empty.
func (s Int64) IsEmpty() bool {
	return len(s) == 0
}

// Clear removes all items from the Int64.
func (s *Int64) Clear() {
	*s = make(map[int64]struct{})
}

// Has judges the specified element whether exists in the Int64.
// it returns true if existed, and false if not.
func (s Int64) Has(element int64) bool {
	_, ok := s[element]
	return ok
}

// HasAll looks for the specified elements to judge
// whether all exist in the Int64.
// it returns true if existed, and false if not.
func (s Int64) HasAll(elements ...int64) bool {
	for _, element := range elements {
		if _, ok := s[element]; !ok {
			return false
		}
	}
	return true
}

// HasAny looks for the specified elements to judge
// whether at least one of the element exists in the Int64.
// it returns true if existed, and false if not.
func (s Int64) HasAny(elements ...int64) bool {
	for _, element := range elements {
		if _, ok := s[element]; ok {
			return true
		}
	}
	return false
}

// List returns the all elements as a slice.
func (s Int64) List() []int64 {
	var dest []int64
	for k := range s {
		dest = append(dest, k)
	}
	return dest
}

// SortedList returns the all elements as a slice sorted by less func.
func (s Int64) SortedList(less func(i, j int64) bool) []int64 {
	dest := s.List()
	sort.Slice(dest, func(i, j int) bool {
		return less(dest[i], dest[j])
	})
	return dest
}

// EachE traverses the elements in the Int64, calling do func for each
// Int64 member. the cycle will be stopped when the do func returns error.
// if err is ErrBreakEach, break the cycle and return nil,
// else, break the cycle and return error.
func (s Int64) EachE(do func(i int64) error) error {
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

// Each traverses the elements in the Int64, calling do func for each
// Int64 member.
func (s Int64) Each(do func(i int64)) {
	for k := range s {
		do(k)
	}
}

// Union returns the union of Int64 s and t.
// For example:
// s = {a, b, c}
// s = {a, c, d, e, f}
// s.Union(s) = {a, b, c, d, e, f}
// s.Union(s) = {a, b, c, d, e, f}
// s.Union(s) = s.Union(s)
func (s Int64) Union(t Int64) Int64 {
	// in order to reduce the number of growing map, copy the largest map here
	var max, min Int64
	if s.Size() > t.Size() {
		max = s
		min = t
	} else {
		max = t
		min = s
	}
	if max.Size() == 0 {
		return NewInt64()
	}
	u := max.Copy()
	for k := range min {
		u[k] = struct{}{}
	}
	return u
}

// Difference returns the difference of Int64 s and t.
// For example:
// s = {a, b, c}
// s = {a, c, d, e, f}
// s.Difference(s) = {b}
// s.Difference(s) = {d, e, f}
func (s Int64) Difference(t Int64) Int64 {
	u := NewInt64()
	for k := range s {
		if !t.Has(k) {
			u.Add(k)
		}
	}
	return u
}

// Intersection returns the intersection of Int64 s and t.
// For example:
// s = {a, b, c}
// s = {a, c, d, e, f}
// s.Intersection(s) = {a, c}
// s.Intersection(s) = {a, c}
// s.Intersection(s) = s.Intersection(s)
func (s Int64) Intersection(t Int64) Int64 {
	var max, min Int64
	if s.Size() > t.Size() {
		max = s
		min = t
	} else {
		max = t
		min = s
	}
	u := NewInt64()
	if min.Size() > 0 {
		for k := range min {
			if max.Has(k) {
				u[k] = struct{}{}
			}
		}
	}
	return u
}

// SymmetricDifference returns a new Int64 with the elements that are either in this Int64
// or in the given Int64, but not in both.
// For example:
// s = {a, c}
// s = {a, b, d}
// s.SymmetricDifference(s) = {c, b, d}
// s.SymmetricDifference(s) = {c, b, d}
// s.SymmetricDifference(s) = s.SymmetricDifference(s)
func (s Int64) SymmetricDifference(t Int64) Int64 {
	return s.Difference(t).Union(t.Difference(s))
}

// IsSubset predicates that tests whether the Int64 s is a subset of Int64 t.
// For example:
// s is subset of s
// s = {a, b, c}
// s = {a, b, c, d}
// s is not subset of s
// s = {a, f}
// s = {a, b, c, d}
func (s Int64) IsSubset(t Int64) bool {
	for k := range s {
		if !t.Has(k) {
			return false
		}
	}
	return true
}

// IsSuperset predicates that tests whether the Int64 s is a super of Int64 t.
// For example:
// s is super of s
// s = {a, b, c, d}
// s = {a, b, c}
// s is not super of s
// s = {a, f}
// s = {a, b, c, d}
func (s Int64) IsSuperset(t Int64) bool {
	return t.IsSubset(s)
}

// Equal predicates that tests whether the Int64 s equals of Int64 t.
// For example:
// s equals of s
// s = {a, b, c}
// s = {a, b, c}
// s does not equal of s
// s = {a, f}
// s = {a, b, c, d}
func (s Int64) Equal(t Int64) bool {
	return len(s) == len(t) && s.IsSubset(t)
}

// Copy returns new Int64 that clones from Int64.
func (s Int64) Copy() Int64 {
	t := NewInt64WithSize(len(s))
	for k := range s {
		t[k] = struct{}{}
	}
	return t
}

// String returns a string representation of Int64
func (s Int64) String() string {
	v := make([]string, 0, s.Size())
	for element := range s {
		v = append(v, fmt.Sprintf("%v", element))
	}
	return fmt.Sprintf("[%s]", strings.Join(v, ", "))
}
