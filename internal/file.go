package mtee

type file interface {
	Write(b []byte) (int, error)
	Close() error
}

/*
type file struct {
	*os.File
}

func (f *file) write(b []byte) (int, error) {
	return f.File.Write(b)
}
*/
