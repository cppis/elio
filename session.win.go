// +build windows

package elio

// GetAddr ...
func (n *Session) GetAddr() (addr string, ok bool) {
	n.ipRemote = n.conn.RemoteAddr().String()
	return n.ipRemote, ok
}

// sysWrite syscall write
func (n *Session) sysWrite(out []byte) (written int, err error) {
	return n.conn.Write(out[written:])
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
