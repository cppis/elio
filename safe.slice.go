// this object referred from dnaeon/gru/utils/slice.go
//
// Copyright (c) 2015-2017 Marin Atanasov Nikolov <dnaeon@gmail.com>
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions
// are met:
//
//  1. Redistributions of source code must retain the above copyright
//     notice, this list of conditions and the following disclaimer
//     in this position and unchanged.
//  2. Redistributions in binary form must reproduce the above copyright
//     notice, this list of conditions and the following disclaimer in the
//     documentation and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE AUTHOR(S) ``AS IS'' AND ANY EXPRESS OR
// IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES
// OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED.
// IN NO EVENT SHALL THE AUTHOR(S) BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT
// NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF
// THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package elio

import "sync"

// DefaultCapacity default capacity
const DefaultCapacity int = 5000

// SafeSlice safe slice
type SafeSlice struct {
	sync.RWMutex
	items    []interface{}
	capacity int
}

// SafeSliceItem contains the index/value pair of an item in a
// concurrent slice
type SafeSliceItem struct {
	Index int
	Value interface{}
}

// NewSafeSlice creates a new safe slice
func NewSafeSlice(c int) *SafeSlice {
	s := &SafeSlice{
		items:    make([]interface{}, 0, c),
		capacity: c,
	}

	return s
}

// Count count
func (s *SafeSlice) Count() int {
	s.Lock()
	defer s.Unlock()

	return len(s.items)
}

// Prepend prepend an item to the concurrent slice
func (s *SafeSlice) Prepend(item interface{}) {
	s.Lock()
	defer s.Unlock()

	s.items = append([]interface{}{item}, s.items...)
}

// Append adds an item to the concurrent slice
func (s *SafeSlice) Append(item interface{}) {
	s.Lock()
	defer s.Unlock()

	s.items = append(s.items, item)
}

// AppendAll adds an items to the concurrent slice
func (s *SafeSlice) AppendAll(items ...interface{}) {
	s.Lock()
	defer s.Unlock()

	s.items = append(s.items, items...)
}

// Paste paste an items to the concurrent slice
func (s *SafeSlice) Paste(items []interface{}) {
	s.Lock()
	defer s.Unlock()

	s.items = append(s.items, items...)
}

// AppendSlice adds an items to the concurrent slice
func (s *SafeSlice) AppendSlice(ss *SafeSlice) {
	s.AppendAll(ss.items...)
}

// OnFilter on filter
type OnFilter func(x interface{}) bool

// FilterOut filter an items out from the concurrent slice
func (s *SafeSlice) FilterOut(f OnFilter) []interface{} {
	s.Lock()
	defer s.Unlock()

	results := s.items[:0]
	var filtered []interface{}
	for _, d := range s.items {
		if false == f(d) {
			results = append(results, d)
		} else {
			filtered = append(filtered, d)
		}
	}

	s.items = results

	return filtered
}

// Fetch fetch slice and renew
func (s *SafeSlice) Fetch() []interface{} {
	s.Lock()
	defer s.Unlock()

	old := s.items
	s.items = make([]interface{}, 0, s.capacity)
	return old
}

// FetchWithLimit fetch slice with limit
func (s *SafeSlice) FetchWithLimit(limit int) []interface{} {
	s.Lock()
	defer s.Unlock()

	var fetch []interface{}

	if limit < len(s.items) {
		// var prev []interface{}
		// fetch = make([]interface{}, 0, limit)
		// prev, s.items = s.items[:limit], s.items[limit:]
		// fetch = append(fetch, prev...)
		fetch, s.items = s.items[:limit], s.items[limit:]
	} else {
		fetch = s.items
		s.items = make([]interface{}, 0, s.capacity)
	}
	return fetch
}

// Shift shift slice and renew
func (s *SafeSlice) Shift() (front interface{}, count int) {
	s.Lock()
	defer s.Unlock()

	count = len(s.items)
	if 0 < count {
		var front interface{}
		front, s.items = s.items[0], s.items[1:]
		return front, (count - 1)
	}

	return nil, 0
}

// Iterate iterates over the items in the concurrent slice
// Each item is sent over a channel, so that
// we can iterate over the slice using the builin range keyword
func (s *SafeSlice) Iterate() <-chan SafeSliceItem {
	c := make(chan SafeSliceItem)

	f := func() {
		s.Lock()
		defer s.Unlock()
		for index, value := range s.items {
			c <- SafeSliceItem{index, value}
		}
		close(c)
	}
	go f()

	return c
}
