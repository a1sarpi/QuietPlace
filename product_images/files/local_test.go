package files

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupLocal(t *testing.T) (*Local, string, func()) {
	// Create a temporary directory using os.MkdirTemp
	dir, err := os.MkdirTemp("", "files")
	if err != nil {
		t.Fatal(err)
	}

	// Assuming 0 is a default value for the int argument
	l, err := NewLocal(dir, 10000) // Pass an int (e.g., 0) and the directory path
	if err != nil {
		t.Fatal(err)
	}

	return l, dir, func() {
		// Cleanup function
		os.RemoveAll(dir)
	}
}

func TestSavesContentsOfReader(t *testing.T) {
	savePath := "/1/test.png"
	fileContents := "Hello World"
	l, dir, cleanup := setupLocal(t)
	defer cleanup()

	err := l.Save(savePath, bytes.NewBuffer([]byte(fileContents)))
	assert.NoError(t, err)

	// Check the file has been correctly written
	f, err := os.Open(filepath.Join(dir, savePath))
	assert.NoError(t, err)

	// Check the contents of the file
	d, err := ioutil.ReadAll(f)
	assert.NoError(t, err)
	assert.Equal(t, fileContents, string(d))
}

func TestGetsContentsAndWritesToWriter(t *testing.T) {
	savePath := "/1/test.png"
	fileContents := "Hello World"
	l, _, cleanup := setupLocal(t)
	defer cleanup()

	// Save a file
	err := l.Save(savePath, bytes.NewBuffer([]byte(fileContents)))
	assert.NoError(t, err)

	// Read the file back
	r, err := l.Get(savePath)
	assert.NoError(t, err)
	defer r.Close()

	// Read the full contents of the reader
	d, err := io.ReadAll(r)
	assert.Equal(t, fileContents, string(d))
}
