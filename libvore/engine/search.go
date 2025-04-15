package engine

import (
	"fmt"

	"github.com/jmeaster30/vore/libvore/ast"
	"github.com/jmeaster30/vore/libvore/bytecode"
	"github.com/jmeaster30/vore/libvore/ds"
	"github.com/jmeaster30/vore/libvore/files"
)

func search(command *bytecode.Command, filename string, reader *files.Reader, mode ReplaceMode) Matches {
	var ci any = command
	switch ci.(type) {
	case bytecode.FindCommand:
		return searchFind(ci.(*bytecode.FindCommand), filename, reader, mode)
	case bytecode.ReplaceCommand:
		return searchReplace(ci.(*bytecode.ReplaceCommand), filename, reader, mode)
	case bytecode.SetCommand:
		return Matches{}
	}
	panic(fmt.Sprintf("Unknown command %T", ci))
}

func findMatches(insts []bytecode.SearchInstruction, all bool, skip int, take int, last int, filename string, reader *files.Reader) Matches {
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
			currentState = matchInstruction(inst, currentState)
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

func searchFind(c *bytecode.FindCommand, filename string, reader *files.Reader, mode ReplaceMode) Matches {
	return findMatches(c.Body, c.All, c.Skip, c.Take, c.Last, filename, reader)
}

func searchReplace(c *bytecode.ReplaceCommand, filename string, reader *files.Reader, mode ReplaceMode) Matches {
	foundMatches := findMatches(c.Body, c.All, c.Skip, c.Take, c.Last, filename, reader)

	replacedMatches := Matches{}
	for _, match := range foundMatches {
		current_state := InitReplacerState(match, len(foundMatches))
		for current_state.programCounter < len(c.Replacer) {
			inst := c.Replacer[current_state.programCounter]
			current_state = executeReplace(inst, current_state)
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

func matchInstruction(i bytecode.SearchInstruction, current_state *SearchEngineState) *SearchEngineState {
	var ii any = i
	switch ii.(type) {
	case bytecode.MatchLiteral:
		return matchLiteral(ii.(bytecode.MatchLiteral), current_state)
	case bytecode.MatchCharClass:
		return matchCharClass(ii.(bytecode.MatchCharClass), current_state)
	case bytecode.MatchVariable:
		return matchVariable(ii.(bytecode.MatchVariable), current_state)
	case bytecode.MatchRange:
		return matchRange(ii.(bytecode.MatchRange), current_state)
	case bytecode.CallSubroutine:
		return matchCallSubroutine(ii.(bytecode.CallSubroutine), current_state)
	case bytecode.Branch:
		return matchBranch(ii.(bytecode.Branch), current_state)
	case bytecode.StartNotIn:
		return matchStartNotIn(ii.(bytecode.StartNotIn), current_state)
	case bytecode.EndNotIn:
		return matchEndNotIn(ii.(bytecode.EndNotIn), current_state)
	case bytecode.FailNotIn:
		return matchFailNotIn(ii.(bytecode.FailNotIn), current_state)
	case bytecode.StartLoop:
		return matchStartLoop(ii.(bytecode.StartLoop), current_state)
	case bytecode.StopLoop:
		return matchStopLoop(ii.(bytecode.StopLoop), current_state)
	case bytecode.StartVarDec:
		return matchStartVarDec(ii.(bytecode.StartVarDec), current_state)
	case bytecode.EndVarDec:
		return matchEndVarDec(ii.(bytecode.EndVarDec), current_state)
	case bytecode.StartSubroutine:
		return matchStartSubroutine(ii.(bytecode.StartSubroutine), current_state)
	case bytecode.EndSubroutine:
		return matchEndSubroutine(ii.(bytecode.EndSubroutine), current_state)
	case bytecode.Jump:
		return matchJump(ii.(bytecode.Jump), current_state)
	}
	panic(fmt.Sprintf("Unknown search instruction %T", ii))
}

func matchLiteral(i bytecode.MatchLiteral, current_state *SearchEngineState) *SearchEngineState {
	next_state := current_state.Copy()
	next_state.MATCH(i.ToFind, i.Not, i.Caseless)
	return next_state
}

func matchCharClass(i bytecode.MatchCharClass, current_state *SearchEngineState) *SearchEngineState {
	next_state := current_state.Copy()
	switch i.Class {
	case ast.ClassAny:
		next_state.MATCHANY(i.Not)
	case ast.ClassWhitespace:
		next_state.MATCHOPTIONS([]string{" ", "\t", "\n", "\r"}, i.Not)
	case ast.ClassDigit:
		next_state.MATCHRANGE("0", "9", i.Not)
	case ast.ClassUpper:
		next_state.MATCHRANGE("A", "Z", i.Not)
	case ast.ClassLower:
		next_state.MATCHRANGE("a", "z", i.Not)
	case ast.ClassLetter:
		next_state.MATCHLETTER(i.Not)
	case ast.ClassFileStart:
		next_state.MATCHFILESTART(i.Not)
	case ast.ClassFileEnd:
		next_state.MATCHFILEEND(i.Not)
	case ast.ClassLineStart:
		next_state.MATCHLINESTART(i.Not)
	case ast.ClassLineEnd:
		next_state.MATCHLINEEND(i.Not)
	case ast.ClassWordStart:
		next_state.MATCHWORDSTART(i.Not)
	case ast.ClassWordEnd:
		next_state.MATCHWORDEND(i.Not)
	case ast.ClassWholeFile:
		next_state.MATCHWHOLEFILE(i.Not)
	case ast.ClassWholeLine:
		next_state.MATCHWHOLELINE(i.Not)
	case ast.ClassWholeWord:
		next_state.MATCHWHOLEWORD(i.Not)
	default:
		panic("Unexpected character class type")
	}
	return next_state
}

func matchVariable(i bytecode.MatchVariable, current_state *SearchEngineState) *SearchEngineState {
	next_state := current_state.Copy()
	next_state.MATCHVAR(i.Name)
	return next_state
}

func matchRange(i bytecode.MatchRange, current_state *SearchEngineState) *SearchEngineState {
	next_state := current_state.Copy()
	next_state.MATCHRANGE(i.From, i.To, i.Not)
	return next_state
}

func matchCallSubroutine(i bytecode.CallSubroutine, current_state *SearchEngineState) *SearchEngineState {
	next_state := current_state.Copy()
	next_state.CALL(i.ToPC, next_state.programCounter+1)
	next_state.JUMP(i.ToPC)
	return next_state
}

func matchBranch(i bytecode.Branch, current_state *SearchEngineState) *SearchEngineState {
	next_state := current_state.Copy()
	flipped := []int{}
	for k := range i.Branches {
		flipped = append(flipped, i.Branches[len(i.Branches)-1-k])
	}

	for _, f := range flipped[:len(flipped)-1] {
		next_state.JUMP(f)
		next_state.CHECKPOINT()
	}

	next_state.JUMP(i.Branches[0])
	return next_state
}

func matchStartNotIn(i bytecode.StartNotIn, current_state *SearchEngineState) *SearchEngineState {
	next_state := current_state.Copy()
	next_state.JUMP(i.NextCheckpointPC)
	next_state.CHECKPOINT()
	next_state.JUMP(current_state.programCounter + 1)
	return next_state
}

func matchFailNotIn(i bytecode.FailNotIn, current_state *SearchEngineState) *SearchEngineState {
	next_state := current_state.Copy()
	next_state.BACKTRACK()
	next_state.BACKTRACK()
	return next_state
}

func matchEndNotIn(i bytecode.EndNotIn, current_state *SearchEngineState) *SearchEngineState {
	next_state := current_state.Copy()
	// TODO this should actually let the rest of the expression backtrack from max size to min size (could just be to 1 since things less than the min are not in)
	cfo := next_state.currentFileOffset
	next_state.CONSUME(i.MaxSize)
	// FIXME: This was added to make it so we don't have an infinite loop when using "not in" in an un-bounded loop
	//        I think a better fix would be to come up with a different way to handle the end of the file
	if cfo == next_state.currentFileOffset {
		next_state.BACKTRACK()
	} else {
		next_state.NEXT()
	}
	return next_state
}

func matchStartLoop(i bytecode.StartLoop, current_state *SearchEngineState) *SearchEngineState {
	next_state := current_state.Copy()

	inited := next_state.INITLOOPSTACK(i.Id, i.Name)
	if !inited {
		if next_state.CHECKZEROMATCHLOOP() {
			next_state.BACKTRACK()
			return next_state
		}
		next_state.INCLOOPSTACK()
	}
	currentIteration := next_state.GETITERATIONSTEP()

	if currentIteration < i.MinLoops {
		next_state.NEXT()
	} else if (i.MaxLoops == -1 || currentIteration <= i.MaxLoops) && i.Fewest {
		next_state.NEXT()
		next_state.CHECKPOINT()
		next_state.POPLOOPSTACK()
		next_state.JUMP(i.ExitLoop + 1)
	} else if (i.MaxLoops == -1 || currentIteration <= i.MaxLoops) && !i.Fewest {
		loop_state := next_state.POPLOOPSTACK()
		pc := next_state.GETPC()
		next_state.JUMP(i.ExitLoop + 1)
		next_state.CHECKPOINT()
		next_state.PUSHLOOPSTACK(loop_state)
		next_state.JUMP(pc + 1)
	} else {
		next_state.BACKTRACK()
	}

	return next_state
}

func matchStopLoop(i bytecode.StopLoop, current_state *SearchEngineState) *SearchEngineState {
	next_state := current_state.Copy()
	next_state.JUMP(i.StartLoop)
	return next_state
}

func matchStartVarDec(i bytecode.StartVarDec, current_state *SearchEngineState) *SearchEngineState {
	next_state := current_state.Copy()
	next_state.STARTVAR(i.Name)
	return next_state
}

func matchEndVarDec(i bytecode.EndVarDec, current_state *SearchEngineState) *SearchEngineState {
	next_state := current_state.Copy()
	next_state.ENDVAR(i.Name)
	return next_state
}

func matchStartSubroutine(i bytecode.StartSubroutine, current_state *SearchEngineState) *SearchEngineState {
	next_state := current_state.Copy()
	next_state.VALIDATECALL(i.Id, i.EndOffset+1)
	next_state.NEXT()
	return next_state
}

func matchEndSubroutine(i bytecode.EndSubroutine, current_state *SearchEngineState) *SearchEngineState {
	next_state := current_state.Copy()

	if len(i.Validate) == 0 {
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
		for _, stmt := range i.Validate {
			pstate = executeStatement(stmt, pstate)
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

func matchJump(i bytecode.Jump, current_state *SearchEngineState) *SearchEngineState {
	next_state := current_state.Copy()
	next_state.JUMP(i.NewProgramCounter)
	return next_state
}

func executeReplace(i bytecode.ReplaceInstruction, current_state *ReplacerState) *ReplacerState {
	var ii any = i
	switch ii.(type) {
	case bytecode.ReplaceString:
		return executeReplaceString(ii.(bytecode.ReplaceString), current_state)
	case bytecode.ReplaceVariable:
		return executeReplaceVariable(ii.(bytecode.ReplaceVariable), current_state)
	case bytecode.ReplaceProcess:
		return executeReplaceProcess(ii.(bytecode.ReplaceProcess), current_state)
	}
	panic(fmt.Sprintf("Unknown replace instruction %T", ii))
}

func executeReplaceString(i bytecode.ReplaceString, current_state *ReplacerState) *ReplacerState {
	next_state := current_state.Copy()
	next_state.WRITESTRING(i.Value)
	next_state.NEXT()
	return next_state
}

func executeReplaceVariable(i bytecode.ReplaceVariable, current_state *ReplacerState) *ReplacerState {
	next_state := current_state.Copy()
	next_state.WRITEVAR(i.Name)
	next_state.NEXT()
	return next_state
}

func executeReplaceProcess(i bytecode.ReplaceProcess, current_state *ReplacerState) *ReplacerState {
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
	for _, stmt := range i.Process {
		pstate = executeStatement(stmt, pstate)
		if pstate.status == RETURNING {
			final_value = pstate.currentValue
			break
		}
	}

	next_state.WRITESTRING(final_value.getString())
	next_state.NEXT()
	return next_state
}
