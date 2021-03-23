English | [简体中文](./README-zh_CN.md)

## Set [![GoDoc](https://pkg.go.dev/badge/github.com/SeananXu/go-set?utm_source=godoc)](https://godoc.org/github.com/SeananXu/go-set) [![Go Report Card](https://goreportcard.com/badge/github.com/SeananXu/go-set)](https://goreportcard.com/report/github.com/SeananXu/go-set)
Set is an abstract data type that can store unique values, without any particular order implementation in Go.

Except for reference document [Set (abstract data type)](https://en.wikipedia.org/wiki/Set_(abstract_data_type)) to define functions
, It also implements some useful functions. for example: Each, EachE, List, SortedList and so on.

## Install
Using `go get` command to get the latest version
```bash
go get github.com/SeananXu/go-set
```
Import it with:
```go
import "github.com/SeananXu/go-set"
```
and use set as the package name inside the code.

## Example
#### Initialization
```go
// interface set
interfaceSet := set.NewInterface()

// int set
intSet := set.NewInt()

// int8 set
int8Set := set.NewInt8()

// int16 set
int16Set := set.NewInt16()

// int32 set
int32Set := set.NewInt32()

// int64 set
int64Set := set.NewInt64()

// float32 set
float32Set := set.NewFloat32()

// float64 set
float64Set := set.NewFloat64()

// uint set
uintSet := set.NewUint()

// uint8 set
uint8Set := set.NewUint8()

// uint16 set
uint16Set := set.NewUint16()

// uint32 set
uint32Set := set.NewUint32()

// uint32 set
uint64Set := set.NewUint64()

// string set
stringSet := set.NewString()
```
#### Basic Operations
```go
// add element
s.Add("element")

// remove element
s.Remove("element")

// pop element
k, ok := s.Pop()

// clear set
s.Clear()

// the number of elements
size := s.Size()

// clone new set
t := s.Copy()

// returns the all elements as a slice
t := s.List()

// returns the all elements as a slice sorted by less func
s.SortedList(func(i, j interface{}) bool {
    return false
})

// returns string
str := s.String()
```
#### Iterator Operations
```go
s.Each(func(i interface{}) {
    log.Println(i)
})

s.EachE(func(i interface{}) error {
    log.Println(i)
    return set.ErrBreakEach
})
```
#### Check Operations
```go
// whether the set is Empty
b := s.IsEmpty()

// whether element exists in set
b := s.Has("element")

// whether all elements exist in set
b := s.HasAll("element1", "element2")

// whether at least one of the element exists in the set
b := s.HasAny("element1", "element2")

// predicates that tests whether the set s is a subset of set t
b := s.IsSuperset(t)

// predicates that tests whether the set s is a super of set t
b := s.IsSubset(t)

// predicates that tests whether the set s equals of set t
b := s.Equal(t)
```
#### Set Operations
```go
// returns the union of sets s and t
s.Union(t)

// returns the intersection of sets s and t
s.Intersection(t)

// returns the difference of sets s and t
s.Difference(t)

// returns the symmetric difference of sets s and t
s.SymmetricDifference(t)
```
more case click [here](./examples/README.md)

## Setgen
The `Setgen` command is used to generate source code for a set class given a type.
It supports the following flags.
- `-s`: Set name, default: element type add 's'.
- `-i`: Import element package, default: don't import package.
- `-p`: Generated go file package, default: directory name.
- `-t`: Set storage element type, this options must be set.
- `-o`: Output file name, default: set name add '.go'.
- `-l`: Whether go file imports 'ErrBreakEach' of 'github.com/SeananXu/go-set', default: import.
- `-h`: Help document.

Install
```
go get github.com/SeananXu/go-set/setgen
```
For example:
```
setgen -t Example
```
## License

The MIT License (MIT) - see [LICENSE](LICENSE) for more details