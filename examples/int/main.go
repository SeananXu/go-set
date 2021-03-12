package main

import (
	"fmt"
	"log"

	"github.com/SeananXu/go-set"
)

func main() {
	log.Println("init int set")
	s := set.NewInt()
	log.Printf("set: %v\n", s)

	log.Println("###### basic operations ######")
	log.Println("set add test, test1, test2, pop1, case1, case2 and case3")
	s.Add(0)
	s.Add(1, 2)
	s.Add(11)
	s.Add(21, 22, 23)
	log.Printf("set: %v\n", s)

	log.Println("set removes test, test1 and test2")
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
	log.Println("set add copy1 and copy2")
	s.Add(41, 42)
	log.Printf("set copy: %v\n", s.Copy())
	log.Printf("set string: %s", s.String())
	log.Printf("set list: %v", s.List())
	s.SortedList(func(i, j int) bool {
		return i < j
	})

	log.Println("###### iterator operations ######")
	s.Each(func(i int) {
		log.Println(i)
	})
	s.EachE(func(i int) error {
		return set.ErrBreakEach
	})

	log.Println("###### check operations ######")
	log.Printf("set: %s\n", s)
	log.Printf("set is empty: %v\n", s.IsEmpty())
	log.Printf("set has e1: %v\n", s.Has(51))
	log.Printf("set has any e1 and e2: %v\n", s.HasAny(51))
	log.Printf("set has all e1 and e2: %v\n", s.HasAll(51))
	t := set.NewInt(61, 62)
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
