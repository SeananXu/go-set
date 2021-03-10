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

// Float32 is a float32 collection that contains no duplicate elements, without any particular order.
// It supports typical set operations: Core set-theoretical operations, Static sets, Dynamic
// sets, Additional operations.
type Float32 map[float32]struct{}

// NewFloat32 initializes a new Float32.
func NewFloat32(elements ...float32) Float32 {
	s := Float32{}
	s.Add(elements...)
	return s
}

// NewFloat32WithSize initializes a new Float32 with the specified size.
func NewFloat32WithSize(size int) Float32 {
	return make(map[float32]struct{}, size)
}

// Add adds the elements to Float32, if it is not present already.
func (s Float32) Add(elements ...float32) {
	for _, element := range elements {
		s[element] = struct{}{}
	}
}

// Remove removes the element from Float32, if it is present.
func (s Float32) Remove(elements ...float32) {
	for _, element := range elements {
		delete(s, element)
	}
}

// Pop returns an arbitrary element of Float32, deleting it from Float32.
// The second value is a bool that is true if the elements existed in
// the Float32, and false if not.
func (s Float32) Pop() (float32, bool) {
	for k := range s {
		delete(s, k)
		return k, true
	}
	return 0, false
}

// Size returns the number of elements in Float32.
func (s Float32) Size() int {
	return len(s)
}

// IsEmpty returns whether the Float32 is Empty.
func (s Float32) IsEmpty() bool {
	return len(s) == 0
}

// Clear removes all items from the Float32.
func (s *Float32) Clear() {
	*s = make(map[float32]struct{})
}

// Has judges the specified element whether exists in the Float32.
// it returns true if existed, and false if not.
func (s Float32) Has(element float32) bool {
	_, ok := s[element]
	return ok
}

// HasAll looks for the specified elements to judge
// whether all exist in the Float32.
// it returns true if existed, and false if not.
func (s Float32) HasAll(elements ...float32) bool {
	for _, element := range elements {
		if _, ok := s[element]; !ok {
			return false
		}
	}
	return true
}

// HasAny looks for the specified elements to judge
// whether at least one of the element exists in the Float32.
// it returns true if existed, and false if not.
func (s Float32) HasAny(elements ...float32) bool {
	for _, element := range elements {
		if _, ok := s[element]; ok {
			return true
		}
	}
	return false
}

// List returns the all elements as a slice.
func (s Float32) List() []float32 {
	var dest []float32
	for k := range s {
		dest = append(dest, k)
	}
	return dest
}

// SortedList returns the all elements as a slice sorted by less func.
func (s Float32) SortedList(less func(i, j float32) bool) []float32 {
	dest := s.List()
	sort.Slice(dest, func(i, j int) bool {
		return less(dest[i], dest[j])
	})
	return dest
}

// EachE traverses the elements in the Float32, calling do func for each
// Float32 member. the cycle will be stopped when the do func returns error.
// if err is ErrBreakEach, break the cycle and return nil,
// else, break the cycle and return error.
func (s Float32) EachE(do func(i float32) error) error {
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

// Each traverses the elements in the Float32, calling do func for each
// Float32 member.
func (s Float32) Each(do func(i float32)) {
	for k := range s {
		do(k)
	}
}

// Union returns the union of Float32 s and t.
// For example:
// s = {a, b, c}
// s = {a, c, d, e, f}
// s.Union(s) = {a, b, c, d, e, f}
// s.Union(s) = {a, b, c, d, e, f}
// s.Union(s) = s.Union(s)
func (s Float32) Union(t Float32) Float32 {
	// in order to reduce the number of growing map, copy the largest map here
	var max, min Float32
	if s.Size() > t.Size() {
		max = s
		min = t
	} else {
		max = t
		min = s
	}
	if max.Size() == 0 {
		return NewFloat32()
	}
	u := max.Copy()
	for k := range min {
		u[k] = struct{}{}
	}
	return u
}

// Difference returns the difference of Float32 s and t.
// For example:
// s = {a, b, c}
// s = {a, c, d, e, f}
// s.Difference(s) = {b}
// s.Difference(s) = {d, e, f}
func (s Float32) Difference(t Float32) Float32 {
	u := NewFloat32()
	for k := range s {
		if !t.Has(k) {
			u.Add(k)
		}
	}
	return u
}

// Intersection returns the intersection of Float32 s and t.
// For example:
// s = {a, b, c}
// s = {a, c, d, e, f}
// s.Intersection(s) = {a, c}
// s.Intersection(s) = {a, c}
// s.Intersection(s) = s.Intersection(s)
func (s Float32) Intersection(t Float32) Float32 {
	var max, min Float32
	if s.Size() > t.Size() {
		max = s
		min = t
	} else {
		max = t
		min = s
	}
	u := NewFloat32()
	if min.Size() > 0 {
		for k := range min {
			if max.Has(k) {
				u[k] = struct{}{}
			}
		}
	}
	return u
}

// SymmetricDifference returns a new Float32 with the elements that are either in this Float32
// or in the given Float32, but not in both.
// For example:
// s = {a, c}
// s = {a, b, d}
// s.SymmetricDifference(s) = {c, b, d}
// s.SymmetricDifference(s) = {c, b, d}
// s.SymmetricDifference(s) = s.SymmetricDifference(s)
func (s Float32) SymmetricDifference(t Float32) Float32 {
	return s.Difference(t).Union(t.Difference(s))
}

// IsSubset predicates that tests whether the Float32 s is a subset of Float32 t.
// For example:
// s is subset of s
// s = {a, b, c}
// s = {a, b, c, d}
// s is not subset of s
// s = {a, f}
// s = {a, b, c, d}
func (s Float32) IsSubset(t Float32) bool {
	for k := range s {
		if !t.Has(k) {
			return false
		}
	}
	return true
}

// IsSuperset predicates that tests whether the Float32 s is a super of Float32 t.
// For example:
// s is super of s
// s = {a, b, c, d}
// s = {a, b, c}
// s is not super of s
// s = {a, f}
// s = {a, b, c, d}
func (s Float32) IsSuperset(t Float32) bool {
	return t.IsSubset(s)
}

// Equal predicates that tests whether the Float32 s equals of Float32 t.
// For example:
// s equals of s
// s = {a, b, c}
// s = {a, b, c}
// s does not equal of s
// s = {a, f}
// s = {a, b, c, d}
func (s Float32) Equal(t Float32) bool {
	return len(s) == len(t) && s.IsSubset(t)
}

// Copy returns new Float32 that clones from Float32.
func (s Float32) Copy() Float32 {
	t := NewFloat32WithSize(len(s))
	for k := range s {
		t[k] = struct{}{}
	}
	return t
}

// String returns a string representation of Float32
func (s Float32) String() string {
	v := make([]string, 0, s.Size())
	for element := range s {
		v = append(v, fmt.Sprintf("%v", element))
	}
	return fmt.Sprintf("[%s]", strings.Join(v, ", "))
}
