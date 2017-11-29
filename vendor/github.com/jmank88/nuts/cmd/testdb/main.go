/*
Command testdb recursively walks the directory given as the first (and only) argument, and copies paths from .txt files
into .db BoltDB database files by the same name.
*/
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/jmank88/nuts"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("too few arguments")
	}
	dir := os.Args[1]
	if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".txt") {
			if err := copyToBolt(path); err != nil {
				log.Printf("failed to copy %q to boltdb: %s\n", path, err)
			}
		}
		return nil
	}); err != nil {
		log.Fatal("failed to walk files", err)
	}
}

const batchSize = 10000

func copyToBolt(txt string) (err error) {
	path := strings.TrimSuffix(txt, ".txt") + ".db"

	db, err := bolt.Open(path, 0666, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err2 := db.Close(); err2 != nil && err == nil {
			err = err2
		}
	}()

	f, err := os.Open(txt)
	if err != nil {
		return err
	}
	defer f.Close()

	// Use default, ScanLines
	s := bufio.NewScanner(f)
	done := false
	for !done {
		if err := db.Update(func(tx *bolt.Tx) error {
			b, err := tx.CreateBucketIfNotExists([]byte("paths"))
			if err != nil {
				return err
			}

			for batch := 0; batch < batchSize; batch++ {
				if !s.Scan() {
					done = true
					return s.Err()
				}
				path := s.Bytes()
				if k, _ := nuts.SeekPathConflict(b.Cursor(), path); k != nil {
					return fmt.Errorf("path %q conflicts with existing %q", string(path), string(k))
				}
				if err := b.Put(path, []byte{}); err != nil {
					return err
				}
			}
			return nil
		}); err != nil {
			return err
		}
	}
	return nil
}
