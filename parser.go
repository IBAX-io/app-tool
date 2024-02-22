package main

import (
	parser "github.com/IBAX-io/needle/grammar"
)

var filePath string

func ParserGrammarBase(data string) error {
	p, err := parser.NewParserBase([]byte(data), "")
	if err != nil {
		return err
	}
	p.Parse()
	p.PrintlnError()
	return nil
}

func ParserGrammarFile(filePath string) error {
	p, err := parser.NewParserFile(filePath)
	if err != nil {
		return err
	}
	p.Parse()
	p.PrintlnError()
	return nil
}
