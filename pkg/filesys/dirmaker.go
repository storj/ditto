package filesys

import "os"

func DirMaker() FsMkdir {
	return osMkdir(os.MkdirAll)
}
