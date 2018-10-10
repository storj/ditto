package filesys

import (
	"os"
	"time"
	"io"
)

//------------------------------------------------------------------------------------
type MockFileInfo struct {
	isDir bool
}

func (fi *MockFileInfo) Name() string {
	return ""
}

func (fi *MockFileInfo) Size() int64 {
	return 0
}

func (fi *MockFileInfo) Mode() os.FileMode {
	return os.FileMode(0)
}

func (fi *MockFileInfo) ModTime() time.Time {
	return time.Now()
}

func (fi *MockFileInfo) IsDir() bool {
	return fi.isDir
}

func (fi *MockFileInfo) Sys() interface{} {
	return nil
}
//-----------------------------------------------------------------------------------

type MockFile struct {
	Err error
}

func (f *MockFile) Close() error {
	return f.Err
}

func (f *MockFile) Read([]byte) (int, error) {
	return 0, f.Err
}

func (f *MockFile) ReadAt([]byte, int64) (int, error) {
	return 0, f.Err
}

func (f *MockFile) Seek(offset int64, whence int) (int64, error) {
	return 0, f.Err
}

func (f *MockFile) Stat() (os.FileInfo, error) {
	if f.Err != nil {
		return nil, f.Err
	}

	return &MockFileInfo{false}, nil
}
//------------------------------------------------------------------------------------

type MockFsOpen func(string) (File, error)

func (fo MockFsOpen) Open(lpath string) (File, error) {
	return fo(lpath)
}
//------------------------------------------------------------------------------------

type MockFsStat = osStat
type MockFsReadDir = osReadDir
type MockDirChecker = osDirChecker
//------------------------------------------------------------------------------------

type MockFileReader func(string) (FReader, error)

func (r MockFileReader) ReadFile(lpath string) (FReader, error) {
	return r(lpath)
}
//------------------------------------------------------------------------------------

type MockFReader struct {}

func (r *MockFReader) Read(b []byte) (int, error) {
	return 0, nil
}

func (r *MockFReader) FileInfo() os.FileInfo {
	return &MockFileInfo{}
}

func (r *MockFReader) Reader() (io.Reader) {
	return r
}