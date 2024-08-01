package main

import (
	"io/fs"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

type File struct {
	Seq           int
	AudioFilePath string
	TextFilePath  string
}

func collectFiles(root string) ([]File, error) {
	r := regexp.MustCompile(`^([0-9]+)_.*\.wav`)

	files := []File{}

	err := filepath.WalkDir(root, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		name := info.Name()

		matches := r.FindAllStringSubmatch(name, -1)
		if len(matches) < 1 || len(matches[0]) != 2 {
			return nil
		}

		seqRaw := matches[0][1]
		seq, err := strconv.Atoi(seqRaw)
		if err != nil {
			return nil
		}

		textFilePath := replaceExtension(path, ".txt")

		files = append(files, File{
			Seq:           seq,
			AudioFilePath: path,
			TextFilePath:  textFilePath,
		})
		return nil
	})
	return files, err
}

func sortFiles(files []File) {
	slices.SortFunc(files, func(a File, b File) int {
		if a.Seq > b.Seq {
			return 1
		} else if a.Seq < b.Seq {
			return -1
		}
		return 0
	})
}

func replaceExtension(path, extension string) string {
	ext := filepath.Ext(path)
	bare := strings.TrimSuffix(path, ext)
	return bare + extension
}
