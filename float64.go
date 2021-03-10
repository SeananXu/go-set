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

// Float64 is a float64 collection that contains no duplicate elements, without any particular order.
// It supports typical set operations: Core set-theoretical operations, Static sets, Dynamic
// sets, Additional operations.
type Float64 map[float64]struct{}

// NewFloat64 initializes a new Float64.
func NewFloat64(elements ...float64) Float64 {
	s := Float64{}
	s.Add(elements...)
	return s
}

// NewFloat64WithSize initializes a new Float64 with the specified size.
func NewFloat64WithSize(size int) Float64 {
	return make(map[float64]struct{}, size)
}

// Add adds the elements to Float64, if it is not present already.
func (s Float64) Add(elements ...float64) {
	for _, element := range elements {
		s[element] = struct{}{}
	}
}

// Remove removes the element from Float64, if it is present.
func (s Float64) Remove(elements ...float64) {
	for _, element := range elements {
		delete(s, element)
	}
}

// Pop returns an arbitrary element of Float64, deleting it from Float64.
// The second value is a bool that is true if the elements existed in
// the Float64, and false if not.
func (s Float64) Pop() (float64, bool) {
	for k := range s {
		delete(s, k)
		return k, true
	}
	return 0, false
}

// Size returns the number of elements in Float64.
func (s Float64) Size() int {
	return len(s)
}

// IsEmpty returns whether the Float64 is Empty.
func (s Float64) IsEmpty() bool {
	return len(s) == 0
}

// Clear removes all items from the Float64.
func (s *Float64) Clear() {
	*s = make(map[float64]struct{})
}

// Has judges the specified element whether exists in the Float64.
// it returns true if existed, and false if not.
func (s Float64) Has(element float64) bool {
	_, ok := s[element]
	return ok
}

// HasAll looks for the specified elements to judge
// whether all exist in the Float64.
// it returns true if existed, and false if not.
func (s Float64) HasAll(elements ...float64) bool {
	for _, element := range elements {
		if _, ok := s[element]; !ok {
			return false
		}
	}
	return true
}

// HasAny looks for the specified elements to judge
// whether at least one of the element exists in the Float64.
// it returns true if existed, and false if not.
func (s Float64) HasAny(elements ...float64) bool {
	for _, element := range elements {
		if _, ok := s[element]; ok {
			return true
		}
	}
	return false
}

// List returns the all elements as a slice.
func (s Float64) List() []float64 {
	var dest []float64
	for k := range s {
		dest = append(dest, k)
	}
	return dest
}

// SortedList returns the all elements as a slice sorted by less func.
func (s Float64) SortedList(less func(i, j float64) bool) []float64 {
	dest := s.List()
	sort.Slice(dest, func(i, j int) bool {
		return less(dest[i], dest[j])
	})
	return dest
}

// EachE traverses the elements in the Float64, calling do func for each
// Float64 member. the cycle will be stopped when the do func returns error.
// if err is ErrBreakEach, break the cycle and return nil,
// else, break the cycle and return error.
func (s Float64) EachE(do func(i float64) error) error {
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

// Each traverses the elements in the Float64, calling do func for each
// Float64 member.
func (s Float64) Each(do func(i float64)) {
	for k := range s {
		do(k)
	}
}

// Union returns the union of Float64 s and t.
// For example:
// s = {a, b, c}
// s = {a, c, d, e, f}
// s.Union(s) = {a, b, c, d, e, f}
// s.Union(s) = {a, b, c, d, e, f}
// s.Union(s) = s.Union(s)
func (s Float64) Union(t Float64) Float64 {
	// in order to reduce the number of growing map, copy the largest map here
	var max, min Float64
	if s.Size() > t.Size() {
		max = s
		min = t
	} else {
		max = t
		min = s
	}
	if max.Size() == 0 {
		return NewFloat64()
	}
	u := max.Copy()
	for k := range min {
		u[k] = struct{}{}
	}
	return u
}

// Difference returns the difference of Float64 s and t.
// For example:
// s = {a, b, c}
// s = {a, c, d, e, f}
// s.Difference(s) = {b}
// s.Difference(s) = {d, e, f}
func (s Float64) Difference(t Float64) Float64 {
	u := NewFloat64()
	for k := range s {
		if !t.Has(k) {
			u.Add(k)
		}
	}
	return u
}

// Intersection returns the intersection of Float64 s and t.
// For example:
// s = {a, b, c}
// s = {a, c, d, e, f}
// s.Intersection(s) = {a, c}
// s.Intersection(s) = {a, c}
// s.Intersection(s) = s.Intersection(s)
func (s Float64) Intersection(t Float64) Float64 {
	var max, min Float64
	if s.Size() > t.Size() {
		max = s
		min = t
	} else {
		max = t
		min = s
	}
	u := NewFloat64()
	if min.Size() > 0 {
		for k := range min {
			if max.Has(k) {
				u[k] = struct{}{}
			}
		}
	}
	return u
}

// SymmetricDifference returns a new Float64 with the elements that are either in this Float64
// or in the given Float64, but not in both.
// For example:
// s = {a, c}
// s = {a, b, d}
// s.SymmetricDifference(s) = {c, b, d}
// s.SymmetricDifference(s) = {c, b, d}
// s.SymmetricDifference(s) = s.SymmetricDifference(s)
func (s Float64) SymmetricDifference(t Float64) Float64 {
	return s.Difference(t).Union(t.Difference(s))
}

// IsSubset predicates that tests whether the Float64 s is a subset of Float64 t.
// For example:
// s is subset of s
// s = {a, b, c}
// s = {a, b, c, d}
// s is not subset of s
// s = {a, f}
// s = {a, b, c, d}
func (s Float64) IsSubset(t Float64) bool {
	for k := range s {
		if !t.Has(k) {
			return false
		}
	}
	return true
}

// IsSuperset predicates that tests whether the Float64 s is a super of Float64 t.
// For example:
// s is super of s
// s = {a, b, c, d}
// s = {a, b, c}
// s is not super of s
// s = {a, f}
// s = {a, b, c, d}
func (s Float64) IsSuperset(t Float64) bool {
	return t.IsSubset(s)
}

// Equal predicates that tests whether the Float64 s equals of Float64 t.
// For example:
// s equals of s
// s = {a, b, c}
// s = {a, b, c}
// s does not equal of s
// s = {a, f}
// s = {a, b, c, d}
func (s Float64) Equal(t Float64) bool {
	return len(s) == len(t) && s.IsSubset(t)
}

// Copy returns new Float64 that clones from Float64.
func (s Float64) Copy() Float64 {
	t := NewFloat64WithSize(len(s))
	for k := range s {
		t[k] = struct{}{}
	}
	return t
}

// String returns a string representation of Float64
func (s Float64) String() string {
	v := make([]string, 0, s.Size())
	for element := range s {
		v = append(v, fmt.Sprintf("%v", element))
	}
	return fmt.Sprintf("[%s]", strings.Join(v, ", "))
}
