package filesys

import (
	"os"
	"io"
)

type FsOpen interface {
	Open(lpath string) (File, error)
}

type FsStat interface {
	Stat(lpath string) (os.FileInfo, error)
}

type FsReadDir interface {
	ReadDir(lpath string) ([]os.FileInfo, error)
}

type FileSystem interface {
	FsOpen
	FsStat
}

type File interface {
	io.Closer
	io.Reader
	io.ReaderAt
	io.Seeker
	Stat() (os.FileInfo, error)
}

type Dir interface {
	os.FileInfo
	Path() string
	Files() []os.FileInfo
	Dirs() []os.FileInfo
}

//------------------------------------------------
type osOpen func(string) (*os.File, error)

func (f osOpen) Open(lpath string) (File, error) {
	return f(lpath)
}
//------------------------------------------------

type osStat func(string) (os.FileInfo, error)

func(s osStat) Stat(lpath string) (os.FileInfo, error) {
	return s(lpath)
}
//------------------------------------------------

type osReadDir func(string) ([]os.FileInfo, error)

func (r osReadDir) ReadDir(lpath string) ([]os.FileInfo, error) {
	return r(lpath)
}
//------------------------------------------------

type dir struct {
	lpath string
	os.FileInfo
	files []os.FileInfo
	dirs []os.FileInfo
}

func (d *dir) Path() string {
	return d.lpath
}

func (d *dir) Files() []os.FileInfo {
	return d.files
}

func (d *dir) Dirs() []os.FileInfo {
	return d.dirs
}
//------------------------------------------------