package countreader

import "io"
import "sync"

type CountReader struct {
	r     io.Reader
	count int
	mu    sync.Mutex
}

func (cr *CountReader) Read(p []byte) (n int, err error) {
	n, err = cr.r.Read(p)
	cr.mu.Lock()
	defer cr.mu.Unlock()
	cr.count += n
	return
}

func (cr *CountReader) Count() int {
	cr.mu.Lock()
	defer cr.mu.Unlock()
	return cr.count
}

