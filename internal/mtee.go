package mtee

import (
	"bufio"
	"fmt"
	"os"
)

// mtee represents the application mtee and its configurations
type mtee struct {
	out []*os.File
	in  *os.File
}

// teeResult represents the result of a mtee goroutine
type teeResult struct {
	ok  bool
	err error
}

func newTeeResult(err error) teeResult {
	return teeResult{
		err == nil,
		err,
	}
}

func (m *mtee) init(files []string, modeAppend bool) error {
	// set outputs
	numOut := 1 + len(files)
	m.out = make([]*os.File, numOut)
	m.in = os.Stdin

	if len(files) > 0 {
		err := m.setFiles(files, modeAppend)
		if err != nil {
			return err
		}
	}

	// last output is stdout
	m.out[numOut-1] = os.Stdout
	return nil

}

func (m *mtee) setFiles(files []string, modeAppend bool) error {
	for i, v := range files {
		err := m.setOut(v, i, modeAppend)
		if err != nil {
			return err
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

	m.out[index] = f
	return nil
}

// tee scans text from in and writes it to out
func (m *mtee) tee() error {
	numOuts := len(m.out)
	results := make(chan teeResult, numOuts)

	scanner := bufio.NewScanner(m.in)
	scanner.Scan()
	text := fmt.Sprintf("%s\n", scanner.Text())

	for _, v := range m.out {
		go func(v *os.File) {
			_, err := v.Write([]byte(text))
			results <- newTeeResult(err)
		}(v)
	}

	for i := 0; i < numOuts; i++ {
		result := <-results
		if !result.ok {
			return result.err
		}
	}

	return nil
}

func (m *mtee) run() error {
	for {
		err := m.tee()
		if err != nil {
			return err
		}
	}

}

// Run runs mtee from a file string and mode string
func Run(files []string, mode bool) error {
	m := &mtee{}

	err := m.init(files, mode)
	if err != nil {
		return err
	}

	// if we opened a file, we need to close it
	for _, f := range m.out {
		defer f.Close()
	}

	return (m.run())
}
