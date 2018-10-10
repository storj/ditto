package filesys

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestDirChecker(t *testing.T) {
	dchecker := MockDirChecker(func(string) (bool, error) {
		return true, nil
	})

	isDir, err := dchecker.CheckIfDir("path")
	assert.NoError(t, err)
	assert.Equal(t, true, isDir)
}

func TestNewDirChecker(t *testing.T) {
	dchecker := NewDirChecker()
	assert.NotNil(t, dchecker)
}