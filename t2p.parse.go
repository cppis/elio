package elio

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

const delimiterLfLf string = "\n\n"
const delimiterCrLfCrLf string = "\r\n\r\n"

// T2pOnParse control on parse
func T2pOnParse(in []byte) (lenParsed, lenDelimit int, ok bool) {
	ok = true

	if lenParsed = bytes.Index(in, []byte(delimiterLfLf)); 0 <= lenParsed {
		lenDelimit = len(delimiterLfLf)
	} else if lenParsed = bytes.Index(in, []byte(delimiterCrLfCrLf)); 0 <= lenParsed {
		lenDelimit = len(delimiterCrLfCrLf)
	} else {
		ok = false
	}

	return lenParsed, lenDelimit, ok
}

// T2pParseCommand parse command
func T2pParseCommand(in []byte) (out []byte) {
	lenIn := len(in)
	if 0 < lenIn {
		out = in[:lenIn]
		var c = 0

		// todo: process request in here
		for i := 0; i < lenIn; i++ {
			if '\b' == in[i] {
				if 0 < c {
					c--
				}
			} else {
				out[c] = in[i]
				c++
			}
		}
		out = out[:c]
	}

	return out
}

// T2PParse T2P parse
func T2PParse(in []byte) (out []string) {
	buffer := bytes.NewBuffer(in)
	reader := bufio.NewReader(buffer)

	for {
		line, _, err := reader.ReadLine()
		if io.EOF == err {
			break
		}

		l := strings.TrimSpace(string(line))
		if "" != l {
			out = append(out, strings.TrimSpace(string(l)))
		}
	}

	return out
}
