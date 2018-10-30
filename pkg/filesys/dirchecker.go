package filesys

import "os"

func DirChecker() FsCheckDir {
	return osDirChecker(CheckIfDir)
}

func CheckIfDir(lpath string) (isDir bool, err error) {
	fi, err := os.Stat(lpath)
	if err != nil {
		return
	}

	return fi.IsDir(), err
}
