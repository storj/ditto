package filesys

import "os"

type DirChecker interface {
	CheckIfDir(string) (bool, error)
}

type osDirChecker func(string) (bool, error)

func NewDirChecker() DirChecker {
	return osDirChecker(CheckIfDir)
}

func(d osDirChecker) CheckIfDir(lpath string) (bool, error) {
	return d(lpath)
}

func CheckIfDir(lpath string) (isDir bool, err error) {
	fi, err := os.Stat(lpath)
	if err != nil {
		return
	}

	return fi.IsDir(), err
}
