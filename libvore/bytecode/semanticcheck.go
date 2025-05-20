package bytecode

import (
	"github.com/jmeaster30/vore/libvore/ast"
	"github.com/jmeaster30/vore/libvore/ds"
)

type ProcessContext int

const (
	PREDICATE ProcessContext = iota
	TRANSFORMATION
)

type ProcessTypeInfo struct {
	currentType ds.Optional[ValueType]
	context     ProcessContext
	environment map[string]ValueType
	inLoop      bool
}

func checkStatement(s *ast.AstProcessStatement, info ProcessTypeInfo) (ProcessTypeInfo, error) {
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
	case *ast.AstProcessBreak:
		return checkBreak(ps, info)
	case *ast.AstProcessContinue:
		return checkContinue(ps, info)
	case *ast.AstProcessExpression:
		return checkExpression(ps, info)
	case *ast.AstProcessDebug:
		return checkDebug(ps, info)
	}
	return info, NewSemanticError(*s, "Unknown AstProcessStatment type")
}

func checkSet(s *ast.AstProcessSet, info ProcessTypeInfo) (ProcessTypeInfo, error) {
	valueInfo, err := checkExpression(&s.Expr, info)
	if err != nil {
		return info, NewSemanticError(s, "")
	}
	valueInfo.environment[s.Name] = valueInfo.currentType.GetValue()
	valueInfo.currentType = ds.None[ValueType]()
	return valueInfo, nil
}

func checkReturn(s *ast.AstProcessReturn, info ProcessTypeInfo) (ProcessTypeInfo, error) {
	valueInfo, err := checkExpression(&s.Expr, info)
	if err != nil {
		return info, err
	}

	if valueInfo.context == PREDICATE && valueInfo.currentType != ds.Some(ValueType_Boolean) {
		return info, NewSemanticError(s, "Since we are in the predicate of a pattern, return values must be a boolean")
	} else if valueInfo.context == TRANSFORMATION && valueInfo.currentType != ds.Some(ValueType_String) && valueInfo.currentType != ds.Some(ValueType_Number) {
		return info, NewSemanticError(s, "Since we are in a transform function, return values must be a string or a number")
	} else {
		valueInfo.currentType = ds.None[ValueType]()
	}

	return valueInfo, nil
}

func checkIf(s *ast.AstProcessIf, info ProcessTypeInfo) (ProcessTypeInfo, error) {
	valueInfo, err := checkExpression(&s.Condition, info)
	if err != nil {
		return valueInfo, err
	}

	if valueInfo.currentType != ds.Some(ValueType_Boolean) {
		return valueInfo, NewSemanticError(s, "Condition of an if statement must be a boolean.")
	}

	for _, stmt := range s.TrueBody {
		valueInfo, err = checkStatement(&stmt, valueInfo)
		if err != nil {
			return valueInfo, err
		}
	}

	for _, stmt := range s.FalseBody {
		valueInfo, err = checkStatement(&stmt, valueInfo)
		if err != nil {
			return valueInfo, err
		}
	}

	return valueInfo, nil
}

func checkDebug(s *ast.AstProcessDebug, info ProcessTypeInfo) (ProcessTypeInfo, error) {
	valueInfo, err := checkExpression(&s.Expr, info)
	if err != nil {
		return valueInfo, err
	}

	valueInfo.currentType = ds.None[ValueType]()
	return valueInfo, nil
}

func checkLoop(s *ast.AstProcessLoop, info ProcessTypeInfo) (ProcessTypeInfo, error) {
	info.inLoop = true
	for _, stmt := range s.Body {
		valueInfo, err := checkStatement(&stmt, info)
		if err != nil {
			return valueInfo, err
		}
	}
	info.inLoop = false
	return info, nil
}

func checkContinue(s *ast.AstProcessContinue, info ProcessTypeInfo) (ProcessTypeInfo, error) {
	if !info.inLoop {
		return info, NewSemanticError(s, "Cannot use 'continue' outside of a loop.")
	}
	return info, nil
}

func checkBreak(s *ast.AstProcessBreak, info ProcessTypeInfo) (ProcessTypeInfo, error) {
	if !info.inLoop {
		return info, NewSemanticError(s, "Cannot use 'break' outside of a loop.")
	}
	return info, nil
}

func checkExpression(s *ast.AstProcessExpression, info ProcessTypeInfo) (ProcessTypeInfo, error) {
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
	return info, NewSemanticError(*s, "Unknown expression")
}

func checkBinaryExpr(s *ast.AstProcessBinaryExpression, info ProcessTypeInfo) (ProcessTypeInfo, error) {
	lhsinfo, lerr := checkExpression(&s.Lhs, info)
	rhsinfo, rerr := checkExpression(&s.Rhs, info)

	if lerr != nil {
		return lhsinfo, lerr
	}
	if rerr != nil {
		return rhsinfo, rerr
	}

	// super basic need to expand on this
	if lhsinfo.currentType == ds.Some(ValueType_String) && s.Op == ast.PLUS {
		lhsinfo.currentType = ds.Some(ValueType_String)
	} else if lhsinfo.currentType == ds.Some(ValueType_String) && (s.Op == ast.DEQUAL || s.Op == ast.NEQUAL || s.Op == ast.LESS || s.Op == ast.GREATER || s.Op == ast.LESSEQ || s.Op == ast.GREATEREQ) {
		lhsinfo.currentType = ds.Some(ValueType_Boolean)
	} else if lhsinfo.currentType == ds.Some(ValueType_Boolean) && (s.Op == ast.AND || s.Op == ast.OR || s.Op == ast.DEQUAL || s.Op == ast.NEQUAL || s.Op == ast.LESS || s.Op == ast.GREATER || s.Op == ast.LESSEQ || s.Op == ast.GREATEREQ) {
		lhsinfo.currentType = ds.Some(ValueType_Boolean)
	} else if lhsinfo.currentType == ds.Some(ValueType_Number) && (s.Op == ast.DEQUAL || s.Op == ast.NEQUAL || s.Op == ast.LESS || s.Op == ast.GREATER || s.Op == ast.LESSEQ || s.Op == ast.GREATEREQ) {
		lhsinfo.currentType = ds.Some(ValueType_Boolean)
	} else if lhsinfo.currentType == ds.Some(ValueType_Number) && (s.Op == ast.PLUS || s.Op == ast.MINUS || s.Op == ast.MULT || s.Op == ast.DIV || s.Op == ast.MOD) {
		lhsinfo.currentType = ds.Some(ValueType_Number)
	} else if lhsinfo.currentType == ds.Some(ValueType_String) && rhsinfo.currentType == ds.Some(ValueType_Number) && (s.Op == ast.PLUS || s.Op == ast.MINUS || s.Op == ast.MULT || s.Op == ast.DIV || s.Op == ast.MOD) {
		lhsinfo.currentType = ds.Some(ValueType_Number)
	} else {
		return lhsinfo, NewSemanticError(s, "Operator not defined for type.")
	}

	return lhsinfo, nil
}

func checkUnaryExpr(s *ast.AstProcessUnaryExpression, info ProcessTypeInfo) (ProcessTypeInfo, error) {
	next_info, err := checkExpression(&s.Expr, info)
	if err != nil {
		return next_info, err
	}

	if next_info.currentType == ds.Some(ValueType_Boolean) && s.Op == ast.NOT {
		next_info.currentType = ds.Some(ValueType_Boolean)
	} else if next_info.currentType == ds.Some(ValueType_String) && (s.Op == ast.HEAD || s.Op == ast.TAIL) {
		next_info.currentType = ds.Some(ValueType_String)
	} else {
		return next_info, NewSemanticError(s, "This operator is not valid on this expression")
	}
	return next_info, nil
}

func checkString(s *ast.AstProcessString, info ProcessTypeInfo) (ProcessTypeInfo, error) {
	info.currentType = ds.Some(ValueType_String)
	return info, nil
}

func checkNumber(s *ast.AstProcessNumber, info ProcessTypeInfo) (ProcessTypeInfo, error) {
	info.currentType = ds.Some(ValueType_Number)
	return info, nil
}

func checkBoolean(s *ast.AstProcessBoolean, info ProcessTypeInfo) (ProcessTypeInfo, error) {
	info.currentType = ds.Some(ValueType_Boolean)
	return info, nil
}

func checkVariable(s *ast.AstProcessVariable, info ProcessTypeInfo) (ProcessTypeInfo, error) {
	t, prs := info.environment[s.Name]
	if prs {
		info.currentType = ds.Some(t)
	} else {
		info.currentType = ds.Some(ValueType_String)
	}
	return info, nil
}
