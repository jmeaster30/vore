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
	iterationStep      int
	iterationDirection int
	engineState        EngineState
}

type VariableRecord struct {
	name        string
	startOffset int
}

type EngineState struct {
	loopStack     *Stack[LoopState]
	backtrack     *Queue[EngineState]
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
}

func (es *EngineState) SEEK() {
	_, serr := es.file.Seek(int64(es.currentFileOffset), 0)
	if serr != nil {
		panic(serr)
	}
}

func (es *EngineState) READ(length int) string {
	es.SEEK()
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

func (es *EngineState) MATCH(value string) {
	if value == es.READ(len(value)) {
		es.CONSUME(len(value))
		es.NEXT()
	} else {
		es.BACKTRACK()
	}
}

func (es *EngineState) NEXT() {
	es.programCounter += 1
}

func CreateState(filename string, file *os.File, fileOffset int, lineNumber int, columnNumber int) *EngineState {
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
}
