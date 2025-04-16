package bytecode

import (
	"fmt"
	"math/rand"

	"github.com/jmeaster30/vore/libvore/ast"
)

type Bytecode struct {
	Bytecode []Command
}

type GeneratedPattern struct {
	search []SearchInstruction
	// FIXME replace the AstProcessStatement with actual bytecode
	validate []ast.AstProcessStatement
}

// TODO why??
type AstProcessProgram []ast.AstProcessStatement

type GenState struct {
	variables             map[string]int
	globalSubroutines     map[string]GeneratedPattern
	globalVariables       map[string]int
	globalTransformations map[string]AstProcessProgram
}

func GenerateBytecode(a *ast.Ast) (*Bytecode, error) {
	bytecode := []Command{}
	gen_state := &GenState{
		globalSubroutines:     make(map[string]GeneratedPattern),
		globalVariables:       make(map[string]int),
		globalTransformations: make(map[string]AstProcessProgram),
	}
	for _, ast_comm := range a.Commands() {
		byte_comm, gen_error := generateCommand(&ast_comm, gen_state)
		if gen_error != nil {
			return nil, NewGenError(gen_error.Error())
		}
		bytecode = append(bytecode, byte_comm)
	}
	return &Bytecode{bytecode}, nil
}

func generateCommand(com *ast.AstCommand, state *GenState) (Command, error) {
	var icom any = *com
	switch c := icom.(type) {
	case *ast.AstFind:
		return generateFindCommand(c, state)
	case *ast.AstReplace:
		return generateReplaceCommand(c, state)
	case *ast.AstSet:
		return generateSetCommand(c, state)
	}
	return nil, NewGenError(fmt.Sprintf("Unknown command %T", icom))
}

func generateFindCommand(f *ast.AstFind, state *GenState) (Command, error) {
	result := FindCommand{
		All:  f.All,
		Skip: f.Skip,
		Take: f.Take,
		Last: f.Last,
		Body: []SearchInstruction{},
	}

	state.variables = make(map[string]int)

	offset := 0
	for _, expr := range f.Body {
		expr_insts, gen_error := generateSearchInstruction(&expr, offset, state)
		if gen_error != nil {
			return nil, gen_error
		}
		offset = offset + len(expr_insts)
		result.Body = append(result.Body, expr_insts...)
	}

	return result, nil
}

func generateReplaceCommand(r *ast.AstReplace, state *GenState) (Command, error) {
	result := ReplaceCommand{
		All:      r.All,
		Skip:     r.Skip,
		Take:     r.Take,
		Last:     r.Last,
		Body:     []SearchInstruction{},
		Replacer: []ReplaceInstruction{},
	}

	state.variables = make(map[string]int)

	offset := 0
	for _, expr := range r.Body {
		expr_insts, gen_error := generateSearchInstruction(&expr, offset, state)
		if gen_error != nil {
			return nil, gen_error
		}
		offset = offset + len(expr_insts)
		result.Body = append(result.Body, expr_insts...)
	}

	offset = 0
	for _, expr := range r.Result {
		expr_insts, gen_error := generateReplaceInstruction(&expr, offset, state)
		if gen_error != nil {
			return nil, gen_error
		}
		offset = offset + len(expr_insts)
		result.Replacer = append(result.Replacer, expr_insts...)
	}

	return result, nil
}

func generateSetCommand(s *ast.AstSet, state *GenState) (Command, error) {
	state.variables = make(map[string]int)

	body, err := generateSetBody(&s.Body, state, s.Id)
	if err != nil {
		return nil, err
	}
	return SetCommand{
		Body: body,
		Id:   s.Id,
	}, nil
}

func generateSetBody(s *ast.AstSetBody, state *GenState, id string) (SetCommandBody, error) {
	var si any = *s
	switch sb := si.(type) {
	case *ast.AstSetTransform:
		return generateSetTransform(*sb, state, id)
	case *ast.AstSetPattern:
		return generateSetPattern(*sb, state, id)
	case *ast.AstSetMatches:
		return generateSetMatches(*sb, state, id)
	}
	return nil, NewGenError(fmt.Sprintf("Unexpected set body %T", si))
}

func generateSetTransform(s ast.AstSetTransform, state *GenState, id string) (SetCommandBody, error) {
	state.globalTransformations[id] = s.Statements

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
	for _, stmt := range s.Statements {
		info = checkStatement(stmt, info)
		if info.currentType == PTERROR {
			// TODO add more info to the gen error like the statement that failed
			return nil, NewGenError(info.errorMessage)
		}
	}

	return SetCommandTransform{s.Statements}, nil
}

func generateSetPattern(s ast.AstSetPattern, state *GenState, id string) (SetCommandBody, error) {
	state.variables = make(map[string]int)

	searchInstructions := []SearchInstruction{}
	offset := 0
	for _, val := range s.Pattern {
		part, err := generateSearchInstruction(&val, offset, state)
		if err != nil {
			return nil, err
		}

		offset += len(part)
		searchInstructions = append(searchInstructions, part...)
	}

	state.globalSubroutines[id] = GeneratedPattern{searchInstructions, s.Body}

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
	for _, stmt := range s.Body {
		info = checkStatement(stmt, info)
		if info.currentType == PTERROR {
			// TODO add more info to the gen error like the statement that failed
			return nil, NewGenError(info.errorMessage)
		}
	}

	return &SetCommandExpression{
		Instructions: searchInstructions,
		Validate:     s.Body,
	}, nil
}

func generateSetMatches(s ast.AstSetMatches, state *GenState, id string) (SetCommandBody, error) {
	command, err := generateCommand(&s.Command, state)
	if err != nil {
		return nil, err
	}
	return &SetCommandMatches{
		Command: command,
	}, nil
}

func generateSearchInstruction(l *ast.AstExpression, offset int, state *GenState) ([]SearchInstruction, error) {
	var il any = *l
	switch si := il.(type) {
	case *ast.AstLoop:
		return generateLoop(si, offset, state)
	case *ast.AstBranch:
		return generateBranch(si, offset, state)
	case *ast.AstDec:
		return generateVarDec(si, offset, state)
	case *ast.AstSub:
		return generateSubroutine(si, offset, state)
	case *ast.AstList:
		return generateList(si, offset, state)
	case *ast.AstRange:
		return generateRange(si, offset, state)
	case *ast.AstPrimary:
		return generatePrimary(si, offset, state)
	}
	return nil, NewGenError(fmt.Sprintf("Unknown expression of type '%T'", il))
}

func generateLoop(l *ast.AstLoop, offset int, state *GenState) ([]SearchInstruction, error) {
	result := []SearchInstruction{}

	current_offset := offset
	if l.Min > 0 && l.Name == "" {
		for i := 0; i < l.Min; i++ {
			// I kinda hate generating this everytime but I also hate the other way where we have to adjust offset values to keep pointers in the body lined up
			body, gen_error := generateSearchInstruction(&l.Body, current_offset, state)
			if gen_error != nil {
				return []SearchInstruction{}, gen_error
			}
			result = append(result, body...)
			current_offset += len(body)
		}
	}

	if l.Min == l.Max && l.Name == "" {
		return result, nil
	}

	body, gen_error := generateSearchInstruction(&l.Body, current_offset+1, state)
	if gen_error != nil {
		return []SearchInstruction{}, gen_error
	}

	newMin := l.Min
	if l.Min > 0 && l.Name == "" {
		newMin = 0
	}
	newMax := l.Max
	if l.Max > 0 && l.Name == "" {
		newMax = l.Max - l.Min
	}

	id := rand.Int63()

	start := StartLoop{
		Id:       id,
		Name:     l.Name,
		MinLoops: newMin,
		MaxLoops: newMax,
		ExitLoop: current_offset + len(body) + 1,
		Fewest:   l.Fewest,
	}

	stop := StopLoop{
		Id:        id,
		Name:      l.Name,
		MinLoops:  newMin,
		MaxLoops:  newMax,
		StartLoop: current_offset,
		Fewest:    l.Fewest,
	}

	result = append(result, start)
	result = append(result, body...)
	result = append(result, stop)
	return result, nil
}

func generateBranch(l *ast.AstBranch, offset int, state *GenState) ([]SearchInstruction, error) {
	left, left_error := generateLiteral(&l.Left, offset+1, state)
	if left_error != nil {
		return []SearchInstruction{}, left_error
	}
	right, right_error := generateSearchInstruction(&l.Right, offset+2+len(left), state)
	if right_error != nil {
		return []SearchInstruction{}, right_error
	}

	b := Branch{
		Branches: []int{
			offset + 1,
			offset + len(left) + 2,
		},
	}

	insts := []SearchInstruction{b}
	insts = append(insts, left...)
	insts = append(insts, Jump{
		NewProgramCounter: offset + len(left) + len(right) + 3,
	})
	insts = append(insts, right...)
	insts = append(insts, Jump{
		NewProgramCounter: offset + len(left) + len(right) + 3,
	})
	return insts, nil
}

func generateVarDec(l *ast.AstDec, offset int, state *GenState) ([]SearchInstruction, error) {
	insts := []SearchInstruction{}
	// offset
	startVarDec := StartVarDec{
		Name: l.Name,
	}

	bodyinsts, gen_error := generateLiteral(&l.Body, offset+1, state)
	if gen_error != nil {
		return []SearchInstruction{}, gen_error
	}

	endVarDec := EndVarDec{
		Name: l.Name,
	}

	_, prs := state.variables[l.Name]
	if prs {
		return []SearchInstruction{}, NewGenError(fmt.Sprintf("name clash '%s'", l.Name))
	}
	state.variables[l.Name] = -1

	insts = append(insts, startVarDec)
	insts = append(insts, bodyinsts...)
	insts = append(insts, endVarDec)
	return insts, nil
}

func generateSubroutine(l *ast.AstSub, offset int, state *GenState) ([]SearchInstruction, error) {
	insts := []SearchInstruction{}

	_, prs := state.variables[l.Name]
	if prs {
		return []SearchInstruction{}, NewGenError(fmt.Sprintf("name clash '%s'", l.Name))
	}
	state.variables[l.Name] = offset

	bodyinsts := []SearchInstruction{}
	loffset := offset + 1
	for _, expr := range l.Body {
		expr_insts, gen_error := generateSearchInstruction(&expr, loffset, state)
		if gen_error != nil {
			return []SearchInstruction{}, gen_error
		}
		loffset = loffset + len(expr_insts)
		bodyinsts = append(bodyinsts, expr_insts...)
	}

	startVarDec := StartSubroutine{
		Id:        offset,
		Name:      l.Name,
		EndOffset: loffset,
	}

	endVarDec := EndSubroutine{
		Name: l.Name,
	}

	insts = append(insts, startVarDec)
	insts = append(insts, bodyinsts...)
	insts = append(insts, endVarDec)
	return insts, nil
}

func generateList(l *ast.AstList, offset int, state *GenState) ([]SearchInstruction, error) {
	if l.Not {
		return generate_not(l, offset, state)
	} else {
		return generate_not_not(l, offset, state)
	}
}

func generateListable(l *ast.AstListable, offset int, state *GenState) ([]SearchInstruction, error) {
	var il any = *l
	switch li := il.(type) {
	case *ast.AstString:
		return generateString(li, offset, state)
	case *ast.AstCharacterClass:
		return generateCharacterClass(li, offset, state)
	case *ast.AstRange:
		return generateRange(li, offset, state)
	}
	return nil, NewGenError(fmt.Sprintf("Unknown listable '%T'", il))
}

func generate_not_not(l *ast.AstList, offset int, state *GenState) ([]SearchInstruction, error) {
	b := Branch{
		Branches: []int{},
	}

	pc := offset + 1
	branches := [][]SearchInstruction{}
	for _, elem := range l.Contents {
		branch_insts, gen_error := generateListable(&elem, pc, state)
		if gen_error != nil {
			return []SearchInstruction{}, gen_error
		}
		b.Branches = append(b.Branches, pc)
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
			NewProgramCounter: end,
		})
	}

	return insts, nil
}

func generate_not(l *ast.AstList, offset int, state *GenState) ([]SearchInstruction, error) {
	insts := []SearchInstruction{}

	pc := offset
	for _, item := range l.Contents {
		item_insts, err := generateListable(&item, pc+1, state)
		if err != nil {
			return []SearchInstruction{}, err
		}

		insts = append(insts, StartNotIn{
			NextCheckpointPC: pc + len(item_insts) + 2,
		})
		insts = append(insts, item_insts...)
		insts = append(insts, FailNotIn{})
		pc = pc + len(item_insts) + 2
	}

	insts = append(insts, EndNotIn{
		MaxSize: l.GetMaxSize(),
	})
	return insts, nil
}

func generateLiteral(l *ast.AstLiteral, offset int, state *GenState) ([]SearchInstruction, error) {
	var il any = *l
	switch ll := il.(type) {
	case *ast.AstString:
		return generateString(ll, offset, state)
	case *ast.AstSubExpr:
		return generateSubExpression(ll, offset, state)
	case *ast.AstVariable:
		return generateVariable(ll, offset, state)
	case *ast.AstCharacterClass:
		return generateCharacterClass(ll, offset, state)
	}
	return nil, NewGenError(fmt.Sprintf("Unkonwn literal type '%T'", il))
}

func generatePrimary(l *ast.AstPrimary, offset int, state *GenState) ([]SearchInstruction, error) {
	return generateLiteral(&l.Literal, offset, state)
}

func generateRange(l *ast.AstRange, offset int, state *GenState) ([]SearchInstruction, error) {
	result := MatchRange{
		From: l.From.Value,
		To:   l.To.Value,
		Not:  false,
	}
	return []SearchInstruction{result}, nil
}

func generateString(l *ast.AstString, offset int, state *GenState) ([]SearchInstruction, error) {
	result := MatchLiteral{
		ToFind:   l.Value,
		Not:      l.Not,
		Caseless: l.Caseless,
	}
	return []SearchInstruction{result}, nil
}

func generateSubExpression(l *ast.AstSubExpr, offset int, state *GenState) ([]SearchInstruction, error) {
	result := []SearchInstruction{}

	loffset := offset
	for _, expr := range l.Body {
		expr_insts, gen_error := generateSearchInstruction(&expr, loffset, state)
		if gen_error != nil {
			return []SearchInstruction{}, gen_error
		}
		loffset = loffset + len(expr_insts)
		result = append(result, expr_insts...)
	}

	return result, nil
}

func generateVariable(l *ast.AstVariable, offset int, state *GenState) ([]SearchInstruction, error) {
	val, prs := state.variables[l.Name]
	if !prs {
		// we don't have a variable check the subroutines
		globalSub, globalPrs := state.globalSubroutines[l.Name]
		if !globalPrs {
			return []SearchInstruction{}, NewGenError(fmt.Sprintf("identifier '%s' is not defined", l.Name))
		}

		state.variables[l.Name] = offset

		bodyinsts := []SearchInstruction{}
		loffset := offset + 1
		for _, expr := range globalSub.search {
			inst := expr.adjust(offset+1, state)
			loffset += 1
			bodyinsts = append(bodyinsts, inst)
		}

		jumpToSubroutine := StartSubroutine{
			Id:        offset,
			Name:      l.Name,
			EndOffset: loffset,
		}

		returnFromSubroutine := EndSubroutine{
			Name:     l.Name,
			Validate: globalSub.validate,
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
			Name: l.Name,
		}
	} else {
		result = CallSubroutine{
			Name: l.Name,
			ToPC: val,
		}
	}
	return []SearchInstruction{result}, nil
}

func generateCharacterClass(l *ast.AstCharacterClass, offset int, state *GenState) ([]SearchInstruction, error) {
	result := MatchCharClass{
		Class: l.ClassType,
		Not:   l.Not,
	}
	return []SearchInstruction{result}, nil
}

func generateReplaceInstruction(l *ast.AstAtom, offset int, state *GenState) ([]ReplaceInstruction, error) {
	var il any = *l
	switch ri := il.(type) {
	case *ast.AstString:
		return generateReplaceString(ri, offset, state)
	case *ast.AstVariable:
		return generateReplaceVariable(ri, offset, state)
	}
	return nil, NewGenError(fmt.Sprintf("Unknown replace instruction '%T'", il))
}

func generateReplaceString(l *ast.AstString, offset int, state *GenState) ([]ReplaceInstruction, error) {
	result := ReplaceString{
		Value: l.Value,
	}
	return []ReplaceInstruction{result}, nil
}

func generateReplaceVariable(l *ast.AstVariable, offset int, state *GenState) ([]ReplaceInstruction, error) {
	transform, prs := state.globalTransformations[l.Name]
	var result ReplaceInstruction
	if prs {
		result = ReplaceProcess{
			Process: transform,
		}
	} else {
		result = ReplaceVariable{
			Name: l.Name,
		}
	}
	return []ReplaceInstruction{result}, nil
}
