package mtee

import (
	"bufio"
	"fmt"
	"os"
)

// mtee represents the application mtee and its configurations
type mtee struct {
	append  bool
	file    *file
	scanner *bufio.Scanner
}

func (m *mtee) init(fstr string, modeAppend bool) error {
	switch fstr {
	case "":
		m.file = nil
	default:
		err := m.setFile(fstr, modeAppend)
		if err != nil {
			return err
		}
	}

	m.scanner = bufio.NewScanner(os.Stdin)
	return nil

}

func (m *mtee) setFile(fstr string, modeAppend bool) error {
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

	m.file = &file{f}
	return nil
}

func (m *mtee) tee() error {
	// read from in
	scanner := bufio.NewScanner(m.in)
	scan = s.Scan()
	text := fmt.Sprintf("%s\n", scan.


	if err != nil {
		return fmt.Errorf("failed to read from in: %w", err)
	}

	text := fmt.Sprintf("%s\n", string(b))
	fmt.Println("text is " + text)

	// write to all outs
	for _, v := range m.out {
		_, err := v.Write([]byte(text))
		if err != nil {
			return err
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
func Run(fs string, mode bool) error {
	m := &mtee{}

	err := m.init(fs, mode)
	if err != nil {
		return err
	}

	// if we opened a file, we need to close it
	if m.file != nil {
		defer m.file.Close()
	}

	return (m.run())
}
