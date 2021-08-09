package types

import (
	"strconv"
	"strings"
)

type TokenType string

const (
	// ProgramKeyword ...
	ProgramKeyword TokenType = "program"
	// IsKeyword ...
	IsKeyword TokenType = "is"
	// BeginKeyword ...
	BeginKeyword TokenType = "begin"
	// EndKeyword ...
	EndKeyword TokenType = "end"
	// GlobalKeyword ...
	GlobalKeyword TokenType = "global"
	// ProcedureKeyword ...
	ProcedureKeyword TokenType = "procedure"
	// VariableKeyword ...
	VariableKeyword TokenType = "variable"
	// TypeKeyword ...
	TypeKeyword TokenType = "type"
	// IntegerKeyword ...
	IntegerKeyword TokenType = "integer"
	// FloatKeyword ...
	FloatKeyword TokenType = "float"
	// StringKeyword ...
	StringKeyword TokenType = "string"
	// BoolKeyword ...
	BoolKeyword TokenType = "bool"
	// EnumKeyword ...
	EnumKeyword TokenType = "enum"
	// IfKeyword ...
	IfKeyword TokenType = "if"
	// ThenKeyword ...
	ThenKeyword TokenType = "then"
	// ElseKeyword ...
	ElseKeyword TokenType = "else"
	// ForKeyword ...
	ForKeyword TokenType = "for"
	// ReturnKeyword ...
	ReturnKeyword TokenType = "return"
	// TrueKeyword ...
	TrueKeyword TokenType = "true"
	// FalseKeyword ...
	FalseKeyword TokenType = "false"

	// OpenRoundBracket ...
	OpenRoundBracket TokenType = "("
	// CloseRoundBracket ...
	CloseRoundBracket TokenType = ")"
	// OpenCurlyBracket ...
	OpenCurlyBracket TokenType = "{"
	// CloseCurlyBracket ...
	CloseCurlyBracket TokenType = "}"
	// OpenSquareBracket ...
	OpenSquareBracket TokenType = "["
	// CloseSquareBracket ...
	CloseSquareBracket TokenType = "]"

	// SemiColonSymbol ...
	SemiColonSymbol TokenType = ";"
	// ColonSymbol ...
	ColonSymbol TokenType = ":"
	// CommaSymbol ...
	CommaSymbol TokenType = ","
	// PeriodSymbol ...
	PeriodSymbol TokenType = "."

	// AssignmentOperator ...
	AssignmentOperator TokenType = ":="
	// AndOperator ...
	AndOperator TokenType = "&"
	// OrOperator ...
	OrOperator TokenType = "|"
	// NotOperator ...
	NotOperator TokenType = "not"
	// AdditionOperator ...
	AdditionOperator TokenType = "+"
	// SubtractionOperator ...
	SubtractionOperator TokenType = "-"
	// LessThanOperator ...
	LessThanOperator TokenType = "<"
	// LessThanEqualOperator ...
	LessThanEqualOperator TokenType = "<="
	// GreaterThanOperator ...
	GreaterThanOperator TokenType = ">"
	// GreaterThanEqualOperator ...
	GreaterThanEqualOperator TokenType = ">="
	// EqualOperator ...
	EqualOperator TokenType = "=="
	//NotEqualOperator ...
	NotEqualOperator TokenType = "!="
	// MultiplicationOperator ...
	MultiplicationOperator TokenType = "*"
	// DivisionOperator ...
	DivisionOperator TokenType = "/"

	// IdentifierToken ..
	IdentifierToken TokenType = "IdentifierToken"
	// IntegerToken ...
	IntegerToken TokenType = "IntegerToken"
	// FloatToken ...
	FloatToken TokenType = "FloatToken"
	// StringToken ...
	StringToken TokenType = "StringToken"
)

var KeywordTokenTypeMap = map[string]TokenType{
	"program":   ProgramKeyword,
	"is":        IsKeyword,
	"begin":     BeginKeyword,
	"end":       EndKeyword,
	"global":    GlobalKeyword,
	"procedure": ProcedureKeyword,
	"variable":  VariableKeyword,
	"type":      TypeKeyword,
	"integer":   IntegerKeyword,
	"float":     FloatKeyword,
	"string":    StringKeyword,
	"bool":      BoolKeyword,
	"enum":      EnumKeyword,
	"if":        IfKeyword,
	"then":      ThenKeyword,
	"else":      ElseKeyword,
	"for":       ForKeyword,
	"return":    ReturnKeyword,
	"true":      TrueKeyword,
	"false":     FalseKeyword,
	"not":       NotOperator,
}

var SymbolTokenTypeMap = map[string]TokenType{
	"(":  OpenRoundBracket,
	")":  CloseRoundBracket,
	"{":  OpenCurlyBracket,
	"}":  CloseCurlyBracket,
	"[":  OpenSquareBracket,
	"]":  CloseSquareBracket,
	";":  SemiColonSymbol,
	":":  ColonSymbol,
	",":  CommaSymbol,
	".":  PeriodSymbol,
	":=": AssignmentOperator,
	"&":  AndOperator,
	"|":  OrOperator,
	"+":  AdditionOperator,
	"-":  SubtractionOperator,
	"<":  LessThanOperator,
	"<=": LessThanEqualOperator,
	">":  GreaterThanOperator,
	">=": GreaterThanEqualOperator,
	"==": EqualOperator,
	"!=": NotEqualOperator,
	"*":  MultiplicationOperator,
	"/":  DivisionOperator,
}

type ScanCode int

const (
	// TokenScanCode ...
	TokenScanCode ScanCode = 0
	// StopWithTokenScanCode ...
	StopWithTokenScanCode ScanCode = 1
	// NewlineScanCode ...
	NewlineScanCode ScanCode = 2
	// AtNextByteScanCode ...
	AtNextByteScanCode ScanCode = 3
	// ErrorScanCode ...
	ErrorScanCode ScanCode = 4
	// LineCommentScanCode ...
	LineCommentScanCode ScanCode = 5
	// BlockCommentOpenScanCode ...
	BlockCommentOpenScanCode ScanCode = 6
	// BlockCommentCloseScanCode ...
	BlockCommentCloseScanCode ScanCode = 7
	// StopNoTokenScanCode ...
	StopNoTokenScanCode ScanCode = 8
)

type Token struct {
	LineNumber  int
	TokenType   TokenType
	StringValue string
	IntValue    int64
	FloatValue  float64
}

func BuildToken(lineNumber int, word string, scanCode ScanCode) (Token, ScanCode) {
	token := Token{lineNumber, "", word, -1, -1}
	return token, scanCode
}

func BuildTokenFromAlphaNumeric(lineNumber int, word string, scanCode ScanCode) (Token, ScanCode) {
	var token Token
	tokenType, exists := KeywordTokenTypeMap[word]
	if exists {
		token = Token{lineNumber, tokenType, "", -1, -1}
	} else {
		token = Token{lineNumber, IdentifierToken, word, -1, -1}
	}
	return token, scanCode
}

func BuildTokenFromNumeric(lineNumber int, word string, scanCode ScanCode) (Token, ScanCode) {
	var token Token
	if strings.Index(word, ".") == -1 {
		intValue, err := strconv.ParseInt(word, 10, 64)
		if err != nil {
			token = Token{lineNumber, "", "", -1, -1}
			scanCode = ErrorScanCode
		} else {
			token = Token{lineNumber, IntegerToken, "", intValue, -1}
		}
	} else {
		float64Value, err := strconv.ParseFloat(word, 64)
		if err != nil {
			token = Token{lineNumber, "", "", -1, -1}
			scanCode = ErrorScanCode
		} else {
			token = Token{lineNumber, FloatToken, "", -1, float64Value}
		}
	}
	return token, scanCode
}

func BuildTokenFromString(lineNumber int, word string, scanCode ScanCode) (Token, ScanCode) {
	token := Token{lineNumber, StringToken, word, -1, -1}
	return token, scanCode
}

func BuildTokenFromSymbol(lineNumber int, word string, scanCode ScanCode) (Token, ScanCode) {
	var token Token
	tokenType, exists := SymbolTokenTypeMap[word]
	if exists {
		token = Token{lineNumber, tokenType, "", -1, -1}
	} else {
		token = Token{lineNumber, "", "", -1, -1}
		scanCode = ErrorScanCode
	}
	return token, scanCode
}

type ProductionType string

const (
	// ProgramProd ...
	ProgramProd ProductionType = "<program>"
	// ProgramHeaderProd ...
	ProgramHeaderProd ProductionType = "<program_header>"
	// ProgramBodyProd ...
	ProgramBodyProd ProductionType = "<program_body>"
	// IdentifierProd ...
	IdentifierProd ProductionType = "<identifier>"
	// DeclarationProd ...
	DeclarationProd ProductionType = "<declaration>"
	// StatementProd ...
	StatementProd ProductionType = "<statement>"
	// ProcedureDeclarationProd ...
	ProcedureDeclarationProd ProductionType = "<procedure_declaration>"
	// VariableDeclarationProd ...
	VariableDeclarationProd ProductionType = "<variable_declaration>"
	// // TypeDeclarationProd ...
	// TypeDeclarationProd ProductionType = ""
	// ProcedureHeaderProd ...
	ProcedureHeaderProd ProductionType = "<procedure_header>"
	// ProcedureBodyProd ...
	ProcedureBodyProd ProductionType = "<procedure_body>"
	// TypeMarkProd ...
	TypeMarkProd ProductionType = "<type_mark>"
	// ParamaterListProd ...
	ParamaterListProd ProductionType = "<parameter_list>"
	// ParamaterProd ...
	ParamaterProd ProductionType = "<parameter>"
	// BoundProd ...
	BoundProd ProductionType = "<bound>"
	// NumberProd ...
	NumberProd ProductionType = "<number>"
	// AssignmentStatementProd ..
	AssignmentStatementProd ProductionType = "<assignment_statement>"
	// IfStatementProd ...
	IfStatementProd ProductionType = "<if_statement>"
	// LoopStatementProd ...
	LoopStatementProd ProductionType = "<loop_statement>"
	// ReturnStatementProd ...
	ReturnStatementProd ProductionType = "<return_statement>"
	// ProcedureCallProd ...
	ProcedureCallProd ProductionType = "<procedure_call>"
	// DestinationProd ...
	DestinationProd ProductionType = "<destination>"
	// ExpressionProd ...
	ExpressionProd ProductionType = "<expression>"
	// ExpressionPrimeProd ...
	ExpressionPrimeProd ProductionType = "<expression>"
	// ArithOpProd ...
	ArithOpProd ProductionType = "<arithOp>"
	// ArithOpPrimeProd ...
	ArithOpPrimeProd ProductionType = "<arithOpPrime>"
	// RelationProd ...
	RelationProd ProductionType = "<relation>"
	// RelationPrimeProd ...
	RelationPrimeProd ProductionType = "<relationPrime>"
	// TermProd ...
	TermProd ProductionType = "<term>"
	// TermPrimeProd ...
	TermPrimeProd ProductionType = "<termPrime>"
	// FactorProd ...
	FactorProd ProductionType = "<factor>"
	// NameProd ...
	NameProd ProductionType = "<name>"
	// AgurmentList ...
	ArgumentListProd ProductionType = "<argument_list>"
	// StringProd ...
	StringProd ProductionType = "<string>"
	// KeywordTerminal ...
	KeywordTerminal ProductionType = "<KeywordTerminal>"
	// SymbolTerminal ...
	SymbolTerminal ProductionType = "<SymbolTerminal>"
)

type ParseNode struct {
	Production           ProductionType
	TerminalToken        Token
	ChildNodes           []ParseNode
	ProcLocalSymbolTable map[string]STEntry
}

type STType string

const (
	// STVarInteger ...
	STVarInteger STType = "integer"
	// STVarIntegerArray ...
	STVarIntegerArray STType = "integer_array"
	// STVarFloat ...
	STVarFloat STType = "float"
	// STVarFloatArray ...
	STVarFloatArray STType = "float_array"
	// STVarString ...
	STVarString STType = "string"
	// STVarStringArray ...
	STVarStringArray STType = "string_array"
	// STVarBool ...
	STVarBool STType = "bool"
	// STVarBoolArray ...
	STVarBoolArray STType = "bool_array"
	// STProcedure ...
	STProcedure STType = "procedure"
	STNone      STType = "none"
)

type STEntry struct {
	Identifier          string
	EntryType           STType
	IsArray             bool
	ArraySize           int
	ProcedureArgTypes   []STType
	ProcedureReturnType STType
	Pointer             int
}
