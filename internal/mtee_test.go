package mtee

import (
	"bufio"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func failNow(t *testing.T, err error, s string) {
	if err != nil {
		assert.FailNow(t, s)
	}
}

var testTeeText = "test tee text"

func TestTee(t *testing.T) {
	var afs = afero.NewMemMapFs()
	f, err := afs.Create("test.txt")
	failNow(t, err, "failed to create file")

	_, err = f.Write([]byte(testTeeText))
	failNow(t, err, "failed to write to file")
	f.Close()

	s := bufio.NewScanner(f)
	c := make(chan teeResult, 1)

	outF := newOut(f, defaultBufSize)

	err = tee(s, []*out{outF}, c)
	assert.NoError(t, err, "tee failed")

}
