package libvore

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

type ReplaceMode int

const (
	OVERWRITE ReplaceMode = iota
	CONFIRM
	NEW
	NOTHING
)

type Match struct {
	filename    string
	matchNumber int
	offset      Range
	line        Range
	column      Range
	value       string
	replacement string
	variables   ValueHashMap
}

type ValueType int

const (
	ValueStringType ValueType = iota
	ValueHashMapType
)

// TODO implement good interface for reading the values from these Value objects
type Value interface {
	String() ValueString
	Hashmap() ValueHashMap

	getType() ValueType
	// add an interface that runs provided functions on each type of value
	process(hashmapFunc func(ValueHashMap), stringFunc func(ValueString))
	Copy() Value
}

type ValueString struct {
	Value string
}

func NewValueString(value string) ValueString {
	return ValueString{value}
}

func (v ValueString) String() ValueString {
	return v
}

func (v ValueString) Hashmap() ValueHashMap {
	result := ValueHashMap{
		Value: make(map[string]Value),
	}

	result.Value["value"] = NewValueString(v.Value)
	return result
}

func (v ValueString) Copy() Value {
	return NewValueString(v.Value)
}

func (v ValueString) getType() ValueType {
	return ValueStringType
}

func (v ValueString) process(hashmapFunc func(ValueHashMap), stringFunc func(ValueString)) {
	stringFunc(v)
}

type ValueHashMap struct {
	Value map[string]Value
}

func NewValueHashMap() ValueHashMap {
	return ValueHashMap{
		Value: make(map[string]Value),
	}
}

func (v ValueHashMap) String() ValueString {
	return NewValueString("[ValueHashMap]")
}

func (v ValueHashMap) Hashmap() ValueHashMap {
	return v
}

func (v ValueHashMap) Copy() Value {
	result := NewValueHashMap()
	for k, val := range v.Value {
		result.Add(k, val.Copy())
	}
	return result
}

func (v ValueHashMap) getType() ValueType {
	return ValueHashMapType
}

func (v ValueHashMap) process(hashmapFunc func(ValueHashMap), stringFunc func(ValueString)) {
	hashmapFunc(v)
}

func (v ValueHashMap) Get(name string) (Value, bool) {
	val, found := v.Value[name]
	return val, found
}

func (v ValueHashMap) Add(name string, value Value) {
	v.Value[name] = value
}

func (v ValueHashMap) Len() int {
	return len(v.Value)
}

func (v ValueHashMap) Keys() []string {
	res := []string{}
	for k := range v.Value {
		res = append(res, k)
	}
	return res
}

type Matches []Match

// TODO convert this to clean all characters of the new ValueHashMap
func cleanControlCharacters(s string) string {
	result := ""
	for _, c := range s {
		switch c {
		case '\t':
			result += "\\t"
		case '\r':
			result += "\\r"
		case '\n':
			result += "\\n"
		case '\\':
			result += "\\\\"
		case '\a':
			result += "\\a"
		case '\b':
			result += "\\b"
		case '\f':
			result += "\\f"
		case '\v':
			result += "\\v"
		default:
			result += string(c)
		}
	}
	return result
}

func (m Matches) FormattedJson() string {
	result := "[\n"
	for i, match := range m {
		format := strings.Split(match.FormattedJson(), "\n")
		for j, line := range format {
			result += "\t" + line
			if j < len(format)-1 {
				result += "\n"
			}
		}

		if i < len(m)-1 {
			result += ",\n"
		} else {
			result += "\n"
		}
	}
	result += "]"
	return result
}

func (m Match) FormattedJson() string {
	result := "{\n"

	result += "\t\"filename\": \"" + m.filename + "\",\n"
	result += "\t\"matchNumber\": \"" + strconv.Itoa(m.matchNumber) + "\",\n"
	result += "\t\"offset\": {\n\t\t\"start\": \"" + strconv.Itoa(m.offset.Start) + "\",\n\t\t\"end\": \"" + strconv.Itoa(m.offset.End) + "\"\n\t},\n"
	result += "\t\"line\": {\n\t\t\"start\": \"" + strconv.Itoa(m.line.Start) + "\",\n\t\t\"end\": \"" + strconv.Itoa(m.line.End) + "\"\n\t},\n"
	result += "\t\"column\": {\n\t\t\"start\": \"" + strconv.Itoa(m.column.Start) + "\",\n\t\t\"end\": \"" + strconv.Itoa(m.column.End) + "\"\n\t},\n"
	result += "\t\"value\": \"" + cleanControlCharacters(m.value) + "\",\n"
	result += "\t\"replaced\": \"" + cleanControlCharacters(m.replacement) + "\",\n"
	result += "\t\"variables\": [\n"

	keys := m.variables.Keys()
	sort.Strings(keys)

	vars := []string{}
	for _, k := range keys {
		key := k
		value, _ := m.variables.Get(k)
		// TODO allow for outputing nested values in the hashmap
		vars = append(vars, "\t\t{\n\t\t\t\""+key+"\": \""+cleanControlCharacters(value.String().Value)+"\"\n\t\t}")
	}

	for i, v := range vars {
		result += v
		if i < len(vars)-1 {
			result += ",\n"
		}
	}

	result += "\n\t]\n}"
	return result
}

func (m Matches) Json() string {
	result := "["
	for i, match := range m {
		result += match.Json()
		if i < len(m)-1 {
			result += ","
		}
	}
	result += "]"
	return result
}

func (m Match) Json() string {
	result := "{"
	result += "\"filename\":\"" + m.filename + "\","
	result += "\"matchNumber\":\"" + strconv.Itoa(m.matchNumber) + "\","
	result += "\"offset\":{\"start\":\"" + strconv.Itoa(m.offset.Start) + "\",\"end\":\"" + strconv.Itoa(m.offset.End) + "\"},"
	result += "\"line\":{\"start\":\"" + strconv.Itoa(m.line.Start) + "\",\"end\":\"" + strconv.Itoa(m.line.End) + "\"},"
	result += "\"column\":{\"start\":\"" + strconv.Itoa(m.column.Start) + "\",\"end\":\"" + strconv.Itoa(m.column.End) + "\"},"
	result += "\"value\":\"" + cleanControlCharacters(m.value) + "\","
	result += "\"replaced\":\"" + cleanControlCharacters(m.replacement) + "\","
	result += "\"variables\":["

	keys := m.variables.Keys()
	sort.Strings(keys)

	vars := []string{}
	for _, k := range keys {
		key := k
		value, _ := m.variables.Get(k)
		// TODO allow for outputing nested ValueHashMaps
		vars = append(vars, "{\""+key+"\":\""+cleanControlCharacters(value.String().Value)+"\"}")
	}

	for i, v := range vars {
		result += v
		if i < len(vars)-1 {
			result += ","
		}
	}

	result += "]}"
	return result
}

func (m Matches) Print() {
	fmt.Println("============")
	for _, match := range m {
		match.Print()
		fmt.Println("============")
	}
}

var matchPrintDepth int

func printHashmap(hashmap ValueHashMap) {
	matchPrintDepth += 1
	keys := hashmap.Keys()
	sort.Strings(keys)
	for _, k := range keys {
		v, _ := hashmap.Get(k)
		fmt.Printf("\n%s'%s' = ", strings.Repeat("  ", matchPrintDepth), k)
		v.process(printHashmap, printString)
	}
	matchPrintDepth -= 1
}

func printString(str ValueString) {
	fmt.Printf("'%s'", str.Value)
}

func (m Match) Print() {
	fmt.Printf("Filename: %s\n", m.filename)
	fmt.Printf("MatchNumber: %d\n", m.matchNumber)
	fmt.Printf("Value: %s\n", m.value)
	fmt.Printf("Replaced: %s\n", m.replacement)
	fmt.Printf("Offset: %d %d\n", m.offset.Start, m.offset.End)
	fmt.Printf("Line: %d %d\n", m.line.Start, m.line.End)
	fmt.Printf("Column: %d %d\n", m.column.Start, m.column.End)
	fmt.Println("Variables:")
	fmt.Print("  [key] = [value]")

	matchPrintDepth = 0

	m.variables.process(printHashmap, printString)
	fmt.Println()
}

type Vore struct {
	tokens   []*Token
	commands []AstCommand
	bytecode []Command
}

func Compile(command string) (*Vore, error) {
	return compile("source", strings.NewReader(command))
}

func CompileFile(source string) (*Vore, error) {
	dat, err := os.Open(source)
	if err != nil {
		return nil, err
	}
	return compile(source, dat)
}

func compile(filename string, reader io.Reader) (*Vore, error) {
	lexer := initLexer(reader)

	tokens := lexer.getTokens()
	//for _, token := range tokens {
	//	fmt.Printf("[%s] '%s' \tline: %d, \tstart column: %d, \tend column: %d\n", token.tokenType.pp(), token.lexeme, token.line.Start, token.column.Start, token.column.End)
	//}

	commands, parseError := parse(tokens)
	if parseError.isError {
		return nil, fmt.Errorf("ERROR:  %s\nToken:  '%s'\nTokenType: %d\nLine:   %d - %d\nColumn: %d - %d", parseError.message, parseError.token.lexeme, parseError.token.tokenType, parseError.token.line.Start, parseError.token.line.End, parseError.token.column.Start, parseError.token.column.End)
	}

	bytecode := []Command{}
	gen_state := &GenState{
		globalSubroutines:     make(map[string]GeneratedPattern),
		globalVariables:       make(map[string]int),
		globalTransformations: make(map[string]AstProcessProgram),
	}
	for _, ast_comm := range commands {
		byte_comm, gen_error := ast_comm.generate(gen_state)
		if gen_error != nil {
			return nil, gen_error
		}
		bytecode = append(bytecode, byte_comm)
	}

	return &Vore{tokens, commands, bytecode}, nil
}

func (v *Vore) RunFiles(filenames []string, mode ReplaceMode, processFilenames bool) Matches {
	actualMode := mode
	if processFilenames {
		actualMode = NOTHING
	}
	result := Matches{}
	for _, command := range v.bytecode {
		//command.print()
		for _, filename := range filenames {
			actualFiles := []string{}
			info, err := os.Stat(filename)
			if err != nil {
				panic(err)
			}
			fixedFilename := filename
			if info.IsDir() {
				if filename[len(filename)-1] != '/' || filename[len(filename)-1] != '\\' {
					fixedFilename += "/"
				}
				entries, err := os.ReadDir(filename)
				if err != nil {
					panic(err)
				}
				for _, entry := range entries {
					actualFiles = append(actualFiles, fixedFilename+entry.Name())
				}
			} else {
				actualFiles = append(actualFiles, fixedFilename)
			}
			for _, actualFilename := range actualFiles {
				var reader *VReader
				if processFilenames {
					reader = VReaderFromString(actualFilename)
				} else {
					reader = VReaderFromFile(actualFilename)
				}
				foundMatches := command.execute(actualFilename, reader, actualMode)
				result = append(result, foundMatches...)
				if processFilenames && len(foundMatches) != 0 && len(foundMatches[0].replacement) != 0 {
					err := os.Rename(actualFilename, foundMatches[0].replacement)
					if err != nil {
						os.Stderr.WriteString("Failed to rename file '" + actualFilename + "' to '" + foundMatches[0].replacement + "'\n")
					}
				}
			}
		}
	}
	return result
}

func (v *Vore) Run(searchText string) Matches {
	result := Matches{}
	for _, command := range v.bytecode {
		reader := VReaderFromString(searchText)
		result = append(result, command.execute("text", reader, NOTHING)...)
	}
	return result
}

func (v *Vore) PrintTokens() {
	for _, token := range v.tokens {
		fmt.Printf("[%s] '%s' \tline: %d, \tstart column: %d, \tend column: %d\n", token.tokenType.pp(), token.lexeme, token.line.Start, token.column.Start, token.column.End)
	}
}

func (v *Vore) PrintAST() {
	for _, command := range v.commands {
		command.print()
	}
	fmt.Println()
}
