package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"syscall/js"

	"github.com/jmeaster30/vore/libvore"
)

func main() {
	done := make(chan struct{}, 0)
	root := js.Global().Get("__libvore__")
	root.Set("voreSearch", js.FuncOf(voreSearch))
	root.Set("voreCompile", js.FuncOf(voreCompile))
	<-done
}

func buildError(err *libvore.VoreError) map[string]any {
	return map[string]any{
		"error": map[string]any{
			"message":   err.Message,
			"token":     err.Token.Lexeme,
			"tokenType": err.Token.TokenType.PP(),
			"line": map[string]any{
				"start": err.Token.Line.Start,
				"end":   err.Token.Line.End,
			},
			"column": map[string]any{
				"start": err.Token.Column.Start,
				"end":   err.Token.Column.End,
			},
			"offset": map[string]any{
				"start": err.Token.Offset.Start,
				"end":   err.Token.Offset.End,
			},
		},
	}
}

func buildMatch(match libvore.Match) map[string]interface{} {
	result := map[string]any{
		"filename":    match.Filename,
		"matchNumber": match.MatchNumber,
		"offset": map[string]any{
			"start": match.Offset.Start,
			"end":   match.Offset.End,
		},
		"line": map[string]any{
			"start": match.Line.Start,
			"end":   match.Line.End,
		},
		"column": map[string]any{
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

func buildMatches(input string, matches libvore.Matches) map[string]any {
	var convertedMatches []any

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

	return map[string]any{
		"input":   input,
		"output":  resultString,
		"matches": convertedMatches,
	}
}

func voreSearch(this js.Value, args []js.Value) any {
	source := args[0]
	input := args[1].String()
	resolve := args[2]
	reject := args[3]
	defer func() {
		if r := recover(); r == nil {
			reject.Invoke(js.ValueOf(map[string]any{
				"error": "Aw man :( ... Go panicked",
			}))
		}
	}()

	var vore *libvore.Vore
	if source.Type() == js.TypeString {
		var err error
		vore, err = libvore.Compile(source.String())
		if err != nil {
			// TODO pass in resolve and reject promises so if there is an error we can reject and use the "catch" syntax
			var detailedErr *libvore.VoreError
			if errors.As(err, &detailedErr) {
				reject.Invoke(js.ValueOf(buildError(detailedErr)))
			}
			return nil
		}
	} else if source.Type() == js.TypeObject {
		bytecodeBytes := source.Get("bytecode")
		if bytecodeBytes.IsUndefined() || bytecodeBytes.IsNull() {
			reject.Invoke(js.ValueOf(map[string]any{
				"error": "source's bytecode was null or undefined :(",
			}))
			return nil
		}

		//js.Global().Get("console").Call("log", fmt.Sprintf("STRING: %s", bytecodeBytes.String()))

		buf := []byte(bytecodeBytes.String())
		//transform the byte array to map
		var bytecodeCommands []map[string]any
		err := json.Unmarshal(buf, &bytecodeCommands)
		if err != nil {
			//js.Global().Get("console").Call("log", fmt.Sprintf("ERROR: %+v", err))
			reject.Invoke(js.ValueOf(map[string]any{
				"error": "Failed to unmarshal json object :(",
			}))
			return nil
		}

		js.Global().Get("console").Call("log", fmt.Sprintf("Getting commands... %+v", bytecodeCommands))
		commands, err := libvore.CommandsFromMap(bytecodeCommands)
		if err != nil {
			js.Global().Get("console").Call("log", fmt.Sprintf("ERROR: %+v", err))
			reject.Invoke(js.ValueOf(map[string]any{
				"error": "Failed to unmarshal json object :(",
			}))
			return nil
		}
		js.Global().Get("console").Call("log", fmt.Sprintf("Commands: %+v", commands))

		vore = libvore.Build(commands)
	}

	js.Global().Get("console").Call("log", fmt.Sprintf("Vore: %+v", vore))
	matches := vore.Run(input)
	resolve.Invoke(js.ValueOf(buildMatches(input, matches)))

	return nil
}

func voreCompile(this js.Value, args []js.Value) any {
	// return an object with a member called "bytecode" that is a uint8array
	source := args[0].String()
	resolve := args[1]
	reject := args[2]
	defer func() {
		if r := recover(); r == nil {
			reject.Invoke(js.ValueOf(map[string]any{
				"error": "Aw man :( ... Go panicked",
			}))
		}
	}()

	vore, err := libvore.Compile(source)
	if err != nil {
		// TODO pass in resolve and reject promises so if there is an error we can reject and use the "catch" syntax
		var detailedErr *libvore.VoreError
		if errors.As(err, &detailedErr) {
			reject.Invoke(js.ValueOf(buildError(detailedErr)))
		}
		return nil
	}

	//js.Global().Get("console").Call("log", fmt.Sprintf("VORE: %+v", vore))

	var serializedBytecode []map[string]any
	for _, command := range vore.Bytecode() {
		serializedBytecode = append(serializedBytecode, command.ToMap())
	}

	//js.Global().Get("console").Call("log", fmt.Sprintf("Serialized: %+v", serializedBytecode))

	bytecodeBytes, err := json.Marshal(serializedBytecode)
	if err != nil {
		reject.Invoke(js.ValueOf(map[string]any{
			"error": "Failed to marshal GO object :(",
		}))
		return nil
	}

	//js.Global().Get("console").Call("log", fmt.Sprintf("BYTES: %s", bytecodeBytes))

	resolve.Invoke(js.ValueOf(map[string]any{
		"bytecode": string(bytecodeBytes[:]),
	}))
	return nil
}
