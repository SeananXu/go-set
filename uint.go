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

// Uint is a uint collection that contains no duplicate elements, without any particular order.
// It supports typical set operations: Core set-theoretical operations, Static sets, Dynamic
// sets, Additional operations.
type Uint map[uint]struct{}

// NewUint initializes a new Uint.
func NewUint(elements ...uint) Uint {
	s := Uint{}
	s.Add(elements...)
	return s
}

// NewUintWithSize initializes a new Uint with the specified size.
func NewUintWithSize(size int) Uint {
	return make(map[uint]struct{}, size)
}

// Add adds the elements to Uint, if it is not present already.
func (s Uint) Add(elements ...uint) {
	for _, element := range elements {
		s[element] = struct{}{}
	}
}

// Remove removes the element from Uint, if it is present.
func (s Uint) Remove(elements ...uint) {
	for _, element := range elements {
		delete(s, element)
	}
}

// Pop returns an arbitrary element of Uint, deleting it from Uint.
// The second value is a bool that is true if the elements existed in
// the Uint, and false if not.
func (s Uint) Pop() (uint, bool) {
	for k := range s {
		delete(s, k)
		return k, true
	}
	return 0, false
}

// Size returns the number of elements in Uint.
func (s Uint) Size() int {
	return len(s)
}

// IsEmpty returns whether the Uint is Empty.
func (s Uint) IsEmpty() bool {
	return len(s) == 0
}

// Clear removes all items from the Uint.
func (s *Uint) Clear() {
	*s = make(map[uint]struct{})
}

// Has judges the specified element whether exists in the Uint.
// it returns true if existed, and false if not.
func (s Uint) Has(element uint) bool {
	_, ok := s[element]
	return ok
}

// HasAll looks for the specified elements to judge
// whether all exist in the Uint.
// it returns true if existed, and false if not.
func (s Uint) HasAll(elements ...uint) bool {
	for _, element := range elements {
		if _, ok := s[element]; !ok {
			return false
		}
	}
	return true
}

// HasAny looks for the specified elements to judge
// whether at least one of the element exists in the Uint.
// it returns true if existed, and false if not.
func (s Uint) HasAny(elements ...uint) bool {
	for _, element := range elements {
		if _, ok := s[element]; ok {
			return true
		}
	}
	return false
}

// List returns the all elements as a slice.
func (s Uint) List() []uint {
	var dest []uint
	for k := range s {
		dest = append(dest, k)
	}
	return dest
}

// SortedList returns the all elements as a slice sorted by less func.
func (s Uint) SortedList(less func(i, j uint) bool) []uint {
	dest := s.List()
	sort.Slice(dest, func(i, j int) bool {
		return less(dest[i], dest[j])
	})
	return dest
}

// EachE traverses the elements in the Uint, calling do func for each
// Uint member. the cycle will be stopped when the do func returns error.
// if err is ErrBreakEach, break the cycle and return nil,
// else, break the cycle and return error.
func (s Uint) EachE(do func(i uint) error) error {
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

// Each traverses the elements in the Uint, calling do func for each
// Uint member.
func (s Uint) Each(do func(i uint)) {
	for k := range s {
		do(k)
	}
}

// Union returns the union of Uint s and t.
// For example:
// s = {a, b, c}
// s = {a, c, d, e, f}
// s.Union(s) = {a, b, c, d, e, f}
// s.Union(s) = {a, b, c, d, e, f}
// s.Union(s) = s.Union(s)
func (s Uint) Union(t Uint) Uint {
	// in order to reduce the number of growing map, copy the largest map here
	var max, min Uint
	if s.Size() > t.Size() {
		max = s
		min = t
	} else {
		max = t
		min = s
	}
	if max.Size() == 0 {
		return NewUint()
	}
	u := max.Copy()
	for k := range min {
		u[k] = struct{}{}
	}
	return u
}

// Difference returns the difference of Uint s and t.
// For example:
// s = {a, b, c}
// s = {a, c, d, e, f}
// s.Difference(s) = {b}
// s.Difference(s) = {d, e, f}
func (s Uint) Difference(t Uint) Uint {
	u := NewUint()
	for k := range s {
		if !t.Has(k) {
			u.Add(k)
		}
	}
	return u
}

// Intersection returns the intersection of Uint s and t.
// For example:
// s = {a, b, c}
// s = {a, c, d, e, f}
// s.Intersection(s) = {a, c}
// s.Intersection(s) = {a, c}
// s.Intersection(s) = s.Intersection(s)
func (s Uint) Intersection(t Uint) Uint {
	var max, min Uint
	if s.Size() > t.Size() {
		max = s
		min = t
	} else {
		max = t
		min = s
	}
	u := NewUint()
	if min.Size() > 0 {
		for k := range min {
			if max.Has(k) {
				u[k] = struct{}{}
			}
		}
	}
	return u
}

// SymmetricDifference returns a new Uint with the elements that are either in this Uint
// or in the given Uint, but not in both.
// For example:
// s = {a, c}
// s = {a, b, d}
// s.SymmetricDifference(s) = {c, b, d}
// s.SymmetricDifference(s) = {c, b, d}
// s.SymmetricDifference(s) = s.SymmetricDifference(s)
func (s Uint) SymmetricDifference(t Uint) Uint {
	return s.Difference(t).Union(t.Difference(s))
}

// IsSubset predicates that tests whether the Uint s is a subset of Uint t.
// For example:
// s is subset of s
// s = {a, b, c}
// s = {a, b, c, d}
// s is not subset of s
// s = {a, f}
// s = {a, b, c, d}
func (s Uint) IsSubset(t Uint) bool {
	for k := range s {
		if !t.Has(k) {
			return false
		}
	}
	return true
}

// IsSuperset predicates that tests whether the Uint s is a super of Uint t.
// For example:
// s is super of s
// s = {a, b, c, d}
// s = {a, b, c}
// s is not super of s
// s = {a, f}
// s = {a, b, c, d}
func (s Uint) IsSuperset(t Uint) bool {
	return t.IsSubset(s)
}

// Equal predicates that tests whether the Uint s equals of Uint t.
// For example:
// s equals of s
// s = {a, b, c}
// s = {a, b, c}
// s does not equal of s
// s = {a, f}
// s = {a, b, c, d}
func (s Uint) Equal(t Uint) bool {
	return len(s) == len(t) && s.IsSubset(t)
}

// Copy returns new Uint that clones from Uint.
func (s Uint) Copy() Uint {
	t := NewUintWithSize(len(s))
	for k := range s {
		t[k] = struct{}{}
	}
	return t
}

// String returns a string representation of Uint
func (s Uint) String() string {
	v := make([]string, 0, s.Size())
	for element := range s {
		v = append(v, fmt.Sprintf("%v", element))
	}
	return fmt.Sprintf("[%s]", strings.Join(v, ", "))
}
