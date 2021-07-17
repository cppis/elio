package elio

// Event event interface
type Event interface {
	Handle()
	String() string
}
