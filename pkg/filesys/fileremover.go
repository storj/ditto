package filesys

import "os"

func FileRemover() (FsRemove) {
	return osRemove(os.Remove)
}
