package main

import (
	"fmt"
	"go-set"
	"log"
)

func main() {
	log.Println("init interface set")
	s := set.NewInterface()
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
	s.SortedList(func(i, j interface{}) bool {
		return false
	})

	log.Println("###### iterator operations ######")
	s.Each(func(i interface{}) {
		log.Println(i)
	})
	s.EachE(func(i interface{}) error {
		return set.ErrBreakEach
	})

	log.Println("###### check operations ######")
	log.Printf("set: %s\n", s)
	log.Printf("set is empty: %v\n", s.IsEmpty())
	log.Printf("set has e1: %v\n", s.Has("e1"))
	log.Printf("set has any e1 and e2: %v\n", s.HasAny("e1"))
	log.Printf("set has all e1 and e2: %v\n", s.HasAll("e1"))
	t := set.NewInterface("t1", "t2")
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
