package filesys

import (
	"os"
	"io/ioutil"
)

type DirReader interface {
	ReadDir(lpath string) (Dir, error)
}

type dirReader struct {
	FsStat
	FsReadDir
}

func NewDirReader() DirReader {
	return &dirReader{osStat(os.Stat), osReadDir(ioutil.ReadDir)}
}

func (d *dirReader) ReadDir(lpath string) (Dir, error) {
	info, err := d.Stat(lpath)
	if err != nil {
		return nil, err
	}

	items, err := d.FsReadDir.ReadDir(lpath)
	if err != nil {
		return nil, err
	}

	var files []os.FileInfo
	var dirs []os.FileInfo

	for i := range items {
		item := items[i]
		if item.IsDir() {
			dirs = append(dirs, item)
			continue
		}

		files = append(files, item)
	}

	return &dir{lpath,info, files, dirs}, nil
}
