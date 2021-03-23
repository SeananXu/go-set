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

// Uint8 is a uint8 collection that contains no duplicate elements, without any particular order.
// It supports typical set operations: Core set-theoretical operations, Static sets, Dynamic
// sets, Additional operations.
type Uint8 map[uint8]struct{}

// NewUint8 initializes a new Uint8.
func NewUint8(elements ...uint8) Uint8 {
	s := Uint8{}
	s.Add(elements...)
	return s
}

// NewUint8WithSize initializes a new Uint8 with the specified size.
func NewUint8WithSize(size int) Uint8 {
	return make(map[uint8]struct{}, size)
}

// Add adds the elements to Uint8, if it is not present already.
func (s Uint8) Add(elements ...uint8) {
	for _, element := range elements {
		s[element] = struct{}{}
	}
}

// Remove removes the element from Uint8, if it is present.
func (s Uint8) Remove(elements ...uint8) {
	for _, element := range elements {
		delete(s, element)
	}
}

// Pop returns an arbitrary element of Uint8, deleting it from Uint8.
// The second value is a bool that is true if the elements existed in
// the Uint8, and false if not.
func (s Uint8) Pop() (uint8, bool) {
	for k := range s {
		delete(s, k)
		return k, true
	}
	return 0, false
}

// Size returns the number of elements in Uint8.
func (s Uint8) Size() int {
	return len(s)
}

// IsEmpty returns whether the Uint8 is Empty.
func (s Uint8) IsEmpty() bool {
	return len(s) == 0
}

// Clear removes all items from the Uint8.
func (s *Uint8) Clear() {
	*s = make(map[uint8]struct{})
}

// Has judges the specified element whether exists in the Uint8.
// it returns true if existed, and false if not.
func (s Uint8) Has(element uint8) bool {
	_, ok := s[element]
	return ok
}

// HasAll looks for the specified elements to judge
// whether all exist in the Uint8.
// it returns true if existed, and false if not.
func (s Uint8) HasAll(elements ...uint8) bool {
	for _, element := range elements {
		if _, ok := s[element]; !ok {
			return false
		}
	}
	return true
}

// HasAny looks for the specified elements to judge
// whether at least one of the element exists in the Uint8.
// it returns true if existed, and false if not.
func (s Uint8) HasAny(elements ...uint8) bool {
	for _, element := range elements {
		if _, ok := s[element]; ok {
			return true
		}
	}
	return false
}

// List returns the all elements as a slice.
func (s Uint8) List() []uint8 {
	var dest []uint8
	for k := range s {
		dest = append(dest, k)
	}
	return dest
}

// SortedList returns the all elements as a slice sorted by less func.
func (s Uint8) SortedList(less func(i, j uint8) bool) []uint8 {
	dest := s.List()
	sort.Slice(dest, func(i, j int) bool {
		return less(dest[i], dest[j])
	})
	return dest
}

// EachE traverses the elements in the Uint8, calling do func for each
// Uint8 member. the cycle will be stopped when the do func returns error.
// if err is ErrBreakEach, break the cycle and return nil,
// else, break the cycle and return error.
func (s Uint8) EachE(do func(i uint8) error) error {
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

// Each traverses the elements in the Uint8, calling do func for each
// Uint8 member.
func (s Uint8) Each(do func(i uint8)) {
	for k := range s {
		do(k)
	}
}

// Union returns the union of Uint8 s and t.
// For example:
// s = {a, b, c}
// s = {a, c, d, e, f}
// s.Union(s) = {a, b, c, d, e, f}
// s.Union(s) = {a, b, c, d, e, f}
// s.Union(s) = s.Union(s)
func (s Uint8) Union(t Uint8) Uint8 {
	// in order to reduce the number of growing map, copy the largest map here
	var max, min Uint8
	if s.Size() > t.Size() {
		max = s
		min = t
	} else {
		max = t
		min = s
	}
	if max.Size() == 0 {
		return NewUint8()
	}
	u := max.Copy()
	for k := range min {
		u[k] = struct{}{}
	}
	return u
}

// Difference returns the difference of Uint8 s and t.
// For example:
// s = {a, b, c}
// s = {a, c, d, e, f}
// s.Difference(s) = {b}
// s.Difference(s) = {d, e, f}
func (s Uint8) Difference(t Uint8) Uint8 {
	u := NewUint8()
	for k := range s {
		if !t.Has(k) {
			u.Add(k)
		}
	}
	return u
}

// Intersection returns the intersection of Uint8 s and t.
// For example:
// s = {a, b, c}
// s = {a, c, d, e, f}
// s.Intersection(s) = {a, c}
// s.Intersection(s) = {a, c}
// s.Intersection(s) = s.Intersection(s)
func (s Uint8) Intersection(t Uint8) Uint8 {
	var max, min Uint8
	if s.Size() > t.Size() {
		max = s
		min = t
	} else {
		max = t
		min = s
	}
	u := NewUint8()
	if min.Size() > 0 {
		for k := range min {
			if max.Has(k) {
				u[k] = struct{}{}
			}
		}
	}
	return u
}

// SymmetricDifference returns a new Uint8 with the elements that are either in this Uint8
// or in the given Uint8, but not in both.
// For example:
// s = {a, c}
// s = {a, b, d}
// s.SymmetricDifference(s) = {c, b, d}
// s.SymmetricDifference(s) = {c, b, d}
// s.SymmetricDifference(s) = s.SymmetricDifference(s)
func (s Uint8) SymmetricDifference(t Uint8) Uint8 {
	return s.Difference(t).Union(t.Difference(s))
}

// IsSubset predicates that tests whether the Uint8 s is a subset of Uint8 t.
// For example:
// s is subset of s
// s = {a, b, c}
// s = {a, b, c, d}
// s is not subset of s
// s = {a, f}
// s = {a, b, c, d}
func (s Uint8) IsSubset(t Uint8) bool {
	for k := range s {
		if !t.Has(k) {
			return false
		}
	}
	return true
}

// IsSuperset predicates that tests whether the Uint8 s is a super of Uint8 t.
// For example:
// s is super of s
// s = {a, b, c, d}
// s = {a, b, c}
// s is not super of s
// s = {a, f}
// s = {a, b, c, d}
func (s Uint8) IsSuperset(t Uint8) bool {
	return t.IsSubset(s)
}

// Equal predicates that tests whether the Uint8 s equals of Uint8 t.
// For example:
// s equals of s
// s = {a, b, c}
// s = {a, b, c}
// s does not equal of s
// s = {a, f}
// s = {a, b, c, d}
func (s Uint8) Equal(t Uint8) bool {
	return len(s) == len(t) && s.IsSubset(t)
}

// Copy returns new Uint8 that clones from Uint8.
func (s Uint8) Copy() Uint8 {
	t := NewUint8WithSize(len(s))
	for k := range s {
		t[k] = struct{}{}
	}
	return t
}

// String returns a string representation of Uint8
func (s Uint8) String() string {
	v := make([]string, 0, s.Size())
	for element := range s {
		v = append(v, fmt.Sprintf("%v", element))
	}
	return fmt.Sprintf("[%s]", strings.Join(v, ", "))
}
