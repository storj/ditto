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

type FsCreate interface {
	Create(lpath string) (File, error)
}

type FsRemove interface {
	Remove(lpath string) (error)
}

type FsMkdir interface {
	Mkdir(lpath string) (error)
}

type FsCheckDir interface {
	CheckIfDir(string) (bool, error)
}

type FsReadDir interface {
	ReadDir(lpath string) ([]os.FileInfo, error)
}

type FileSystem interface {
	FsOpen
	FsStat
	FsCreate
	FsRemove
	FsReadDir
}

type File interface {
	io.Closer
	io.Reader
	io.ReaderAt
	io.Seeker
	io.Writer
	io.WriterAt
	Stat() (os.FileInfo, error)
}

type Dir interface {
	os.FileInfo
	Path() string
	Files() []os.FileInfo
	Dirs() []os.FileInfo
}


