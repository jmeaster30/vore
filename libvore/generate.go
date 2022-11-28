package libvore

type GenState struct {
	loopId    int
	variables map[string]int
}

func (f *AstFind) generate() Command {
	result := FindCommand{
		all:  f.all,
		skip: f.skip,
		take: f.take,
		last: f.last,
		body: []Instruction{},
	}

	current_state := &GenState{
		loopId:    0,
		variables: make(map[string]int),
	}

	offset := 0
	for _, expr := range f.body {
		expr_insts := expr.generate(offset, current_state)
		offset = offset + len(expr_insts)
		result.body = append(result.body, expr_insts...)
	}

	return result
}

func (r *AstReplace) generate() Command {
	result := ReplaceCommand{
		all:      r.all,
		skip:     r.skip,
		take:     r.take,
		last:     r.last,
		body:     []Instruction{},
		replacer: []RInstruction{},
	}

	current_state := &GenState{
		loopId:    0,
		variables: make(map[string]int),
	}

	offset := 0
	for _, expr := range r.body {
		expr_insts := expr.generate(offset, current_state)
		offset = offset + len(expr_insts)
		result.body = append(result.body, expr_insts...)
	}

	offset = 0
	for _, expr := range r.result {
		expr_insts := expr.generateReplace(offset, current_state)
		offset = offset + len(expr_insts)
		result.replacer = append(result.replacer, expr_insts...)
	}

	return result
}

func (s *AstSet) generate() Command {
	return SetCommand{}
}

func (l *AstLoop) generate(offset int, state *GenState) []Instruction {
	state.loopId += 1
	body := l.body.generate(offset+1, state)

	start := StartLoop{
		id:       state.loopId,
		minLoops: l.min,
		maxLoops: l.max,
		exitLoop: offset + len(body) + 1,
		fewest:   l.fewest,
	}

	stop := StopLoop{
		id:        state.loopId,
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

func (l *AstBranch) generate(offset int, state *GenState) []Instruction {

	left := l.left.generate(offset+1, state)
	right := l.right.generate(offset+2+len(left), state)

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

func (l *AstDec) generate(offset int, state *GenState) []Instruction {
	insts := []Instruction{}
	// offset
	startVarDec := StartVarDec{
		name: l.name,
	}

	bodyinsts := l.body.generate(offset+1, state)

	endVarDec := EndVarDec{
		name: l.name,
	}

	_, prs := state.variables[l.name]
	if prs {
		panic("Name clash '" + l.name + "'")
	}
	state.variables[l.name] = -1

	insts = append(insts, startVarDec)
	insts = append(insts, bodyinsts...)
	insts = append(insts, endVarDec)
	return insts
}

func (l *AstSub) generate(offset int, state *GenState) []Instruction {
	insts := []Instruction{}

	_, prs := state.variables[l.name]
	if prs {
		panic("Name clash '" + l.name + "'")
	}
	state.variables[l.name] = offset

	bodyinsts := []Instruction{}
	loffset := offset + 1
	for _, expr := range l.body {
		expr_insts := expr.generate(loffset, state)
		loffset = loffset + len(expr_insts)
		bodyinsts = append(bodyinsts, expr_insts...)
	}

	startVarDec := StartSubroutine{
		id:        offset,
		name:      l.name,
		endOffset: loffset,
	}

	endVarDec := EndSubroutine{
		name: l.name,
	}

	insts = append(insts, startVarDec)
	insts = append(insts, bodyinsts...)
	insts = append(insts, endVarDec)
	return insts
}

func (l *AstList) generate(offset int, state *GenState) []Instruction {
	b := Branch{
		branches: []int{},
	}

	pc := offset + 1
	branches := [][]Instruction{}
	for _, elem := range l.contents {
		branch_insts := elem.generate(pc, state)
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

func (l *AstPrimary) generate(offset int, state *GenState) []Instruction {
	return l.literal.generate(offset, state)
}

func (l *AstRange) generate(offset int, state *GenState) []Instruction {
	result := MatchRange{
		from: l.from.value,
		to:   l.to.value,
	}
	return []Instruction{result}
}

func (l *AstString) generate(offset int, state *GenState) []Instruction {
	result := MatchLiteral{
		toFind: l.value,
	}
	return []Instruction{result}
}

func (l *AstSubExpr) generate(offset int, state *GenState) []Instruction {
	result := []Instruction{}

	loffset := offset
	for _, expr := range l.body {
		expr_insts := expr.generate(loffset, state)
		loffset = loffset + len(expr_insts)
		result = append(result, expr_insts...)
	}

	return result
}

func (l *AstVariable) generate(offset int, state *GenState) []Instruction {
	val, prs := state.variables[l.name]
	if !prs {
		panic("Variable not defined")
	}
	var result Instruction
	if val == -1 {
		result = MatchVariable{
			name: l.name,
		}
	} else {
		result = CallSubroutine{
			name: l.name,
			toPC: val,
		}
	}
	return []Instruction{result}
}

func (l *AstCharacterClass) generate(offset int, state *GenState) []Instruction {
	result := MatchCharClass{
		class: l.classType,
	}
	return []Instruction{result}
}

func (l *AstString) generateReplace(offset int, state *GenState) []RInstruction {
	result := ReplaceString{
		value: l.value,
	}
	return []RInstruction{result}
}

func (l *AstVariable) generateReplace(offset int, state *GenState) []RInstruction {
	result := ReplaceVariable{
		name: l.name,
	}
	return []RInstruction{result}
}
