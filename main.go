package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	outputName     string
	inputName      string
	unpackMode     bool
	packMode       bool
	debug          bool
	version        bool
	singleSeparate bool
	importNew      bool
	withGraph      bool
	dirs           = []string{dirSnippet, dirMenu, dirLang, dirTable, dirParam, dirData, dirPage, dirCon}
)

func main() {
	flag.BoolVar(&unpackMode, "unpack", false, "-u, unpacking mode")
	flag.StringVar(&inputName, "input", ".", "-i, path for input files, filename for pack and dirname/ (slashed) for unpack")
	flag.StringVar(&outputName, "output", "output", "-o, output filename for JSON if input file name not pointed")

	// shorthand
	flag.StringVar(&outputName, "o", "output", "-output")
	flag.StringVar(&inputName, "i", ".", "input")
	flag.BoolVar(&unpackMode, "u", false, "-unpack")
	flag.BoolVar(&version, "v", false, "-version")
	flag.BoolVar(&debug, "d", false, "debug")
	flag.BoolVar(&withGraph, "g", false, "make graphical structure in dot-file")
	flag.Parse()

	args := flag.Args()
	if argsCount := len(args); argsCount == 0 {
		// without args run gui
		SimpleGui()
	} else {
		if argsCount == 1 && version {
			fmt.Println(currentTitle)
		} else if argsCount >= 1 {
			if version {
				fmt.Println(currentTitle)
			}
			name := args[0]
			if stats, err := os.Stat(name); err == nil {
				if stats.IsDir() {
					packMode = true
					inputName = name
				} else if stats.Mode() != os.ModeType {
					unpackMode = true
					inputName = name
				}
			} else {
				fmt.Println(err)
			}
		}
		checkOutput()
	}
	if unpackMode {
		unpackJSON(inputName)
	} else if packMode {
		packJSON(inputName)
	}
}

func checkOutput() {
	if outputName == "output" && inputName != "." { // we have only inputname
		parts := strings.Split(inputName, separator)
		pLen := len(parts)
		outputName = parts[pLen-1]
		if unpackMode {
			ext := filepath.Ext(outputName)
			outputName = outputName[:len(outputName)-len(ext)]
			outputName = outputName + separator
		} else {
			if strings.HasSuffix(inputName, separator) {
				outputName = parts[pLen-2]
			}
		}
	}

	if unpackMode {
		if stats, err := os.Stat(inputName); inputName == "." || stats.IsDir() || err != nil {
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(helpMsg)
			return
		}
		if !strings.HasSuffix(outputName, separator) {
			outputName = outputName + separator
		}
		if debug {
			fmt.Println("output dir name:", outputName)
		}
	} else {
		if !strings.HasSuffix(inputName, separator) {
			fmt.Println(helpMsg)
			return
		}
	}
}

func createDir(path string) {
	if _, e := os.Stat(path); os.IsNotExist(e) {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			fmt.Println("error create dir", err)
		}
	}
}

func writeFileString(filename, content string) {
	outFile, err := os.Create(filepath.Join(outputName, filename))
	if err != nil {
		// fmt.Println("error write file:", err)
		return
	}
	defer outFile.Close()
	if _, err := outFile.WriteString(content); err != nil {
		fmt.Println(err)
		return
	}
}

func stringInSlice(arr []string, val string) bool {
	for _, v := range arr {
		if v == val {
			return true
		}
	}
	return false
}
