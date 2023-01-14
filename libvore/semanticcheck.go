package libvore

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

func (s AstProcessSet) check(info ProcessTypeInfo) ProcessTypeInfo {
	valueInfo := s.expr.check(info)
	if valueInfo.currentType == PTERROR {
		return valueInfo
	}
	valueInfo.environment[s.name] = valueInfo.currentType
	valueInfo.currentType = PTOK
	return valueInfo
}

func (s AstProcessReturn) check(info ProcessTypeInfo) ProcessTypeInfo {
	valueInfo := s.expr.check(info)
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

func (s AstProcessIf) check(info ProcessTypeInfo) ProcessTypeInfo {
	valueInfo := s.condition.check(info)
	if valueInfo.currentType == PTERROR {
		return valueInfo
	}

	if valueInfo.currentType != PTBOOLEAN {
		valueInfo.currentType = PTERROR
		valueInfo.errorMessage = "Condition of an if statement must be a boolean."
		return valueInfo
	}

	for _, stmt := range s.trueBody {
		valueInfo = stmt.check(valueInfo)
		if valueInfo.currentType == PTERROR {
			return valueInfo
		}
	}

	for _, stmt := range s.falseBody {
		valueInfo = stmt.check(valueInfo)
		if valueInfo.currentType == PTERROR {
			return valueInfo
		}
	}

	return valueInfo
}

func (s AstProcessDebug) check(info ProcessTypeInfo) ProcessTypeInfo {
	valueInfo := s.expr.check(info)
	if valueInfo.currentType == PTERROR {
		return valueInfo
	}

	valueInfo.currentType = PTOK
	return valueInfo
}

func (s AstProcessLoop) check(info ProcessTypeInfo) ProcessTypeInfo {
	for _, stmt := range s.body {
		info = stmt.check(info)
		if info.currentType == PTERROR {
			return info
		}
	}
	return info
}

func (s AstProcessContinue) check(info ProcessTypeInfo) ProcessTypeInfo {
	if !info.inLoop {
		info.currentType = PTERROR
		info.errorMessage = "Cannot use 'continue' outside of a loop."
	}
	return info
}

func (s AstProcessBreak) check(info ProcessTypeInfo) ProcessTypeInfo {
	if !info.inLoop {
		info.currentType = PTERROR
		info.errorMessage = "Cannot use 'continue' outside of a loop."
	}
	return info
}

func (s AstProcessBinaryExpression) check(info ProcessTypeInfo) ProcessTypeInfo {
	lhsinfo := s.lhs.check(info)
	rhsinfo := s.rhs.check(info)
	// super basic need to expand on this
	if lhsinfo.currentType == PTERROR {
		return lhsinfo
	} else if rhsinfo.currentType == PTERROR {
		return rhsinfo
	} else if lhsinfo.currentType == PTSTRING && s.op == PLUS {
		lhsinfo.currentType = PTSTRING
	} else if lhsinfo.currentType == PTSTRING && (s.op == DEQUAL || s.op == NEQUAL || s.op == LESS || s.op == GREATER || s.op == LESSEQ || s.op == GREATEREQ) {
		lhsinfo.currentType = PTBOOLEAN
	} else if lhsinfo.currentType == PTBOOLEAN && (s.op == AND || s.op == OR || s.op == DEQUAL || s.op == NEQUAL || s.op == LESS || s.op == GREATER || s.op == LESSEQ || s.op == GREATEREQ) {
		lhsinfo.currentType = PTBOOLEAN
	} else if lhsinfo.currentType == PTNUMBER && (s.op == DEQUAL || s.op == NEQUAL || s.op == LESS || s.op == GREATER || s.op == LESSEQ || s.op == GREATEREQ) {
		lhsinfo.currentType = PTBOOLEAN
	} else if lhsinfo.currentType == PTNUMBER && (s.op == PLUS || s.op == MINUS || s.op == MULT || s.op == DIV || s.op == MOD) {
		lhsinfo.currentType = PTNUMBER
	} else if lhsinfo.currentType == PTSTRING && rhsinfo.currentType == PTNUMBER && (s.op == PLUS || s.op == MINUS || s.op == MULT || s.op == DIV || s.op == MOD) {
		lhsinfo.currentType = PTNUMBER
	} else {
		lhsinfo.currentType = PTERROR
		lhsinfo.errorMessage = "Operator not defined for type."
	}

	return lhsinfo
}

func (s AstProcessUnaryExpression) check(info ProcessTypeInfo) ProcessTypeInfo {
	next_info := s.expr.check(info)
	if next_info.currentType == PTBOOLEAN && s.op == NOT {
		next_info.currentType = PTBOOLEAN
	} else if next_info.currentType == PTSTRING && (s.op == HEAD || s.op == TAIL) {
		next_info.currentType = PTSTRING
	} else if next_info.currentType != PTERROR {
		next_info.currentType = PTERROR
		next_info.errorMessage = "This operator is not valid on this expression" // TODO add better error message here
	}
	return next_info
}

func (s AstProcessString) check(info ProcessTypeInfo) ProcessTypeInfo {
	info.currentType = PTSTRING
	return info
}

func (s AstProcessNumber) check(info ProcessTypeInfo) ProcessTypeInfo {
	info.currentType = PTNUMBER
	return info
}

func (s AstProcessVariable) check(info ProcessTypeInfo) ProcessTypeInfo {
	t, prs := info.environment[s.name]
	if prs {
		info.currentType = t
	} else {
		info.currentType = PTSTRING
	}
	return info
}
