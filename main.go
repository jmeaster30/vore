package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"

	"github.com/jmeaster30/vore/libvore"
)

var replaceModeArg = libvore.NEW

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
	debug_arg := flag.Bool("debug", false, "Prints the AST of the supplied command or source file")
	files_arg := flag.String("files", "", "Files to search")
	filenames_arg := flag.Bool("filenames", false, "Process filenames instead of file contents")
	out_json_arg := flag.Bool("json", false, "Output JSON to STDOUT")
	out_fjson_arg := flag.Bool("formatted-json", false, "Output formatted JSON to STDOUT")
	json_file_arg := flag.String("json-file", "", "JSON output file")
	fjson_file_arg := flag.String("formatted-json-file", "", "Formatted JSON output file")
	no_output_arg := flag.Bool("no-output", false, "Do not output any results")
	profile_arg := flag.String("profile", "", "CPU Profile")
	flag.Func("replace-mode", "File mode for replace statements [NEW, NOTHING, OVERWRITE] (default: NEW)", replaceMode)
	flag.Parse()

	source := *source_arg
	files := *files_arg
	process_filenames := *filenames_arg
	json_file := *json_file_arg
	fjson_file := *fjson_file_arg
	out_json := *out_json_arg
	out_fjson := *out_fjson_arg
	no_output := *no_output_arg
	profile_file := *profile_arg
	command := *command_arg
	debug := *debug_arg

	if debug {
		fmt.Printf("source: '%s'\n", source)
		fmt.Printf("command: '%s'\n", command)
		fmt.Printf("files: '%s'\n", files)
		fmt.Printf("replace mode: '%s'\n", replaceModeArg)
		if debug {
			fmt.Print("D ")
		}
		if len(profile_file) != 0 {
			fmt.Print("P ")
		}
		if process_filenames {
			fmt.Print("F ")
		}
		if out_json {
			fmt.Print("j ")
		}
		if out_fjson {
			fmt.Print("J ")
		}
		if no_output {
			fmt.Print("N")
		}
		fmt.Print("\n")

	}

	if profile_file != "" {
		f, err := os.Create(profile_file)
		if err != nil {
			log.Fatal(err)
		}
		err = pprof.StartCPUProfile(f)
		if err != nil {
			log.Fatal(err)
		}
		defer pprof.StopCPUProfile()
	}

	if len(files) == 0 && !debug {
		fmt.Println("Please supply some files to search O.O")
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

	var vore *libvore.Vore
	var compError error
	if len(source) != 0 {
		vore, compError = libvore.CompileFile(source)
	} else {
		vore, compError = libvore.Compile(command)
	}

	if compError != nil {
		log.Fatal(compError)
	}

  if debug {
    vore.PrintAST()
  }

	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
  }

	searchFiles := libvore.ParsePath(*files_arg).GetFileList(currentDir)
	if len(searchFiles) == 0 {
		fmt.Println("No files to search :(")
		return
	}

	results := vore.RunFiles(searchFiles, replaceModeArg, process_filenames)

	if no_output { // skip all output
		return
	}

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
