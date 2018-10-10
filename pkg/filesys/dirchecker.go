package filesys

type DirChecker interface {
	CheckIfDir(string) (bool, error)
}

type BDirChecker func(string) (bool, error)

func(d BDirChecker) CheckIfDir(lpath string) (bool, error) {
	return d(lpath)
}
