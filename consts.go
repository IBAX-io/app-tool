package main

import (
	"os"
	"strings"
)

const (
	currentVersion = "0.9.14"
	currentTitle   = "Applications Packager " + currentVersion

	eSIM  = ".sim"
	ePTL  = ".ptl"
	eJSON = ".json"
	eCSV  = ".csv"

	dirBlock = "blocks"
	dirMenu  = "menus"
	dirLang  = "languages"
	dirTable = "tables"
	dirParam = "parameters"
	dirData  = "data"
	dirPage  = "pages"
	dirCon   = "contracts"

	typeBlock = "blocks"
	typeMenu  = "menu"
	typeLang  = "languages"
	typeTable = "tables"
	typeParam = "app_params"
	typePage  = "pages"
	typeCon   = "contracts"

	defaultCondition  = "ContractConditions(\"MainCondition\")"
	defaultMenu       = "default_menu"
	defaultPermission = "{\"insert\":\"true\",\"update\":\"true\",\"new_column\":\"true\"}"
	configName        = "config.json"
	separator         = string(os.PathSeparator)
	structFileName    = "struct.dot"

	helpMsg = "please choose directory for packing, example:\n    ap dirfiles" + separator + "\nor file to unpacking, example:\n    ap file.json"
)

type configFile struct {
	Name       string       `json:"name,omitempty"`
	Conditions string       `json:"conditions,omitempty"`
	Blocks     []importConf `json:"blocks"`
	Contracts  []importConf `json:"contracts"`
	Menus      []importConf `json:"menus"`
	Pages      []importConf `json:"pages"`
	Tables     []importConf `json:"tables"`
	Params     []importConf `json:"parameters"`
}

type exportFile struct {
	Name       string         `json:"name,omitempty"`
	Conditions string         `json:"conditions,omitempty"`
	Blocks     []importStruct `json:"blocks"`
	Contracts  []importStruct `json:"contracts"`
	Data       []dataStruct   `json:"data"`
	Languages  []importStruct `json:"languages"`
	Menus      []importStruct `json:"menus"`
	Pages      []importStruct `json:"pages"`
	Parameters []importStruct `json:"parameters"`
	Tables     []importStruct `json:"tables"`
}

func (e *exportFile) cleaning() {
	for i := range e.Blocks {
		e.Blocks[i].Type = ""
	}
	for i := range e.Contracts {
		e.Contracts[i].Type = ""
	}
	for i := range e.Languages {
		e.Languages[i].Type = ""
	}
	for i := range e.Menus {
		e.Menus[i].Type = ""
	}
	for i := range e.Pages {
		e.Pages[i].Type = ""
	}
	for i := range e.Parameters {
		e.Parameters[i].Type = ""
	}
	for i := range e.Tables {
		e.Tables[i].Type = ""
	}
}

type importFile struct {
	Blocks     []commonStruct `json:"blocks"`
	Contracts  []commonStruct `json:"contracts"`
	Data       []dataStruct   `json:"data"`
	Languages  []commonStruct `json:"languages"`
	Menus      []commonStruct `json:"menus"`
	Pages      []commonStruct `json:"pages"`
	Parameters []commonStruct `json:"parameters"`
	Tables     []commonStruct `json:"tables"`
	Name       string         `json:"name,omitempty"`
}

func (item *importStruct) dir() string {
	if !strings.HasSuffix(item.Type, "s") {
		return item.Type + "s"
	}
	return item.Type
}
func (item *importStruct) fullName() string {
	return item.Name + item.ext()
}
func (item *importStruct) ext() string {
	switch item.Type {
	case typeBlock:
		return ePTL
	case typeMenu:
		return ePTL
	case typePage:
		return ePTL
	case typeParam:
		return eCSV
	case typeCon:
		return eSIM
	default:
		return eJSON
	}
}

type dataFile struct {
	Name       string         `json:"name"`
	Conditions string         `json:"conditions,omitempty"`
	Data       []importStruct `json:"data"`
}
type dataConf struct {
	Name       string       `json:"name"`
	Conditions string       `json:"conditions,omitempty"`
	Data       []importConf `json:"data"`
}

type importStruct struct {
	Name         string `json:",omitempty"`
	Confirmation string `json:",omitempty"`
	Conditions   string `json:",omitempty"`
	Value        string `json:",omitempty"`
	Trans        string `json:",omitempty"`
	Menu         string `json:",omitempty"`
	Columns      string `json:",omitempty"`
	Permissions  string `json:",omitempty"`
	Type         string `json:",omitempty"`
}
type importConf struct {
	Name         string `json:",omitempty"`
	Confirmation string `json:",omitempty"`
	Conditions   string `json:",omitempty"`
	Menu         string `json:",omitempty"`
	Permissions  string `json:",omitempty"`
	Type         string `json:",omitempty"`
}

type commonStruct struct {
	Name       string
	Value      string
	Conditions string
	Trans      string
	Columns    string
	Table      string
}
type testFormatStruct struct {
	Name       string         `json:"name,omitempty"`
	Conditions string         `json:"conditions,omitempty"`
	Blocks     []importStruct `json:"blocks,omitempty"`
	Contracts  []importStruct `json:"contracts,omitempty"`
	Languages  []importStruct `json:"languages,omitempty"`
	Menus      []importStruct `json:"menus,omitempty"`
	Pages      []importStruct `json:"pages,omitempty"`
	Parameters []importStruct `json:"parameters,omitempty"`
	Tables     []importStruct `json:"tables,omitempty"`
}

func (t *testFormatStruct) len() (l int) {
	if t.Name != "" {
		l++
	}
	if t.Blocks != nil {
		l++
	}
	if t.Contracts != nil {
		l++
	}
	if t.Languages != nil {
		l++
	}
	if t.Menus != nil {
		l++
	}
	if t.Pages != nil {
		l++
	}
	if t.Parameters != nil {
		l++
	}
	if t.Tables != nil {
		l++
	}
	return l
}

type dataStruct struct {
	Table   string
	Columns []string
	Data    [][]string
}

type graphStruct struct {
	Name      string
	Value     string
	Group     string
	Path      string
	Dir       string
	FontColor string
	Color     string
	EdgeLabel string
}
