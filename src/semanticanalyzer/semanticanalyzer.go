package semanticanalyzer

import (
	"compiler/src/types"
	"errors"
	"log"
)

var globalSymbolTable = map[string]types.STEntry{}
var builtinSymbolTable = map[string]types.STEntry{}

func SemanticAnalysis(node *types.ParseNode, parseGlobalSymbolTable map[string]types.STEntry, parseBuiltinSymbolTable map[string]types.STEntry) {
	globalSymbolTable = parseGlobalSymbolTable
	builtinSymbolTable = parseGlobalSymbolTable

	err := CheckNode(node, nil, types.STEntry{})
	if err != nil {
		log.Fatal(err)
	}
}

func CheckNode(node *types.ParseNode, localSymbolTable map[string]types.STEntry, stEntry types.STEntry) error {
	var localST map[string]types.STEntry
	localST = localSymbolTable
	if node.Production == types.ProcedureDeclarationProd {
		localST = node.ProcLocalSymbolTable
	}
	var entry types.STEntry
	entry = stEntry
	var identifier types.ParseNode
	if node.Production == types.ProcedureDeclarationProd {
		header := node.ChildNodes[0]
		// if header.ChildNodes[0].TerminalToken.TokenType == types.GlobalKeyword {
		// 	identifier = node.ChildNodes[2]
		// } else {
		// 	identifier = node.ChildNodes[1]
		// }
		identifier = header.ChildNodes[1]
		entry = localST[identifier.TerminalToken.StringValue]
	}

	for _, child := range node.ChildNodes {
		CheckNode(&child, localST, entry)
	}

	if node.Production == types.AssignmentStatementProd {
		// for rule 14
		return CheckAssignmentStatementNode(node, localST)
	}
	if node.Production == types.LoopStatementProd {
		// for rule 15
		// check assignment statement
		return CheckLoopStatementNode(node, localST, entry)
	}
	if node.Production == types.IfStatementProd {
		// for rule 15
		// check assignment statement
		return CheckIfStatementNode(node, localST, entry)
	}
	if node.Production == types.ReturnStatementProd {
		// for rule 15
		// check assignment statement
		return CheckReturnStatementNode(node, localST, entry)
	}

	return nil
}

func CheckProcedureCallNode(node *types.ParseNode, localSymbolTable map[string]types.STEntry) (types.STType, error) {
	errString := "Semantic Analysis Error: Procedure call argument types do not match procedure declaration parameter types"
	identifier := node.ChildNodes[0].TerminalToken.StringValue
	var stEntry types.STEntry
	stEntryLocal, existsLocal := localSymbolTable[identifier]
	stEntryGlobal, existsGlobal := globalSymbolTable[identifier]
	stEntryBuiltin, existsBuiltin := builtinSymbolTable[identifier]
	if existsLocal {
		stEntry = stEntryLocal
	} else if existsGlobal {
		stEntry = stEntryGlobal
	} else if existsBuiltin {
		stEntry = stEntryBuiltin
	}

	if node.ChildNodes[2].Production == types.ArgumentListProd {
		err := CheckArgumentListNode(&node.ChildNodes[3], localSymbolTable, stEntry)
		if err != nil {
			return types.STNone, errors.New(err.Error())
		}
		// DO SOEMTHING HERE?
	} else if node.ChildNodes[2].TerminalToken.TokenType == types.CloseRoundBracket {
		if len(stEntry.ProcedureArgTypes) != 0 {
			return types.STNone, errors.New(errString)
		}
	}

	return stEntry.ProcedureReturnType, nil
}

func CheckArgumentListNode(node *types.ParseNode, localSymbolTable map[string]types.STEntry, stEntry types.STEntry) error {
	errString := "Semantic Analysis Error: Procedure call argument types do not match procedure declaration parameter types"

	argListSTTypes := []types.STType{}

	for _, child := range node.ChildNodes {
		argType, err := CheckExpressionNode(&child, localSymbolTable)
		if err != nil {
			return errors.New(err.Error())
		}
		argListSTTypes = append(argListSTTypes, argType)
	}

	if len(argListSTTypes) != len(stEntry.ProcedureArgTypes) {
		return errors.New(errString)
	}

	for i := range argListSTTypes {
		if argListSTTypes[i] != stEntry.ProcedureArgTypes[i] {
			return errors.New(errString)
		}
	}

	return nil
}

func CheckAssignmentStatementNode(node *types.ParseNode, localSymbolTable map[string]types.STEntry) error {
	errString := "Semantic Analysis Error: Expression type is not compatible with destination type"

	destSTType, err := CheckDestinationNode(&node.ChildNodes[0], localSymbolTable)
	if err != nil {
		return errors.New(err.Error())
	}

	exprSTType, err := CheckExpressionNode(&node.ChildNodes[2], localSymbolTable)
	if err != nil {
		return errors.New(err.Error())
	}

	if exprSTType != destSTType {
		if !(exprSTType == types.STVarBool && destSTType == types.STVarInteger) &&
			!(exprSTType == types.STVarInteger && destSTType == types.STVarBool) &&
			!(exprSTType == types.STVarInteger && destSTType == types.STVarFloat) &&
			!(exprSTType == types.STVarFloat && destSTType == types.STVarInteger) {
			return errors.New(errString)
		}
	}

	return nil
}

func CheckDestinationNode(node *types.ParseNode, localSymbolTable map[string]types.STEntry) (types.STType, error) {
	errString := "Semantic Analysis Error: Index expression does not evaluate to an integer"

	identifier := node.ChildNodes[0].TerminalToken.StringValue
	var stEntry types.STEntry
	stEntryLocal, existsLocal := localSymbolTable[identifier]
	stEntryGlobal, existsGlobal := globalSymbolTable[identifier]
	if existsLocal {
		stEntry = stEntryLocal
	} else if existsGlobal {
		stEntry = stEntryGlobal
	}

	if len(node.ChildNodes) > 1 {
		exprSTType, err := CheckExpressionNode(&node.ChildNodes[2], localSymbolTable)
		if err != nil {
			return types.STNone, errors.New(err.Error())
		}
		if exprSTType != types.STVarInteger {
			return types.STNone, errors.New(errString)
		}

		if stEntry.EntryType == types.STVarIntegerArray {
			return types.STVarInteger, nil
		} else if stEntry.EntryType == types.STVarFloatArray {
			return types.STVarFloat, nil
		} else if stEntry.EntryType == types.STVarStringArray {
			return types.STVarString, nil
		} else if stEntry.EntryType == types.STVarBoolArray {
			return types.STVarBool, nil
		}
	}

	return stEntry.EntryType, nil
}

func CheckNameNode(node *types.ParseNode, localSymbolTable map[string]types.STEntry) (types.STType, error) {
	errString := "Semantic Analysis Error: Index expression does not evaluate to an integer"

	identifier := node.ChildNodes[0].TerminalToken.StringValue
	var stEntry types.STEntry
	stEntryLocal, existsLocal := localSymbolTable[identifier]
	stEntryGlobal, existsGlobal := globalSymbolTable[identifier]
	if existsLocal {
		stEntry = stEntryLocal
	} else if existsGlobal {
		stEntry = stEntryGlobal
	}

	if len(node.ChildNodes) > 1 {
		exprSTType, err := CheckExpressionNode(&node.ChildNodes[2], localSymbolTable)
		if err != nil {
			return types.STNone, errors.New(err.Error())
		}
		if exprSTType != types.STVarInteger {
			return types.STNone, errors.New(errString)
		}

		if stEntry.EntryType == types.STVarIntegerArray {
			return types.STVarInteger, nil
		} else if stEntry.EntryType == types.STVarFloatArray {
			return types.STVarFloat, nil
		} else if stEntry.EntryType == types.STVarStringArray {
			return types.STVarString, nil
		} else if stEntry.EntryType == types.STVarBoolArray {
			return types.STVarBool, nil
		}
	}

	return stEntry.EntryType, nil
}

func CheckExpressionNode(node *types.ParseNode, localSymbolTable map[string]types.STEntry) (types.STType, error) {
	errString := "Semantic Analysis Error: Incompatible types"

	// hasNot := false
	aop_index := 0
	exp_index := 1
	if node.ChildNodes[0].TerminalToken.TokenType == types.NotOperator {
		// hasNot = true
		aop_index = 1
		exp_index = 2
	}

	stType, err := CheckArithOpNode(&node.ChildNodes[aop_index], localSymbolTable)
	if err != nil {
		return types.STNone, errors.New(err.Error())
	}

	if len(node.ChildNodes) > 1 {
		termPrimeSTType, err := CheckExpressionPrimeNode(&node.ChildNodes[exp_index], localSymbolTable)
		if err != nil {
			return types.STNone, errors.New(err.Error())
		}

		if stType == termPrimeSTType && (stType == types.STVarInteger || stType == types.STVarBool) {
			return stType, nil
		} else {
			return types.STNone, errors.New(errString)
		}
	}

	return stType, nil
}

func CheckExpressionPrimeNode(node *types.ParseNode, localSymbolTable map[string]types.STEntry) (types.STType, error) {
	errString := "Semantic Analysis Error: Incompatible types"

	// hasNot := false
	aop_index := 0
	exp_index := 1
	if node.ChildNodes[0].TerminalToken.TokenType == types.NotOperator {
		// hasNot = true
		aop_index = 1
		exp_index = 2
	}

	stType, err := CheckArithOpNode(&node.ChildNodes[aop_index], localSymbolTable)
	if err != nil {
		return types.STNone, errors.New(err.Error())
	}

	if len(node.ChildNodes) > 1 {
		termPrimeSTType, err := CheckExpressionPrimeNode(&node.ChildNodes[exp_index], localSymbolTable)
		if err != nil {
			return types.STNone, errors.New(err.Error())
		}

		if stType == termPrimeSTType && (stType == types.STVarInteger || stType == types.STVarBool) {
			return stType, nil
		} else {
			return types.STNone, errors.New(errString)
		}
	}

	if stType != types.STVarInteger && stType != types.STVarBool {
		return types.STNone, errors.New(errString)
	}

	return stType, nil
}

func CheckArithOpNode(node *types.ParseNode, localSymbolTable map[string]types.STEntry) (types.STType, error) {
	errString := "Semantic Analysis Error: Incompatible types"
	stType, err := CheckRelationNode(&node.ChildNodes[0], localSymbolTable)
	if err != nil {
		return types.STNone, errors.New(err.Error())
	}

	if len(node.ChildNodes) > 1 {
		termPrimeSTType, err := CheckArithOpPrimeNode(&node.ChildNodes[1], localSymbolTable)
		if err != nil {
			return types.STNone, errors.New(err.Error())
		}

		if stType == termPrimeSTType && (stType == types.STVarInteger || stType == types.STVarFloat) {
			return stType, nil
		} else if stType == types.STVarInteger && termPrimeSTType == types.STVarFloat {
			return termPrimeSTType, nil
		} else if stType == types.STVarFloat && termPrimeSTType == types.STVarInteger {
			return stType, nil
		} else {
			return types.STNone, errors.New(errString)
		}
	}

	return stType, nil
}

func CheckArithOpPrimeNode(node *types.ParseNode, localSymbolTable map[string]types.STEntry) (types.STType, error) {
	errString := "Semantic Analysis Error: Incompatible types"
	stType, err := CheckFactorNode(&node.ChildNodes[1], localSymbolTable)
	if err != nil {
		return types.STNone, errors.New(err.Error())
	}

	if len(node.ChildNodes) > 2 {
		termPrimeSTType, err := CheckTermPrimeNode(&node.ChildNodes[2], localSymbolTable)
		if err != nil {
			return types.STNone, errors.New(err.Error())
		}

		if stType == termPrimeSTType && (stType == types.STVarInteger || stType == types.STVarFloat) {
			return stType, nil
		} else if stType == types.STVarInteger && termPrimeSTType == types.STVarFloat {
			return termPrimeSTType, nil
		} else if stType == types.STVarFloat && termPrimeSTType == types.STVarInteger {
			return stType, nil
		} else {
			return types.STNone, errors.New(errString)
		}
	}

	if stType != types.STVarInteger && stType != types.STVarFloat {
		return types.STNone, errors.New(errString)
	}

	return stType, nil
}

func CheckRelationNode(node *types.ParseNode, localSymbolTable map[string]types.STEntry) (types.STType, error) {
	errString := "Semantic Analysis Error: Incompatible types"
	stType, err := CheckTermNode(&node.ChildNodes[0], localSymbolTable)
	if err != nil {
		return types.STNone, errors.New(err.Error())
	}

	if len(node.ChildNodes) > 1 {
		relPrimeSTType, err := CheckRelationPrimeNode(&node.ChildNodes[1], localSymbolTable)
		if err != nil {
			return types.STNone, errors.New(err.Error())
		}

		if stType == relPrimeSTType && (stType == types.STVarInteger || stType == types.STVarFloat || stType == types.STVarBool) {
			return types.STVarBool, nil
		} else if stType == types.STVarBool && relPrimeSTType == types.STVarInteger {
			return stType, nil
		} else if stType == types.STVarInteger && relPrimeSTType == types.STVarBool {
			return relPrimeSTType, nil
		} else if stType == types.STVarFloat && relPrimeSTType == types.STVarFloat {
			return types.STVarBool, nil
		} else if stType == types.STVarString && relPrimeSTType == types.STVarString {
			if node.ChildNodes[1].ChildNodes[0].TerminalToken.TokenType == types.EqualOperator || node.ChildNodes[1].ChildNodes[0].TerminalToken.TokenType == types.NotEqualOperator {
				return types.STVarBool, nil
			} else {
				return types.STNone, errors.New(errString)
			}
		} else {
			return types.STNone, errors.New(errString)
		}
	}

	return stType, nil
}

func CheckRelationPrimeNode(node *types.ParseNode, localSymbolTable map[string]types.STEntry) (types.STType, error) {
	errString := "Semantic Analysis Error: Incompatible types"
	stType, err := CheckTermNode(&node.ChildNodes[1], localSymbolTable)
	if err != nil {
		return types.STNone, errors.New(err.Error())
	}

	if len(node.ChildNodes) > 2 {
		relPrimeSTType, err := CheckRelationPrimeNode(&node.ChildNodes[2], localSymbolTable)
		if err != nil {
			return types.STNone, errors.New(err.Error())
		}

		if stType == relPrimeSTType && (stType == types.STVarInteger || stType == types.STVarFloat || stType == types.STVarBool) {
			return stType, nil
		} else if stType == types.STVarBool && relPrimeSTType == types.STVarInteger {
			return stType, nil
		} else if stType == types.STVarInteger && relPrimeSTType == types.STVarBool {
			return relPrimeSTType, nil
		} else if stType == types.STVarFloat && relPrimeSTType == types.STVarFloat {
			return stType, nil
		} else if stType == types.STVarString && relPrimeSTType == types.STVarString {
			if node.ChildNodes[1].ChildNodes[0].TerminalToken.TokenType == types.EqualOperator || node.ChildNodes[1].ChildNodes[0].TerminalToken.TokenType == types.NotEqualOperator {
				return types.STVarBool, nil
			} else {
				return types.STNone, errors.New(errString)
			}
		} else {
			return types.STNone, errors.New(errString)
		}
	}

	if stType != types.STVarInteger && stType != types.STVarFloat && stType != types.STVarBool && stType != types.STVarString {
		return types.STNone, errors.New(errString)
	}

	return stType, nil
}

func CheckTermNode(node *types.ParseNode, localSymbolTable map[string]types.STEntry) (types.STType, error) {
	errString := "Semantic Analysis Error: Incompatible types"
	stType, err := CheckFactorNode(&node.ChildNodes[0], localSymbolTable)
	if err != nil {
		return types.STNone, errors.New(err.Error())
	}

	if len(node.ChildNodes) > 1 {
		termPrimeSTType, err := CheckTermPrimeNode(&node.ChildNodes[1], localSymbolTable)
		if err != nil {
			return types.STNone, errors.New(err.Error())
		}

		if stType == termPrimeSTType && (stType == types.STVarInteger || stType == types.STVarFloat) {
			return stType, nil
		} else if stType == types.STVarInteger && termPrimeSTType == types.STVarFloat {
			return termPrimeSTType, nil
		} else if stType == types.STVarFloat && termPrimeSTType == types.STVarInteger {
			return stType, nil
		} else {
			return types.STNone, errors.New(errString)
		}
	}

	return stType, nil
}

func CheckTermPrimeNode(node *types.ParseNode, localSymbolTable map[string]types.STEntry) (types.STType, error) {
	errString := "Semantic Analysis Error: Incompatible types"
	stType, err := CheckFactorNode(&node.ChildNodes[1], localSymbolTable)
	if err != nil {
		return types.STNone, errors.New(err.Error())
	}

	if len(node.ChildNodes) > 2 {
		termPrimeSTType, err := CheckTermPrimeNode(&node.ChildNodes[2], localSymbolTable)
		if err != nil {
			return types.STNone, errors.New(err.Error())
		}

		if stType == termPrimeSTType && (stType == types.STVarInteger || stType == types.STVarFloat) {
			return stType, nil
		} else if stType == types.STVarInteger && termPrimeSTType == types.STVarFloat {
			return termPrimeSTType, nil
		} else if stType == types.STVarFloat && termPrimeSTType == types.STVarInteger {
			return stType, nil
		} else {
			return types.STNone, errors.New(errString)
		}
	}

	if stType != types.STVarInteger && stType != types.STVarFloat {
		return types.STNone, errors.New(errString)
	}

	return stType, nil
}

func CheckFactorNode(node *types.ParseNode, localSymbolTable map[string]types.STEntry) (types.STType, error) {
	errString := "Semantic Analysis Error: Unknown factor type"

	if node.ChildNodes[0].TerminalToken.TokenType == types.SubtractionOperator {
		if node.ChildNodes[1].Production == types.NameProd {
			return CheckNameNode(&node.ChildNodes[1], localSymbolTable)
		} else if node.ChildNodes[1].Production == types.NumberProd {
			if node.ChildNodes[1].TerminalToken.TokenType == types.FloatToken {
				return types.STVarFloat, nil
			} else if node.ChildNodes[1].TerminalToken.TokenType == types.IntegerToken {
				return types.STVarInteger, nil
			}
		}
	} else if node.ChildNodes[0].TerminalToken.TokenType == types.OpenRoundBracket {
		return CheckExpressionNode(&node.ChildNodes[1], localSymbolTable)
	} else if node.ChildNodes[0].Production == types.ProcedureCallProd {
		return CheckProcedureCallNode(&node.ChildNodes[0], localSymbolTable)
	} else if node.ChildNodes[0].Production == types.NameProd {
		return CheckNameNode(&node.ChildNodes[0], localSymbolTable)
	} else if node.ChildNodes[0].Production == types.NumberProd {
		if node.ChildNodes[0].TerminalToken.TokenType == types.FloatToken {
			return types.STVarFloat, nil
		} else if node.ChildNodes[0].TerminalToken.TokenType == types.IntegerToken {
			return types.STVarInteger, nil
		}
	} else if node.ChildNodes[0].Production == types.StringProd {
		return types.STVarString, nil
	} else if node.ChildNodes[0].TerminalToken.TokenType == types.TrueKeyword || node.ChildNodes[0].TerminalToken.TokenType == types.FalseKeyword {
		return types.STVarBool, nil
	}

	return types.STNone, errors.New(errString)
}

func CheckLoopStatementNode(node *types.ParseNode, localSymbolTable map[string]types.STEntry, stEntry types.STEntry) error {
	errString := "Semantic Analysis Error: Loop expression does not evaluate to a boolean value"

	err := CheckAssignmentStatementNode(&node.ChildNodes[2], localSymbolTable)
	if err != nil {
		return errors.New(err.Error())
	}

	stType, err := CheckExpressionNode(&node.ChildNodes[4], localSymbolTable)
	if err != nil {
		return errors.New(err.Error())
	}
	if stType != types.STVarBool || stType != types.STVarInteger {
		return errors.New(errString)
	}

	for _, child := range node.ChildNodes[6:] {
		if child.TerminalToken.TokenType == types.EndKeyword {
			break
		}
		err = CheckNode(&child, localSymbolTable, stEntry)
		if err != nil {
			return errors.New(err.Error())
		}
	}

	return nil
}

func CheckIfStatementNode(node *types.ParseNode, localSymbolTable map[string]types.STEntry, stEntry types.STEntry) error {
	errString := "Semantic Analysis Error: If expression does not evaluate to a boolean value"

	stType, err := CheckExpressionNode(&node.ChildNodes[2], localSymbolTable)
	if err != nil {
		return errors.New(err.Error())
	}
	if stType != types.STVarBool || stType != types.STVarInteger {
		return errors.New(errString)
	}

	for _, child := range node.ChildNodes[5:] {
		if child.TerminalToken.TokenType == types.EndKeyword {
			break
		}
		err = CheckNode(&child, localSymbolTable, stEntry)
		if err != nil {
			return errors.New(err.Error())
		}
	}

	return nil
}

func CheckReturnStatementNode(node *types.ParseNode, localSymbolTable map[string]types.STEntry, stEntry types.STEntry) error {
	errString := "Semantic Analysis Error: Returned value does not match procedure declaration return type"

	stType, err := CheckExpressionNode(&node.ChildNodes[1], localSymbolTable)
	if err != nil {
		return errors.New(err.Error())
	}

	if stType != stEntry.EntryType {
		return errors.New(errString)
	}

	return nil
}
