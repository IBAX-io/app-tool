package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func packJSON(path string) {
	out := packDir(path)
	var arr []string
	arr = append(arr)
	path = filepath.Dir(path)
	files, err := os.ReadDir(path)
	if err != nil {
		return
	}

	for _, f := range files {
		fname := f.Name()
		fpath := filepath.Join(path, fname)
		if debug {
			fmt.Println(fpath)
		}
		sf, err := os.Stat(fpath)
		if err != nil {
			fmt.Println(err)
			return
		}
		if sf.IsDir() {
			dir := packDir(fpath)
			switch fname {
			case dirSnippet:
				out.Snippets = append(out.Snippets, dir.Snippets...)
			case dirMenu:
				out.Menus = append(out.Menus, dir.Menus...)
			case dirLang:
				out.Languages = append(out.Languages, dir.Languages...)
			case dirTable:
				out.Tables = append(out.Tables, dir.Tables...)
			case dirParam:
				out.Parameters = append(out.Parameters, dir.Parameters...)
			case typeParam:
				out.Parameters = append(out.Parameters, dir.Parameters...)
			case dirData:
				out.Data = append(out.Data, dir.Data...)
			case dirPage:
				out.Pages = append(out.Pages, dir.Pages...)
			case dirCon:
				for _, cont := range dir.Contracts {
					err := ParserGrammarFile(cont.FullPath)
					if err != nil {
						panic(err)
					}
				}
				out.Contracts = append(out.Contracts, dir.Contracts...)
			}
		}
	}

	if countEntries(out) > 0 {
		readConfig(&out)
		if len(out.Contracts) > 0 {
			out.Contracts = sortContracts(out.Contracts)
		}
		var result []byte
		if importNew {
			data := dataFile{}
			data.Name = out.Name
			data.Conditions = out.Conditions
			data.Data = append(data.Data, out.Snippets...)
			data.Data = append(data.Data, out.Menus...)
			data.Data = append(data.Data, out.Languages...)
			data.Data = append(data.Data, out.Tables...)
			data.Data = append(data.Data, out.Parameters...)
			data.Data = append(data.Data, out.Pages...)
			data.Data = append(data.Data, out.Contracts...)
			result, _ = _JSONMarshal(data, true)
		} else {
			out.cleaning()
			result, _ = _JSONMarshal(&out, true)
		}

		if !strings.HasSuffix(outputName, eJSON) {
			outputName += eJSON
		}
		outFile, err := os.Create(outputName)
		if err != nil {
			if debug {
				fmt.Println(err)
			}
			return
		}
		defer outFile.Close()
		outFile.WriteString(string(result))

		if abs, err := filepath.Abs(path); err == nil {
			abspath := filepath.Join(abs, structFileName)
			createGraph(abspath)
		}
		fmt.Println("pack complete!\noutput file:", outputName)
	}
}
func packDir(path string) (out exportFile) {
	out.Snippets = []importStruct{}
	out.Contracts = []importStruct{}
	out.Data = []dataStruct{}
	out.Languages = []importStruct{}
	out.Menus = []importStruct{}
	out.Pages = []importStruct{}
	out.Parameters = []importStruct{}
	out.Tables = []importStruct{}

	files, err := os.ReadDir(path)
	if err != nil {
		return
	}

	absdir, _ := filepath.Abs(path)
	absdirParts := strings.Split(absdir, separator)
	fdir := absdirParts[len(absdirParts)-1]
	for _, f := range files {
		fname := f.Name()
		ext := filepath.Ext(fname)
		if debug {
			fmt.Println(fname)
		}

		switch ext {
		case ePTL:
			switch {
			case fdir == dirMenu || fdir == typeMenu:
				el := encodeStd(path, fname)
				el.Type = typeMenu
				out.Menus = append(out.Menus, el)
			case fdir == dirSnippet || fdir == typeSnippet:
				el := encodeStd(path, fname)
				el.Type = typeSnippet
				out.Snippets = append(out.Snippets, el)
			default:
				el := encodePage(path, fname)
				el.Type = typePage
				out.Pages = append(out.Pages, el)
			}
		case eJSON:
			switch {
			case fdir == dirParam || fdir == typeParam:
				el := encodeStd(path, fname)
				el.Type = typeParam
				out.Parameters = append(out.Parameters, el)
			case fdir == dirLang || fdir == typeLang:
				el := encodeLang(path, fname)
				el.Type = typeLang
				out.Languages = append(out.Languages, el)
			case fdir == dirTable || fdir == typeTable:
				el := encodeTable(path, fname)
				el.Type = typeTable
				out.Tables = append(out.Tables, el)
			case fdir == dirData:
				el := encodeData(path, fname)
				out.Data = append(out.Data, el)
			}
		case eCSV:
			switch {
			case fdir == dirParam || fdir == typeParam:
				el := encodeStd(path, fname)
				el.Type = typeParam
				out.Parameters = append(out.Parameters, el)
			}
		case eSIM:
			el := encodeStd(path, fname)
			el.Type = typeCon
			out.Contracts = append(out.Contracts, el)
		}
	}
	return
}

func encodePage(path, fname string) (result importStruct) {
	ext := filepath.Ext(fname)
	name := fname[:len(fname)-len(ext)]
	fpath := filepath.Join(path, fname)
	result.Menu = defaultMenu
	result.Name = name
	result.Value = file2str(fpath)
	result.Conditions = defaultCondition
	return
}
func encodeData(path, fname string) (result dataStruct) {
	ext := filepath.Ext(fname)
	name := fname[:len(fname)-len(ext)]
	fpath := filepath.Join(path, fname)
	result.Table = name
	dataFile := file2data(fpath)
	result.Columns = dataFile.Columns
	result.Data = dataFile.Data
	return
}
func encodeTable(path, fname string) (result importStruct) {
	ext := filepath.Ext(fname)
	name := fname[:len(fname)-len(ext)]
	fpath := filepath.Join(path, fname)
	result.Name = name
	result.Columns = file2str(fpath)
	result.Permissions = defaultPermission
	return
}
func encodeLang(path, fname string) (result importStruct) {
	ext := filepath.Ext(fname)
	name := fname[:len(fname)-len(ext)]
	fpath := filepath.Join(path, fname)
	result.Name = name
	result.Trans = file2str(fpath)
	result.Conditions = ""
	return
}

func encodeStd(path, fname string) (result importStruct) {
	ext := filepath.Ext(fname)
	name := fname[:len(fname)-len(ext)]
	fpath := filepath.Join(path, fname)
	result.FullPath = fpath
	result.Name = name
	result.Value = file2str(fpath)
	result.Conditions = defaultCondition
	return
}

func file2str(filename string) (str string) {
	bs, err := os.ReadFile(filename)
	if err != nil {
		return
	}
	str = string(bs)
	return
}

func file2data(filename string) (result dataStruct) {
	bs, err := os.ReadFile(filename)
	if err != nil {
		return
	}
	json.Unmarshal(bs, &result)
	return
}

func _JSONMarshal(v interface{}, unescape bool) ([]byte, error) {
	b, err := json.MarshalIndent(v, "", "    ")

	if unescape {
		b = bytes.Replace(b, []byte("\\u003c"), []byte("<"), -1)
		b = bytes.Replace(b, []byte("\\u003e"), []byte(">"), -1)
		b = bytes.Replace(b, []byte("\\u0026"), []byte("&"), -1)
	}
	return b, err
}
func sortContracts(c []importStruct) []importStruct {
	loops := map[string]string{}
	lenC := len(c)
	notSwapped := true
	for n := 0; n < lenC; n++ {
		notSwapped = true
		for i := lenC - 1; i > 0; i-- {
			for j := i - 1; j >= 0; j-- {
				if textContainsContract(c[j].Value, c[i].Name) {
					if textContainsContract(c[i].Value, c[j].Name) { // detect call contract loop
						if _, ok := loops[c[j].Name]; !ok {
							loops[c[i].Name] = c[j].Name
						}
					}
					c[i], c[j] = c[j], c[i]
					notSwapped = false
				}
			}
		}
		if notSwapped {
			break
		}
	}
	if len(loops) > 0 {
		fmt.Println("loops:")
		for key, val := range loops {
			fmt.Printf("%v <=> %v\n", key, val)
		}
	}
	return c
}

func textContainsContract(text, name string) bool {
	re := regexp.MustCompile(name + "\\s*\\(")

	lines := strings.Split(text, "\n")
	for _, l := range lines {
		line := strings.Trim(l, " ")
		if !strings.HasPrefix(line, "//") && re.MatchString(line) {
			return true
		}
	}
	return false
}
func countEntries(file exportFile) (count int) {
	return len(file.Snippets) +
		len(file.Contracts) +
		len(file.Data) +
		len(file.Languages) +
		len(file.Menus) +
		len(file.Pages) +
		len(file.Parameters) +
		len(file.Tables)
}
