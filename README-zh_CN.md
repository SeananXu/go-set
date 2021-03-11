[English](./README.md) | 简体中文

## Set
Set 是用 go 语言实现存储**唯一**和**无序**元素的数据结构.

除了 [Set (abstract data type)](https://en.wikipedia.org/wiki/Set_(abstract_data_type)) 文档定义功能之外, 它还
实现很多好用的功能, 例如: Each, EachE, List, SortedList...

## 安装
使用 `go get` 指令获取最新代码
```bash
go get github.com/SeananXu/go-set
```
`improt` 依赖:
```go
import "github.com/SeananXu/go-set"
```
使用 `set` 作为包名使用代码

## 例子
#### Initialization
```go
// interface set
interfaceSet := set.NewInterface()

// int set
intSet := set.NewInt()

// int64 set
int64Set := set.NewInt64()

// float32 set
float32Set := set.NewFloat32()

// float64 set
float64Set := set.NewFloat64()

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
更多点击[这里](./README-zh_CN.md)

## License

The MIT License (MIT) - see [LICENSE](./LISENCE) for more details