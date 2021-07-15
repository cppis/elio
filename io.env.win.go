// +build windows

package elio

// GetMaxInCount get max in count
func GetMaxInCount(io string) (c int) {
	c = DefaultDefIOInCount

	AppInfo().Msgf("get.max.incount %d of model:'%s'", c, IoDefault.String())
	return c
}

// GetCurrentIO get current io
func GetCurrentIO(io string) IOs {
	return IoDefault
}

// GenIO gen IO
func GenIO(io string) IO {
	AppInfo().Msgf("gen iomodel:'%s'", IoDefault.String())
	return NewIoDefault()
}

// SetLimit set limit
func SetLimit() (err error) {
	AppInfo().Msgf("set.limit not supported")
	return nil
}
