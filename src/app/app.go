package app

import (
	"compiler/src/codegen"
	"compiler/src/parser"
	"compiler/src/scanner"
	"compiler/src/semanticanalyzer"
)

// App ...
func App(inputFile string) {
	tokenList := scanner.ScanFile(inputFile)
	scanner.PrintTokenList(tokenList)
	parseTreeRoot := parser.Parse(tokenList)
	parser.PrintParseNodes(&parseTreeRoot, 0)
	semanticanalyzer.SemanticAnalysis(&parseTreeRoot, parser.GetGlobalSymbolTable(), parser.GetBuiltinSymbolTable())
	codegen.GenerateC(&parseTreeRoot, parser.GetGlobalSymbolTable(), parser.GetBuiltinSymbolTable())
}
