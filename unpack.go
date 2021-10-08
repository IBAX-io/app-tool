package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/spkg/bom"
)

func unpackJSON(filename string) {
	bs, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	bs = bom.Clean(bs)
	test := testFormatStruct{}
	if err := json.Unmarshal(bs, &test); err != nil {
		fmt.Println("unmarshal file test:", err)
		return
	}
	if test.len() == 1 && len(test.Name) > 0 {
		importNew = true
		file := dataFile{}
		if err := json.Unmarshal(bs, &file); err != nil {
			fmt.Println("unmarshal file error_2:", err)
			return
		}
		unpackDataFile(file.Data)

	} else {
		file := importFile{}
		if err := json.Unmarshal(bs, &file); err != nil {
			fmt.Println("unmarshal file error_1:", err)
			return
		}
		unpackStruct(file.Contracts, eSIM, dirCon)
		unpackStruct(file.Menus, ePTL, dirMenu)
		unpackStruct(file.Blocks, ePTL, dirBlock)
		unpackStruct(file.Pages, ePTL, dirPage)
		unpackStruct(file.Tables, eJSON, dirTable)
		unpackStruct(file.Parameters, eCSV, dirParam)
		unpackStruct(file.Languages, eJSON, dirLang)

		if len(file.Data) > 0 {
			createDir(filepath.Join(outputName, dirData))
			for _, item := range file.Data {
				name := item.Table + eJSON
				name = filepath.Join(dirData, name)
				result, _ := json.MarshalIndent(item, "", "    ")
				writeFileString(name, string(result))
			}
		}
	}
	writeConfig(bs)
	if abs, err := filepath.Abs(outputName); err == nil {
		abspath := filepath.Join(abs, structFileName)
		createGraph(abspath)
	}
}

func unpackStruct(items []commonStruct, tail, dir string) {
	if len(items) > 0 {
		createDir(filepath.Join(outputName, dir))
		for _, item := range items {
			value := item.Value
			fmt.Println(value)
			switch dir {
			case dirTable:
				value = item.Columns
			case dirLang:
				value = item.Trans
			}
			name := item.Name
			if len(item.Table) > 0 {
				name = item.Table
			}
			fullName := name + tail
			fullName = filepath.Join(dir, fullName)
			writeFileString(fullName, value)
		}
	}
}
func unpackDataFile(items []importStruct) {
	for _, item := range items {
		createDir(filepath.Join(outputName, item.dir()))
		value := item.Value
		switch item.dir() {
		case dirTable:
			value = item.Columns
		case dirLang:
			value = item.Trans
		}
		fullName := filepath.Join(item.dir(), item.fullName())
		writeFileString(fullName, value)
	}
}
