package mtee

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
)

const (
	defaultBufSize = 4096
	stdOutBufSize  = 1
)

// write t to v, store results in c
func writeAndStore(b []byte, w io.Writer, c chan teeResult) {
	_, err := w.Write(b)
	c <- newTeeResult(err)
}

// tee scans from a scanner and writes to all outs, storing
// the results in channel
func tee(s *bufio.Scanner, w []*out, c chan teeResult) error {
	s.Scan()
	b := s.Bytes()
	b = append(b, '\n')

	for _, v := range w {
		go writeAndStore(b, v, c)
	}

	for i := 0; i < len(c); i++ {
		result := <-c
		if !result.ok {
			return result.err
		}
	}
	return nil
}

// mtee is the application mtee and its configurations
type mtee struct {
	out     []*out
	in      *os.File
	scanner *bufio.Scanner
	results chan teeResult
}

// teeResult is the result of a mtee goroutine
type teeResult struct {
	ok  bool
	err error
}

// newTeeResult takes an error and wraps it in a teeResult
func newTeeResult(err error) teeResult {
	return teeResult{
		err == nil,
		err,
	}
}

func (m *mtee) init(files []string, modeAppend bool) error {
	// set outputs
	numOut := 1 + len(files)
	m.out = make([]*out, numOut)
	m.in = os.Stdin
	m.scanner = bufio.NewScanner(m.in)

	if len(files) > 0 {
		err := m.setFiles(files, modeAppend)
		if err != nil {
			return err
		}
	}

	// last output is stdout
	m.out[numOut-1] = newOut(os.Stdout, stdOutBufSize)
	m.results = make(chan teeResult, len(files))
	return nil

}

func (m *mtee) setFiles(files []string, modeAppend bool) error {
	numFiles := len(files)
	results := make(chan teeResult, numFiles)

	for i, v := range files {
		go func(i int, v string) {
			err := m.setOut(v, i, modeAppend)
			results <- newTeeResult(err)
		}(i, v)
	}

	for i := 0; i < numFiles; i++ {
		result := <-results
		if !result.ok {
			return result.err
		}
	}
	return nil
}

func (m *mtee) setOut(fstr string, index int, modeAppend bool) error {
	var f *os.File
	var err error

	// open or create the file
	switch modeAppend {
	case true:
		f, err = os.OpenFile(fstr, os.O_APPEND|os.O_WRONLY, 0755)
	case false:
		f, err = os.Create(fstr)
	}

	if err != nil {
		return err
	}

	m.out[index] = newOut(f, defaultBufSize)

	return nil
}

// tee scans text from in and writes it to all outs
func (m *mtee) tee() error {
	if len(m.out) < 1 {
		return fmt.Errorf("FATAL: tee() - requires at least one out")
	}
	if m.scanner == nil {
		return fmt.Errorf("FATAL tee() - no scanner")
	}

	return tee(m.scanner, m.out, m.results)
}

func (m *mtee) run() error {
	done := false
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		done = true
	}()
	for !done {
		err := m.tee()
		if err != nil {
			return err
		}
	}
	return nil
}

// Run runs mtee from a file string and mode string
func Run(files []string, mode bool) error {
	m := &mtee{}

	err := m.init(files, mode)
	defer m.in.Close()
	if err != nil {
		return err
	}

	err = m.run()
	// if we opened a file, we need to close it
	for _, f := range m.out {
		f.Close()
	}
	return err

}
