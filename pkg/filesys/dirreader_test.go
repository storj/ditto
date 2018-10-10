package filesys

import (
	"testing"
	"os"
	"github.com/stretchr/testify/assert"
	"github.com/pkg/errors"
)

func TestNewDirReader(t *testing.T) {
	dreader := NewDirReader()
	assert.NotNil(t, dreader)
}

func TestDirReader(t *testing.T) {
	testError := errors.New("test error")

	cases := []struct{
		testName string
		testFunc func(t *testing.T)
	}{
		{
			"No error",
			func(t *testing.T) {
				dreader := dirReader{
					MockFsStat(func(s string) (os.FileInfo, error) {
						return &MockFileInfo{}, nil
					}),
					MockFsReadDir(func(s string) ([]os.FileInfo, error) {
						list := []os.FileInfo{&MockFileInfo{isDir:true}, &MockFileInfo{}, &MockFileInfo{}}
						return list, nil
					}),
				}

				_, err := dreader.ReadDir("path")
				assert.NoError(t, err)
			},
		},
		{
			"Stat error",
			func(t *testing.T) {
				dreader := dirReader{
					MockFsStat(func(s string) (os.FileInfo, error) {
						return nil, testError
					}),
					MockFsReadDir(func(s string) ([]os.FileInfo, error) {
						return nil, nil
					}),
				}

				_, err := dreader.ReadDir("path")
				assert.Error(t, err)
			},
		},
		{
			"Read dir error",
			func(t *testing.T) {
				dreader := dirReader{
					MockFsStat(func(s string) (os.FileInfo, error) {
						return &MockFileInfo{}, nil
					}),
					MockFsReadDir(func(s string) ([]os.FileInfo, error) {
						return nil, testError
					}),
				}

				_, err := dreader.ReadDir("path")
				assert.Error(t, err)
			},
		},
	}

	for _, c := range cases {
		t.Run(c.testName, c.testFunc)
	}
}
