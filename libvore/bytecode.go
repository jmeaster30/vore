package libvore

import (
	"fmt"
)

type Command interface {
	execute(string, *VReader) []Match
	print()
}

type FindCommand struct {
	all  bool
	skip int
	take int
	body []Instruction
}

func (c FindCommand) print() {
	fmt.Println("FIND COMMAND")
	for _, inst := range c.body {
		inst.print()
	}
}

func (c FindCommand) execute(filename string, reader *VReader) []Match {
	matches := []Match{}
	matchNumber := 0
	fileOffset := 0
	lineNumber := 1
	columnNumber := 1
	//fmt.Printf("searching %s %d\n", filename, filesize)
	//fmt.Printf("%t %d %d\n", c.all, c.skip, c.take)
	for c.all || matchNumber < c.skip+c.take {
		currentState := CreateState(filename, reader, fileOffset, lineNumber, columnNumber)
		for currentState.status == INPROCESS {
			inst := c.body[currentState.programCounter]
			//fmt.Printf("inst: %d\n", currentState.programCounter)
			//inst.print()
			//fmt.Printf("BEFORE = PC: %d\tBTK: %d\n", currentState.programCounter, currentState.backtrack.Size())
			currentState = inst.execute(currentState)

			//fmt.Printf("AFTER  = PC: %d\tBTK: %d\n", currentState.programCounter, currentState.backtrack.Size())
			if currentState.status == INPROCESS && currentState.programCounter >= len(c.body) {
				currentState.SUCCESS()
			}
		}

		if currentState.status == SUCCESS && len(currentState.currentMatch) != 0 && matchNumber >= c.skip {
			//fmt.Println("SUCCESS ====================================================")
			foundMatch := currentState.MakeMatch(matchNumber + 1)
			matches = append(matches, foundMatch)
			fileOffset = currentState.currentFileOffset
			lineNumber = currentState.currentLineNum
			columnNumber = currentState.currentColumnNum
			matchNumber += 1
		} else {
			//fmt.Println("FAIL =======================================================")
			if currentState.status == SUCCESS && len(currentState.currentMatch) != 0 {
				matchNumber += 1
			}
			skipC := reader.ReadAt(1, fileOffset)
			if len(skipC) != 1 {
				panic("WOW THAT IS NOT GOOD :(")
			}
			fileOffset += 1
			columnNumber += 1
			if rune(skipC[0]) == rune('\n') {
				lineNumber += 1
				columnNumber = 1
			}
		}

		if fileOffset >= reader.size {
			break
		}
	}

	return matches
}

type Instruction interface {
	execute(*EngineState) *EngineState
	print()
}

type MatchLiteral struct {
	toFind string
}

func (i MatchLiteral) print() {
	fmt.Printf("MATCH LITERAL '%s'\n", i.toFind)
}

func (i MatchLiteral) execute(current_state *EngineState) *EngineState {
	next_state := current_state.Copy()
	next_state.MATCH(i.toFind)
	return next_state
}

type MatchCharClass struct {
	class AstCharacterClassType
}

func (i MatchCharClass) print() {
	fmt.Println("MATCH CLASS")
}

func (i MatchCharClass) execute(current_state *EngineState) *EngineState {
	next_state := current_state.Copy()
	switch i.class {
	case ClassAny:
		next_state.MATCHANY()
	case ClassWhitespace:
		next_state.MATCHOPTIONS([]string{" ", "\t", "\n", "\r"})
	case ClassDigit:
		next_state.MATCHRANGE("0", "9")
	case ClassUpper:
		next_state.MATCHRANGE("A", "Z")
	case ClassLower:
		next_state.MATCHRANGE("a", "z")
	case ClassLetter:
		next_state.MATCHLETTER()
	case ClassFileStart:
		next_state.MATCHFILESTART()
	case ClassFileEnd:
		next_state.MATCHFILEEND()
	case ClassLineStart:
		next_state.MATCHLINESTART()
	case ClassLineEnd:
		next_state.MATCHLINEEND()
	default:
		panic("Unexpected character class type")
	}
	return next_state
}

type MatchVariable struct {
	name string
}

func (i MatchVariable) print() {
	fmt.Printf("MATCH VAR '%s'\n", i.name)
}

func (i MatchVariable) execute(current_state *EngineState) *EngineState {
	next_state := current_state.Copy()
	next_state.MATCHVAR(i.name)
	return next_state
}

type MatchRange struct {
	from string
	to   string
}

func (i MatchRange) print() {
	fmt.Printf("MATCH RANGE '%s' to '%s'\n", i.from, i.to)
}

func (i MatchRange) execute(current_state *EngineState) *EngineState {
	next_state := current_state.Copy()
	next_state.MATCHRANGE(i.from, i.to)
	return next_state
}

type CallSubroutine struct {
	name string
}

func (i CallSubroutine) print() {
	fmt.Println("CALL SUB")
}

func (i CallSubroutine) execute(current_state *EngineState) *EngineState {
	panic("unimplemented subroutines")
}

type Branch struct {
	branches []int
}

func (i Branch) print() {
	fmt.Print("BRANCH ")
	for _, b := range i.branches {
		fmt.Printf("%d\t", b)
	}
	fmt.Println()
}

func (i Branch) execute(current_state *EngineState) *EngineState {
	next_state := current_state.Copy()
	flipped := []int{}
	for k := range i.branches {
		flipped = append(flipped, i.branches[len(i.branches)-1-k])
	}

	for _, f := range flipped[:len(flipped)-1] {
		next_state.JUMP(f)
		next_state.CHECKPOINT()
	}

	next_state.JUMP(i.branches[0])
	return next_state
}

type StartLoop struct {
	id       int
	minLoops int
	maxLoops int
	fewest   bool
	exitLoop int
}

func (i StartLoop) print() {
	fmt.Printf("START LOOP %d %d %t %d\n", i.minLoops, i.maxLoops, i.fewest, i.exitLoop)
}

func (i StartLoop) execute(current_state *EngineState) *EngineState {
	next_state := current_state.Copy()

	inited := next_state.INITLOOPSTACK(i.id)
	if !inited {
		next_state.INCLOOPSTACK()
	}
	currentIteration := next_state.GETITERATIONSTEP()

	if currentIteration < i.minLoops-1 {
		next_state.NEXT()
	} else if (i.maxLoops == -1 || currentIteration <= i.maxLoops) && i.fewest {
		next_state.NEXT()
		next_state.CHECKPOINT()
		next_state.POPLOOPSTACK()
		next_state.JUMP(i.exitLoop + 1)
	} else if (i.maxLoops == -1 || currentIteration <= i.maxLoops) && !i.fewest {
		loop_state := next_state.POPLOOPSTACK()
		pc := next_state.GETPC()
		next_state.JUMP(i.exitLoop + 1)
		next_state.CHECKPOINT()
		next_state.PUSHLOOPSTACK(loop_state)
		next_state.JUMP(pc + 1)
	} else {
		next_state.BACKTRACK()
	}

	return next_state
}

type StopLoop struct {
	id        int
	minLoops  int
	maxLoops  int
	fewest    bool
	startLoop int
}

func (i StopLoop) print() {
	fmt.Printf("END LOOP %d %d %t %d\n", i.minLoops, i.maxLoops, i.fewest, i.startLoop)
}

func (i StopLoop) execute(current_state *EngineState) *EngineState {
	next_state := current_state.Copy()
	next_state.JUMP(i.startLoop)
	return next_state
}

type StartVarDec struct {
	name string
}

func (i StartVarDec) print() {
	fmt.Printf("START VARDEC '%s'\n", i.name)
}

func (i StartVarDec) execute(current_state *EngineState) *EngineState {
	next_state := current_state.Copy()
	next_state.STARTVAR(i.name)
	return next_state
}

type EndVarDec struct {
	name string
}

func (i EndVarDec) print() {
	fmt.Printf("END VARDEC '%s'\n", i.name)
}

func (i EndVarDec) execute(current_state *EngineState) *EngineState {
	next_state := current_state.Copy()
	next_state.ENDVAR(i.name)
	return next_state
}

type Jump struct {
	newProgramCounter int
}

func (i Jump) print() {
	fmt.Printf("JUMP %d\n", i.newProgramCounter)
}

func (i Jump) execute(current_state *EngineState) *EngineState {
	next_state := current_state.Copy()
	next_state.JUMP(i.newProgramCounter)
	return next_state
}
