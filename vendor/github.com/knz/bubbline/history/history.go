package history

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"os"
)

const cookie = "_HiStOrY_V2_"

// LoadHistory loadsa a history from a file loaded from the specified
// path. The file must be in the same format as used by libedit.
func LoadHistory(fileName string) ([]string, error) {
	f, err := os.Open(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer func() { _ = f.Close() }()
	return loadHistoryFromFile(f)
}

func loadHistoryFromFile(f io.Reader) ([]string, error) {
	var buf [len(cookie) + 1]byte
	n, err := f.Read(buf[:])
	if err == io.EOF {
		// empty file.
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	sl := buf[:n]
	if !bytes.Equal(sl, []byte(cookie+"\n")) {
		// Cookie not recognized. No-op.
		return nil, nil
	}
	// Read the remainder of the file.
	contents, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	if len(contents) > 0 && contents[len(contents)-1] == '\n' {
		contents = contents[:len(contents)-1]
	}
	lines := bytes.Split(contents, []byte("\n"))
	hist := make([]string, 0, len(lines))
	for _, line := range lines {
		// Unescape octal codes.
		resultEnd := 0
		for c := 0; c < len(line); c++ {
			foundEscape := false
			if line[c] == '\\' && c+3 < len(line) {
				foundEscape = true
				var b byte
				for i := 1; i <= 3; i++ {
					digit := line[c+i]
					if digit < '0' || digit > '7' {
						foundEscape = false
						break
					}
					b = (b << 3) | (digit - '0')
				}
				if foundEscape {
					line[resultEnd] = b
					c += 3
				}
			}
			if !foundEscape {
				line[resultEnd] = line[c]
			}
			resultEnd++
		}
		result := line[:resultEnd]
		hist = append(hist, string(result))
	}
	return hist, nil
}

// SaveHistory saves a history to the specified file.
// The file will be written in the same format as used by libedit.
func SaveHistory(h []string, fileName string) (retErr error) {
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer func() {
		closeErr := f.Close()
		if retErr == nil {
			retErr = closeErr
		}
	}()

	return saveHistoryToFile(h, f)
}

func saveHistoryToFile(h []string, f io.Writer) error {
	w := bufio.NewWriter(f)
	_, err := w.Write([]byte(cookie + "\n"))
	if err != nil {
		return err
	}
	for _, entry := range h {
		var buf bytes.Buffer
		for c := 0; c < len(entry); c++ {
			if b := entry[c]; b == ' ' || b == '\t' || b == '\n' || b == '\\' {
				buf.WriteByte('\\')
				buf.WriteByte((b>>6)&7 + '0')
				buf.WriteByte((b>>3)&7 + '0')
				buf.WriteByte((b>>0)&7 + '0')
			} else {
				buf.WriteByte(b)
			}
		}
		buf.WriteByte('\n')
		_, err := w.Write(buf.Bytes())
		if err != nil {
			return err
		}
	}
	return w.Flush()
}
