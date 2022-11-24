package libvore

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Match struct {
	filename    string
	matchNumber int
	offset      Range
	line        Range
	column      Range
	value       string
	variables   map[string]string
}

type Matches []Match

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
			result += "\\"
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
	result += "\t\"variables\": [\n"

	keys := make([]string, 0, len(m.variables))
	for k := range m.variables {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	vars := []string{}
	for _, k := range keys {
		key := k
		value := m.variables[k]
		vars = append(vars, "\t\t{\n\t\t\t\""+key+"\": \""+cleanControlCharacters(value)+"\"\n\t\t}")
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
	result += "\"variables\":["

	keys := make([]string, 0, len(m.variables))
	for k := range m.variables {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	vars := []string{}
	for _, k := range keys {
		key := k
		value := m.variables[k]
		vars = append(vars, "{\""+key+"\":\""+cleanControlCharacters(value)+"\"}")
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

func (m Match) Print() {
	fmt.Printf("Filename: %s\n", m.filename)
	fmt.Printf("MatchNumber: %d\n", m.matchNumber)
	fmt.Printf("Value: %s\n", m.value)
	fmt.Printf("Offset: %d %d\n", m.offset.Start, m.offset.End)
	fmt.Printf("Line: %d %d\n", m.line.Start, m.line.End)
	fmt.Printf("Column: %d %d\n", m.column.Start, m.column.End)
	fmt.Println("Variables:")
	fmt.Println("\t[key] = [value]")
	for key, value := range m.variables {
		fmt.Printf("\t%s = %s\n", key, value)
	}
}

type Vore struct {
	tokens   []*Token
	commands []AstCommand
	bytecode []Command
}

func Compile(command string) Vore {
	return compile("source", strings.NewReader(command))
}

func CompileFile(source string) Vore {
	dat, err := os.Open(source)
	if err != nil {
		panic(err)
	}
	return compile(source, dat)
}

func compile(filename string, reader io.Reader) Vore {
	lexer := initLexer(reader)

	tokens := lexer.getTokens()
	commands, parseError := parse(tokens)
	if parseError.isError {
		panic(fmt.Sprintf("\nERROR:  %s\nToken:  '%s'\nLine:   %d - %d\nColumn: %d - %d\n", parseError.message, parseError.token.lexeme, parseError.token.line.Start, parseError.token.line.End, parseError.token.column.Start, parseError.token.column.End))
	}

	bytecode := []Command{}
	for _, ast_comm := range commands {
		byte_comm := ast_comm.generate()
		bytecode = append(bytecode, byte_comm)
	}

	return Vore{tokens, commands, bytecode}
}

func (v *Vore) Run(filenames []string) Matches {
	result := []Match{}
	for _, command := range v.bytecode {
		//command.print()
		for _, filename := range filenames {
			result = append(result, command.execute(filename)...)
		}
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