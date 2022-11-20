package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jmeaster30/vore/libvore"
)

func main() {
	source_arg := flag.String("src", "", "Vore source file to run on search files")
	command_arg := flag.String("com", "", "Vore command to run on search files")
	files_arg := flag.String("files", "", "Files to search")
	ide_arg := flag.Bool("ide", false, "Open source and files in vore ide")
	flag.Parse()

	source := *source_arg
	files := *files_arg
	ide := *ide_arg
	command := *command_arg

	if ide {
		fmt.Println("Sorry repl mode is not implemented :(")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if len(files) == 0 {
		fmt.Println("Please supply some files to search.")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if len(source) != 0 && len(command) != 0 {
		fmt.Println("Cannot use both a source file and a command at the same time.")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if len(source) == 0 && len(command) == 0 {
		fmt.Println("Must supply either a source file or a command.")
		flag.PrintDefaults()
		os.Exit(1)
	}

	var vore libvore.Vore
	if len(source) != 0 {
		vore = libvore.CompileFile(source)
	} else {
		vore = libvore.Compile(command)
	}

	//vore.PrintTokens()
	//vore.PrintAST()
	results := vore.Run([]string{*files_arg})
	for _, match := range results {
		match.Print()
	}
}
