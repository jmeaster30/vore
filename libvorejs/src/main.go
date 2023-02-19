package main

import (
	"syscall/js"

	"github.com/jmeaster30/vore/libvore"
)

func main() {
	done := make(chan struct{}, 0)
	root := js.Global().Get("__libvore__")
	root.Set("voreSearch", js.FuncOf(voreSearch))
	<-done
}

// TODO I would like to add the "compile" and "run" functions so you don't have to compile the source each search
// TODO I would also like to use promises for asynchronous code and a little better error handling

func buildError(err *libvore.VoreError) map[string]interface{} {
	return map[string]interface{}{
		"error": map[string]interface{}{
			"message":   err.Message,
			"token":     err.Token.Lexeme,
			"tokenType": err.Token.TokenType.PP(),
			"line": map[string]interface{}{
				"start": err.Token.Line.Start,
				"end":   err.Token.Line.End,
			},
			"column": map[string]interface{}{
				"start": err.Token.Column.Start,
				"end":   err.Token.Column.End,
			},
			"offset": map[string]interface{}{
				"start": err.Token.Offset.Start,
				"end":   err.Token.Offset.End,
			},
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
		resultString += input[inputIndex:len(input)]
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
		// TODO pass in resolve and reject promises so if there is an error we can reject and use the "catch" syntax
		if detailedErr, ok := err.(*libvore.VoreError); ok {
			reject.Invoke(js.ValueOf(buildError(detailedErr)))
		} else {
			reject.Invoke(js.ValueOf(map[string]interface{}{
				"error": err.Error(),
			}))
		}
		return nil
	}
	matches := vore.Run(input)
	resolve.Invoke(js.ValueOf(buildMatches(input, matches)))
	return nil
}
