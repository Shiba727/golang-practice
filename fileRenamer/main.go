package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

type matchResult struct {
	base  string
	index int
	ext   string
}

var (
	re            = regexp.MustCompile("^(.+?) ([0-9]{4}) [(]([0-9]+) of ([0-9]+)[)][.](.+?)$")
	replaceString = "$2 - $1 - $3 of $4.$5"
	walkDir       string
	dryrun        bool
)

func init() {
	flag.StringVar(&walkDir, "walkdir", "sample", "the directory where you save all the files that need to rename")
	flag.BoolVar(&dryrun, "dry", true, "whether or not should the tool change the names")
}

func main() {
	flag.Parse()

	// toRename := make(map[string][]string)
	var toRename []string

	err := filepath.Walk(walkDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("failed access path: %s ", path)
		}
		if info.IsDir() {
			return nil
		}
		if _, err := match(info.Name()); err == nil {
			toRename = append(toRename, path)
		}
		return nil
	})
	if err != nil {
		log.Fatal("Failed to walk through the dir", walkDir)
	}

	for _, path := range toRename {
		dir := filepath.Dir(path)
		filename := filepath.Base(path)
		newFilename, _ := match(filename)
		newPath := filepath.Join(dir, newFilename)

		fmt.Printf("mv %s => %s\n", path, newPath)
		if !dryrun {
			err = os.Rename(path, newPath)
			if err != nil {
				log.Fatal("failed to rename", err)
			}
		}
	}
}

//match returns the new file name, or an error if the file name did't match our pattern.
func match(filename string) (string, error) {
	if !re.MatchString(filename) {
		return "", fmt.Errorf("%s did not match our pattern", filename)
	}
	return re.ReplaceAllString(filename, replaceString), nil
}
