package complete

import (
	"path/filepath"
)

type FileCompleter struct {
	Root string
}

func NewFileCompleter(s string) Completer {
	return &FileCompleter{Root: s}
}

func (f *FileCompleter) Complete(s string) ([]string, error) {
	n, err := filepath.Glob(filepath.Join(f.Root, s+"*"))
	if err != nil || len(n) == 0 {
		return n, err
	}
	files := make([]string, len(n))
	for i := range n {
		files[i] = filepath.Base(n[i])
	}
	return files, err
}
