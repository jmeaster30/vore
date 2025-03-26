package libvore

import (
	"fmt"

	"github.com/jmeaster30/vore/libvore/ds"
	"github.com/jmeaster30/vore/libvore/files"
)

type Command interface {
	execute(string, *files.Reader, ReplaceMode) Matches
	String() string
}

type FindCommand struct {
	all  bool
	skip int
	take int
	last int
	body []SearchInstruction
}

func (f FindCommand) String() string {
	return fmt.Sprintf("(find (all %t) (min %d max %d) (last %d) %s)", f.all, f.skip, f.take, f.last, f.body)
}

func findMatches(insts []SearchInstruction, all bool, skip int, take int, last int, filename string, reader *files.Reader) Matches {
	matches := ds.NewQueue[Match]()
	matchNumber := 0
	fileOffset := 0
	lineNumber := 1
	columnNumber := 1

	if reader.Size() == 0 {
		return Matches{}
	}

	for all || matchNumber < skip+take {
		currentState := CreateState(filename, reader, fileOffset, lineNumber, columnNumber)
		for currentState.status == INPROCESS {
			inst := insts[currentState.programCounter]
			currentState = inst.execute(currentState)
			// fmt.Printf("PC: %d INST: %+v STATE: %+v\n", currentState.programCounter, inst, currentState)
			if currentState.status == INPROCESS && currentState.programCounter >= len(insts) {
				currentState.SUCCESS()
			}
		}

		if currentState.status == SUCCESS && len(currentState.currentMatch) != 0 && matchNumber >= skip {
			// fmt.Println("====== SUCCESS ======")
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
			// fmt.Println("====== FAILED  ======")
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

		if fileOffset >= reader.Size() {
			break
		}
	}

	return matches.Contents()
}

func (c FindCommand) execute(filename string, reader *files.Reader, mode ReplaceMode) Matches {
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

func (r ReplaceCommand) String() string {
	return fmt.Sprintf("(replace (%t %d %d %d))", r.all, r.skip, r.take, r.last)
}

func (c ReplaceCommand) execute(filename string, reader *files.Reader, mode ReplaceMode) Matches {
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

	var writer *files.Writer
	replaceReader := reader
	if mode == NEW {
		writer = files.WriterFromFile(filename + ".vored")
	} else if mode == OVERWRITE {
		// If we are overwriting the file we have to load the original
		// into memory since we will be writing over areas of text that
		// we need to read from
		replaceReader = files.ReaderFromFileToMemory(filename)
		writer = files.WriterFromFile(filename)
	} else if mode == NOTHING {
		writer = files.WriterFromMemory()
	}

	lastReaderOffset := 0
	currentWriterOffset := 0
	for i := 0; i < len(replacedMatches); i++ {
		// read from where we left off to the next replacedMatch
		currentReaderLength := replacedMatches[i].Offset.Start - lastReaderOffset
		orig := replaceReader.ReadAt(currentReaderLength, lastReaderOffset)
		writer.WriteAt(currentWriterOffset, orig)
		currentWriterOffset += currentReaderLength
		lastReaderOffset += currentReaderLength

		// write the replacement. We have to update the lastReaderOffset with the part of the string that was matched
		writer.WriteAt(currentWriterOffset, replacedMatches[i].Replacement.GetValueOrDefault(""))
		currentWriterOffset += len(replacedMatches[i].Replacement.GetValueOrDefault(""))
		lastReaderOffset += len(replacedMatches[i].Value)
	}
	if lastReaderOffset < replaceReader.Size() {
		outputValue := replaceReader.ReadAt(reader.Size()-lastReaderOffset, lastReaderOffset)
		writer.WriteAt(currentWriterOffset, outputValue)
	}

	writer.Close()
	replaceReader.Close()

	return replacedMatches
}

type SetCommand struct {
	body SetCommandBody
	id   string
}

func (s SetCommand) String() string {
	return fmt.Sprintf("(set (id %s) (%v)", s.id, s.body)
}

func (c SetCommand) execute(filename string, reader *files.Reader, mode ReplaceMode) Matches {
	return Matches{}
}

type SetCommandBody interface {
	execute(state *GlobalState, id string) *GlobalState
}

type SetCommandExpression struct {
	instructions []SearchInstruction
	validate     []AstProcessStatement
}

func (s SetCommandExpression) execute(state *GlobalState, id string) *GlobalState {
	return state
}

type SetCommandMatches struct {
	command Command
}

func (s SetCommandMatches) execute(state *GlobalState, id string) *GlobalState {
	// TODO run through the command we have and store the matches in variables
	// TODO or do the matches in the generate function
	return state
}

type SetCommandTransform struct {
	statements []AstProcessStatement
}

func (s SetCommandTransform) execute(state *GlobalState, id string) *GlobalState {
	return state
}

type SearchInstruction interface {
	execute(*SearchEngineState) *SearchEngineState
	adjust(offset int, state *GenState) SearchInstruction
	String() string
}

type ReplaceInstruction interface {
	execute(*ReplacerState) *ReplacerState
}

type MatchLiteral struct {
	not      bool
	toFind   string
	caseless bool
}

func (i MatchLiteral) String() string {
	return fmt.Sprintf("(literal (not %t) (caseless %t) '%s')", i.not, i.caseless, i.toFind)
}

func (i MatchLiteral) adjust(offset int, state *GenState) SearchInstruction {
	return i
}

func (i MatchLiteral) execute(current_state *SearchEngineState) *SearchEngineState {
	next_state := current_state.Copy()
	next_state.MATCH(i.toFind, i.not, i.caseless)
	return next_state
}

type MatchCharClass struct {
	not   bool
	class AstCharacterClassType
}

func (i MatchCharClass) String() string {
	return fmt.Sprintf("(class (not %t) %s)", i.not, i.class)
}

func (i MatchCharClass) adjust(offset int, state *GenState) SearchInstruction {
	return i
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
	case ClassWordStart:
		next_state.MATCHWORDSTART(i.not)
	case ClassWordEnd:
		next_state.MATCHWORDEND(i.not)
	case ClassWholeFile:
		next_state.MATCHWHOLEFILE(i.not)
	case ClassWholeLine:
		next_state.MATCHWHOLELINE(i.not)
	case ClassWholeWord:
		next_state.MATCHWHOLEWORD(i.not)
	default:
		panic("Unexpected character class type")
	}
	return next_state
}

type MatchVariable struct {
	name string
}

func (i MatchVariable) String() string {
	return fmt.Sprintf("(var '%s')", i.name)
}

func (i MatchVariable) adjust(offset int, state *GenState) SearchInstruction {
	return i
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

func (i MatchRange) String() string {
	return fmt.Sprintf("(range (not %t) (from '%s') (to '%s'))", i.not, i.from, i.to)
}

func (i MatchRange) adjust(offset int, state *GenState) SearchInstruction {
	return i
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

func (i CallSubroutine) String() string {
	return fmt.Sprintf("(call '%s' %d)", i.name, i.toPC)
}

func (i CallSubroutine) adjust(offset int, state *GenState) SearchInstruction {
	i.toPC += offset
	return i
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

func (i Branch) String() string {
	return fmt.Sprintf("(branch %v)", i.branches)
}

func (i Branch) adjust(offset int, state *GenState) SearchInstruction {
	for idx := range i.branches {
		i.branches[idx] += offset
	}
	return i
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

func (i StartNotIn) String() string {
	return fmt.Sprintf("(startNotIn %d)", i.nextCheckpointPC)
}

func (i StartNotIn) adjust(offset int, state *GenState) SearchInstruction {
	i.nextCheckpointPC += offset
	return i
}

func (i StartNotIn) execute(current_state *SearchEngineState) *SearchEngineState {
	next_state := current_state.Copy()
	next_state.JUMP(i.nextCheckpointPC)
	next_state.CHECKPOINT()
	next_state.JUMP(current_state.programCounter + 1)
	return next_state
}

type FailNotIn struct{}

func (i FailNotIn) String() string {
	return "(failNotIn)"
}

func (i FailNotIn) adjust(offset int, state *GenState) SearchInstruction {
	return i
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

func (i EndNotIn) String() string {
	return fmt.Sprintf("(endNotIn %d)", i.maxSize)
}

func (i EndNotIn) adjust(offset int, state *GenState) SearchInstruction {
	return i
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
	id       int64
	minLoops int
	maxLoops int
	fewest   bool
	exitLoop int
	name     string
}

func (i StartLoop) String() string {
	return fmt.Sprintf("(startLoop '%s' (min %d max %d) (lazy %t) %d %d)", i.name, i.minLoops, i.maxLoops, i.fewest, i.id, i.exitLoop)
}

func (i StartLoop) adjust(offset int, state *GenState) SearchInstruction {
	i.exitLoop += offset
	return i
}

func (i StartLoop) execute(current_state *SearchEngineState) *SearchEngineState {
	next_state := current_state.Copy()

	inited := next_state.INITLOOPSTACK(i.id, i.name)
	if !inited {
		if next_state.CHECKZEROMATCHLOOP() {
			next_state.BACKTRACK()
			return next_state
		}
		next_state.INCLOOPSTACK()
	}
	currentIteration := next_state.GETITERATIONSTEP()

	if currentIteration < i.minLoops {
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
	id        int64
	minLoops  int
	maxLoops  int
	fewest    bool
	startLoop int
	name      string
}

func (i StopLoop) String() string {
	return fmt.Sprintf("(stopLoop '%s' (min %d max %d) (lazy %t) %d %d)", i.name, i.minLoops, i.maxLoops, i.fewest, i.id, i.startLoop)
}

func (i StopLoop) adjust(offset int, state *GenState) SearchInstruction {
	i.startLoop += offset
	return i
}

func (i StopLoop) execute(current_state *SearchEngineState) *SearchEngineState {
	next_state := current_state.Copy()
	next_state.JUMP(i.startLoop)
	return next_state
}

type StartVarDec struct {
	name string
}

func (i StartVarDec) String() string {
	return fmt.Sprintf("(startVarDec '%s')", i.name)
}

func (i StartVarDec) adjust(offset int, state *GenState) SearchInstruction {
	return i
}

func (i StartVarDec) execute(current_state *SearchEngineState) *SearchEngineState {
	next_state := current_state.Copy()
	next_state.STARTVAR(i.name)
	return next_state
}

type EndVarDec struct {
	name string
}

func (i EndVarDec) String() string {
	return fmt.Sprintf("(endVarDec '%s')", i.name)
}

func (i EndVarDec) adjust(offset int, state *GenState) SearchInstruction {
	return i
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

func (i StartSubroutine) String() string {
	return fmt.Sprintf("(startSub '%s' %d %d)", i.name, i.id, i.endOffset)
}

func (i StartSubroutine) adjust(offset int, state *GenState) SearchInstruction {
	i.endOffset += offset
	return i
}

func (i StartSubroutine) execute(current_state *SearchEngineState) *SearchEngineState {
	next_state := current_state.Copy()
	next_state.VALIDATECALL(i.id, i.endOffset+1)
	next_state.NEXT()
	return next_state
}

type EndSubroutine struct {
	name     string
	validate []AstProcessStatement
}

func (i EndSubroutine) String() string {
	return fmt.Sprintf("(endSub '%s')", i.name)
}

func (i EndSubroutine) adjust(offset int, state *GenState) SearchInstruction {
	return i
}

func (i EndSubroutine) execute(current_state *SearchEngineState) *SearchEngineState {
	next_state := current_state.Copy()

	if len(i.validate) == 0 {
		next_state.RETURN()
	} else {
		env := make(map[string]ProcessValue)
		subMatch := current_state.currentMatch[current_state.callStack.Peek().startMatchOffset:]
		env["match"] = ProcessValueString{subMatch}
		env["matchLength"] = ProcessValueNumber{len(subMatch)}
		// TODO add more variables here!

		pstate := ProcessState{
			currentValue: ProcessValueString{""},
			environment:  env,
			status:       NEXT,
		}
		var final_value ProcessValue = ProcessValueBoolean{true}
		for _, stmt := range i.validate {
			pstate = stmt.execute(pstate)
			if pstate.status == RETURNING {
				final_value = pstate.currentValue
				break
			}
		}

		if final_value.getBoolean() {
			next_state.RETURN()
		} else {
			next_state.BACKTRACK()
		}
	}

	return next_state
}

type Jump struct {
	newProgramCounter int
}

func (i Jump) String() string {
	return fmt.Sprintf("(jump %d)", i.newProgramCounter)
}

func (i Jump) adjust(offset int, state *GenState) SearchInstruction {
	i.newProgramCounter += offset
	return i
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

type ReplaceProcess struct {
	process []AstProcessStatement
}

func (i ReplaceProcess) execute(current_state *ReplacerState) *ReplacerState {
	next_state := current_state.Copy()

	// execute AST
	env := make(map[string]ProcessValue)
	keys := current_state.variables.Keys()
	for _, key := range keys {
		value, _ := current_state.variables.Get(key)
		if value.getType() == ValueStringType {
			env[key] = ProcessValueString{value.String().Value}
		}
		// TODO Need to add process hash maps or merge into the main Values
	}

	env["match"] = ProcessValueString{next_state.match.Value}
	env["matchLength"] = ProcessValueNumber{len(next_state.match.Value)}
	env["matchNumber"] = ProcessValueNumber{next_state.match.MatchNumber}

	pstate := ProcessState{
		currentValue: ProcessValueString{""},
		environment:  env,
		status:       NEXT,
	}
	var final_value ProcessValue = ProcessValueBoolean{true}
	for _, stmt := range i.process {
		pstate = stmt.execute(pstate)
		if pstate.status == RETURNING {
			final_value = pstate.currentValue
			break
		}
	}

	next_state.WRITESTRING(final_value.getString())

	next_state.NEXT()
	return next_state
}
