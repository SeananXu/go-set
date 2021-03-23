[English](./README.md) | 简体中文

## Set [![GoDoc](https://pkg.go.dev/badge/github.com/SeananXu/go-set?utm_source=godoc)](https://godoc.org/github.com/SeananXu/go-set) [![Go Report Card](https://goreportcard.com/badge/github.com/SeananXu/go-set)](https://goreportcard.com/report/github.com/SeananXu/go-set) 
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
更多点击[这里](./examples/README-zh_CN.md)

## Setgen
`Setgen` 根据指定的文件自动生成对应 `Set` 文件的命令行工具
- `-s`: Set name, default: element type add 's'.
- `-i`: Import element package, default: don't import package.
- `-p`: Generated go file package, default: directory name.
- `-t`: Set storage element type, this options must be set.
- `-o`: Output file name, default: set name add '.go'.
- `-l`: Whether go file imports 'ErrBreakEach' of 'github.com/SeananXu/go-set', default: import.
- `-h`: Help document.

安装
```
go get github.com/SeananXu/go-set/setgen
```
例如:
```
setgen -t Example
```

## License

The MIT License (MIT) - see [LICENSE](LICENSE) for more details