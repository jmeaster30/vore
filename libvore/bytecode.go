package libvore

type Command interface {
	execute(string, *VReader, ReplaceMode) Matches
}

type FindCommand struct {
	all  bool
	skip int
	take int
	last int
	body []SearchInstruction
}

func findMatches(insts []SearchInstruction, all bool, skip int, take int, last int, filename string, reader *VReader) Matches {
	matches := NewQueue[Match]()
	matchNumber := 0
	fileOffset := 0
	lineNumber := 1
	columnNumber := 1

	//for i, inst := range insts {
	//	fmt.Printf("[%d] %v\n", i, inst)
	//}

	for all || matchNumber < skip+take {
		currentState := CreateState(filename, reader, fileOffset, lineNumber, columnNumber)
		for currentState.status == INPROCESS {
			inst := insts[currentState.programCounter]
			currentState = inst.execute(currentState)
			//fmt.Printf("pc: %d inst: %v\n", currentState.programCounter, inst)

			if currentState.status == INPROCESS && currentState.programCounter >= len(insts) {
				currentState.SUCCESS()
			}
		}

		if currentState.status == SUCCESS && len(currentState.currentMatch) != 0 && matchNumber >= skip {
			//fmt.Println("====== SUCCESS ======")
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
			//fmt.Println("====== FAILED  ======")
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
	body     []SearchInstruction
	replacer []ReplaceInstruction
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
	body         SetCommandBody
	isSubroutine bool
	isMatches    bool
	id           string
}

func (c SetCommand) execute(filename string, reader *VReader, mode ReplaceMode) Matches {
	return Matches{}
}

type SetCommandBody interface {
	execute(state *GlobalState, id string) *GlobalState
}

type SetCommandExpression struct {
	instructions []SearchInstruction
}

func (s SetCommandExpression) execute(state *GlobalState, id string) *GlobalState {
	return state
}

type SetCommandMatches struct {
	command Command
}

func (s SetCommandMatches) execute(state *GlobalState, id string) *GlobalState {
	// run through the command we have and store the matches in variables
	return state
}

type SearchInstruction interface {
	execute(*SearchEngineState) *SearchEngineState
	adjust(offset int, state *GenState) (SearchInstruction, int)
}

type ReplaceInstruction interface {
	execute(*ReplacerState) *ReplacerState
}

type MatchLiteral struct {
	not    bool
	toFind string
}

func (i MatchLiteral) adjust(offset int, state *GenState) (SearchInstruction, int) {
	return i, state.loopId
}
func (i MatchLiteral) execute(current_state *SearchEngineState) *SearchEngineState {
	next_state := current_state.Copy()
	next_state.MATCH(i.toFind, i.not)
	return next_state
}

type MatchCharClass struct {
	not   bool
	class AstCharacterClassType
}

func (i MatchCharClass) adjust(offset int, state *GenState) (SearchInstruction, int) {
	return i, state.loopId
}
func (i MatchCharClass) execute(current_state *SearchEngineState) *SearchEngineState {
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

func (i MatchVariable) adjust(offset int, state *GenState) (SearchInstruction, int) {
	return i, state.loopId
}
func (i MatchVariable) execute(current_state *SearchEngineState) *SearchEngineState {
	next_state := current_state.Copy()
	next_state.MATCHVAR(i.name)
	return next_state
}

type MatchRange struct {
	not  bool
	from string
	to   string
}

func (i MatchRange) adjust(offset int, state *GenState) (SearchInstruction, int) {
	return i, state.loopId
}
func (i MatchRange) execute(current_state *SearchEngineState) *SearchEngineState {
	next_state := current_state.Copy()
	next_state.MATCHRANGE(i.from, i.to, i.not)
	return next_state
}

type CallSubroutine struct {
	name string
	toPC int
}

func (i CallSubroutine) adjust(offset int, state *GenState) (SearchInstruction, int) {
	i.toPC += offset
	return i, state.loopId
}
func (i CallSubroutine) execute(current_state *SearchEngineState) *SearchEngineState {
	next_state := current_state.Copy()
	next_state.CALL(i.toPC, next_state.programCounter+1)
	next_state.JUMP(i.toPC)
	return next_state
}

type Branch struct {
	branches []int
}

func (i Branch) adjust(offset int, state *GenState) (SearchInstruction, int) {
	for idx := range i.branches {
		i.branches[idx] += offset
	}
	return i, state.loopId
}
func (i Branch) execute(current_state *SearchEngineState) *SearchEngineState {
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

type StartNotIn struct {
	nextCheckpointPC int
}

func (i StartNotIn) adjust(offset int, state *GenState) (SearchInstruction, int) {
	i.nextCheckpointPC += offset
	return i, state.loopId
}
func (i StartNotIn) execute(current_state *SearchEngineState) *SearchEngineState {
	next_state := current_state.Copy()
	next_state.JUMP(i.nextCheckpointPC)
	next_state.CHECKPOINT()
	next_state.JUMP(current_state.programCounter + 1)
	return next_state
}

type FailNotIn struct{}

func (i FailNotIn) adjust(offset int, state *GenState) (SearchInstruction, int) {
	return i, state.loopId
}
func (i FailNotIn) execute(current_state *SearchEngineState) *SearchEngineState {
	next_state := current_state.Copy()
	next_state.BACKTRACK()
	next_state.BACKTRACK()
	return next_state
}

type EndNotIn struct {
	maxSize int
}

func (i EndNotIn) adjust(offset int, state *GenState) (SearchInstruction, int) {
	return i, state.loopId
}
func (i EndNotIn) execute(current_state *SearchEngineState) *SearchEngineState {
	next_state := current_state.Copy()
	// TODO this should actually let the rest of the expression backtrack from max size to min size (could just be to 1 since things less than the min are not in)
	cfo := next_state.currentFileOffset
	next_state.CONSUME(i.maxSize)
	// FIXME: This was added to make it so we don't have an infinite loop when using "not in" in an un-bounded loop
	//        I think a better fix would be to come up with a different way to handle the end of the file
	if cfo == next_state.currentFileOffset {
		next_state.BACKTRACK()
	} else {
		next_state.NEXT()
	}
	return next_state
}

type StartLoop struct {
	id       int
	minLoops int
	maxLoops int
	fewest   bool
	exitLoop int
}

func (i StartLoop) adjust(offset int, state *GenState) (SearchInstruction, int) {
	i.exitLoop += offset
	i.id += state.loopId
	return i, state.loopId
}
func (i StartLoop) execute(current_state *SearchEngineState) *SearchEngineState {
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

func (i StopLoop) adjust(offset int, state *GenState) (SearchInstruction, int) {
	i.id += state.loopId
	i.startLoop += offset
	return i, state.loopId
}
func (i StopLoop) execute(current_state *SearchEngineState) *SearchEngineState {
	next_state := current_state.Copy()
	next_state.JUMP(i.startLoop)
	return next_state
}

type StartVarDec struct {
	name string
}

func (i StartVarDec) adjust(offset int, state *GenState) (SearchInstruction, int) {
	return i, state.loopId
}
func (i StartVarDec) execute(current_state *SearchEngineState) *SearchEngineState {
	next_state := current_state.Copy()
	next_state.STARTVAR(i.name)
	return next_state
}

type EndVarDec struct {
	name string
}

func (i EndVarDec) adjust(offset int, state *GenState) (SearchInstruction, int) {
	return i, state.loopId
}
func (i EndVarDec) execute(current_state *SearchEngineState) *SearchEngineState {
	next_state := current_state.Copy()
	next_state.ENDVAR(i.name)
	return next_state
}

type StartSubroutine struct {
	id        int
	name      string
	endOffset int
}

func (i StartSubroutine) adjust(offset int, state *GenState) (SearchInstruction, int) {
	i.endOffset += offset
	return i, state.loopId
}
func (i StartSubroutine) execute(current_state *SearchEngineState) *SearchEngineState {
	next_state := current_state.Copy()
	next_state.VALIDATECALL(i.id, i.endOffset+1)
	next_state.NEXT()
	return next_state
}

type EndSubroutine struct {
	name string
}

func (i EndSubroutine) adjust(offset int, state *GenState) (SearchInstruction, int) {
	return i, state.loopId
}
func (i EndSubroutine) execute(current_state *SearchEngineState) *SearchEngineState {
	next_state := current_state.Copy()
	next_state.RETURN()
	return next_state
}

type Jump struct {
	newProgramCounter int
}

func (i Jump) adjust(offset int, state *GenState) (SearchInstruction, int) {
	i.newProgramCounter += offset
	return i, state.loopId
}
func (i Jump) execute(current_state *SearchEngineState) *SearchEngineState {
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
