module github.com/jmeaster30/vore/libvore

go 1.19

require (
	github.com/jmeaster30/vore/libvore/ds v0.0.0
	github.com/jmeaster30/vore/libvore/files v0.0.0 // indirect
	github.com/jmeaster30/vore/libvore/testutils v0.0.0
	github.com/jmeaster30/vore/libvore/algo v0.0.0
	github.com/jmeaster30/vore/libvore/ast v0.0.0
	github.com/jmeaster30/vore/libvore/bytecode v0.0.0
	github.com/jmeaster30/vore/libvore/engine v0.0.0
)

replace (
	github.com/jmeaster30/vore/libvore/algo => ./algo
	github.com/jmeaster30/vore/libvore/ast => ./ast
	github.com/jmeaster30/vore/libvore/bytecode => ./bytecode
	github.com/jmeaster30/vore/libvore/ds => ./ds
	github.com/jmeaster30/vore/libvore/engine => ./engine
	github.com/jmeaster30/vore/libvore/files => ./files
	github.com/jmeaster30/vore/libvore/testutils => ./testutils
)
