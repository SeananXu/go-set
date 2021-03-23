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

package main

import (
	"fmt"
	"log"

	"github.com/SeananXu/go-set"
)

func main() {
	log.Println("init uint set")
	s := set.NewUint()
	log.Printf("set: %v\n", s)

	log.Println("###### basic operations ######")
	log.Println("set add 0, 1, 2, 11, 21, 22 and 23")
	s.Add(0)
	s.Add(1, 2)
	s.Add(11)
	s.Add(21, 22, 23)
	log.Printf("set: %v\n", s)

	log.Println("set removes 0, 1 and 2")
	s.Remove(0)
	s.Remove(1, 2)
	log.Printf("set: %v\n", s)
	k, ok := s.Pop()
	if ok {
		log.Printf("set pop: %d\n", k)
		log.Printf("set: %v\n", s)
	}
	log.Println("set clear")
	s.Clear()
	log.Printf("set size: %d\n", s.Size())
	log.Println("set add 41 and 42")
	s.Add(41, 42)
	log.Printf("set copy: %v\n", s.Copy())
	log.Printf("set string: %s", s.String())
	log.Printf("set list: %v", s.List())
	s.SortedList(func(i, j uint) bool {
		return i < j
	})

	log.Println("###### iterator operations ######")
	s.Each(func(i uint) {
		log.Println(i)
	})
	s.EachE(func(i uint) error {
		return set.ErrBreakEach
	})

	log.Println("###### check operations ######")
	log.Printf("set: %s\n", s)
	log.Printf("set is empty: %v\n", s.IsEmpty())
	log.Printf("set has 51: %v\n", s.Has(51))
	log.Printf("set has any 51 and 52: %v\n", s.HasAny(51, 52))
	log.Printf("set has all 51 and 52: %v\n", s.HasAll(51, 52))
	t := set.NewUint(61, 62)
	log.Printf("t set: %v\n", t)
	s.IsSuperset(t)
	s.IsSubset(t)
	s.Equal(t)

	fmt.Println("###### set operations ######")
	s.Union(t)
	s.Intersection(t)
	s.Difference(t)
	s.SymmetricDifference(t)
}
