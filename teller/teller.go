// Stub Seek() method so that current position can be retrieved from Readers
// or Writers that don't implement Seeker.
package teller

import "errors"
import "io"

type writer struct {
	io.Writer
	seeker io.Seeker // Supplied Writer cast to a Seeker or nil
	offset int64
}

// Error returned when a non-nop seek was attempted on a wrapped Reader or Writer
// that wraps a non-seekable Reader of Writer.
var ErrSeekUnsupported = errors.New("teller: Nonzero seeks are unsupported")

// Wraps a Writer into a WriteSeeker. On Seek we first attempt to seek using
// original Writer; if that fails and a nop seek was requested, we return
// position based on counting bytes written through the wrapped Writer.
func NewWriter(w io.Writer) io.WriteSeeker {
	ws := &writer{Writer: w}
	if s, ok := w.(io.WriteSeeker); ok {
		ws.seeker = s
	}
	return ws
}

func (w *writer) Write(p []byte) (n int, err error) {
	n, err = w.Writer.Write(p)
	w.offset += int64(n)
	return
}

func (w *writer) Seek(offset int64, whence int) (ret int64, err error) {
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

type reader struct {
	io.Reader
	seeker io.Seeker // Supplied Reader cast to a Seeker or nil
	offset int64
}

// Wraps a Reader into a ReadSeeker. On Seek we first attempt to seek using
// original Reader; if that fails and a nop seek was requested, we return
// position based on counting bytes read through the wrapped Reader.
func NewReader(r io.Reader) io.ReadSeeker {
	rs := &reader{Reader: r}
	if s, ok := r.(io.ReadSeeker); ok {
		rs.seeker = s
	}
	return rs
}

func (r *reader) Read(p []byte) (n int, err error) {
	n, err = r.Reader.Read(p)
	r.offset += int64(n)
	return
}

func (r *reader) Seek(offset int64, whence int) (ret int64, err error) {
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
