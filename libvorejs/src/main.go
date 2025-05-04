package main

import (
	"syscall/js"

	"github.com/jmeaster30/vore/libvore"
	"github.com/jmeaster30/vore/libvore/ds"
)

func main() {
	done := make(chan struct{}, 0)
	root := js.Global().Get("__libvore__")
	root.Set("voreSearch", js.FuncOf(voreSearch))
	<-done
}

// TODO I would like to add the "compile" and "run" functions so you don't have to compile the source each search
// TODO I would also like to use promises for asynchronous code and a little better error handling

func buildRange(value ds.Range) map[string]any {
	return map[string]any{
		"start": value.Start,
		"end":   value.End,
	}
}

func buildLexError(err libvore.LexError) map[string]any {
	return map[string]any{
		"error": map[string]any{
			"type":    "LexError",
			"message": err.Message(),
			"token": map[string]any{
				"lexeme": err.Token().Lexeme,
				"type":   err.Token().TokenType.PP(),
				"offset": buildRange(err.Token().Offset),
				"line":   buildRange(err.Token().Line),
				"column": buildRange(err.Token().Column),
			},
		},
	}
}

func buildParseError(err libvore.ParseError) map[string]any {
	return map[string]any{
		"error": map[string]any{
			"type":    "ParseError",
			"message": err.Message(),
			"token": map[string]any{
				"lexeme": err.Token().Lexeme,
				"type":   err.Token().TokenType.PP(),
				"offset": buildRange(err.Token().Offset),
				"line":   buildRange(err.Token().Line),
				"column": buildRange(err.Token().Column),
			},
		},
	}
}

func buildGenError(err libvore.GenError) map[string]any {
	return map[string]any{
		"error": map[string]any{
			"type":    "GenError",
			"message": err.Message(),
		},
	}
}

func buildSemanticError(err libvore.SemanticError) map[string]any {
	return map[string]any{
		"error": map[string]any{
			"type":    "SemanticError",
			"message": err.Message(),
		},
	}
}

func buildExecError(err libvore.ExecError) map[string]any {
	return map[string]any{
		"error": map[string]any{
			"type":    "ExecError",
			"message": err.Message(),
		},
	}
}

func buildError(err error) map[string]any {
	switch e := err.(type) {
	case *libvore.LexError:
		return buildLexError(e)
	case *libvore.ParseError:
		return buildParseError(e)
	case *libvore.GenError:
		return buildGenError(e)
	case *libvore.SemanticError:
		return buildSemanticError(e)
	case *libvore.ExecError:
		return buildExecError(e)
	}

	return map[string]any{
		"error": map[string]any{
			"message": err.Error(),
		},
	}
}

func buildMatch(match libvore.Match) map[string]interface{} {
	result := map[string]interface{}{
		"filename":    match.Filename,
		"matchNumber": match.MatchNumber,
		"offset": map[string]interface{}{
			"start": match.Offset.Start,
			"end":   match.Offset.End,
		},
		"line": map[string]interface{}{
			"start": match.Line.Start,
			"end":   match.Line.End,
		},
		"column": map[string]interface{}{
			"start": match.Column.Start,
			"end":   match.Column.End,
		},
		"value": match.Value,
	}
	if match.Replacement.HasValue() {
		result["replacement"] = match.Replacement.GetValue()
	}
	return result
}

func buildMatches(input string, matches libvore.Matches) map[string]interface{} {
	convertedMatches := []interface{}{}

	// I think the libvore.Vore.Run function should ultimately return the resulting string but not quite sure if I like that
	resultString := ""
	inputIndex := 0

	for _, match := range matches {
		convertedMatches = append(convertedMatches, buildMatch(match))
		startSlice := input[inputIndex:match.Offset.Start]
		resultString += startSlice
		resultString += match.Replacement.GetValueOrDefault(match.Value)
		inputIndex = match.Offset.End
	}

	if inputIndex < len(input) {
		resultString += input[inputIndex:]
	}

	return map[string]interface{}{
		"input":   input,
		"output":  resultString,
		"matches": convertedMatches,
	}
}

func voreSearch(this js.Value, args []js.Value) any {
	source := args[0].String()
	input := args[1].String()
	resolve := args[2]
	reject := args[3]
	defer func() {
		if r := recover(); r == nil {
			reject.Invoke(js.ValueOf(map[string]interface{}{
				"error": "Aw man :( ... Go paniced",
			}))
		}
	}()

	vore, err := libvore.Compile(source)
	if err != nil {
		reject.Invoke(js.ValueOf(buildError(err)))
		return nil
	}
	matches := vore.Run(input)
	resolve.Invoke(js.ValueOf(buildMatches(input, matches)))
	return nil
}
