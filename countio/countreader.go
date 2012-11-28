package countreader

import "io"
import "sync"

type countReader struct {
	r     io.Reader
	count int
	mu    sync.Mutex
}

func (cr *countReader) Read(p []byte) (n int, err error) {
	n, err = cr.r.Read(p)
	cr.mu.Lock()
	cr.count += n
	cr.mu.Unlock()
	return
}
