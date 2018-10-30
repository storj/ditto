package filesys

import "os"

//ForceFileCreator wraps default os creation that truncates the file if one exists
func ForceFileCreator() (FsCreate) {
	return osCreate(os.Create)
}

//FileCreator default behavior is to create only of none exists
func FileCreator() (FsCreate) {
	return osOpenFile(os.OpenFile)
}
