package main

import (
	"fmt"
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

func voreSearch(this js.Value, args []js.Value) any {
	fmt.Println("HERE!!!!")
	returnObject := make(map[string]interface{})
	vore, err := libvore.Compile(args[0].String())
	if err != nil {
		fmt.Println(err)
		returnObject["error"] = err.Error()
		return js.ValueOf(returnObject)
	}
	matches := vore.Run(args[1].String())
	returnObject["matches"] = matches
	fmt.Printf("THERE WERE %d MATCHES", len(matches))
	return returnObject
}
