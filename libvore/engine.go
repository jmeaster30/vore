package libvore

import (
	"strconv"
)

type Status int

const (
	SUCCESS Status = iota
	FAILED
	INPROCESS
)

type LoopState struct {
	loopId        int
	callLevel     int
	iterationStep int
}

type VariableRecord struct {
	name        string
	startOffset int
}

type CallState struct {
	id           int
	returnOffset int
}

type SearchEngineState struct {
	loopStack     *Stack[LoopState]
	backtrack     *Stack[SearchEngineState]
	variableStack *Stack[VariableRecord]
	callStack     *Stack[CallState]
	environment   map[string]string

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
		if from <= value && value <= to {
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

func (es *SearchEngineState) MATCH(value string, not bool) {
	comp := es.READ(len(value))
	//fmt.Printf("is(%d) '%s' == '%s'\n", es.currentFileOffset, value, comp)
	if value == comp {
		//fmt.Println("YEAH!!")
		es.CONSUME(len(value))
		es.NEXT()
	} else {
		//fmt.Println("no :(")
		es.BACKTRACK()
	}
}

func (es *SearchEngineState) MATCHVAR(name string) {
	value, found := es.environment[name]
	if !found {
		es.BACKTRACK()
	} else {
		es.MATCH(value, false)
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

func (es *SearchEngineState) INITLOOPSTACK(loopId int) bool {
	top := es.loopStack.Peek()
	if es.loopStack.IsEmpty() || top.loopId != loopId || top.callLevel != int(es.callStack.Size()) {
		es.loopStack.Push(LoopState{
			loopId:        loopId,
			callLevel:     int(es.callStack.Size()),
			iterationStep: 0,
		})
		return true
	}
	return false
}

func (es *SearchEngineState) INCLOOPSTACK() {
	if es.loopStack.IsEmpty() {
		panic("oh crap :(")
	}
	es.loopStack.Peek().iterationStep += 1
}

func (es *SearchEngineState) GETITERATIONSTEP() int {
	if es.loopStack.IsEmpty() {
		panic("oh crap :(")
	}
	return es.loopStack.Peek().iterationStep
}

func (es *SearchEngineState) POPLOOPSTACK() LoopState {
	return *es.loopStack.Pop()
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
	es.environment[name] = value
	es.NEXT()
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
	//fmt.Println("CALLSTACK")
	//for i, s := range es.callStack.store {
	//	fmt.Printf("(%d) %d - %d\n", i, s.id, s.returnOffset)
	//}
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
		environment:       make(map[string]string),
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
	envCopy := make(map[string]string)
	for k, v := range es.environment {
		envCopy[k] = v
	}

	return &SearchEngineState{
		loopStack:         es.loopStack.Copy(),
		backtrack:         es.backtrack.Copy(),
		variableStack:     es.variableStack.Copy(),
		callStack:         es.callStack.Copy(),
		environment:       envCopy,
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
	result := Match{
		filename:    es.filename,
		matchNumber: matchNumber,
		offset:      *NewRange(es.startFileOffset, es.currentFileOffset),
		line:        *NewRange(es.startLineNum, es.currentLineNum),
		column:      *NewRange(es.startColumnNum, es.currentColumnNum),
		value:       es.currentMatch,
		variables:   make(map[string]string),
	}

	for key, value := range es.environment {
		result.variables[key] = value
	}

	return result
}

type ReplacerState struct {
	variables      map[string]string
	match          Match
	programCounter int
}

func InitReplacerState(match Match, totalMatches int) *ReplacerState {
	variables := make(map[string]string)
	for key, value := range match.variables {
		variables[key] = value
	}

	variables["totalMatches"] = strconv.Itoa(totalMatches)
	variables["matchNumber"] = strconv.Itoa(match.matchNumber)
	variables["startOffset"] = strconv.Itoa(match.offset.Start)
	variables["endOffset"] = strconv.Itoa(match.offset.Start)
	variables["lineNumber"] = strconv.Itoa(match.line.Start)
	variables["columnNumber"] = strconv.Itoa(match.column.Start)
	variables["value"] = strconv.Itoa(match.offset.Start)
	variables["filename"] = strconv.Itoa(match.offset.Start)
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
	rs.match.replacement += value
}

func (rs *ReplacerState) WRITEVAR(name string) {
	value, found := rs.variables[name]
	if found {
		rs.match.replacement += value
	}
}

func (rs *ReplacerState) Copy() *ReplacerState {
	varsCopy := make(map[string]string)
	for k, v := range rs.variables {
		varsCopy[k] = v
	}

	return &ReplacerState{
		programCounter: rs.programCounter,
		match:          rs.match,
		variables:      varsCopy,
	}
}

func (rs *ReplacerState) Set(from *ReplacerState) {
	rs.variables = from.variables
	rs.match = from.match
	rs.programCounter = from.programCounter
}

type GlobalState struct {
	subroutines map[string][]SearchInstruction
	matches     map[string]Matches
}
