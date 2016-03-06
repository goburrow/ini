// Package ini provides helpers for reading ini file
package ini

import (
	"bufio"
	"bytes"
	"io"
	"os"
)

// WalkFunc is the type of the function called for each key visited by Walk.
type WalkFunc func(section, name, value []byte) error

// Walk walks all items from r, calling walkFn for each item.
func Walk(r io.Reader, walkFn WalkFunc) error {
	var section []byte
	sep := []byte{'='}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 || line[0] == '#' || line[0] == ';' {
			// Skip comments
			continue
		}
		if line[0] == '[' && line[len(line)-1] == ']' {
			// Section
			if cap(section) < len(line)-2 {
				section = make([]byte, len(line)-2)
			} else {
				section = section[:len(line)-2]
			}
			copy(section, line[1:len(line)-1])
			continue
		}

		key := bytes.SplitN(line, sep, 2)
		if len(key) != 2 {
			continue
		}
		name := bytes.TrimSpace(key[0])
		value := bytes.TrimSpace(key[1])
		err := walkFn(section, name, value)
		if err != nil {
			return err
		}
	}
	return nil
}

// WalkFile is similar to Walk but takes file path for the reader.
func WalkFile(path string, walkFn WalkFunc) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return Walk(f, walkFn)
}