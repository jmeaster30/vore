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

// func buildMatches(matches libvore.Matches) map[string]interface{} {

// }

func voreSearch(this js.Value, args []js.Value) interface{} {
	vore, err := libvore.Compile(args[0].String())
	if err != nil {
		if detailedErr, ok := err.(*libvore.VoreError); ok {
			return js.ValueOf(buildError(detailedErr))
		} else {
			return js.ValueOf(map[string]interface{}{
				"error": err.Error(),
			})
		}

	}
	matches := vore.Run(args[1].String())
	return js.ValueOf(map[string]interface{}{"numberOfMatches": len(matches)})
}
