package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

func getFileSizeString(file os.FileInfo) string {
	if file.Size() == 0 {
		return "(empty)"
	}
	return "(" + strconv.FormatInt(file.Size(), 10) + "b)"
}

func dirTreePrefix(out io.Writer, path string, printFiles bool,
	prefix string) error {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return fmt.Errorf("error reading dir %s: %s",
			path, err)
	}

	dirsOnly := make([]os.FileInfo, 0)
	if !printFiles {
		for _, file := range files {
			if file.IsDir() {
				dirsOnly = append(dirsOnly, file)
			}
		}
		files = dirsOnly
	}
	for i, file := range files {
		if file.IsDir() {
			var newPrefix string
			newPath := filepath.Join(path, file.Name())

			if i == len(files)-1 {
				fmt.Fprintf(out, "%s└───%s\n", prefix, file.Name())
				newPrefix = prefix + `	`
			} else {
				fmt.Fprintf(out, "%s├───%s\n", prefix, file.Name())
				newPrefix = prefix + `│	`
			}

			dirTreePrefix(out, newPath, printFiles, newPrefix)
		} else if printFiles {
			if i == len(files)-1 {
				fmt.Fprintf(out, "%s└───%s %s\n",
					prefix, file.Name(), getFileSizeString(file))
			} else {
				fmt.Fprintf(out, "%s├───%s %s\n",
					prefix, file.Name(), getFileSizeString(file))
			}
		}
	}

	return nil
}
func dirTree(out io.Writer, path string, printFiles bool) error {
	return dirTreePrefix(out, path, printFiles, "")
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)

	if err != nil {
		panic(err.Error())
	}
}
