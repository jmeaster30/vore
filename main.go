package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/jmeaster30/vore/libvore"
)

var replaceModeArg libvore.ReplaceMode

func replaceMode(value string) error {
	switch value {
	case "OVERWRITE":
		replaceModeArg = libvore.OVERWRITE
	case "NOTHING":
		replaceModeArg = libvore.NOTHING
	case "":
		fallthrough
	case "NEW":
		replaceModeArg = libvore.NEW
	default:
		return errors.New("Expected [NEW, NOTHING, OVERWRITE] but got '" + value + "'.")
	}
	return nil
}

func main() {
	source_arg := flag.String("src", "", "Vore source file to run on search files")
	command_arg := flag.String("com", "", "Vore command to run on search files")
	files_arg := flag.String("files", "", "Files to search")
	out_json_arg := flag.Bool("json", false, "JSON output file")
	out_fjson_arg := flag.Bool("formatted-json", false, "Formatted JSON output file")
	json_file_arg := flag.String("json-file", "", "JSON output file")
	fjson_file_arg := flag.String("formatted-json-file", "", "Formatted JSON output file")
	ide_arg := flag.Bool("ide", false, "Open source and files in vore ide")
	flag.Func("replace-mode", "File mode for replace statements [OVERWRITE, NEW] (default: NEW)", replaceMode)
	flag.Parse()

	source := *source_arg
	files := *files_arg
	json_file := *json_file_arg
	fjson_file := *fjson_file_arg
	out_json := *out_json_arg
	out_fjson := *out_fjson_arg
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

	if out_json && out_fjson {
		fmt.Println("Can't output both json and formatted json to stdout.")
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
	results := vore.RunFiles([]string{*files_arg}, replaceModeArg)
	if len(results) == 0 {
		fmt.Println("There were no matches :(")
	} else {
		if !out_json && !out_fjson {
			fmt.Printf("There were %d matches :)\n", len(results))
		}
		if len(json_file) != 0 {
			f := OpenFile(json_file)
			Truncate(f)
			f.WriteString(results.Json())
		}
		if len(fjson_file) != 0 {
			f := OpenFile(fjson_file)
			Truncate(f)
			f.WriteString(results.FormattedJson())
		}
		if out_json {
			fmt.Println(results.Json())
		} else if out_fjson {
			fmt.Println(results.FormattedJson())
		} else {
			results.Print()
		}
	}
}

func OpenFile(filename string) *os.File {
	f, err := os.OpenFile(filename, os.O_CREATE, os.ModeAppend)
	if err != nil {
		panic(err)
	}
	return f
}

func Truncate(f *os.File) {
	terr := f.Truncate(0)
	if terr != nil {
		panic(terr)
	}
	_, serr := f.Seek(0, 0)
	if serr != nil {
		panic(serr)
	}
}
