// +build windows

package elio

import "syscall"

const (
	ShutRd   = 0x0
	ShutWr   = 0x1
	ShutRdWr = 0x2

	SigInt  = syscall.SIGINT
	SigTerm = syscall.SIGTERM
)

type Sockaddr interface{}
