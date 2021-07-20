package main

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
)

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

func dirTree(out io.Writer, path string, printFiles bool) (err error) {
	err = getDir(out, path, "", printFiles)
	return
}

func getDir(out io.Writer, path string, startLine string, printFiles bool) (err error) {
	files, err := os.ReadDir(path)

	var correctedFiles []fs.DirEntry
	for _, file := range files {
		if (!printFiles && file.IsDir()) || printFiles {
			correctedFiles = append(correctedFiles, file)
		}
	}

	var size int = len(correctedFiles)
	isLast := false

	for idx, corrFile := range correctedFiles {
		_ = idx
		if idx == size-1 {
			isLast = true
		}

		var levelSymbol string
		if printFiles || corrFile.IsDir() {
			if isLast {
				levelSymbol = "└"
			} else {
				levelSymbol = "├"
			}

			info, _ := corrFile.Info()
			formattedLine := getFormattedOutputLine(startLine, levelSymbol, corrFile.Name(), info)
			out.Write([]byte(formattedLine))
			// fmt.Println(formattedLine)
		}

		if corrFile.IsDir() {
			var newStartLine string
			if isLast {
				newStartLine = startLine + "	"
			} else {
				newStartLine = startLine + "│	"
			}

			nextPath := path + string(filepath.Separator) + corrFile.Name()
			getDir(out, nextPath, newStartLine, printFiles)
		}
	}

	return err
}

func getFormattedOutputLine(startLine string, graphicSymbol string, name string, info os.FileInfo) (result string) {

	var size string
	if info.Size() == 0 {
		size = "empty"
	} else {
		size = strconv.FormatInt(info.Size(), 10) + "b"
	}

	if !info.IsDir() {
		result = startLine + graphicSymbol + "───" + name + " " + "(" + size + ")" + "\n"
	} else {
		result = startLine + graphicSymbol + "───" + name + "\n"
	}

	return
}
