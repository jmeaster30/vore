package libvore

import "fmt"

type GenState struct {
	loopId    int
	variables map[string]int
}

func (f *AstFind) generate() (Command, error) {
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
		expr_insts, gen_error := expr.generate(offset, current_state)
		if gen_error != nil {
			return nil, gen_error
		}
		offset = offset + len(expr_insts)
		result.body = append(result.body, expr_insts...)
	}

	return result, nil
}

func (r *AstReplace) generate() (Command, error) {
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
		expr_insts, gen_error := expr.generate(offset, current_state)
		if gen_error != nil {
			return nil, gen_error
		}
		offset = offset + len(expr_insts)
		result.body = append(result.body, expr_insts...)
	}

	offset = 0
	for _, expr := range r.result {
		expr_insts, gen_error := expr.generateReplace(offset, current_state)
		if gen_error != nil {
			return nil, gen_error
		}
		offset = offset + len(expr_insts)
		result.replacer = append(result.replacer, expr_insts...)
	}

	return result, nil
}

func (s *AstSet) generate() (Command, error) {
	return SetCommand{}, nil
}

func (l *AstLoop) generate(offset int, state *GenState) ([]Instruction, error) {
	state.loopId += 1
	body, gen_error := l.body.generate(offset+1, state)
	if gen_error != nil {
		return []Instruction{}, gen_error
	}

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
	return result, nil
}

func (l *AstBranch) generate(offset int, state *GenState) ([]Instruction, error) {

	left, left_error := l.left.generate(offset+1, state)
	if left_error != nil {
		return []Instruction{}, left_error
	}
	right, right_error := l.right.generate(offset+2+len(left), state)
	if right_error != nil {
		return []Instruction{}, right_error
	}

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
	return insts, nil
}

func (l *AstDec) generate(offset int, state *GenState) ([]Instruction, error) {
	insts := []Instruction{}
	// offset
	startVarDec := StartVarDec{
		name: l.name,
	}

	bodyinsts, gen_error := l.body.generate(offset+1, state)
	if gen_error != nil {
		return []Instruction{}, gen_error
	}

	endVarDec := EndVarDec{
		name: l.name,
	}

	_, prs := state.variables[l.name]
	if prs {
		return []Instruction{}, fmt.Errorf("Name clash '%s'", l.name)
	}
	state.variables[l.name] = -1

	insts = append(insts, startVarDec)
	insts = append(insts, bodyinsts...)
	insts = append(insts, endVarDec)
	return insts, nil
}

func (l *AstSub) generate(offset int, state *GenState) ([]Instruction, error) {
	insts := []Instruction{}

	_, prs := state.variables[l.name]
	if prs {
		return []Instruction{}, fmt.Errorf("Name clash '%s'", l.name)
	}
	state.variables[l.name] = offset

	bodyinsts := []Instruction{}
	loffset := offset + 1
	for _, expr := range l.body {
		expr_insts, gen_error := expr.generate(loffset, state)
		if gen_error != nil {
			return []Instruction{}, gen_error
		}
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
	return insts, nil
}

func (l *AstList) generate(offset int, state *GenState) ([]Instruction, error) {
	if l.not {
		return l.generate_not(offset, state)
	} else {
		return l.generate_not_not(offset, state)
	}
}

func (l *AstList) generate_not_not(offset int, state *GenState) ([]Instruction, error) {
	b := Branch{
		branches: []int{},
	}

	pc := offset + 1
	branches := [][]Instruction{}
	for _, elem := range l.contents {
		branch_insts, gen_error := elem.generate(pc, state)
		if gen_error != nil {
			return []Instruction{}, gen_error
		}
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

	return insts, nil
}

func (l *AstList) generate_not(offset int, state *GenState) ([]Instruction, error) {
	insts := []Instruction{}

	pc := offset
	for _, item := range l.contents {
		item_insts, err := item.generate(pc+1, state)
		if err != nil {
			return []Instruction{}, err
		}

		insts = append(insts, StartNotIn{
			nextCheckpointPC: pc + len(item_insts) + 2,
		})
		insts = append(insts, item_insts...)
		insts = append(insts, FailNotIn{})
		pc = pc + len(item_insts) + 2
	}

	insts = append(insts, EndNotIn{
		maxSize: l.getMaxSize(),
	})
	return insts, nil
}

func (l *AstPrimary) generate(offset int, state *GenState) ([]Instruction, error) {
	return l.literal.generate(offset, state)
}

func (l *AstRange) generate(offset int, state *GenState) ([]Instruction, error) {
	result := MatchRange{
		from: l.from.value,
		to:   l.to.value,
		not:  false,
	}
	return []Instruction{result}, nil
}

func (l *AstString) generate(offset int, state *GenState) ([]Instruction, error) {
	result := MatchLiteral{
		toFind: l.value,
		not:    l.not,
	}
	return []Instruction{result}, nil
}

func (l *AstSubExpr) generate(offset int, state *GenState) ([]Instruction, error) {
	result := []Instruction{}

	loffset := offset
	for _, expr := range l.body {
		expr_insts, gen_error := expr.generate(loffset, state)
		if gen_error != nil {
			return []Instruction{}, gen_error
		}
		loffset = loffset + len(expr_insts)
		result = append(result, expr_insts...)
	}

	return result, nil
}

func (l *AstVariable) generate(offset int, state *GenState) ([]Instruction, error) {
	val, prs := state.variables[l.name]
	if !prs {
		return []Instruction{}, fmt.Errorf("Variable '%s' is not defined", l.name)
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
	return []Instruction{result}, nil
}

func (l *AstCharacterClass) generate(offset int, state *GenState) ([]Instruction, error) {
	result := MatchCharClass{
		class: l.classType,
		not:   l.not,
	}
	return []Instruction{result}, nil
}

func (l *AstString) generateReplace(offset int, state *GenState) ([]RInstruction, error) {
	result := ReplaceString{
		value: l.value,
	}
	return []RInstruction{result}, nil
}

func (l *AstVariable) generateReplace(offset int, state *GenState) ([]RInstruction, error) {
	result := ReplaceVariable{
		name: l.name,
	}
	return []RInstruction{result}, nil
}
