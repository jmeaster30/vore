package bytecode

import (
	"fmt"

	"github.com/jmeaster30/vore/libvore/ast"
)

type Command interface {
	// execute(string, *files.Reader, ReplaceMode) Matches
	IsCommand()
	String() string
}

type FindCommand struct {
	All  bool
	Skip int
	Take int
	Last int
	Body []SearchInstruction
}

func (f FindCommand) IsCommand() {}

func (f FindCommand) String() string {
	return fmt.Sprintf("(find (all %t) (min %d max %d) (last %d) %s)", f.All, f.Skip, f.Take, f.Last, f.Body)
}

/*
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
	return findMatches(c.Body, c.All, c.Skip, c.Take, c.Last, filename, reader)
}
*/

type ReplaceCommand struct {
	All      bool
	Skip     int
	Take     int
	Last     int
	Body     []SearchInstruction
	Replacer []ReplaceInstruction
}

func (r ReplaceCommand) IsCommand() {}

func (r ReplaceCommand) String() string {
	return fmt.Sprintf("(replace (%t %d %d %d))", r.All, r.Skip, r.Take, r.Last)
}

/*
func (c ReplaceCommand) execute(filename string, reader *files.Reader, mode ReplaceMode) Matches {
	foundMatches := findMatches(c.Body, c.All, c.Skip, c.Take, c.Last, filename, reader)

	replacedMatches := Matches{}
	for _, match := range foundMatches {
		current_state := InitReplacerState(match, len(foundMatches))
		for current_state.programCounter < len(c.Replacer) {
			inst := c.Replacer[current_state.programCounter]
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
*/

type SetCommand struct {
	Body SetCommandBody
	Id   string
}

func (s SetCommand) IsCommand() {}

func (s SetCommand) String() string {
	return fmt.Sprintf("(set (id %s) (%v)", s.Id, s.Body)
}

//func (c SetCommand) execute(filename string, reader *files.Reader, mode ReplaceMode) Matches {
//	return Matches{}
//}

type SetCommandBody interface {
	IsSetCommandBody()
	// execute(state *GlobalState, id string) *GlobalState
}

type SetCommandExpression struct {
	Instructions []SearchInstruction
	Validate     []ast.AstProcessStatement
}

func (s SetCommandExpression) IsSetCommandBody() {}

/*
func (s SetCommandExpression) execute(state *GlobalState, id string) *GlobalState {
	return state
}
*/

type SetCommandMatches struct {
	Command Command
}

func (s SetCommandMatches) IsSetCommandBody() {}

/*
func (s SetCommandMatches) execute(state *GlobalState, id string) *GlobalState {
	// TODO run through the command we have and store the matches in variables
	// TODO or do the matches in the generate function
	return state
}
*/

type SetCommandTransform struct {
	Statements []ast.AstProcessStatement
}

func (s SetCommandTransform) IsSetCommandBody() {}

//func (s SetCommandTransform) execute(state *GlobalState, id string) *GlobalState {
//	return state
//}

type SearchInstruction interface {
	// execute(*SearchEngineState) *SearchEngineState
	IsSearchInstruction()
	adjust(offset int, state *GenState) SearchInstruction
	String() string
}

type ReplaceInstruction interface {
	IsReplaceInstruction()
	// execute(*ReplacerState) *ReplacerState
}

type MatchLiteral struct {
	Not      bool
	ToFind   string
	Caseless bool
}

func (i MatchLiteral) IsSearchInstruction() {}

func (i MatchLiteral) String() string {
	return fmt.Sprintf("(literal (not %t) (caseless %t) '%s')", i.Not, i.Caseless, i.ToFind)
}

func (i MatchLiteral) adjust(offset int, state *GenState) SearchInstruction {
	return i
}

//func (i MatchLiteral) execute(current_state *SearchEngineState) *SearchEngineState {
//	next_state := current_state.Copy()
//	next_state.MATCH(i.toFind, i.not, i.caseless)
//	return next_state
//}

type MatchCharClass struct {
	Not   bool
	Class ast.AstCharacterClassType
}

func (i MatchCharClass) IsSearchInstruction() {}

func (i MatchCharClass) String() string {
	return fmt.Sprintf("(class (not %t) %s)", i.Not, i.Class)
}

func (i MatchCharClass) adjust(offset int, state *GenState) SearchInstruction {
	return i
}

/*
func (i MatchCharClass) execute(current_state *SearchEngineState) *SearchEngineState {
	next_state := current_state.Copy()
	switch i.Class {
	case ClassAny:
		next_state.MATCHANY(i.Not)
	case ClassWhitespace:
		next_state.MATCHOPTIONS([]string{" ", "\t", "\n", "\r"}, i.Not)
	case ClassDigit:
		next_state.MATCHRANGE("0", "9", i.Not)
	case ClassUpper:
		next_state.MATCHRANGE("A", "Z", i.Not)
	case ClassLower:
		next_state.MATCHRANGE("a", "z", i.Not)
	case ClassLetter:
		next_state.MATCHLETTER(i.Not)
	case ClassFileStart:
		next_state.MATCHFILESTART(i.Not)
	case ClassFileEnd:
		next_state.MATCHFILEEND(i.Not)
	case ClassLineStart:
		next_state.MATCHLINESTART(i.Not)
	case ClassLineEnd:
		next_state.MATCHLINEEND(i.Not)
	case ClassWordStart:
		next_state.MATCHWORDSTART(i.Not)
	case ClassWordEnd:
		next_state.MATCHWORDEND(i.Not)
	case ClassWholeFile:
		next_state.MATCHWHOLEFILE(i.Not)
	case ClassWholeLine:
		next_state.MATCHWHOLELINE(i.Not)
	case ClassWholeWord:
		next_state.MATCHWHOLEWORD(i.Not)
	default:
		panic("Unexpected character class type")
	}
	return next_state
}
*/

type MatchVariable struct {
	Name string
}

func (i MatchVariable) IsSearchInstruction() {}

func (i MatchVariable) String() string {
	return fmt.Sprintf("(var '%s')", i.Name)
}

func (i MatchVariable) adjust(offset int, state *GenState) SearchInstruction {
	return i
}

//func (i MatchVariable) execute(current_state *SearchEngineState) *SearchEngineState {
//	next_state := current_state.Copy()
//	next_state.MATCHVAR(i.Name)
//	return next_state
//}

type MatchRange struct {
	Not  bool
	From string
	To   string
}

func (i MatchRange) IsSearchInstruction() {}

func (i MatchRange) String() string {
	return fmt.Sprintf("(range (not %t) (from '%s') (to '%s'))", i.Not, i.From, i.To)
}

func (i MatchRange) adjust(offset int, state *GenState) SearchInstruction {
	return i
}

//func (i MatchRange) execute(current_state *SearchEngineState) *SearchEngineState {
//	next_state := current_state.Copy()
//	next_state.MATCHRANGE(i.From, i.To, i.Not)
//	return next_state
//}

type CallSubroutine struct {
	Name string
	ToPC int
}

func (i CallSubroutine) IsSearchInstruction() {}

func (i CallSubroutine) String() string {
	return fmt.Sprintf("(call '%s' %d)", i.Name, i.ToPC)
}

func (i CallSubroutine) adjust(offset int, state *GenState) SearchInstruction {
	i.ToPC += offset
	return i
}

//func (i CallSubroutine) execute(current_state *SearchEngineState) *SearchEngineState {
//	next_state := current_state.Copy()
//	next_state.CALL(i.ToPC, next_state.programCounter+1)
//	next_state.JUMP(i.ToPC)
//	return next_state
//}

type Branch struct {
	Branches []int
}

func (i Branch) IsSearchInstruction() {}

func (i Branch) String() string {
	return fmt.Sprintf("(branch %v)", i.Branches)
}

func (i Branch) adjust(offset int, state *GenState) SearchInstruction {
	for idx := range i.Branches {
		i.Branches[idx] += offset
	}
	return i
}

// func (i Branch) execute(current_state *SearchEngineState) *SearchEngineState {
//	next_state := current_state.Copy()
//	flipped := []int{}
//	for k := range i.Branches {
//		flipped = append(flipped, i.Branches[len(i.Branches)-1-k])
//	}

//	for _, f := range flipped[:len(flipped)-1] {
//		next_state.JUMP(f)
//		next_state.CHECKPOINT()
//	}

//	next_state.JUMP(i.Branches[0])
//	return next_state
//}

type StartNotIn struct {
	NextCheckpointPC int
}

func (i StartNotIn) IsSearchInstruction() {}

func (i StartNotIn) String() string {
	return fmt.Sprintf("(startNotIn %d)", i.NextCheckpointPC)
}

func (i StartNotIn) adjust(offset int, state *GenState) SearchInstruction {
	i.NextCheckpointPC += offset
	return i
}

//func (i StartNotIn) execute(current_state *SearchEngineState) *SearchEngineState {
//	next_state := current_state.Copy()
//	next_state.JUMP(i.NextCheckpointPC)
//	next_state.CHECKPOINT()
//	next_state.JUMP(current_state.programCounter + 1)
//	return next_state
//}

type FailNotIn struct{}

func (i FailNotIn) IsSearchInstruction() {}

func (i FailNotIn) String() string {
	return "(failNotIn)"
}

func (i FailNotIn) adjust(offset int, state *GenState) SearchInstruction {
	return i
}

//func (i FailNotIn) execute(current_state *SearchEngineState) *SearchEngineState {
//	next_state := current_state.Copy()
//	next_state.BACKTRACK()
//	next_state.BACKTRACK()
//	return next_state
//}

type EndNotIn struct {
	MaxSize int
}

func (i EndNotIn) IsSearchInstruction() {}

func (i EndNotIn) String() string {
	return fmt.Sprintf("(endNotIn %d)", i.MaxSize)
}

func (i EndNotIn) adjust(offset int, state *GenState) SearchInstruction {
	return i
}

//func (i EndNotIn) execute(current_state *SearchEngineState) *SearchEngineState {
//	next_state := current_state.Copy()
// TODO this should actually let the rest of the expression backtrack from max size to min size (could just be to 1 since things less than the min are not in)
//	cfo := next_state.currentFileOffset
//	next_state.CONSUME(i.MaxSize)
// FIXME: This was added to make it so we don't have an infinite loop when using "not in" in an un-bounded loop
//        I think a better fix would be to come up with a different way to handle the end of the file
//	if cfo == next_state.currentFileOffset {
//		next_state.BACKTRACK()
//	} else {
//		next_state.NEXT()
//	}
//	return next_state
//}

type StartLoop struct {
	Id       int64
	MinLoops int
	MaxLoops int
	Fewest   bool
	ExitLoop int
	Name     string
}

func (i StartLoop) IsSearchInstruction() {}

func (i StartLoop) String() string {
	return fmt.Sprintf("(startLoop '%s' (min %d max %d) (lazy %t) %d %d)", i.Name, i.MinLoops, i.MaxLoops, i.Fewest, i.Id, i.ExitLoop)
}

func (i StartLoop) adjust(offset int, state *GenState) SearchInstruction {
	i.ExitLoop += offset
	return i
}

// func (i StartLoop) execute(current_state *SearchEngineState) *SearchEngineState {
//	next_state := current_state.Copy()

//	inited := next_state.INITLOOPSTACK(i.Id, i.Name)
//	if !inited {
//		if next_state.CHECKZEROMATCHLOOP() {
//			next_state.BACKTRACK()
//			return next_state
//		}
//		next_state.INCLOOPSTACK()
//	}
//	currentIteration := next_state.GETITERATIONSTEP()

//	if currentIteration < i.MinLoops {
//		next_state.NEXT()
//	} else if (i.MaxLoops == -1 || currentIteration <= i.MaxLoops) && i.Fewest {
//		next_state.NEXT()
//		next_state.CHECKPOINT()
//		next_state.POPLOOPSTACK()
//		next_state.JUMP(i.ExitLoop + 1)
//	} else if (i.MaxLoops == -1 || currentIteration <= i.MaxLoops) && !i.Fewest {
//		loop_state := next_state.POPLOOPSTACK()
//		pc := next_state.GETPC()
//		next_state.JUMP(i.ExitLoop + 1)
//		next_state.CHECKPOINT()
//		next_state.PUSHLOOPSTACK(loop_state)
//		next_state.JUMP(pc + 1)
//	} else {
//		next_state.BACKTRACK()
//	}

//	return next_state
//}

type StopLoop struct {
	Id        int64
	MinLoops  int
	MaxLoops  int
	Fewest    bool
	StartLoop int
	Name      string
}

func (i StopLoop) IsSearchInstruction() {}

func (i StopLoop) String() string {
	return fmt.Sprintf("(stopLoop '%s' (min %d max %d) (lazy %t) %d %d)", i.Name, i.MinLoops, i.MaxLoops, i.Fewest, i.Id, i.StartLoop)
}

func (i StopLoop) adjust(offset int, state *GenState) SearchInstruction {
	i.StartLoop += offset
	return i
}

//func (i StopLoop) execute(current_state *SearchEngineState) *SearchEngineState {
//	next_state := current_state.Copy()
//	next_state.JUMP(i.StartLoop)
//	return next_state
//}

type StartVarDec struct {
	Name string
}

func (i StartVarDec) IsSearchInstruction() {}

func (i StartVarDec) String() string {
	return fmt.Sprintf("(startVarDec '%s')", i.Name)
}

func (i StartVarDec) adjust(offset int, state *GenState) SearchInstruction {
	return i
}

//func (i StartVarDec) execute(current_state *SearchEngineState) *SearchEngineState {
//	next_state := current_state.Copy()
//	next_state.STARTVAR(i.Name)
//	return next_state
//}

type EndVarDec struct {
	Name string
}

func (i EndVarDec) IsSearchInstruction() {}

func (i EndVarDec) String() string {
	return fmt.Sprintf("(endVarDec '%s')", i.Name)
}

func (i EndVarDec) adjust(offset int, state *GenState) SearchInstruction {
	return i
}

//func (i EndVarDec) execute(current_state *SearchEngineState) *SearchEngineState {
//	next_state := current_state.Copy()
//	next_state.ENDVAR(i.Name)
//	return next_state
//}

type StartSubroutine struct {
	Id        int
	Name      string
	EndOffset int
}

func (i StartSubroutine) IsSearchInstruction() {}

func (i StartSubroutine) String() string {
	return fmt.Sprintf("(startSub '%s' %d %d)", i.Name, i.Id, i.EndOffset)
}

func (i StartSubroutine) adjust(offset int, state *GenState) SearchInstruction {
	i.EndOffset += offset
	return i
}

//func (i StartSubroutine) execute(current_state *SearchEngineState) *SearchEngineState {
//	next_state := current_state.Copy()
//	next_state.VALIDATECALL(i.Id, i.EndOffset+1)
//	next_state.NEXT()
//	return next_state
//}

type EndSubroutine struct {
	Name     string
	Validate []ast.AstProcessStatement
}

func (i EndSubroutine) IsSearchInstruction() {}

func (i EndSubroutine) String() string {
	return fmt.Sprintf("(endSub '%s')", i.Name)
}

func (i EndSubroutine) adjust(offset int, state *GenState) SearchInstruction {
	return i
}

// func (i EndSubroutine) execute(current_state *SearchEngineState) *SearchEngineState {
//	next_state := current_state.Copy()

//	if len(i.Validate) == 0 {
//		next_state.RETURN()
//	} else {
//		env := make(map[string]ProcessValue)
//		subMatch := current_state.currentMatch[current_state.callStack.Peek().startMatchOffset:]
//		env["match"] = ProcessValueString{subMatch}
//		env["matchLength"] = ProcessValueNumber{len(subMatch)}
// TODO add more variables here!

//		pstate := ProcessState{
//			currentValue: ProcessValueString{""},
//			environment:  env,
//			status:       NEXT,
//		}
//		var final_value ProcessValue = ProcessValueBoolean{true}
//		for _, stmt := range i.Validate {
//			pstate = stmt.execute(pstate)
//			if pstate.status == RETURNING {
//				final_value = pstate.currentValue
//				break
//			}
//		}

//		if final_value.getBoolean() {
//			next_state.RETURN()
//		} else {
//			next_state.BACKTRACK()
//		}
//	}

//	return next_state
//}

type Jump struct {
	NewProgramCounter int
}

func (i Jump) IsSearchInstruction() {}

func (i Jump) String() string {
	return fmt.Sprintf("(jump %d)", i.NewProgramCounter)
}

func (i Jump) adjust(offset int, state *GenState) SearchInstruction {
	i.NewProgramCounter += offset
	return i
}

//func (i Jump) execute(current_state *SearchEngineState) *SearchEngineState {
//	next_state := current_state.Copy()
//	next_state.JUMP(i.NewProgramCounter)
//	return next_state
//}

type ReplaceString struct {
	Value string
}

func (i ReplaceString) IsReplaceInstruction() {}

//func (i ReplaceString) execute(current_state *ReplacerState) *ReplacerState {
//	next_state := current_state.Copy()
//	next_state.WRITESTRING(i.Value)
//	next_state.NEXT()
//	return next_state
//}

type ReplaceVariable struct {
	Name string
}

func (i ReplaceVariable) IsReplaceInstruction() {}

//func (i ReplaceVariable) execute(current_state *ReplacerState) *ReplacerState {
//	next_state := current_state.Copy()
//	next_state.WRITEVAR(i.Name)
//	next_state.NEXT()
//	return next_state
//}

type ReplaceProcess struct {
	Process []ast.AstProcessStatement
}

func (i ReplaceProcess) IsReplaceInstruction() {}

// func (i ReplaceProcess) execute(current_state *ReplacerState) *ReplacerState {
//	next_state := current_state.Copy()

// execute AST
//	env := make(map[string]ProcessValue)
//	keys := current_state.variables.Keys()
//	for _, key := range keys {
//		value, _ := current_state.variables.Get(key)
//		if value.getType() == ValueStringType {
//			env[key] = ProcessValueString{value.String().Value}
//		}
// TODO Need to add process hash maps or merge into the main Values
//	}

//	env["match"] = ProcessValueString{next_state.match.Value}
//	env["matchLength"] = ProcessValueNumber{len(next_state.match.Value)}
//	env["matchNumber"] = ProcessValueNumber{next_state.match.MatchNumber}

//	pstate := ProcessState{
//		currentValue: ProcessValueString{""},
//		environment:  env,
//		status:       NEXT,
//	}
//	var final_value ProcessValue = ProcessValueBoolean{true}
//	for _, stmt := range i.Process {
//		pstate = stmt.execute(pstate)
//		if pstate.status == RETURNING {
//			final_value = pstate.currentValue
//			break
//		}
//	}

//	next_state.WRITESTRING(final_value.getString())

//	next_state.NEXT()
//	return next_state
//}
