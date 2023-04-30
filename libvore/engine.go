package libvore

import (
	"strconv"
	"strings"
)

type Status int

const (
	SUCCESS Status = iota
	FAILED
	INPROCESS
)

type LoopState struct {
	loopId              int64
	callLevel           int
	iterationStep       int
	name                string
	loopMatchIndexStart int
	variables           ValueHashMap
}

type VariableRecord struct {
	name        string
	startOffset int
}

type CallState struct {
	id               int
	returnOffset     int
	startMatchOffset int
}

type SearchEngineState struct {
	loopStack     *Stack[LoopState]
	backtrack     *Stack[SearchEngineState]
	variableStack *Stack[VariableRecord]
	callStack     *Stack[CallState]
	environment   ValueHashMap

	status            Status
	programCounter    int
	currentFileOffset int
	currentMatch      string
	currentLineNum    int
	currentColumnNum  int
	startFileOffset   int
	startLineNum      int
	startColumnNum    int
	reader            *VReader
	filename          string
}

func (es *SearchEngineState) SEEK() {
	es.reader.Seek(es.currentFileOffset)
}

func (es *SearchEngineState) SEEKTO(offset int) {
	es.reader.Seek(offset)
}

func (es *SearchEngineState) READ(length int) string {
	es.SEEK()
	return es.reader.Read(length)
}

func (es *SearchEngineState) READAT(offset int, length int) string {
	es.SEEKTO(offset)
	return es.reader.Read(length)
}

func (es *SearchEngineState) CONSUME(amount int) {
	value := es.READ(amount)
	es.currentMatch += value
	es.currentFileOffset += len(value)
	for _, c := range value {
		es.currentColumnNum += 1
		if c == rune('\n') {
			es.currentLineNum += 1
			es.currentColumnNum = 1
		}
	}
}

func (es *SearchEngineState) BACKTRACK() {
	if es.backtrack.Size() == 0 {
		es.FAIL()
	} else {
		next_state := es.backtrack.Pop()
		es.Set(next_state)
	}
}

func (es *SearchEngineState) FAIL() {
	es.status = FAILED
}

func (es *SearchEngineState) SUCCESS() {
	es.status = SUCCESS
}

func (es *SearchEngineState) MATCHFILESTART(not bool) {
	if es.currentFileOffset == 0 {
		if not {
			es.BACKTRACK()
		} else {
			es.NEXT()
		}
	} else {
		if not {
			es.NEXT()
		} else {
			es.BACKTRACK()
		}
	}
}

func (es *SearchEngineState) MATCHFILEEND(not bool) {
	if es.currentFileOffset == es.reader.size {
		if not {
			es.BACKTRACK()
		} else {
			es.NEXT()
		}
	} else {
		if not {
			es.NEXT()
		} else {
			es.BACKTRACK()
		}
	}
}

func (es *SearchEngineState) MATCHLINESTART(not bool) {
	if es.currentFileOffset == 0 {
		if not {
			es.BACKTRACK()
		} else {
			es.NEXT()
		}
		return
	}

	value := es.READAT(es.currentFileOffset-1, 1)
	if value == "\n" {
		if not {
			es.BACKTRACK()
		} else {
			es.NEXT()
		}
	} else {
		if not {
			es.NEXT()
		} else {
			es.BACKTRACK()
		}
	}
}

func (es *SearchEngineState) MATCHLINEEND(not bool) {
	nextChar := es.READ(1)
	nextTwoChar := es.READ(2)
	if nextChar == "\n" || nextTwoChar == "\r\n" || es.currentFileOffset == es.reader.size {
		if not {
			es.BACKTRACK()
		} else {
			es.NEXT()
		}
	} else {
		if not {
			es.NEXT()
		} else {
			es.BACKTRACK()
		}
	}
}

func IsLetter(value string) bool {
	return ("a" <= value && value <= "z") || ("A" <= value && value <= "Z") || ("0" <= value && value <= "9") || value == "_"
}

func (es *SearchEngineState) MATCHWORDSTART(not bool) {
	if es.currentFileOffset == es.reader.size {
		if not {
			es.BACKTRACK()
		} else {
			es.NEXT()
		}
		return
	}

	current := es.READ(1)
	if es.currentFileOffset == 0 {
		if IsLetter(current) {
			if not {
				es.BACKTRACK()
			} else {
				es.NEXT()
			}
		} else {
			if not {
				es.NEXT()
			} else {
				es.BACKTRACK()
			}
		}
		return
	}

	previous := es.READAT(es.currentFileOffset-1, 1)
	if IsLetter(current) && !IsLetter(previous) {
		if not {
			es.BACKTRACK()
		} else {
			es.NEXT()
		}
	} else {
		if not {
			es.NEXT()
		} else {
			es.BACKTRACK()
		}
	}
}

func (es *SearchEngineState) MATCHWORDEND(not bool) {
	if es.currentFileOffset == 0 {
		if not {
			es.BACKTRACK()
		} else {
			es.NEXT()
		}
		return
	}

	current := es.READ(1)
	if es.currentFileOffset == es.reader.size {
		if !IsLetter(current) {
			if not {
				es.BACKTRACK()
			} else {
				es.NEXT()
			}
		} else {
			if not {
				es.NEXT()
			} else {
				es.BACKTRACK()
			}
		}
		return
	}

	previous := es.READAT(es.currentFileOffset-1, 1)

	if !IsLetter(current) && IsLetter(previous) {
		if not {
			es.BACKTRACK()
		} else {
			es.NEXT()
		}
	} else {
		if not {
			es.NEXT()
		} else {
			es.BACKTRACK()
		}
	}
}

func (es *SearchEngineState) MATCHWHOLEFILE(not bool) {
	if es.currentFileOffset != 0 {
		if not {
			// TODO Should this be a zero match or a single char match?
			// I think zero match is better so it keeps in line with the other file anchors
			es.NEXT()
		} else {
			es.BACKTRACK()
		}
		return
	}

	if not {
		es.BACKTRACK()
		return
	}

	// TODO This is probably going to be a performance concern
	es.CONSUME(es.reader.size)
	es.NEXT()
}

func (es *SearchEngineState) MATCHWHOLELINE(not bool) {
	if (es.currentFileOffset != 0 && es.READAT(es.currentFileOffset-1, 1) != "\n") || es.currentFileOffset == es.reader.size {
		if not {
			es.NEXT()
		} else {
			es.BACKTRACK()
		}
		return
	}

	if not {
		es.BACKTRACK()
		return
	}

	// we know we are at the start of a line and we can read a character
	for {
		es.CONSUME(1)
		if es.currentFileOffset == es.reader.size {
			break
		}

		nextChar := es.READ(1)
		nextTwoChar := es.READ(2)
		if nextChar == "\n" || nextTwoChar == "\r\n" || es.currentFileOffset == es.reader.size {
			break
		}
	}

	es.NEXT()
}

func (es *SearchEngineState) MATCHWHOLEWORD(not bool) {
	if (es.currentFileOffset != 0 && (!IsLetter(es.READ(1)) || IsLetter(es.READAT(es.currentFileOffset-1, 1)))) || es.currentFileOffset == es.reader.size {
		if not {
			es.NEXT()
		} else {
			es.BACKTRACK()
		}
		return
	}

	if not {
		es.BACKTRACK()
		return
	}

	// we know we are at the start of a line and we can read a character
	for {
		es.CONSUME(1)
		if es.currentFileOffset == es.reader.size {
			break
		}

		current := es.READ(1)
		previous := es.READAT(es.currentFileOffset-1, 1)

		if !IsLetter(current) && IsLetter(previous) {
			break
		}
	}

	es.NEXT()
}

func (es *SearchEngineState) MATCHANY(not bool) {
	if not {
		es.BACKTRACK()
		return
	}
	value := es.READ(1)
	if value == "" {
		es.BACKTRACK()
	} else {
		es.CONSUME(1)
		es.NEXT()
	}
}

func (es *SearchEngineState) MATCHRANGE(from string, to string, not bool) {
	min := len(from)
	max := len(to)

	for i := max; i >= min; i-- {
		value := es.READ(i)
		if (from <= value && value <= to && !not) || ((from > value || value > to) && not) {
			es.CONSUME(i)
			es.NEXT()
			return
		}
	}

	es.BACKTRACK()
}

func (es *SearchEngineState) MATCHLETTER(not bool) {
	// TODO I would prefer if I had a generic way to do these multirange searches
	value := es.READ(1)
	if ("a" <= value && value <= "z") || ("A" <= value && value <= "Z") {
		if not {
			es.BACKTRACK()
		} else {
			es.CONSUME(1)
			es.NEXT()
		}
	} else {
		if not {
			es.CONSUME(1)
			es.NEXT()
		} else {
			es.BACKTRACK()
		}
	}
}

func (es *SearchEngineState) MATCHOPTIONS(options []string, not bool) {
	value := es.READ(1)
	if value == "" {
		es.BACKTRACK()
		return
	}

	for _, opt := range options {
		if value == opt {
			if not {
				es.BACKTRACK()
				return
			} else {
				es.CONSUME(1)
				es.NEXT()
				return
			}
		}
	}

	if not {
		es.CONSUME(1)
		es.NEXT()
	} else {
		es.BACKTRACK()
	}
}

func compare(a string, b string, caseless bool) bool {
	if caseless {
		return strings.EqualFold(a, b)
	} else {
		return a == b
	}
}

func (es *SearchEngineState) MATCH(value string, not bool, caseless bool) {
	comp := es.READ(len(value))

	if len(comp) == 0 {
		es.BACKTRACK()
		return
	}

	if !not && compare(value, comp, caseless) {
		es.CONSUME(len(value))
		es.NEXT()
	} else if not && !compare(value, comp, caseless) {
		es.CONSUME(len(value))
		es.NEXT()
	} else {
		es.BACKTRACK()
	}
}

func (es *SearchEngineState) MATCHVAR(name string) {
	value, found := es.environment.Get(name)
	if !found {
		es.BACKTRACK()
	} else if value.getType() == ValueHashMapType {
		// TODO add syntax for indexing hash maps but also I want something a bit better than just failing here
		es.BACKTRACK()
	} else {
		es.MATCH(value.String().Value, false, false)
	}
}

func (es *SearchEngineState) NEXT() {
	es.programCounter += 1
}

func (es *SearchEngineState) JUMP(pc int) {
	es.programCounter = pc
}

func (es *SearchEngineState) GETPC() int {
	return es.programCounter
}

func (es *SearchEngineState) INITLOOPSTACK(loopId int64, name string) bool {
	top := es.loopStack.Peek()
	if es.loopStack.IsEmpty() || top.loopId != loopId || top.callLevel != int(es.callStack.Size()) {
		lstate := LoopState{
			loopId:              loopId,
			name:                name,
			callLevel:           int(es.callStack.Size()),
			iterationStep:       0,
			loopMatchIndexStart: len(es.currentMatch),
			variables:           NewValueHashMap(),
		}
		lstate.variables.Add("0", NewValueHashMap())
		es.loopStack.Push(lstate)
		return true
	}
	return false
}

func (es *SearchEngineState) INCLOOPSTACK() {
	if es.loopStack.IsEmpty() {
		panic("oh crap :(")
	}
	es.loopStack.Peek().iterationStep += 1
	es.loopStack.Peek().loopMatchIndexStart = len(es.currentMatch)
	es.loopStack.Peek().variables.Add(strconv.Itoa(es.loopStack.Peek().iterationStep), NewValueHashMap())
}

func (es *SearchEngineState) GETITERATIONSTEP() int {
	if es.loopStack.IsEmpty() {
		panic("oh crap :(")
	}
	return es.loopStack.Peek().iterationStep
}

func (es *SearchEngineState) CHECKZEROMATCHLOOP() bool {
	if es.loopStack.IsEmpty() {
		panic("Loop stack is empty :(")
	}
	return es.loopStack.Peek().loopMatchIndexStart == len(es.currentMatch)
}

func (es *SearchEngineState) POPLOOPSTACK() LoopState {
	if es.loopStack.IsEmpty() {
		panic("oh crap :(")
	}
	top := es.loopStack.Peek()
	es.loopStack.Pop()
	if top.name != "" {
		es.INSERTVARIABLE(top.name, top.variables)
	}
	return *top
}

func (es *SearchEngineState) PUSHLOOPSTACK(loopState LoopState) {
	es.loopStack.Push(loopState)
}

func (es *SearchEngineState) STARTVAR(name string) {
	record := VariableRecord{
		name:        name,
		startOffset: len(es.currentMatch),
	}
	es.variableStack.Push(record)
	es.NEXT()
}

func (es *SearchEngineState) ENDVAR(name string) {
	record := es.variableStack.Pop()
	if record.name != name {
		panic("UHOH BAD INSTRUCTIONS I TRIED RESOLVING A VARIABLE THAT I WASN'T EXPECTING")
	}
	value := es.currentMatch[record.startOffset:]
	es.INSERTVARIABLE(name, NewValueString(value))
	es.NEXT()
}

func (es *SearchEngineState) INSERTVARIABLE(name string, value Value) {
	var lowestScope *LoopState = nil
	for i := int(es.loopStack.Size()) - 1; i >= 0; i-- {
		lowestScope = es.loopStack.Index(i)
		if lowestScope.name != "" {
			break
		}
	}

	if lowestScope == nil || lowestScope.name == "" {
		es.environment.Add(name, value)
	} else {
		index := strconv.Itoa(lowestScope.iterationStep)
		v, prs := lowestScope.variables.Get(index)
		if !prs {
			m := NewValueHashMap()
			m.Add(name, value)
			lowestScope.variables.Add(index, m)
		} else {
			m := v.Hashmap()
			m.Add(name, value)
			lowestScope.variables.Add(index, m)
		}
	}
}

func (es *SearchEngineState) VALIDATECALL(id int, returnOffset int) {
	top := es.callStack.Peek()
	if top == nil || top.id != id {
		es.CALL(id, returnOffset)
	}
}

func (es *SearchEngineState) CALL(id int, returnOffset int) {
	es.callStack.Push(CallState{
		id:           id,
		returnOffset: returnOffset,
	})
}

func (es *SearchEngineState) RETURN() {
	top := es.callStack.Pop()
	if top == nil {
		panic("BAD CALL STACK :(")
	}
	es.programCounter = top.returnOffset
}

func (es *SearchEngineState) CHECKPOINT() {
	checkpoint := es.Copy()
	es.backtrack.Push(*checkpoint)
}

func CreateState(filename string, reader *VReader, fileOffset int, lineNumber int, columnNumber int) *SearchEngineState {
	return &SearchEngineState{
		loopStack:         NewStack[LoopState](),
		backtrack:         NewStack[SearchEngineState](),
		variableStack:     NewStack[VariableRecord](),
		callStack:         NewStack[CallState](),
		environment:       NewValueHashMap(),
		status:            INPROCESS,
		programCounter:    0,
		currentFileOffset: fileOffset,
		currentLineNum:    lineNumber,
		currentColumnNum:  columnNumber,
		startFileOffset:   fileOffset,
		startLineNum:      lineNumber,
		startColumnNum:    columnNumber,
		reader:            reader,
		filename:          filename,
	}
}

func (es *SearchEngineState) Copy() *SearchEngineState {
	return &SearchEngineState{
		loopStack:         es.loopStack.Copy(),
		backtrack:         es.backtrack.Copy(),
		variableStack:     es.variableStack.Copy(),
		callStack:         es.callStack.Copy(),
		environment:       es.environment,
		status:            es.status,
		programCounter:    es.programCounter,
		currentFileOffset: es.currentFileOffset,
		currentMatch:      es.currentMatch,
		currentLineNum:    es.currentLineNum,
		currentColumnNum:  es.currentColumnNum,
		startFileOffset:   es.startFileOffset,
		startLineNum:      es.startLineNum,
		startColumnNum:    es.startColumnNum,
		reader:            es.reader,
		filename:          es.filename,
	}
}

func (es *SearchEngineState) Set(value *SearchEngineState) {
	es.loopStack = value.loopStack
	es.backtrack = value.backtrack
	es.variableStack = value.variableStack
	es.callStack = value.callStack
	es.environment = value.environment
	es.status = value.status
	es.programCounter = value.programCounter
	es.currentFileOffset = value.currentFileOffset
	es.currentMatch = value.currentMatch
	es.currentLineNum = value.currentLineNum
	es.currentColumnNum = value.currentColumnNum
	es.startFileOffset = value.startFileOffset
	es.startLineNum = value.startLineNum
	es.startColumnNum = value.startColumnNum
	es.reader = value.reader
	es.filename = value.filename
}

func (es *SearchEngineState) MakeMatch(matchNumber int) Match {
	return Match{
		Filename:    es.filename,
		MatchNumber: matchNumber,
		Offset:      *NewRange(es.startFileOffset, es.currentFileOffset),
		Line:        *NewRange(es.startLineNum, es.currentLineNum),
		Column:      *NewRange(es.startColumnNum, es.currentColumnNum),
		Value:       es.currentMatch,
		Variables:   es.environment,
	}
}

type ReplacerState struct {
	variables      ValueHashMap
	match          Match
	programCounter int
}

func InitReplacerState(match Match, totalMatches int) *ReplacerState {
	variables := match.Variables.Copy().Hashmap()

	variables.Add("totalMatches", NewValueString(strconv.Itoa(totalMatches)))
	variables.Add("matchNumber", NewValueString(strconv.Itoa(match.MatchNumber)))
	variables.Add("startOffset", NewValueString(strconv.Itoa(match.Offset.Start)))
	variables.Add("endOffset", NewValueString(strconv.Itoa(match.Offset.End)))
	variables.Add("lineNumber", NewValueString(strconv.Itoa(match.Line.Start)))
	variables.Add("columnNumber", NewValueString(strconv.Itoa(match.Column.Start)))
	variables.Add("value", NewValueString(match.Value))
	variables.Add("filename", NewValueString(match.Filename))
	return &ReplacerState{
		variables:      variables,
		match:          match,
		programCounter: 0,
	}
}

func (rs *ReplacerState) NEXT() {
	rs.programCounter += 1
}

func (rs *ReplacerState) WRITESTRING(value string) {
	rs.match.Replacement = Some(rs.match.Replacement.GetValueOrDefault("") + value)
}

func (rs *ReplacerState) WRITEVAR(name string) {
	value, found := rs.variables.Get(name)
	if found && value.getType() == ValueStringType {
		rs.match.Replacement = Some(rs.match.Replacement.GetValueOrDefault("") + value.String().Value)
	}
}

func (rs *ReplacerState) Copy() *ReplacerState {
	return &ReplacerState{
		programCounter: rs.programCounter,
		match:          rs.match,
		variables:      rs.variables,
	}
}

func (rs *ReplacerState) Set(from *ReplacerState) {
	rs.variables = from.variables
	rs.match = from.match
	rs.programCounter = from.programCounter
}

type GlobalState struct {
	//subroutines map[string][]SearchInstruction
	//matches     map[string]Matches
}
