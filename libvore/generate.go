package libvore

var loop_id int

func (f *AstFind) generate() Command {
	loop_id = 0
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
	loop_id += 1
	body := l.body.generate(offset + 1)

	start := StartLoop{
		id:       loop_id,
		minLoops: l.min,
		maxLoops: l.max,
		exitLoop: offset + len(body) + 1,
		fewest:   l.fewest,
	}

	stop := StopLoop{
		id:        loop_id,
		minLoops:  l.min,
		maxLoops:  l.max,
		startLoop: offset,
		fewest:    l.fewest,
	}

	result := []Instruction{start}
	result = append(result, body...)
	result = append(result, stop)
	return result
}

func (l *AstBranch) generate(offset int) []Instruction {

	left := l.left.generate(offset + 1)
	right := l.right.generate(offset + 2 + len(left))

	b := Branch{
		branches: []int{
			offset + 1,
			offset + len(left) + 2,
		},
	}

	insts := []Instruction{b}
	insts = append(insts, left...)
	insts = append(insts, Jump{
		newProgramCounter: offset + len(left) + len(right) + 3,
	})
	insts = append(insts, right...)
	insts = append(insts, Jump{
		newProgramCounter: offset + len(left) + len(right) + 3,
	})
	return insts
}

func (l *AstDec) generate(offset int) []Instruction {
	insts := []Instruction{}
	if l.isSubroutine {
		panic("subroutines aren't generated yet")
		// start sub routine
		// sub routine body
		// end sub routine
	} else {
		// offset
		startVarDec := StartVarDec{
			name: l.name,
		}

		bodyinsts := l.body.generate(offset + 1)

		endVarDec := EndVarDec{
			name: l.name,
		}

		insts = append(insts, startVarDec)
		insts = append(insts, bodyinsts...)
		insts = append(insts, endVarDec)
	}
	return insts
}

func (l *AstList) generate(offset int) []Instruction {
	b := Branch{
		branches: []int{},
	}

	pc := offset + 1
	branches := [][]Instruction{}
	for _, elem := range l.contents {
		branch_insts := elem.generate(pc)
		b.branches = append(b.branches, pc)
		pc += len(branch_insts) + 1
		branches = append(branches, branch_insts)
	}

	end := offset + 1
	for _, instss := range branches {
		end += len(instss) + 1
	}

	insts := []Instruction{b}
	for _, instss := range branches {
		insts = append(insts, instss...)
		insts = append(insts, Jump{
			newProgramCounter: end,
		})
	}

	return insts
}

func (l *AstPrimary) generate(offset int) []Instruction {
	return l.literal.generate(offset)
}

func (l *AstRange) generate(offset int) []Instruction {
	result := MatchRange{
		from: l.from.value,
		to:   l.to.value,
	}
	return []Instruction{result}
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
	result := MatchVariable{
		name: l.name,
	}
	return []Instruction{result}
}

func (l *AstCharacterClass) generate(offset int) []Instruction {
	result := MatchCharClass{
		class: l.classType,
	}
	return []Instruction{result}
}
