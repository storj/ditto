package filesys

import "os"

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

type osCreate func(string) (*os.File, error)

func (c osCreate) Create(lpath string) (File, error) {
	return c(lpath)
}
//------------------------------------------------

type osOpenFile func(string, int, os.FileMode) (*os.File, error)

func (o osOpenFile) Create(lpath string) (File, error) {
	return o(lpath, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
}
//------------------------------------------------

type osRemove func(string) (error)

func (r osRemove) Remove (lpath string) (error) {
	return r(lpath)
}
//------------------------------------------------

type osMkdir func(string, os.FileMode) (error)

func (mk osMkdir) Mkdir(lpath string) (error) {
	return mk(lpath, os.ModePerm)
}
//------------------------------------------------

type osDirChecker func(string) (bool, error)

func(d osDirChecker) CheckIfDir(lpath string) (bool, error) {
	return d(lpath)
}
//------------------------------------------------

type osReadDir func(string) ([]os.FileInfo, error)

func (r osReadDir) ReadDir(lpath string) ([]os.FileInfo, error) {
	return r(lpath)
}
//------------------------------------------------
