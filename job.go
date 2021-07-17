package elio

import (
	"fmt"
)

// Job job interface
type Job interface {
	String() string
	Work() error
}

// WriteJob write job
type WriteJob struct {
	session *Session
	out     []byte
}

// String string
func (j *WriteJob) String() string {
	return fmt.Sprintf("WriteJob::%p", j)
}

// Work work
func (j *WriteJob) Work() (err error) {
	defer func() {
		j.session.DecRef()
	}()

	var out []byte
	var written int

	outs := j.session.outQueue.Fetch()

	if 0 < len(outs) {
		for _, o := range outs {
			out = append(out, o.(*ByteBuffer).Bytes()...)
			PutByteBuffer(o.(*ByteBuffer))
		}

		written, err = j.session.ioCore.io.Write(j.session, out)
		if len(out) == written {
		} else {
			AppError().Str(LogSession, j.session.String()).
				Err(err).Msgf("writing failed with fd:%v written:%d", j.session.fd, written)
		}
	}

	return err
}
