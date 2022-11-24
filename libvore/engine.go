package libvore

import (
	"os"
)

type Status int

const (
	SUCCESS Status = iota
	FAILED
	INPROCESS
)

type LoopState struct {
	loopId        int
	iterationStep int
}

type VariableRecord struct {
	name        string
	startOffset int
}

type EngineState struct {
	loopStack     *Stack[LoopState]
	backtrack     *Queue[EngineState] // TODO this may need to be a stack
	variableStack *Stack[VariableRecord]
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
	file              *os.File
	filename          string
	filesize          int
}

func (es *EngineState) SEEK() {
	_, serr := es.file.Seek(int64(es.currentFileOffset), 0)
	if serr != nil {
		panic(serr)
	}
}

func (es *EngineState) SEEKTO(offset int) {
	_, serr := es.file.Seek(int64(offset), 0)
	if serr != nil {
		panic(serr)
	}
}

func (es *EngineState) READ(length int) string {
	es.SEEK()
	if es.currentFileOffset+length-1 >= es.filesize {
		return ""
	}
	currentString := make([]byte, length)
	n, err := es.file.Read(currentString)
	if err != nil {
		panic(err)
	}
	if n != length {
		return ""
	}
	return string(currentString)
}

func (es *EngineState) READAT(offset int, length int) string {
	es.SEEKTO(offset)
	if offset+length-1 >= es.filesize {
		return ""
	}
	currentString := make([]byte, length)
	n, err := es.file.Read(currentString)
	if err != nil {
		panic(err)
	}
	if n != length {
		return ""
	}
	return string(currentString)
}

func (es *EngineState) CONSUME(amount int) {
	value := es.READ(amount)
	es.currentMatch += value
	es.currentFileOffset += amount
	for _, c := range value {
		es.currentColumnNum += 1
		if c == rune('\n') {
			es.currentLineNum += 1
			es.currentColumnNum = 1
		}
	}
}

func (es *EngineState) BACKTRACK() {
	if es.backtrack.Size() == 0 {
		es.FAIL()
	} else {
		next_state := es.backtrack.Pop()
		es.Set(next_state)
		// TODO we may need some stuff here to better do backtracking
	}
}

func (es *EngineState) FAIL() {
	es.status = FAILED
}

func (es *EngineState) SUCCESS() {
	es.status = SUCCESS
}

func (es *EngineState) MATCHFILESTART() {
	if es.currentFileOffset == 0 {
		es.NEXT()
	} else {
		es.BACKTRACK()
	}
}

func (es *EngineState) MATCHFILEEND() {
	if es.currentFileOffset == es.filesize {
		es.NEXT()
	} else {
		es.BACKTRACK()
	}
}

func (es *EngineState) MATCHLINESTART() {
	if es.currentFileOffset == 0 {
		es.NEXT()
		return
	}

	value := es.READAT(es.currentFileOffset-1, 1)
	if value == "\n" {
		es.NEXT()
	} else {
		es.BACKTRACK()
	}
}

func (es *EngineState) MATCHLINEEND() {
	value := es.READ(1)
	if value == "\n" || es.currentFileOffset == es.filesize {
		es.NEXT()
	} else {
		es.BACKTRACK()
	}
}

func (es *EngineState) MATCHANY() {
	value := es.READ(1)
	if value == "" {
		es.BACKTRACK()
	} else {
		es.CONSUME(1)
		es.NEXT()
	}
}

func (es *EngineState) MATCHRANGE(from string, to string) {
	//? Is it possible to extend this to fit our need of range matches?
	value := es.READ(1)
	if from <= value && value <= to {
		es.CONSUME(1)
		es.NEXT()
	} else {
		es.BACKTRACK()
	}
}

func (es *EngineState) MATCHLETTER() {
	// TODO I would prefer if I had a generic way to do these multirange searches
	value := es.READ(1)
	if ("a" <= value && value <= "z") || ("A" <= value && value <= "Z") {
		es.CONSUME(1)
		es.NEXT()
	} else {
		es.BACKTRACK()
	}
}

func (es *EngineState) MATCHOPTIONS(options []string) {
	value := es.READ(1)
	if value == "" {
		es.BACKTRACK()
		return
	}

	for _, opt := range options {
		if value == opt {
			es.CONSUME(1)
			es.NEXT()
			return
		}
	}

	es.BACKTRACK()
}

func (es *EngineState) MATCH(value string) {
	if value == es.READ(len(value)) {
		es.CONSUME(len(value))
		es.NEXT()
	} else {
		es.BACKTRACK()
	}
}

func (es *EngineState) MATCHVAR(name string) {
	value, found := es.environment[name]
	if !found {
		es.BACKTRACK()
	} else {
		es.MATCH(value)
	}
}

func (es *EngineState) NEXT() {
	es.programCounter += 1
}

func (es *EngineState) JUMP(pc int) {
	es.programCounter = pc
}

func (es *EngineState) GETPC() int {
	return es.programCounter
}

func (es *EngineState) INITLOOPSTACK(loopId int) bool {
	if es.loopStack.IsEmpty() || es.loopStack.Peek().loopId != loopId {
		es.loopStack.Push(LoopState{
			loopId:        loopId,
			iterationStep: 0,
		})
		return true
	}
	return false
}

func (es *EngineState) INCLOOPSTACK() {
	if es.loopStack.IsEmpty() {
		panic("oh crap :(")
	}
	es.loopStack.Peek().iterationStep += 1
}

func (es *EngineState) GETITERATIONSTEP() int {
	if es.loopStack.IsEmpty() {
		panic("oh crap :(")
	}
	return es.loopStack.Peek().iterationStep
}

func (es *EngineState) POPLOOPSTACK() LoopState {
	return *es.loopStack.Pop()
}

func (es *EngineState) PUSHLOOPSTACK(loopState LoopState) {
	es.loopStack.Push(loopState)
}

func (es *EngineState) STARTVAR(name string) {
	record := VariableRecord{
		name:        name,
		startOffset: len(es.currentMatch),
	}
	es.variableStack.Push(record)
	es.NEXT()
}

func (es *EngineState) ENDVAR(name string) {
	record := es.variableStack.Pop()
	if record.name != name {
		panic("UHOH BAD INSTRUCTIONS I TRIED RESOLVING A VARIABLE THAT I WASN'T EXPECTING")
	}
	value := es.currentMatch[record.startOffset:]
	es.environment[name] = value
	es.NEXT()
}

func (es *EngineState) CHECKPOINT() {
	checkpoint := es.Copy()
	es.backtrack.PushFront(*checkpoint)
}

func CreateState(filename string, filesize int, file *os.File, fileOffset int, lineNumber int, columnNumber int) *EngineState {
	return &EngineState{
		loopStack:         NewStack[LoopState](),
		backtrack:         NewQueue[EngineState](),
		variableStack:     NewStack[VariableRecord](),
		environment:       make(map[string]string),
		status:            INPROCESS,
		programCounter:    0,
		currentFileOffset: fileOffset,
		currentLineNum:    lineNumber,
		currentColumnNum:  columnNumber,
		startFileOffset:   fileOffset,
		startLineNum:      lineNumber,
		startColumnNum:    columnNumber,
		file:              file,
		filename:          filename,
		filesize:          filesize,
	}
}

func (es *EngineState) Copy() *EngineState {
	envCopy := make(map[string]string)
	for k, v := range es.environment {
		envCopy[k] = v
	}

	return &EngineState{
		loopStack:         es.loopStack.Copy(),
		backtrack:         es.backtrack.Copy(),
		variableStack:     es.variableStack.Copy(),
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
		file:              es.file,
		filename:          es.filename,
		filesize:          es.filesize,
	}
}

func (es *EngineState) Set(value *EngineState) {
	es.loopStack = value.loopStack
	es.backtrack = value.backtrack
	es.variableStack = value.variableStack
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
	es.file = value.file
	es.filename = value.filename
	es.filesize = value.filesize
}

func (es *EngineState) MakeMatch(matchNumber int) Match {
	result := Match{
		filename:     es.filename,
		matchNumber:  matchNumber,
		fileOffset:   *NewRange(uint64(es.startFileOffset), uint64(es.currentFileOffset)),
		lineNumber:   *NewRange(uint64(es.startLineNum), uint64(es.currentLineNum)),
		columnNumber: *NewRange(uint64(es.startColumnNum), uint64(es.currentColumnNum)),
		value:        es.currentMatch,
		variables:    make(map[string]string),
	}

	for key, value := range es.environment {
		result.variables[key] = value
	}

	return result
}
