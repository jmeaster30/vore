package libvore

import (
	"fmt"
	"os"
)

func readFile(filename string) (*os.File, int64) {
	dat, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	fstat, ferr := dat.Stat()
	if ferr != nil {
		panic(ferr)
	}
	return dat, fstat.Size()
}

type Command interface {
	execute(string) []Match
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

func (c FindCommand) execute(filename string) []Match {
	fmt.Println("FIND COMMAND")

	file, filesize := readFile(filename)
	matches := []Match{}
	matchNumber := 0
	fileOffset := 0
	lineNumber := 1
	columnNumber := 1
	fmt.Printf("searching %s %d\n", filename, filesize)
	fmt.Printf("%t %d %d\n", c.all, c.skip, c.take)
	for c.all || matchNumber < c.skip+c.take {
		currentState := CreateState(filename, file, fileOffset, lineNumber, columnNumber)
		for currentState.status == INPROCESS {
			inst := c.body[currentState.programCounter]
			fmt.Printf("inst: %d\n", currentState.programCounter)
			inst.print()
			currentState = inst.execute(currentState)
			if currentState.status == INPROCESS && currentState.programCounter >= len(c.body) {
				currentState.SUCCESS()
			}
		}

		if currentState.status == SUCCESS && len(currentState.currentMatch) != 0 {
			fmt.Println("SUCCESS!!!!")
			foundMatch := Match{
				filename:     filename,
				matchNumber:  matchNumber + 1,
				fileOffset:   *NewRange(uint64(currentState.startFileOffset), uint64(currentState.currentFileOffset)),
				lineNumber:   *NewRange(uint64(currentState.startLineNum), uint64(currentState.currentLineNum)),
				columnNumber: *NewRange(uint64(currentState.startColumnNum), uint64(currentState.currentColumnNum)),
				value:        currentState.currentMatch,
				variables:    make(map[string]string),
			}
			matches = append(matches, foundMatch)
			fileOffset = currentState.currentFileOffset
			lineNumber = currentState.currentLineNum
			columnNumber = currentState.currentColumnNum
			matchNumber += 1
		} else {
			skipC := make([]byte, 1)
			n, err := file.ReadAt(skipC, int64(fileOffset))
			if n != 1 || err != nil {
				panic("WOW THAT IS NOT GOOD :(")
			}
			fileOffset += 1
			columnNumber += 1
			if rune(skipC[0]) == rune('\n') {
				lineNumber += 1
				columnNumber = 1
			}
		}

		if int64(fileOffset) >= filesize {
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
	return &EngineState{}
}

type MatchVariable struct {
	name string
}

func (i MatchVariable) print() {
	fmt.Println("MATCH VAR")
}

func (i MatchVariable) execute(current_state *EngineState) *EngineState {
	return &EngineState{}
}

type CallSubroutine struct {
	name string
}

func (i CallSubroutine) print() {
	fmt.Println("CALL SUB")
}

func (i CallSubroutine) execute(current_state *EngineState) *EngineState {
	return &EngineState{}
}

type Branch struct {
	left  int
	right int
}

func (i Branch) print() {
	fmt.Println("BRANCH")
}

func (i Branch) execute(current_state *EngineState) *EngineState {
	return &EngineState{}
}

type StartLoop struct {
	minLoops  int
	maxLoopes int
	loopBody  int
	exitLoop  int
}

func (i StartLoop) print() {
	fmt.Println("START LOOP")
}

func (i StartLoop) execute(current_state *EngineState) *EngineState {
	return &EngineState{}
}

type StopLoop struct {
	minLoops  int
	maxLoopes int
	startLoop int
}

func (i StopLoop) print() {
	fmt.Println("STOP LOOP")
}

func (i StopLoop) execute(current_state *EngineState) *EngineState {
	return &EngineState{}
}

type StartVarDec struct {
	name string
}

func (i StartVarDec) print() {
	fmt.Println("START VARDEC")
}

func (i StartVarDec) execute(current_state *EngineState) *EngineState {
	return &EngineState{}
}

type EndVarDec struct {
	name string
}

func (i EndVarDec) print() {
	fmt.Println("END VARDEC")
}

func (i EndVarDec) execute(current_state *EngineState) *EngineState {
	return &EngineState{}
}

type Jump struct {
	newProgramCounter int
}

func (i Jump) print() {
	fmt.Println("JUMP")
}

func (i Jump) execute(current_state *EngineState) *EngineState {
	return &EngineState{}
}
