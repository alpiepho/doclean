package walk_example

import (
	"io/fs"
	"path/filepath"
)

func walk(s string, d fs.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if !d.IsDir() {
		println(s)
	}
	return nil
}

func main() {
	filepath.Walk("..", walk)
}
