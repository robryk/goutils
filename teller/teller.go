package teller

import "errors"
import "io"

type Writer struct {
	io.Writer
	seeker io.Seeker // Supplied Writer as a Seeker or nil
	offset int64
}

var ErrSeekUnsupported = errors.New("teller: Nonzero seeks are unsupported")

func NewWriter(w io.Writer) io.WriteSeeker {
	ws := &Writer{Writer: w}
	if s, ok := w.(io.WriteSeeker); ok {
		ws.seeker = s
	}
	return ws
}

func (w *Writer) Write(p []byte) (n int, err error) {
	n, err = w.Writer.Write(p)
	w.offset += int64(n)
	return
}

func (w *Writer) Seek(offset int64, whence int) (ret int64, err error) {
	err = ErrSeekUnsupported
	if w.seeker != nil {
		ret, err = w.seeker.Seek(offset, whence)
		if err == nil {
			w.offset = ret
			return
		}
	}
	if offset != 0 || whence != 1 {
		// We can't seek or have failed to seek and the seek wasn't a nop.
		return
	}
	return w.offset, nil
}

type Reader struct {
	io.Reader
	seeker io.Seeker // Supplied Reader as a Seeker or nil
	offset int64
}

func NewReader(r io.Reader) io.ReadSeeker {
	rs := &Reader{Reader: r}
	if s, ok := r.(io.ReadSeeker); ok {
		rs.seeker = s
	}
	return rs
}

func (r *Reader) Read(p []byte) (n int, err error) {
	n, err = r.Reader.Read(p)
	r.offset += int64(n)
	return
}

func (r *Reader) Seek(offset int64, whence int) (ret int64, err error) {
	err = ErrSeekUnsupported
	if r.seeker != nil {
		ret, err = r.seeker.Seek(offset, whence)
		if err == nil {
			r.offset = ret
			return
		}
	}
	if offset != 0 || whence != 1 {
		// We can't seek or have failed to seek and the seek wasn't a nop.
		return
	}
	return r.offset, nil
}
