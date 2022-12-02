package libvore

type Command interface {
	execute(string, *VReader, ReplaceMode) Matches
}

type FindCommand struct {
	all  bool
	skip int
	take int
	last int
	body []Instruction
}

func findMatches(insts []Instruction, all bool, skip int, take int, last int, filename string, reader *VReader) Matches {
	matches := NewQueue[Match]()
	matchNumber := 0
	fileOffset := 0
	lineNumber := 1
	columnNumber := 1

	for all || matchNumber < skip+take {
		currentState := CreateState(filename, reader, fileOffset, lineNumber, columnNumber)
		for currentState.status == INPROCESS {
			inst := insts[currentState.programCounter]
			currentState = inst.execute(currentState)

			if currentState.status == INPROCESS && currentState.programCounter >= len(insts) {
				currentState.SUCCESS()
			}
		}

		if currentState.status == SUCCESS && len(currentState.currentMatch) != 0 && matchNumber >= skip {
			foundMatch := currentState.MakeMatch(matchNumber + 1)
			matches.PushBack(foundMatch)
			if last != 0 {
				matches.Limit(last)
			}
			fileOffset = currentState.currentFileOffset
			lineNumber = currentState.currentLineNum
			columnNumber = currentState.currentColumnNum
			matchNumber += 1
		} else {
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

	return matches.Contents()
}

func (c FindCommand) execute(filename string, reader *VReader, mode ReplaceMode) Matches {
	return findMatches(c.body, c.all, c.skip, c.take, c.last, filename, reader)
}

type ReplaceCommand struct {
	all      bool
	skip     int
	take     int
	last     int
	body     []Instruction
	replacer []RInstruction
}

func (c ReplaceCommand) execute(filename string, reader *VReader, mode ReplaceMode) Matches {
	foundMatches := findMatches(c.body, c.all, c.skip, c.take, c.last, filename, reader)

	replacedMatches := Matches{}
	for _, match := range foundMatches {
		current_state := InitReplacerState(match, len(foundMatches))
		for current_state.programCounter < len(c.replacer) {
			inst := c.replacer[current_state.programCounter]
			current_state = inst.execute(current_state)
		}
		replacedMatches = append(replacedMatches, current_state.match)
	}

	var writer *VWriter
	if mode == NEW {
		writer = NewVWriter(filename + ".vored")
	} else if mode == OVERWRITE {
		writer = NewVWriter(filename)
	} else if mode == NOTHING {
		writer = DummyVWriter()
	}

	lastReaderOffset := 0
	currentWriterOffset := 0
	for i := 0; i < len(replacedMatches); i++ {
		currentReaderLength := replacedMatches[i].offset.Start - lastReaderOffset
		orig := reader.ReadAt(currentReaderLength, lastReaderOffset)
		writer.WriteAt(currentWriterOffset, orig)
		currentWriterOffset += currentReaderLength
		lastReaderOffset += currentReaderLength
		writer.WriteAt(currentWriterOffset, replacedMatches[i].replacement)
		currentWriterOffset += len(replacedMatches[i].replacement)
		lastReaderOffset += len(replacedMatches[i].value)
	}
	if lastReaderOffset < reader.size {
		outputValue := reader.ReadAt(reader.size-lastReaderOffset, lastReaderOffset)
		writer.WriteAt(currentWriterOffset, outputValue)
	}

	return replacedMatches
}

type SetCommand struct {
}

func (c SetCommand) execute(filename string, reader *VReader, mode ReplaceMode) Matches {
	return Matches{}
}

type Instruction interface {
	execute(*EngineState) *EngineState
}

type RInstruction interface {
	execute(*ReplacerState) *ReplacerState
}

type MatchLiteral struct {
	not    bool
	toFind string
}

func (i MatchLiteral) execute(current_state *EngineState) *EngineState {
	next_state := current_state.Copy()
	next_state.MATCH(i.toFind, i.not)
	return next_state
}

type MatchCharClass struct {
	not   bool
	class AstCharacterClassType
}

func (i MatchCharClass) execute(current_state *EngineState) *EngineState {
	next_state := current_state.Copy()
	switch i.class {
	case ClassAny:
		next_state.MATCHANY(i.not)
	case ClassWhitespace:
		next_state.MATCHOPTIONS([]string{" ", "\t", "\n", "\r"}, i.not)
	case ClassDigit:
		next_state.MATCHRANGE("0", "9", i.not)
	case ClassUpper:
		next_state.MATCHRANGE("A", "Z", i.not)
	case ClassLower:
		next_state.MATCHRANGE("a", "z", i.not)
	case ClassLetter:
		next_state.MATCHLETTER(i.not)
	case ClassFileStart:
		next_state.MATCHFILESTART(i.not)
	case ClassFileEnd:
		next_state.MATCHFILEEND(i.not)
	case ClassLineStart:
		next_state.MATCHLINESTART(i.not)
	case ClassLineEnd:
		next_state.MATCHLINEEND(i.not)
	default:
		panic("Unexpected character class type")
	}
	return next_state
}

type MatchVariable struct {
	name string
}

func (i MatchVariable) execute(current_state *EngineState) *EngineState {
	next_state := current_state.Copy()
	next_state.MATCHVAR(i.name)
	return next_state
}

type MatchRange struct {
	not  bool
	from string
	to   string
}

func (i MatchRange) execute(current_state *EngineState) *EngineState {
	next_state := current_state.Copy()
	next_state.MATCHRANGE(i.from, i.to, i.not)
	return next_state
}

type CallSubroutine struct {
	name string
	toPC int
}

func (i CallSubroutine) execute(current_state *EngineState) *EngineState {
	next_state := current_state.Copy()
	next_state.CALL(i.toPC, next_state.programCounter+1)
	next_state.JUMP(i.toPC)
	return next_state
}

type Branch struct {
	branches []int
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

func (i StartLoop) execute(current_state *EngineState) *EngineState {
	next_state := current_state.Copy()

	inited := next_state.INITLOOPSTACK(i.id)
	if !inited {
		next_state.INCLOOPSTACK()
	}
	currentIteration := next_state.GETITERATIONSTEP()

	if currentIteration < i.minLoops {
		//fmt.Println("Less than min")
		next_state.NEXT()
	} else if (i.maxLoops == -1 || currentIteration <= i.maxLoops) && i.fewest {
		//fmt.Println("All or less than max FEWEST")
		next_state.NEXT()
		next_state.CHECKPOINT()
		next_state.POPLOOPSTACK()
		next_state.JUMP(i.exitLoop + 1)
	} else if (i.maxLoops == -1 || currentIteration <= i.maxLoops) && !i.fewest {
		//fmt.Println("All or less than max")
		loop_state := next_state.POPLOOPSTACK()
		pc := next_state.GETPC()
		next_state.JUMP(i.exitLoop + 1)
		next_state.CHECKPOINT()
		next_state.PUSHLOOPSTACK(loop_state)
		next_state.JUMP(pc + 1)
	} else {
		//fmt.Printf("FAIL! current: %d min: %d max: %d\n", currentIteration, i.minLoops, i.maxLoops)
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

func (i StopLoop) execute(current_state *EngineState) *EngineState {
	next_state := current_state.Copy()
	next_state.JUMP(i.startLoop)
	return next_state
}

type StartVarDec struct {
	name string
}

func (i StartVarDec) execute(current_state *EngineState) *EngineState {
	next_state := current_state.Copy()
	next_state.STARTVAR(i.name)
	return next_state
}

type EndVarDec struct {
	name string
}

func (i EndVarDec) execute(current_state *EngineState) *EngineState {
	next_state := current_state.Copy()
	next_state.ENDVAR(i.name)
	return next_state
}

type StartSubroutine struct {
	id        int
	name      string
	endOffset int
}

func (i StartSubroutine) execute(current_state *EngineState) *EngineState {
	next_state := current_state.Copy()
	next_state.VALIDATECALL(i.id, i.endOffset+1)
	next_state.NEXT()
	return next_state
}

type EndSubroutine struct {
	name string
}

func (i EndSubroutine) execute(current_state *EngineState) *EngineState {
	next_state := current_state.Copy()
	next_state.RETURN()
	return next_state
}

type Jump struct {
	newProgramCounter int
}

func (i Jump) execute(current_state *EngineState) *EngineState {
	next_state := current_state.Copy()
	next_state.JUMP(i.newProgramCounter)
	return next_state
}

type ReplaceString struct {
	value string
}

func (i ReplaceString) execute(current_state *ReplacerState) *ReplacerState {
	next_state := current_state.Copy()
	next_state.WRITESTRING(i.value)
	next_state.NEXT()
	return next_state
}

type ReplaceVariable struct {
	name string
}

func (i ReplaceVariable) execute(current_state *ReplacerState) *ReplacerState {
	next_state := current_state.Copy()
	next_state.WRITEVAR(i.name)
	next_state.NEXT()
	return next_state
}
