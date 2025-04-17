package bytecode

import (
	"fmt"

	"github.com/jmeaster30/vore/libvore/ast"
)

type ProcessType int

const (
	PTUNKNOWN ProcessType = iota
	PTOK
	PTERROR
	PTSTRING
	PTNUMBER
	PTBOOLEAN
)

type ProcessContext int

const (
	PREDICATE ProcessContext = iota
	TRANSFORMATION
)

type ProcessTypeInfo struct {
	currentType  ProcessType
	errorMessage string
	context      ProcessContext
	environment  map[string]ProcessType
	inLoop       bool
}

func checkStatement(s *ast.AstProcessStatement, info ProcessTypeInfo) ProcessTypeInfo {
	var si any = *s
	switch ps := si.(type) {
	case *ast.AstProcessSet:
		return checkSet(ps, info)
	case *ast.AstProcessReturn:
		return checkReturn(ps, info)
	case *ast.AstProcessIf:
		return checkIf(ps, info)
	case *ast.AstProcessLoop:
		return checkLoop(ps, info)
	case ast.AstProcessBreak:
		return checkBreak(info)
	case ast.AstProcessContinue:
		return checkContinue(info)
	case *ast.AstProcessExpression:
		return checkExpression(ps, info)
	case *ast.AstProcessDebug:
		return checkDebug(ps, info)
	}
	info.currentType = PTERROR
	info.errorMessage = fmt.Sprintf("Unknown expression '%T'", si)
	return info
}

func checkSet(s *ast.AstProcessSet, info ProcessTypeInfo) ProcessTypeInfo {
	valueInfo := checkExpression(&s.Expr, info)
	if valueInfo.currentType == PTERROR {
		return valueInfo
	}
	valueInfo.environment[s.Name] = valueInfo.currentType
	valueInfo.currentType = PTOK
	return valueInfo
}

func checkReturn(s *ast.AstProcessReturn, info ProcessTypeInfo) ProcessTypeInfo {
	valueInfo := checkExpression(&s.Expr, info)
	if valueInfo.currentType == PTERROR {
		return valueInfo
	}

	if valueInfo.context == PREDICATE && valueInfo.currentType != PTBOOLEAN {
		valueInfo.currentType = PTERROR
		valueInfo.errorMessage = "Since we are in the predicate of a pattern, return values must be a boolean"
	} else if valueInfo.context == TRANSFORMATION && valueInfo.currentType != PTSTRING && valueInfo.currentType != PTNUMBER {
		valueInfo.currentType = PTERROR
		valueInfo.errorMessage = "Since we are in a transform function, return values must be a string or a number"
	} else {
		valueInfo.currentType = PTOK
	}

	return valueInfo
}

func checkIf(s *ast.AstProcessIf, info ProcessTypeInfo) ProcessTypeInfo {
	valueInfo := checkExpression(&s.Condition, info)
	if valueInfo.currentType == PTERROR {
		return valueInfo
	}

	if valueInfo.currentType != PTBOOLEAN {
		valueInfo.currentType = PTERROR
		valueInfo.errorMessage = "Condition of an if statement must be a boolean."
		return valueInfo
	}

	for _, stmt := range s.TrueBody {
		valueInfo = checkStatement(&stmt, valueInfo)
		if valueInfo.currentType == PTERROR {
			return valueInfo
		}
	}

	for _, stmt := range s.FalseBody {
		valueInfo = checkStatement(&stmt, valueInfo)
		if valueInfo.currentType == PTERROR {
			return valueInfo
		}
	}

	return valueInfo
}

func checkDebug(s *ast.AstProcessDebug, info ProcessTypeInfo) ProcessTypeInfo {
	valueInfo := checkExpression(&s.Expr, info)
	if valueInfo.currentType == PTERROR {
		return valueInfo
	}

	valueInfo.currentType = PTOK
	return valueInfo
}

func checkLoop(s *ast.AstProcessLoop, info ProcessTypeInfo) ProcessTypeInfo {
	info.inLoop = true
	for _, stmt := range s.Body {
		info = checkStatement(&stmt, info)
		if info.currentType == PTERROR {
			return info
		}
	}
	info.inLoop = false
	return info
}

func checkContinue(info ProcessTypeInfo) ProcessTypeInfo {
	if !info.inLoop {
		info.currentType = PTERROR
		info.errorMessage = "Cannot use 'continue' outside of a loop."
	}
	return info
}

func checkBreak(info ProcessTypeInfo) ProcessTypeInfo {
	if !info.inLoop {
		info.currentType = PTERROR
		info.errorMessage = "Cannot use 'break' outside of a loop."
	}
	return info
}

func checkExpression(s *ast.AstProcessExpression, info ProcessTypeInfo) ProcessTypeInfo {
	var si any = *s
	switch pe := si.(type) {
	case ast.AstProcessBinaryExpression:
		return checkBinaryExpr(&pe, info)
	case ast.AstProcessUnaryExpression:
		return checkUnaryExpr(&pe, info)
	case ast.AstProcessString:
		return checkString(&pe, info)
	case ast.AstProcessNumber:
		return checkNumber(&pe, info)
	case ast.AstProcessBoolean:
		return checkBoolean(&pe, info)
	case ast.AstProcessVariable:
		return checkVariable(&pe, info)
	}
	info.currentType = PTERROR
	info.errorMessage = fmt.Sprintf("Unknown expression '%T'", si)
	return info
}

func checkBinaryExpr(s *ast.AstProcessBinaryExpression, info ProcessTypeInfo) ProcessTypeInfo {
	lhsinfo := checkExpression(&s.Lhs, info)
	rhsinfo := checkExpression(&s.Rhs, info)
	// super basic need to expand on this
	if lhsinfo.currentType == PTERROR {
		return lhsinfo
	} else if rhsinfo.currentType == PTERROR {
		return rhsinfo
	} else if lhsinfo.currentType == PTSTRING && s.Op == ast.PLUS {
		lhsinfo.currentType = PTSTRING
	} else if lhsinfo.currentType == PTSTRING && (s.Op == ast.DEQUAL || s.Op == ast.NEQUAL || s.Op == ast.LESS || s.Op == ast.GREATER || s.Op == ast.LESSEQ || s.Op == ast.GREATEREQ) {
		lhsinfo.currentType = PTBOOLEAN
	} else if lhsinfo.currentType == PTBOOLEAN && (s.Op == ast.AND || s.Op == ast.OR || s.Op == ast.DEQUAL || s.Op == ast.NEQUAL || s.Op == ast.LESS || s.Op == ast.GREATER || s.Op == ast.LESSEQ || s.Op == ast.GREATEREQ) {
		lhsinfo.currentType = PTBOOLEAN
	} else if lhsinfo.currentType == PTNUMBER && (s.Op == ast.DEQUAL || s.Op == ast.NEQUAL || s.Op == ast.LESS || s.Op == ast.GREATER || s.Op == ast.LESSEQ || s.Op == ast.GREATEREQ) {
		lhsinfo.currentType = PTBOOLEAN
	} else if lhsinfo.currentType == PTNUMBER && (s.Op == ast.PLUS || s.Op == ast.MINUS || s.Op == ast.MULT || s.Op == ast.DIV || s.Op == ast.MOD) {
		lhsinfo.currentType = PTNUMBER
	} else if lhsinfo.currentType == PTSTRING && rhsinfo.currentType == PTNUMBER && (s.Op == ast.PLUS || s.Op == ast.MINUS || s.Op == ast.MULT || s.Op == ast.DIV || s.Op == ast.MOD) {
		lhsinfo.currentType = PTNUMBER
	} else {
		lhsinfo.currentType = PTERROR
		lhsinfo.errorMessage = "Operator not defined for type."
	}

	return lhsinfo
}

func checkUnaryExpr(s *ast.AstProcessUnaryExpression, info ProcessTypeInfo) ProcessTypeInfo {
	next_info := checkExpression(&s.Expr, info)
	if next_info.currentType == PTBOOLEAN && s.Op == ast.NOT {
		next_info.currentType = PTBOOLEAN
	} else if next_info.currentType == PTSTRING && (s.Op == ast.HEAD || s.Op == ast.TAIL) {
		next_info.currentType = PTSTRING
	} else if next_info.currentType != PTERROR {
		next_info.currentType = PTERROR
		next_info.errorMessage = "This operator is not valid on this expression" // TODO add better error message here
	}
	return next_info
}

func checkString(s *ast.AstProcessString, info ProcessTypeInfo) ProcessTypeInfo {
	info.currentType = PTSTRING
	return info
}

func checkNumber(s *ast.AstProcessNumber, info ProcessTypeInfo) ProcessTypeInfo {
	info.currentType = PTNUMBER
	return info
}

func checkBoolean(s *ast.AstProcessBoolean, info ProcessTypeInfo) ProcessTypeInfo {
	info.currentType = PTBOOLEAN
	return info
}

func checkVariable(s *ast.AstProcessVariable, info ProcessTypeInfo) ProcessTypeInfo {
	t, prs := info.environment[s.Name]
	if prs {
		info.currentType = t
	} else {
		info.currentType = PTSTRING
	}
	return info
}
