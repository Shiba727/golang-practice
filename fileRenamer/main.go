package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	var toRename []string
	fileCount := 0
	dir := "sample"

	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("failed access path: %s ", path)
		}
		if info.IsDir() {
			return nil
		}
		if _, err := match(info.Name(), 0); err == nil {
			fileCount++
			toRename = append(toRename, info.Name())
		}
		return nil
	})
	if err != nil {
		log.Fatal("Failed to walk through the dir", dir)
	}

	for _, originName := range toRename {
		newName, err := match(originName, fileCount)
		if err != nil {
			log.Fatal("no match ", err)
		}
		newPath := filepath.Join(dir, newName)
		oldPath := filepath.Join(dir, originName)

		err = os.Rename(oldPath, newPath)
		if err != nil {
			log.Fatal("failed to rename ", err)
		}
		fmt.Printf("mv %s => %s\n", oldPath, newPath)
	}
}

//match returns the new file name, or an error if the file name did't match our pattern.
func match(filename string, total int) (string, error) {
	pieces := strings.Split(filename, ".")
	ext := pieces[len(pieces)-1]
	tmp := strings.Join(pieces[0:len(pieces)-1], ".")
	pieces = strings.Split(tmp, "_")
	name := strings.Join(pieces[0:len(pieces)-1], "_")
	number, err := strconv.Atoi(pieces[len(pieces)-1])
	if err != nil {
		return "", fmt.Errorf("%s did not match our pattern", filename)
	}

	return fmt.Sprintf("%s - %d of %d.%s", strings.Title(name), number, total, ext), nil
}
