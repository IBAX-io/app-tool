package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/tmc/dot"
)

var (
	graphMap      = map[string][]string{}
	dirsGraph     = []string{dirMenu, dirPage, dirSnippet, dirCon, dirTable}
	graphDot      = dot.NewGraph("G")
	contractsList = []string{}
	labelType     = "label"
	nodeColors    = map[string]string{
		dirPage:    "green",
		dirCon:     "red",
		dirMenu:    "blue",
		dirSnippet: "green",
	}
	nodeShapes = map[string]string{
		dirPage:    "record",
		dirCon:     "record",
		dirMenu:    "record",
		dirSnippet: "record",
	}
	page2Contr     = regexp.MustCompile("\\(.*?Contract:\\s*@?1?(\\w+)")
	page2Page      = regexp.MustCompile("\\(.*?Page:\\s*@?1?(\\w+)")
	tableWrite     = regexp.MustCompile("(?:DBInsert|DBUpdate|DBUpdateExt)\\(\\s*[\"]@?1?(\\w+?)[\"]")
	tableRead      = regexp.MustCompile("(?:DBFind|DBRow)\\(.*?[\"\x60']?@?1?(\\w+)['\x60\",)]")
	page2Table     = regexp.MustCompile("(?:DBFind)\\(.*?[\"\x60']?@?1?(\\w+)['\x60\",)]")
	includeSnippet = regexp.MustCompile("Include\\(\\s*Name:\\s*@?1?(.*?)[,\\s\\)]|Include\\(\\s*@?1?([^:]*?)[\\),\\s]")
)

func createGraph(filename string) {
	if !withGraph {
		return
	}
	graphDot.SetType(dot.DIGRAPH)
	graphDot.Set("rankdir", "LR")
	graphDot.Set("fontsize", "24")
	graphDot.Set("nojustify", "true")
	graphDot.Set("size", "30")
	graphDot.Set("ordering", "out")

	label := strings.Trim(outputName, separator)
	label = strings.TrimSuffix(label, ".json")
	labelGraph := fmt.Sprintf("%s\n%s", label, time.Now().Format(time.RFC850))
	graphDot.Set(labelType, labelGraph)

	graphList := []graphStruct{}
	dir := filepath.Dir(filename)
	dirAbs, _ := filepath.Abs(dir)
	files, err := os.ReadDir(dirAbs)
	if err != nil {
		return
	}

	for _, f := range files {
		fname := f.Name()
		fpath := filepath.Join(dirAbs, fname)
		if debug {
			fmt.Println(fpath)
		}
		sf, err := os.Stat(fpath)
		if err != nil {
			fmt.Println(err)
			return
		}

		if sf.IsDir() && stringInSlice(dirsGraph, fname) {
			graphList = append(graphList, dirToGraph(fpath)...)
		}
	}
	for _, dGraph := range dirsGraph {
		for _, gs := range graphList {
			if gs.Dir == dGraph {
				createNodeWithEdges(&gs)
			}
		}
	}
	writeGraph(filename)
}

func dirToGraph(path string) (out []graphStruct) {
	files, err := os.ReadDir(path)
	if err != nil {
		return
	}

	dirAbs, _ := filepath.Abs(path)
	dirAbsParts := strings.Split(dirAbs, separator)
	dir := dirAbsParts[len(dirAbsParts)-1]
	for _, f := range files {
		nameExt := f.Name()
		ext := filepath.Ext(nameExt)
		name := nameExt[:len(nameExt)-len(ext)]
		fileAbs, _ := filepath.Abs(filepath.Join(path, nameExt))

		if stringInSlice(dirsGraph, dir) {
			gs := graphStruct{}
			gs.Name = name
			if dir != dirTable {
				gs.Value = file2str(fileAbs)
			}
			gs.Group = parseGroup(name)
			gs.Dir = dir
			out = append(out, gs)

			if dir == dirCon {
				contractsList = append(contractsList, name)
			}
		}
	}
	return
}

func createNodeWithEdges(gs *graphStruct) {
	node := dot.NewNode(getNodeName(gs.Name, gs.Dir))
	settingsNode(node, gs.Dir)
	revert := false
	switch gs.Dir {
	case dirCon:
		createContractNodes(node, gs, dirCon)
		createNodes(node, tableWrite, gs, dirTable, revert)
		createNodes(node, tableRead, gs, dirTable, true)
	case dirPage:
		createNodes(node, page2Contr, gs, dirCon, revert)
		createNodes(node, page2Table, gs, dirTable, true)
		createNodes(node, page2Page, gs, dirPage, revert)
		createNodes(node, includeSnippet, gs, dirSnippet, revert)
	case dirSnippet:
		createNodes(node, page2Contr, gs, dirCon, revert)
		createNodes(node, page2Table, gs, dirTable, true)
		createNodes(node, page2Page, gs, dirPage, revert)
	case dirMenu:
		createNodes(node, page2Page, gs, dirPage, revert)
	}

	graphDot.AddNode(node)
}

func createNodes(parentNode *dot.Node, pat *regexp.Regexp, gs *graphStruct, dir string, revert bool) {
	s := strings.Replace(gs.Value, "`", `"`, -1)
	arr := pat.FindAllStringSubmatch(s, -1)
	for _, match := range arr {
		for i := range match {
			if i > 0 {
				if match[i] != "" {
					name := match[i]
					if !stringInSlice(graphMap[parentNode.Name()], name) { // check exist graph heads
						createNode(parentNode, name, dir, gs, revert)
					}
				}
			}
		}
	}
}

func settingsNode(node *dot.Node, dir string) {
	node.Set("fontcolor", nodeColors[dir])
	node.Set("color", nodeColors[dir])
	node.Set("rank", "same")
	node.Set("group", dir)
	node.Set("shape", nodeShapes[dir])
}

func createNode(parentNode *dot.Node, nameOrig, dir string, gs *graphStruct, revert bool) {
	name := getNodeName(nameOrig, dir)
	parentName := parentNode.Name()
	node := dot.NewNode(name)
	settingsNode(node, dir)

	if _, ok := graphMap[parentName]; !ok {
		graphMap[parentName] = []string{}
	}

	edge := dot.NewEdge(parentNode, node)
	edge.Set("color", nodeColors[dir])
	if revert {
		edge = dot.NewEdge(node, parentNode)
		// edge.Set("color", nodeColors[gs.Dir])
	}

	if dir == dirSnippet {
		edge.Set(labelType, "include")
	}

	graphDot.AddEdge(edge)
	graphMap[parentName] = append(graphMap[parentName], nameOrig)
}

func createContractNodes(parentNode *dot.Node, gs *graphStruct, dir string) {
	s := strings.Replace(gs.Value, "`", `"`, -1)
	for _, name := range contractsList {
		if name != gs.Name && strings.Contains(s, name) {
			if !stringInSlice(graphMap[parentNode.Name()], name) { // check exist graph heads
				createNode(parentNode, name, dir, gs, false)
			}
		}
	}
}

func getNodeName(name, dir string) (_name string) {
	_name = fmt.Sprintf("%s\n%s", name, strings.TrimSuffix(dir, "s"))
	if strings.Contains(_name, ",") {
		_name = strings.Join(strings.Split(_name, ","), "\n")
	}
	_name = strings.Replace(_name, `"`, "", -1)
	_name = strings.Replace(_name, "`", "", -1)
	return
}

func writeGraph(name string) {
	outFile, err := os.Create(name)
	if err != nil {
		fmt.Println("error write file:", err)
		return
	}
	defer outFile.Close()
	if _, err := outFile.WriteString(graphDot.String()); err != nil {
		fmt.Println(err)
		return
	}
}

func parseGroup(n string) string {
	name := underscore(n)
	if strings.Contains(name, "_") {
		parts := strings.Split(name, "_")
		return strings.ToLower(parts[0])
	}
	return "0"
}

var camel = regexp.MustCompile("(^[^A-Z0-9]*|[A-Z0-9]*)([A-Z0-9][^A-Z]+|$)")

func underscore(s string) string {
	var a []string
	for _, sub := range camel.FindAllStringSubmatch(s, -1) {
		if sub[1] != "" {
			a = append(a, sub[1])
		}
		if sub[2] != "" {
			a = append(a, sub[2])
		}
	}
	return strings.ToLower(strings.Join(a, "_"))
}
