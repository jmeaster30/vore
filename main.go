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
	var _ = *files_arg // TODO
	ide := *ide_arg
	command := *command_arg

	if ide {
		fmt.Println("Sorry repl mode is not implemented :(")
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

	if len(source) != 0 {
		// Use source file
		test := libvore.CompileFile(source)
		test.PrintTokens()
	} else {
		// Use command
		test := libvore.Compile(command)
		test.PrintTokens()
	}
}
