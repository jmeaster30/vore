# Getting Started

## Installation
---

I don't know how go packaging and installation works so just install go and run the following commands to build the executable from source. `libvore` uses external dependencies other than `go` itself so you shouldn't have to worry about making sure anything else is installed.

```bash
git clone https://github.com/jmeaster30/vore.git
cd ./vore
go build .
```

## Basics of Command Line Interface
---

The command line interface was meant to be designed for simplicity so there aren't any esoteric flags that control the things you need. Each flag has a name that describes what the flag is for.

In the command line, run:
```bash
./vore -help
```

Output:
```
Usage of vore:
  -com string
        Vore command to run on search files
  -files string
        Files to search
  -formatted-json
        Formatted JSON output file
  -formatted-json-file string
        Formatted JSON output file
  -ide
        Open source and files in vore ide
  -json
        JSON output file
  -json-file string
        JSON output file
  -replace-mode value
        File mode for replace statements [NEW, NOTHING, OVERWRITE] (default: NEW)
  -src string
        Vore source file to run on search files
```


### Hello World Vore Command
---

Open a text editor for the file we want to search and add:

```text
Hello, Lilith
```

Save the file as "HelloLilith.txt"

In the command line, run:

```bash
./vore -com "find all 'Hello, ' (at least 1 letter) = name" -files "HelloLilith.txt"
```

### Hello World Vore Source Script
---

Open a text editor for the file we want to search and add:

```text
Hello, Lilith
```

Save the file as "HelloLilith.txt"

Open your text editor for the source file that will contain our vore script:

```vore
find all
  'Hello, '
  (at least 1 letter) = name
```

Save the file as "HelloName.vore"

In the command line, run:

```bash
./vore -src "HelloName.vore" -files "HelloLilith.txt"
```

### JSON Output To Console
---

Follwing the section [Hello World Vore Source Script](#hello-world-vore-source-script).

In the command line, run:

```bash
./vore -src "HelloName.vore" -files "HelloLilith.txt" -json
```

### Formatted JSON Output To Console
---

Follwing the section [Hello World Vore Source Script](#hello-world-vore-source-script).

In the command line, run:

```bash
./vore -src "HelloName.vore" -files "HelloLilith.txt" -formatted-json
```

### JSON Output To File
---

Follwing the section [Hello World Vore Source Script](#hello-world-vore-source-script).

In the command line, run:

```bash
./vore -src "HelloName.vore" -files "HelloLilith.txt" -json-file "lilith.output"
```

### Formatted JSON Output To File
---

Follwing the section [Hello World Vore Source Script](#hello-world-vore-source-script).

In the command line, run:

```bash
./vore -src "HelloName.vore" -files "HelloLilith.txt" -formatted-json-file "lilith.formatted.output"
```
