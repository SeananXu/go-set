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
	log.Println("init string set")
	s := set.NewString()
	log.Printf("set: %v\n", s)

	log.Println("###### basic operations ######")
	log.Println("set add test, test1, test2, pop1, case1, case2 and case3")
	s.Add("test")
	s.Add("test1", "test2")
	s.Add("pop1")
	s.Add("case1", "case2", "case3")
	log.Printf("set: %v\n", s)

	log.Println("set removes test, test1 and test2")
	s.Remove("test")
	s.Remove("test1", "test2")
	log.Printf("set: %v\n", s)
	k, ok := s.Pop()
	if ok {
		log.Printf("set pop: %s\n", k)
		log.Printf("set: %v\n", s)
	}
	log.Println("set clear")
	s.Clear()
	log.Printf("set size: %d\n", s.Size())
	log.Println("set add copy1 and copy2")
	s.Add("copy1", "copy2")
	log.Printf("set copy: %v\n", s.Copy())
	log.Printf("set string: %s", s.String())
	log.Printf("set list: %v", s.List())
	s.SortedList(func(i, j string) bool {
		return false
	})

	log.Println("###### iterator operations ######")
	s.Each(func(i string) {
		log.Println(i)
	})
	s.EachE(func(i string) error {
		return set.ErrBreakEach
	})

	log.Println("###### check operations ######")
	log.Printf("set: %s\n", s)
	log.Printf("set is empty: %v\n", s.IsEmpty())
	log.Printf("set has e1: %v\n", s.Has("e1"))
	log.Printf("set has any e1 and e2: %v\n", s.HasAny("e1", "e2"))
	log.Printf("set has all e1 and e2: %v\n", s.HasAll("e1", "e2"))
	t := set.NewString("t1", "t2")
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
