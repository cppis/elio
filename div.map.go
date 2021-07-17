package elio

import (
	"fmt"
)

const (
	// InvalidDivIndex partition invalid index
	InvalidDivIndex uint32 = ^uint32(0)
)

// DivMap division map
type DivMap struct {
	division []*UnsafeMap
	current  uint32
}

// NewDivMap new division map
func NewDivMap(l uint32) *DivMap {
	s := new(DivMap)
	if nil != s {
		s.division = make([]*UnsafeMap, l, l)

		for i := uint32(0); i < l; i++ {
			s.division[i] = NewUnsafeMap()
		}
	}

	return s
}

// IsValidIndex is valid index
func (d *DivMap) IsValidIndex(i uint32) bool {
	if i <= InvalidDivIndex {
		return false
	}

	return true
}

// GetCurrent get current partition
func (d *DivMap) GetCurrent() (c uint32, s *UnsafeMap) {
	l := len(d.division)
	if l <= 0 {
		return InvalidDivIndex, nil
	}

	defer func() {
		if 0 < l {
			d.current++
			d.current = d.current % uint32(l)
		}
	}()

	c = d.current

	return c, d.division[c]
}

// GetLeast get least partition
func (d *DivMap) GetLeast() (i uint32, s *UnsafeMap) {
	l := 0

	for c, p := range d.division {
		count := d.Count()

		if 0 == l {
			i = uint32(c)
			s = p
			l = count
		} else if count < l {
			i = uint32(c)
			s = p
			l = count
		}
	}

	return i, s
}

// Get get partition
func (d *DivMap) Get(i uint32) (s *UnsafeMap, err error) {
	if len(d.division) < int(i) {
		err = fmt.Errorf("invalid index %d", i)

	} else {
		s = d.division[i]
	}

	return s, err
}

// Count get count
func (d *DivMap) Count() int {
	return len(d.division)
}

// GetCounts get partitions count
func (d *DivMap) GetCounts() []int {
	var c []int
	for _, s := range d.division {
		c = append(c, s.Count())
	}

	return c
}

// Set set to partition
func (d *DivMap) Set(k uint64, v interface{}) (uint32, *UnsafeMap) {
	i, m := d.GetCurrent()
	if InvalidDivIndex != i {
		m.Set(k, v)
	}

	return i, m
}

// Del del from partition
func (d *DivMap) Del(i uint32, k uint64) (v interface{}, ok bool) {
	m, err := d.Get(i)
	if nil == err {
		ok = true
		v, ok = m.Get(k)
		if ok {
			m.Del(k)
		}
	}

	return v, ok
}
