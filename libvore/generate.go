package libvore

func (f *AstFind) generate() Command {
	result := FindCommand{
		all:  f.all,
		skip: f.skip,
		take: f.take,
		body: []Instruction{},
	}

	offset := 0
	for _, expr := range f.body {
		expr_insts := expr.generate(offset)
		offset = offset + len(expr_insts)
		result.body = append(result.body, expr_insts...)
	}

	return result
}

func (r *AstReplace) generate() Command {
	panic(":(")
}

func (s *AstSet) generate() Command {
	panic(":(")
}

func (l *AstLoop) generate(offset int) []Instruction {
	panic(":(")
}

func (l *AstOptional) generate(offset int) []Instruction {
	panic(":(")
}

func (l *AstBranch) generate(offset int) []Instruction {
	panic(":(")
}

func (l *AstDec) generate(offset int) []Instruction {
	panic(":(")
}

func (l *AstList) generate(offset int) []Instruction {
	panic(":(")
}

func (l *AstPrimary) generate(offset int) []Instruction {
	return l.literal.generate(offset)
}

func (l *AstRange) generate(offset int) []Instruction {
	panic(":(")
}

func (l *AstString) generate(offset int) []Instruction {
	result := MatchLiteral{
		toFind: l.value,
	}

	return []Instruction{result}
}

func (l *AstSubExpr) generate(offset int) []Instruction {
	result := []Instruction{}

	loffset := offset
	for _, expr := range l.body {
		expr_insts := expr.generate(loffset)
		loffset = loffset + len(expr_insts)
		result = append(result, expr_insts...)
	}

	return result
}

func (l *AstVariable) generate(offset int) []Instruction {
	panic(":(")
}

func (l *AstCharacterClass) generate(offset int) []Instruction {
	result := MatchCharClass{
		class: l.classType,
	}
	return []Instruction{result}
}
