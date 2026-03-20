package bitset

import (
	"iter"
	"math/bits"
)

// BitSet128 is a bitset implementation for integers in the range [0, 127],
// based on two `uint64` bitmasks.
type BitSet128 struct {
	lo uint64
	hi uint64
}

func (s *BitSet128) Add(item int) {
	if item < 64 {
		s.lo |= 1 << uint(item)
	} else {
		s.hi |= 1 << uint(item-64)
	}
}

func (s *BitSet128) Remove(item int) {
	if item < 64 {
		s.lo &= ^(1 << uint(item))
	} else {
		s.hi &= ^(1 << uint(item-64))
	}
}

func (s BitSet128) Contains(item int) bool {
	if item < 64 {
		return s.lo&(1<<uint(item)) != 0
	}
	return s.hi&(1<<uint(item-64)) != 0
}

func (s BitSet128) Size() int {
	return bits.OnesCount64(s.lo) + bits.OnesCount64(s.hi)
}

func (s BitSet128) Empty() bool {
	return s.lo == 0 && s.hi == 0
}

func (s BitSet128) Union(a BitSet128) BitSet128 {
	return BitSet128{s.lo | a.lo, s.hi | a.hi}
}

func (s BitSet128) Intersection(a BitSet128) BitSet128 {
	return BitSet128{s.lo & a.lo, s.hi & a.hi}
}

func (s BitSet128) Difference(a BitSet128) BitSet128 {
	return BitSet128{s.lo &^ a.lo, s.hi &^ a.hi}
}

func (s BitSet128) All() iter.Seq[int] {
	return func(yield func(int) bool) {
		lo := s.lo
		for lo != 0 {
			lsb := bits.TrailingZeros64(lo)
			if !yield(lsb) {
				return
			}
			lo &= lo - 1
		}
		hi := s.hi
		for hi != 0 {
			lsb := bits.TrailingZeros64(hi)
			if !yield(lsb + 64) {
				return
			}
			hi &= hi - 1
		}
	}
}

func (s BitSet128) Values() []int {
	values := make([]int, 0, s.Size())
	for k := range s.All() {
		values = append(values, k)
	}
	return values
}
