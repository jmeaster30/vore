module github.com/jmeaster30/vore

go 1.19

require (
	github.com/jmeaster30/vore/libvore v0.0.0
	github.com/jmeaster30/vore/libvore/engine v0.0.0
)

replace (
	github.com/jmeaster30/vore/libvore => ./libvore
	github.com/jmeaster30/vore/libvore/engine => ./libvore/engine
)
