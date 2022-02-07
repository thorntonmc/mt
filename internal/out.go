package mtee

import (
	"bufio"
	"os"
	"sync"
)

// out wraps the output's file, buffer, and mutex
type out struct {
	file *os.File
	buf  *bufio.Writer
	mu   sync.Mutex
}

// Write locks the files mutex, defers the unlock, and writes to the buffer
func (o *out) Write(b []byte) (n int, err error) {
	o.mu.Lock()
	defer o.mu.Unlock()
	return o.buf.Write(b)
}

// Close flushes the file's buffer and closes the file
func (o *out) Close() error {
	err := o.buf.Flush()
	if err != nil {
		return err
	}
	return o.file.Close()
}

// newOut takes a file and a buffer size, and returns an out
func newOut(f *os.File, n int) *out {
	return &out{
		file: f,
		buf:  bufio.NewWriterSize(f, n),
	}
}
