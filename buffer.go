package elio

// Buffer is a helper type for managing input streams from inside the Data event.
type Buffer struct{ b []byte }

// Bytes get bytes of buffer.
func (b *Buffer) Bytes() []byte {
	return b.b
}

// Len get length of buffer.
func (b *Buffer) Len() int {
	return len(b.b)
}

// Begin accepts a new packet and returns a working sequence of unprocessed bytes.
func (b *Buffer) Begin(pos int) []byte {
	return b.b[pos:]
}

// Get get contents
func (b *Buffer) Get(begin, end int) []byte {
	return b.b[begin:end]
}

// GetFrom get contents from begin
func (b *Buffer) GetFrom(begin int) []byte {
	return b.b[begin:]
}

// GetTo get contents to end
func (b *Buffer) GetTo(end int) []byte {
	return b.b[:end]
}

// Put put contents to buffer
func (b *Buffer) Put(data []byte) ([]byte, int) {
	if len(data) > 0 {
		b.b = append(b.b, data...)
	}

	return b.b[0:], len(b.b)
}

// Clear clear buffer
func (b *Buffer) Clear(pos int) {
	if pos < 0 {
		b.b = b.b[:0]
	} else {
		b.b = b.b[pos:]
	}
}
