// +build linux

package elio

import (
	"syscall"

	"golang.org/x/sys/unix"
)

const (
	// ShutRd shut rd
	ShutRd = unix.SHUT_RD
	// ShutWr shut wr
	ShutWr = unix.SHUT_WR
	// ShutRdWr shut rdwr
	ShutRdWr = unix.SHUT_RDWR

	// SigInt sig int
	SigInt = syscall.SIGINT
	// SigTerm sig term
	SigTerm = syscall.SIGTERM
)

// Sockaddr sockaddr type definition
type Sockaddr unix.Sockaddr
