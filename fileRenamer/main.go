package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type matchResult struct {
	base  string
	index int
	ext   string
}

func main() {
	toRename := make(map[string][]string)
	var walkDir string
	var dryrun bool
	flag.StringVar(&walkDir, "walkdir", "sample", "the directory where you save all the files that need to rename")
	flag.BoolVar(&dryrun, "dry", true, "whether or not should the tool change the names")
	flag.Parse()

	err := filepath.Walk(walkDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("failed access path: %s ", path)
		}
		if info.IsDir() {
			return nil
		}
		if res, err := match(info.Name()); err == nil {
			curDir := filepath.Dir(path)
			key := filepath.Join(curDir, fmt.Sprintf("%s.%s", res.base, res.ext))
			toRename[key] = append(toRename[key], info.Name())
		}
		return nil
	})
	if err != nil {
		log.Fatal("Failed to walk through the dir", walkDir)
	}

	for key, files := range toRename {
		dir := filepath.Dir(key)
		count := len(files)
		sort.Strings(files)
		for i, filename := range files {

			res, err := match(filename)
			if err != nil {
				log.Fatal("no match", err)
			}

			newName := fmt.Sprintf("%s-%d of %d.%s", res.base, (i + 1), count, res.ext)
			newPath := filepath.Join(dir, newName)
			oldPath := filepath.Join(dir, filename)

			fmt.Printf("mv %s => %s\n", oldPath, newPath)
			if !dryrun {
				err = os.Rename(oldPath, newPath)
				if err != nil {
					log.Fatal("failed to rename", err)
				}
			}
		}
	}
}

//match returns the new file name, or an error if the file name did't match our pattern.
func match(filename string) (*matchResult, error) {
	pieces := strings.Split(filename, ".")
	ext := pieces[len(pieces)-1]
	tmp := strings.Join(pieces[0:len(pieces)-1], ".")
	pieces = strings.Split(tmp, "_")
	name := strings.Join(pieces[0:len(pieces)-1], "_")
	number, err := strconv.Atoi(pieces[len(pieces)-1])
	if err != nil {
		return nil, fmt.Errorf("%s did not match our pattern", filename)
	}
	return &matchResult{strings.Title(name), number, ext}, nil
}
