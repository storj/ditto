package filesys

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"errors"
)

func TestFileReader(t *testing.T) {
	testError := errors.New("test error")

	cases := []struct{
		testName string
		testFunc func(t *testing.T)
	}{
		{
			"No error",
			func(t *testing.T) {
				freader := &baseFileReader{
					MockFsOpen(func(lpath string) (File, error) {
						return &MockFile{}, nil
					}),
				}

				file, err := freader.ReadFile("path")
				assert.NoError(t, err)
				assert.NotNil(t, file)
			},
		},
		{
			"Error open",
			func(t *testing.T) {
				freader := &baseFileReader{
					MockFsOpen(func(lpath string) (File, error) {
						return nil, testError
					}),
				}

				file, err := freader.ReadFile("path")
				assert.Nil(t, file)
				assert.Error(t, err)
				assert.Equal(t, testError, err)
			},
		},
		{

			"Error stat",
			func(t *testing.T) {
				freader := &baseFileReader{
					MockFsOpen(func(lpath string) (File, error) {
						return &MockFile{testError}, nil
					}),
				}

				file, err := freader.ReadFile("path")
				assert.Nil(t, file)
				assert.Error(t, err)
				assert.Equal(t, testError, err)
			},
		},
	}

	for _, c := range cases {
		t.Run(c.testName, c.testFunc)
	}
}

func TestHashReader(t *testing.T) {
	cases := []struct{
		testName string
		testFunc func(t *testing.T)
	}{
		{
			"No error",
			func(t *testing.T) {
				hreader := &hashFileReader{
					MockFileReader(func(lpath string) (FReader, error) {
						return &MockFReader{}, nil
					}),
				}

				rd, err := hreader.ReadFileH("path")
				assert.NoError(t, err)
				assert.NotNil(t, rd)
			},
		},
	}

	for _, c := range cases {
		t.Run(c.testName, c.testFunc)
	}
}
