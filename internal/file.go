package mtee

import "os"

type file struct {
	*os.File
}

func (f *file) write(b []byte) (int, error) {
	return f.File.Write(b)
}
