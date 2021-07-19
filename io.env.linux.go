// +build linux

package elio

import (
	"golang.org/x/sys/unix"
)

// GetMaxInCount get max in count
func GetMaxInCount(io string) (c int) {
	c = DefaultPollInCount
	if IoDefault.String() == io {
		c = DefaultDefIOInCount
	}

	//AppInfo().Msgf("get.max.incount %d of model:'%s'", c, io)
	return c
}

// GetCurrentIO get current io
func GetCurrentIO(io string) IOs {
	if IoDefault.String() == io {
		return IoDefault
	}

	return IoPoll
}

// GenIO gen io
func GenIO(io string) IoModel {
	if IoDefault.String() == io {
		AppInfo().Msgf("gen io.model:'%s'", IoDefault.String())
		return NewIoDefault()
	}

	AppInfo().Msgf("gen io.model:%s'", IoPoll.String())
	return NewIoPoll()
}

// SetLimit set limit
func SetLimit() (err error) {
	var rLimit unix.Rlimit
	if err = unix.Getrlimit(unix.RLIMIT_NOFILE, &rLimit); nil == err {
		rLimit.Cur = rLimit.Max
		if err = unix.Setrlimit(unix.RLIMIT_NOFILE, &rLimit); nil == err {
		}
	}

	AppInfo().Msgf("set.limit set current limit:%d", rLimit.Cur)

	if nil != err {
		AppError().Err(err).Msgf("failed to set limit")
	}

	return err
}
