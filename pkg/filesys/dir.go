package filesys

import "os"

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
