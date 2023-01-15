package libvore

import (
	"fmt"
	"math/rand"
)

type GeneratedPattern struct {
	search   []SearchInstruction
	validate []AstProcessStatement
}

type AstProcessProgram []AstProcessStatement

type GenState struct {
	variables             map[string]int
	globalSubroutines     map[string]GeneratedPattern
	globalVariables       map[string]int
	globalTransformations map[string]AstProcessProgram
}

func (f *AstFind) generate(state *GenState) (Command, error) {
	result := FindCommand{
		all:  f.all,
		skip: f.skip,
		take: f.take,
		last: f.last,
		body: []SearchInstruction{},
	}

	state.variables = make(map[string]int)

	offset := 0
	for _, expr := range f.body {
		expr_insts, gen_error := expr.generate(offset, state)
		if gen_error != nil {
			return nil, gen_error
		}
		offset = offset + len(expr_insts)
		result.body = append(result.body, expr_insts...)
	}

	return result, nil
}

func (r *AstReplace) generate(state *GenState) (Command, error) {
	result := ReplaceCommand{
		all:      r.all,
		skip:     r.skip,
		take:     r.take,
		last:     r.last,
		body:     []SearchInstruction{},
		replacer: []ReplaceInstruction{},
	}

	state.variables = make(map[string]int)

	offset := 0
	for _, expr := range r.body {
		expr_insts, gen_error := expr.generate(offset, state)
		if gen_error != nil {
			return nil, gen_error
		}
		offset = offset + len(expr_insts)
		result.body = append(result.body, expr_insts...)
	}

	offset = 0
	for _, expr := range r.result {
		expr_insts, gen_error := expr.generateReplace(offset, state)
		if gen_error != nil {
			return nil, gen_error
		}
		offset = offset + len(expr_insts)
		result.replacer = append(result.replacer, expr_insts...)
	}

	return result, nil
}

func (s *AstSet) generate(state *GenState) (Command, error) {
	state.variables = make(map[string]int)

	body, err := s.body.generate(state, s.id)
	if err != nil {
		return nil, err
	}
	return SetCommand{
		body: body,
		id:   s.id,
	}, nil
}

func (s AstSetTransform) generate(state *GenState, id string) (SetCommandBody, error) {
	state.globalTransformations[id] = s.statements

	// semantic check
	env := make(map[string]ProcessType)
	env["match"] = PTSTRING
	env["matchLength"] = PTNUMBER
	// TODO pull variables from search pattern and add them here

	info := ProcessTypeInfo{
		currentType:  PTOK,
		errorMessage: "",
		context:      TRANSFORMATION,
		environment:  env,
		inLoop:       false,
	}
	for _, stmt := range s.statements {
		info = stmt.check(info)
		if info.currentType == PTERROR {
			return nil, fmt.Errorf("%s", info.errorMessage)
		}
	}

	return SetCommandTransform{s.statements}, nil
}

func (s AstSetPattern) generate(state *GenState, id string) (SetCommandBody, error) {
	state.variables = make(map[string]int)

	searchInstructions, err := s.pattern.generate(0, state)
	if err != nil {
		return nil, err
	}

	state.globalSubroutines[id] = GeneratedPattern{searchInstructions, s.body}

	// semantic check
	env := make(map[string]ProcessType)
	env["match"] = PTSTRING
	env["matchLength"] = PTNUMBER
	// TODO pull variables from search pattern and add them here

	info := ProcessTypeInfo{
		currentType:  PTOK,
		errorMessage: "",
		context:      PREDICATE,
		environment:  env,
		inLoop:       false,
	}
	for _, stmt := range s.body {
		info = stmt.check(info)
		if info.currentType == PTERROR {
			return nil, fmt.Errorf("%s", info.errorMessage)
		}
	}

	return &SetCommandExpression{
		instructions: searchInstructions,
		validate:     s.body,
	}, nil
}

func (s AstSetMatches) generate(state *GenState, id string) (SetCommandBody, error) {
	command, err := s.command.generate(state)
	if err != nil {
		return nil, err
	}
	return &SetCommandMatches{
		command: command,
	}, nil
}

func (l *AstLoop) generate(offset int, state *GenState) ([]SearchInstruction, error) {
	body, gen_error := l.body.generate(offset+1, state)
	if gen_error != nil {
		return []SearchInstruction{}, gen_error
	}

	id := rand.Int63()

	start := StartLoop{
		id:       id,
		minLoops: l.min,
		maxLoops: l.max,
		exitLoop: offset + len(body) + 1,
		fewest:   l.fewest,
	}

	stop := StopLoop{
		id:        id,
		minLoops:  l.min,
		maxLoops:  l.max,
		startLoop: offset,
		fewest:    l.fewest,
	}

	result := []SearchInstruction{start}
	result = append(result, body...)
	result = append(result, stop)
	return result, nil
}

func (l *AstBranch) generate(offset int, state *GenState) ([]SearchInstruction, error) {

	left, left_error := l.left.generate(offset+1, state)
	if left_error != nil {
		return []SearchInstruction{}, left_error
	}
	right, right_error := l.right.generate(offset+2+len(left), state)
	if right_error != nil {
		return []SearchInstruction{}, right_error
	}

	b := Branch{
		branches: []int{
			offset + 1,
			offset + len(left) + 2,
		},
	}

	insts := []SearchInstruction{b}
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

func (l *AstDec) generate(offset int, state *GenState) ([]SearchInstruction, error) {
	insts := []SearchInstruction{}
	// offset
	startVarDec := StartVarDec{
		name: l.name,
	}

	bodyinsts, gen_error := l.body.generate(offset+1, state)
	if gen_error != nil {
		return []SearchInstruction{}, gen_error
	}

	endVarDec := EndVarDec{
		name: l.name,
	}

	_, prs := state.variables[l.name]
	if prs {
		return []SearchInstruction{}, fmt.Errorf("name clash '%s'", l.name)
	}
	state.variables[l.name] = -1

	insts = append(insts, startVarDec)
	insts = append(insts, bodyinsts...)
	insts = append(insts, endVarDec)
	return insts, nil
}

func (l *AstSub) generate(offset int, state *GenState) ([]SearchInstruction, error) {
	insts := []SearchInstruction{}

	_, prs := state.variables[l.name]
	if prs {
		return []SearchInstruction{}, fmt.Errorf("name clash '%s'", l.name)
	}
	state.variables[l.name] = offset

	bodyinsts := []SearchInstruction{}
	loffset := offset + 1
	for _, expr := range l.body {
		expr_insts, gen_error := expr.generate(loffset, state)
		if gen_error != nil {
			return []SearchInstruction{}, gen_error
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

func (l *AstList) generate(offset int, state *GenState) ([]SearchInstruction, error) {
	if l.not {
		return l.generate_not(offset, state)
	} else {
		return l.generate_not_not(offset, state)
	}
}

func (l *AstList) generate_not_not(offset int, state *GenState) ([]SearchInstruction, error) {
	b := Branch{
		branches: []int{},
	}

	pc := offset + 1
	branches := [][]SearchInstruction{}
	for _, elem := range l.contents {
		branch_insts, gen_error := elem.generate(pc, state)
		if gen_error != nil {
			return []SearchInstruction{}, gen_error
		}
		b.branches = append(b.branches, pc)
		pc += len(branch_insts) + 1
		branches = append(branches, branch_insts)
	}

	end := offset + 1
	for _, instss := range branches {
		end += len(instss) + 1
	}

	insts := []SearchInstruction{b}
	for _, instss := range branches {
		insts = append(insts, instss...)
		insts = append(insts, Jump{
			newProgramCounter: end,
		})
	}

	return insts, nil
}

func (l *AstList) generate_not(offset int, state *GenState) ([]SearchInstruction, error) {
	insts := []SearchInstruction{}

	pc := offset
	for _, item := range l.contents {
		item_insts, err := item.generate(pc+1, state)
		if err != nil {
			return []SearchInstruction{}, err
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

func (l *AstPrimary) generate(offset int, state *GenState) ([]SearchInstruction, error) {
	return l.literal.generate(offset, state)
}

func (l *AstRange) generate(offset int, state *GenState) ([]SearchInstruction, error) {
	result := MatchRange{
		from: l.from.value,
		to:   l.to.value,
		not:  false,
	}
	return []SearchInstruction{result}, nil
}

func (l *AstString) generate(offset int, state *GenState) ([]SearchInstruction, error) {
	result := MatchLiteral{
		toFind: l.value,
		not:    l.not,
	}
	return []SearchInstruction{result}, nil
}

func (l *AstSubExpr) generate(offset int, state *GenState) ([]SearchInstruction, error) {
	result := []SearchInstruction{}

	loffset := offset
	for _, expr := range l.body {
		expr_insts, gen_error := expr.generate(loffset, state)
		if gen_error != nil {
			return []SearchInstruction{}, gen_error
		}
		loffset = loffset + len(expr_insts)
		result = append(result, expr_insts...)
	}

	return result, nil
}

func (l *AstVariable) generate(offset int, state *GenState) ([]SearchInstruction, error) {
	val, prs := state.variables[l.name]
	if !prs {
		// we don't have a variable check the subroutines
		globalSub, globalPrs := state.globalSubroutines[l.name]
		if !globalPrs {
			return []SearchInstruction{}, fmt.Errorf("identifier '%s' is not defined", l.name)
		}

		state.variables[l.name] = offset

		bodyinsts := []SearchInstruction{}
		loffset := offset + 1
		for _, expr := range globalSub.search {
			inst := expr.adjust(offset+1, state)
			loffset += 1
			bodyinsts = append(bodyinsts, inst)
		}

		jumpToSubroutine := StartSubroutine{
			id:        offset,
			name:      l.name,
			endOffset: loffset,
		}

		returnFromSubroutine := EndSubroutine{
			name:     l.name,
			validate: globalSub.validate,
		}

		insts := []SearchInstruction{}
		insts = append(insts, jumpToSubroutine)
		insts = append(insts, bodyinsts...)
		insts = append(insts, returnFromSubroutine)
		return insts, nil
	}
	var result SearchInstruction
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
	return []SearchInstruction{result}, nil
}

func (l *AstCharacterClass) generate(offset int, state *GenState) ([]SearchInstruction, error) {
	result := MatchCharClass{
		class: l.classType,
		not:   l.not,
	}
	return []SearchInstruction{result}, nil
}

func (l *AstString) generateReplace(offset int, state *GenState) ([]ReplaceInstruction, error) {
	result := ReplaceString{
		value: l.value,
	}
	return []ReplaceInstruction{result}, nil
}

func (l *AstVariable) generateReplace(offset int, state *GenState) ([]ReplaceInstruction, error) {
	transform, prs := state.globalTransformations[l.name]
	var result ReplaceInstruction
	if prs {
		result = ReplaceProcess{
			process: transform,
		}
	} else {
		result = ReplaceVariable{
			name: l.name,
		}
	}
	return []ReplaceInstruction{result}, nil
}
