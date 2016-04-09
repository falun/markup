package buildindex

import ()

type FileEntry struct {
	Root string
	Name string
	Dir  bool
}

type ByFilepath []FileEntry

func (f ByFilepath) Len() int {
	return len(f)
}

func (f ByFilepath) Less(i, j int) bool {
	if f[i].Root > f[j].Root {
		return false
	}

	if (f[i].Root == f[j].Root) && (f[i].Name >= f[j].Name) {
		return false
	}

	return true
}

func (f ByFilepath) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}
