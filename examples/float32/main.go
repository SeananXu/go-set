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
	log.Println("init float32 set")
	s := set.NewFloat32()
	log.Printf("set: %v\n", s)

	log.Println("###### basic operations ######")
	log.Println("set add -0.1, -0.2, 5.9, 3.1, 7.8, 4.2 and 5.9")
	s.Add(-0.1)
	s.Add(-0.2, 5.9)
	s.Add(3.1)
	s.Add(7.8, 4.2, 5.9)
	log.Printf("set: %v\n", s)

	log.Println("set removes -0.1, -0.2 and 5.9")
	s.Remove(-0.1)
	s.Remove(-0.2, 5.9)
	log.Printf("set: %v\n", s)
	k, ok := s.Pop()
	if ok {
		log.Printf("set pop: %f\n", k)
		log.Printf("set: %v\n", s)
	}
	log.Println("set clear")
	s.Clear()
	log.Printf("set size: %d\n", s.Size())
	log.Println("set add 11.5 and 17.9")
	s.Add(11.5, 17.9)
	log.Printf("set copy: %v\n", s.Copy())
	log.Printf("set string: %s", s.String())
	log.Printf("set list: %v", s.List())
	s.SortedList(func(i, j float32) bool {
		return false
	})

	log.Println("###### iterator operations ######")
	s.Each(func(i float32) {
		log.Println(i)
	})
	s.EachE(func(i float32) error {
		return set.ErrBreakEach
	})

	log.Println("###### check operations ######")
	log.Printf("set: %s\n", s)
	log.Printf("set is empty: %v\n", s.IsEmpty())
	log.Printf("set has 1: %v\n", s.Has(1))
	log.Printf("set has any 1 and 1.2: %v\n", s.HasAny(1, 1.2))
	log.Printf("set has all 1 and 1.2: %v\n", s.HasAll(1, 1.2))
	t := set.NewFloat32(3.0, 8.3)
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
