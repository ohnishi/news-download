package core

type stringIncluder interface {
	Include(str string) bool
}

// StringSet はstringの集合を表す。
type StringSet map[string]struct{}

// NewStringSet はStringSetを新たに生成する。
func NewStringSet(strs ...string) *StringSet {
	s := StringSet{}
	s.Add(strs...)
	return &s
}

// Add は要素を追加する。
func (s *StringSet) Add(strs ...string) {
	for _, str := range strs {
		(*s)[str] = struct{}{}
	}
}

// Include はstrが含まれてばtrueを返す。
func (s *StringSet) Include(str string) bool {
	_, ok := (*s)[str]
	return ok
}

// Size は含まれる要素数を返す。
func (s *StringSet) Size() int {
	return len(*s)
}

// Merge はs2の要素を追加する。
func (s *StringSet) Merge(s2 *StringSet) {
	for k := range *s2 {
		(*s)[k] = struct{}{}
	}
}

// Slice はstringの集合をスライスとして返す。
func (s *StringSet) Slice() []string {
	strs := []string{}
	for k := range *s {
		strs = append(strs, k)
	}
	return strs
}

// Subtract はsに含まれていてs2に含まれていない要素の集合を返す。
func (s *StringSet) Subtract(s2 stringIncluder) StringSet {
	res := StringSet{}
	for k := range *s {
		if !s2.Include(k) {
			res.Add(k)
		}
	}
	return res
}

// StringOrderedSet は要素を追加した順序を保持するstringの集合を表す
type StringOrderedSet struct {
	orderedStrings []string
	strings        StringSet
}

// NewStringOrderedSet はStringOrderSetを新たに生成する。
func NewStringOrderedSet(strs ...string) *StringOrderedSet {
	s := StringOrderedSet{
		orderedStrings: nil,
		strings:        StringSet{},
	}
	s.Add(strs...)
	return &s
}

// Add は要素を追加する（既に含まれる場合は何もしない）
func (s *StringOrderedSet) Add(strs ...string) {
	for _, str := range strs {
		if !s.strings.Include(str) {
			s.orderedStrings = append(s.orderedStrings, str)
			s.strings.Add(str)
		}
	}
}

// Include はstrが含まれていればtrueを返す。
func (s *StringOrderedSet) Include(str string) bool {
	return s.strings.Include(str)
}

// Size は含まれる要素数を返す。
func (s *StringOrderedSet) Size() int {
	return s.strings.Size()
}

// Merge はs2の要素を追加する。
func (s *StringOrderedSet) Merge(s2 *StringOrderedSet) {
	s.Add(s2.orderedStrings...)
}

// Slice はstringの集合をスライスとして返す。
func (s *StringOrderedSet) Slice() []string {
	res := make([]string, len(s.orderedStrings))
	copy(res, s.orderedStrings)
	return res
}

// Subtract はsに含まれていてs2に含まれていない要素の集合を返す。
func (s *StringOrderedSet) Subtract(s2 stringIncluder) StringOrderedSet {
	res := NewStringOrderedSet()
	for _, str := range s.orderedStrings {
		if !s2.Include(str) {
			res.Add(str)
		}
	}
	return *res
}

type int64Includer interface {
	Include(i int64) bool
}

// Int64Set はint64の集合を表す。
type Int64Set map[int64]struct{}

// NewInt64Set はInt64Setを新たに生成する。
func NewInt64Set(ints ...int64) *Int64Set {
	s := Int64Set{}
	s.Add(ints...)
	return &s
}

// Add は要素を追加する。
func (s *Int64Set) Add(ints ...int64) {
	for _, i := range ints {
		(*s)[i] = struct{}{}
	}
}

// Include はiが含まれてばtrueを返す。
func (s *Int64Set) Include(i int64) bool {
	_, ok := (*s)[i]
	return ok
}

// Size はInt64Setに含まれる要素数を返す。
func (s *Int64Set) Size() int {
	return len(*s)
}

// Merge はs2の要素を追加する。
func (s *Int64Set) Merge(s2 *Int64Set) {
	for k := range *s2 {
		(*s)[k] = struct{}{}
	}
}

// Slice はint64の集合をスライスとして返す。
func (s *Int64Set) Slice() []int64 {
	ints := []int64{}
	for k := range *s {
		ints = append(ints, k)
	}
	return ints
}

// Subtract はsに含まれていてs2に含まれていない要素の集合を返す。
func (s *Int64Set) Subtract(s2 int64Includer) Int64Set {
	res := Int64Set{}
	for k := range *s {
		if !s2.Include(k) {
			res.Add(k)
		}
	}
	return res
}

// Int64OrderedSet は要素を追加した順序を保持するint64の集合を表す
type Int64OrderedSet struct {
	orderedInts []int64
	ints        Int64Set
}

// NewInt64OrderedSet はInt64OrderSetを新たに生成する
func NewInt64OrderedSet(ints ...int64) *Int64OrderedSet {
	s := Int64OrderedSet{
		orderedInts: nil,
		ints:        Int64Set{},
	}
	s.Add(ints...)
	return &s
}

// Add は要素を追加する（既に含まれる場合は何もしない）
func (s *Int64OrderedSet) Add(ints ...int64) {
	for _, i := range ints {
		if !s.ints.Include(i) {
			s.orderedInts = append(s.orderedInts, i)
			s.ints.Add(i)
		}
	}
}

// Include はiが含まれていればtrueを返す。
func (s *Int64OrderedSet) Include(i int64) bool {
	return s.ints.Include(i)
}

// Size は含まれる要素数を返す。
func (s *Int64OrderedSet) Size() int {
	return s.ints.Size()
}

// Merge はs2の要素を追加する。
func (s *Int64OrderedSet) Merge(s2 *Int64OrderedSet) {
	s.Add(s2.orderedInts...)
}

// Slice はint64の集合をスライスとして返す。
func (s *Int64OrderedSet) Slice() []int64 {
	res := make([]int64, len(s.orderedInts))
	copy(res, s.orderedInts)
	return res
}

// Subtract はsに含まれていてs2に含まれていない要素の集合を返す。
func (s *Int64OrderedSet) Subtract(s2 int64Includer) Int64OrderedSet {
	res := NewInt64OrderedSet()
	for _, i := range s.orderedInts {
		if !s2.Include(i) {
			res.Add(i)
		}
	}
	return *res
}

type intIncluder interface {
	Include(i int) bool
}

// IntSet はintの集合を表す。
type IntSet map[int]struct{}

// NewIntSet はIntSetを新たに生成する。
func NewIntSet(ints ...int) *IntSet {
	s := IntSet{}
	s.Add(ints...)
	return &s
}

// Add は要素を追加する。
func (s *IntSet) Add(ints ...int) {
	for _, i := range ints {
		(*s)[i] = struct{}{}
	}
}

// Include はiが含まれてばtrueを返す。
func (s *IntSet) Include(i int) bool {
	_, ok := (*s)[i]
	return ok
}

// Size はIntSetに含まれる要素数を返す。
func (s *IntSet) Size() int {
	return len(*s)
}

// Merge はs2の要素を追加する。
func (s *IntSet) Merge(s2 *IntSet) {
	for k := range *s2 {
		(*s)[k] = struct{}{}
	}
}

// Slice はintの集合をスライスとして返す。
func (s *IntSet) Slice() []int {
	ints := []int{}
	for k := range *s {
		ints = append(ints, k)
	}
	return ints
}

// Subtract はsに含まれていてs2に含まれていない要素の集合を返す。
func (s *IntSet) Subtract(s2 intIncluder) IntSet {
	res := IntSet{}
	for k := range *s {
		if !s2.Include(k) {
			res.Add(k)
		}
	}
	return res
}

// IntOrderedSet は要素を追加した順序を保持するintの集合を表す
type IntOrderedSet struct {
	orderedInts []int
	ints        IntSet
}

// NewIntOrderedSet はInt64OrderSetを新たに生成する
func NewIntOrderedSet(ints ...int) *IntOrderedSet {
	s := IntOrderedSet{
		orderedInts: nil,
		ints:        IntSet{},
	}
	s.Add(ints...)
	return &s
}

// Add は要素を追加する（既に含まれる場合は何もしない）
func (s *IntOrderedSet) Add(ints ...int) {
	for _, i := range ints {
		if !s.ints.Include(i) {
			s.orderedInts = append(s.orderedInts, i)
			s.ints.Add(i)
		}
	}
}

// Include はintが含まれていればtrueを返す。
func (s *IntOrderedSet) Include(i int) bool {
	return s.ints.Include(i)
}

// Size は含まれる要素数を返す。
func (s *IntOrderedSet) Size() int {
	return s.ints.Size()
}

// Merge はs2の要素を追加する。
func (s *IntOrderedSet) Merge(s2 *IntOrderedSet) {
	s.Add(s2.orderedInts...)
}

// Slice はintの集合をスライスとして返す。
func (s *IntOrderedSet) Slice() []int {
	res := make([]int, len(s.orderedInts))
	copy(res, s.orderedInts)
	return res
}

// Subtract はsに含まれていてs2に含まれていない要素の集合を返す。
func (s *IntOrderedSet) Subtract(s2 intIncluder) IntOrderedSet {
	res := NewIntOrderedSet()
	for _, i := range s.orderedInts {
		if !s2.Include(i) {
			res.Add(i)
		}
	}
	return *res
}
