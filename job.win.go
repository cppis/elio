// +build windows

package elio

// // Work work
// func (j *WriteJob) Work() error {
// 	var err error
// 	var out []byte
// 	var written int

// 	outs := j.session.outQueue.Fetch()
// 	//fmt.Printf("outs %d\n", len(outs))

// 	for _, o := range outs {
// 		out = append(out, o.(*ByteBuffer).Bytes()...)
// 		PutByteBuffer(o.(*ByteBuffer))
// 	}

// 	//fmt.Printf("out len %d\n", len(out))

// 	l := len(out)
// 	for written < l {
// 		var w int
// 		w, err = j.session.service.io.Write(j.session, out[written:])
// 		if 0 < w {
// 			//fmt.Printf("-%d", w)
// 			written += w

// 		} else {
// 			//AppError().Str(LogObject, j.String()).Str(LogSession, j.session.String()).
// 			AppError().Str(LogSession, j.session.String()).
// 				Err(err).Msgf("writing failed with fd:%v written:%d w:%d count.outqueue:%d", j.session.fd, written, w, l)

// 			break
// 		}
// 	}

// 	return err
// }
