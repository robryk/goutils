package teller

import "errors"
import "io"

type Writer struct {
	io.Writer
	offset int64
}

var ErrSeekUnsupported = errors.New("teller: Nonzero seeks are unsupported")

func NewWriter(w io.Writer) io.WriteSeeker {
	if ws, ok := w.(io.WriteSeeker); ok {
		return ws
	}
	return &Writer{Writer: w}
}

func (w *Writer) Write(p []byte) (n int, err error) {
	n, err = w.Writer.Write(p)
	w.offset += int64(n)
	return
}

func (w *Writer) Seek(offset int64, whence int) (ret int64, err error) {
	ret = w.offset
	if offset != 0 || whence != 1 {
		err = ErrSeekUnsupported
	}
	return
}

type Reader struct {
	io.Reader
	offset int64
}

func NewReader(r io.Reader) io.ReadSeeker {
	if rs, ok := r.(io.ReadSeeker); ok {
		return rs
	}
	return &Reader{Reader: r}
}

func (r *Reader) Read(p []byte) (n int, err error) {
	n, err = r.Reader.Read(p)
	r.offset += int64(n)
	return
}

func (r *Reader) Seek(offset int64, whence int) (ret int64, err error) {
	ret = r.offset
	if offset != 0 || whence != 1 {
		err = ErrSeekUnsupported
	}
	return
}
