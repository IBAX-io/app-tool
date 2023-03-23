package main

import (
	"encoding/json"
	"fmt"
	"github.com/spkg/bom"
	"os"
	"path/filepath"
	"strconv"
)

func unpackJSON(filename string) {
	bs, err := os.ReadFile(filename)
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
		if split {
			splitDataFile(file)
			return
		} else {
			unpackDataFile(file.Data)
		}
	} else {
		file := importFile{}
		if err := json.Unmarshal(bs, &file); err != nil {
			fmt.Println("unmarshal file error_1:", err)
			return
		}
		unpackStruct(file.Contracts, eSIM, dirCon)
		unpackStruct(file.Menus, ePTL, dirMenu)
		unpackStruct(file.Snippets, ePTL, dirSnippet)
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

func splitDataFile(items dataFile) {
	var tp = make(map[string][]importStruct)
	for _, out := range items.Data {
		var path string
		ext := filepath.Ext(outputName)
		path = outputName[:len(outputName)-len(ext)] + "." + out.dir() + ext
		tp[path] = append(tp[path], out)
	}
	for path, out := range tp {
		segments := splitArray(out, int64(number))
		for i := 0; i < len(segments); i++ {
			split := segments[i]
			data := dataFile{}
			data.Name = items.Name
			data.Conditions = items.Conditions
			data.Data = append(data.Data, split...)
			result, _ := _JSONMarshal(data, true)
			ext := filepath.Ext(path)
			path := path[:len(path)-len(ext)] + strconv.Itoa(i+1) + ext
			writeFileString2(path, string(result))
		}
	}

}

func splitArray(arr []importStruct, num int64) [][]importStruct {
	max := int64(len(arr))

	if max <= num || num <= 0 {
		return [][]importStruct{arr}
	}
	var quantity int64
	if max%num == 0 {
		quantity = max / num
	} else {
		quantity = (max / num) + 1
	}
	var segments = make([][]importStruct, 0)
	var start, end, i int64
	for i = 1; i <= quantity; i++ {
		end = i * num
		if i != quantity {
			segments = append(segments, arr[start:end])
		} else {
			segments = append(segments, arr[start:])
		}
		start = i * num
	}
	return segments
}
