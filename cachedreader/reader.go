package cachedreader

import "io"
import "sync"

// CachedReader transforms a Reader into a ReaderAt by caching reader's output.
type CachedReader struct {
	source io.Reader
	buffer []byte
	err    error
	mu     sync.Mutex
}

func (ci *CachedReader) ReadAt(p []byte, off int) (int, error) {
	// If need be, this can be changes to an rwlock and locked for writing only when
	// we need to read more (or rather, released and relocked for writing).
	ci.mu.Lock()
	defer ci.mu.Unlock()

	// If there is not enough data in buffer and we haven't failed to read more
	if len(ci.buffer) < off+len(p) && ci.err == nil {
		if cap(ci.buffer) < off+len(p) {
			newCap := off + len(p)
			if newCap < 2*cap(ci.buffer) {
				newCap = 2*cap(ci.buffer)
			}
			newBuffer := make([]byte, 0, newCap)
			newBuffer = append(newBuffer, ci.buffer...)
			ci.buffer = newBuffer
		}
		var n int
		n, ci.err = ci.source.Read(ci.buffer[len(ci.buffer) : off+len(p)])
		ci.buffer = ci.buffer[:len(ci.buffer)+n]
	}

	n := copy(p, ci.buffer[off:])
	if n < len(p) {
		if ci.err == nil {
			panic("CachedReader is short on data and has no errors")
		}
		return n, ci.err
	}
	// If we have the data to satisfy the request, we return no errors
	return n, nil
}

// Returns the memory usage of the buffer in bytes.
func (ci *CachedReader) MemoryUsage() int {
	return cap(ci.buffer)
}

// Creates a CachedReader from the given reader. The source reader will be read from
// nonconcurrently, but in various goroutines.
func NewCachedReader(source io.Reader) *CachedReader {
	return &CachedReader{
		source: source,
		buffer: nil,
		err: nil,
	}
}

