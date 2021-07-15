// +build linux

package elio

import (
	"fmt"
	"net"

	"golang.org/x/sys/unix"
)

// // onRelease
// func (n *Session) onRelease(m interface{}) {
// 	poll := m.(*ioPoll)

// 	m.poll.DeleteR(n.fd)
// 	m.Close(n)

// 	m.GetService().sessionCmap.Remove(fmt.Sprintf("0x%08x", n.fd))
// }

// GetAddr ...
func (n *Session) GetAddr() (string, bool) {
	// No attempt at error reporting because there are no possible errors,
	// and the caller won't report them anyway.
	//var sa unix.Sockaddr
	//sa, _ := unix.Getsockname(n.fd)
	sa, _ := unix.Getpeername(n.fd)
	switch sa := sa.(type) {
	case *unix.SockaddrInet4:
		//return &sa.Addr, sa.Port, true
		ip := net.IP{sa.Addr[0], sa.Addr[1], sa.Addr[2], sa.Addr[3]}
		n.ipRemote = fmt.Sprintf("%s:%d", ip, sa.Port)
	case *unix.SockaddrInet6:
		//return &sa.Addr, sa.Port, true
		ip := net.IP{
			sa.Addr[0], sa.Addr[1], sa.Addr[2], sa.Addr[3],
			sa.Addr[4], sa.Addr[5], sa.Addr[6], sa.Addr[7],
			sa.Addr[8], sa.Addr[9], sa.Addr[10], sa.Addr[11],
			sa.Addr[12], sa.Addr[13], sa.Addr[14], sa.Addr[15],
		}
		n.ipRemote = fmt.Sprintf("%s:%d", ip, sa.Port)
	}
	return n.ipRemote, true
}

// sysWrite syscall write
func (n *Session) sysWrite(out []byte) (written int, err error) {
	return unix.Write(n.fd, out[written:])
}

// sysWriteAll syscall write all
func (n *Session) sysWriteAll(out []byte) (written int, err error) {
	var w int
	l := len(out)
	for written < l {
		w, err = n.sysWrite(out[written:])
		if 0 < w {
			written += w

		} else {
			break
		}
	}

	return written, err
}
