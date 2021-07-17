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

// DefaultCapacity default capacity
//const DefaultCapacity int = 5000

// UnsafeSlice safe slice
type UnsafeSlice struct {
	items    []interface{}
	capacity int
}

// UnsafeSliceItem contains the index/value pair of an item in a
// concurrent slice
type UnsafeSliceItem struct {
	Index int
	Value interface{}
}

// NewUnsafeSlice creates a new concurrent slice
func NewUnsafeSlice(c int) *UnsafeSlice {
	s := &UnsafeSlice{
		items:    make([]interface{}, 0, c),
		capacity: c,
	}

	return s
}

// Count count
func (s *UnsafeSlice) Count() int {
	return len(s.items)
}

// Prepend prepend an item to the concurrent slice
func (s *UnsafeSlice) Prepend(item interface{}) {
	s.items = append([]interface{}{item}, s.items...)
}

// Append adds an item to the concurrent slice
func (s *UnsafeSlice) Append(item interface{}) {
	s.items = append(s.items, item)
}

// AppendAll adds an items to the concurrent slice
func (s *UnsafeSlice) AppendAll(items ...interface{}) {
	s.items = append(s.items, items...)
}

// Paste paste an items to the concurrent slice
func (s *UnsafeSlice) Paste(items []interface{}) {
	s.items = append(s.items, items...)
}

// AppendSlice adds an items to the concurrent slice
func (s *UnsafeSlice) AppendSlice(ss *UnsafeSlice) {
	s.AppendAll(ss.items...)
}

// FilterOut filter an items out from the concurrent slice
func (s *UnsafeSlice) FilterOut(f OnFilter) []interface{} {
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
func (s *UnsafeSlice) Fetch() []interface{} {
	old := s.items
	s.items = make([]interface{}, 0, s.capacity)
	return old
}

// FetchWithLimit fetch slice with limit
func (s *UnsafeSlice) FetchWithLimit(limit int) []interface{} {
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
func (s *UnsafeSlice) Shift() (front interface{}, count int) {
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
func (s *UnsafeSlice) Iterate() <-chan UnsafeSliceItem {
	c := make(chan UnsafeSliceItem)

	f := func() {
		for index, value := range s.items {
			c <- UnsafeSliceItem{index, value}
		}
		close(c)
	}
	go f()

	return c
}
