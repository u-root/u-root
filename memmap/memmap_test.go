package memmap

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"
	"unicode"
)

func TestSysMemmap(t *testing.T) {
	var m string
	mm, err := Ranges()
	if err != nil {
		t.Fatalf("MemMap: got %v, want nil", err)
	}
	for _, v := range mm {
		m = m + v.String()
	}
	// We could just kick off a cat here but I hate to depend on running programs and it would
	// require a shell for globbing. ewww.
	var fs string
	err = filepath.Walk("/sys/firmware/memmap", func(name string, fi os.FileInfo, err error) error {
		// be dumb. If it's one of our names, read it and throw it onto the end.
		// this should never happen, unless they add weird non-directory things in the future.
		if !fi.IsDir() || !unicode.IsDigit(rune(fi.Name()[0])) {
			return nil
		}
		for _, v := range []string{"start", "end", "type"} {
			s, err := ioutil.ReadFile(path.Join(name, v))
			if err != nil {
				return err
			}
			fs = fs + string(s)
		}
		return err

	})
	if len(fs) != len(m) {
		t.Errorf("File system string len != string length of memmap: fs is %d, memmap is %d", len(fs), len(m))
	}
	if fs != m {
		t.Fatalf("fs is %s and m.String is %s: they're different and need to be the same", fs, m)
	}
	t.Logf("fs is %s and m.String is %s: OK", fs, m)

}
