package elio

// IOs io models type definition
type IOs int

const (
	// IoAuto IO auto
	IoAuto IOs = iota
	// IoDefault IO default(golang method)
	IoDefault
	// IoPoll IO poll
	IoPoll
)

const (
	// DefaultDefIOInCount default def IO in count
	DefaultDefIOInCount = 1
	// DefaultPollInCount default poll in count
	DefaultPollInCount = 4
)

// String string
func (i IOs) String() string {
	return [...]string{"auto", "default", "poll"}[i]
}

// IOsFromString IO models from string
func IOsFromString(m string) IOs {
	switch m {
	case IoAuto.String():
		return IoAuto
	case IoPoll.String():
		return IoPoll
	default:
		return IoDefault
	}
}
