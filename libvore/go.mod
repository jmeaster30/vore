module github.com/jmeaster30/vore/libvore

go 1.19

require (
	github.com/jmeaster30/vore/libvore/ds v0.0.0
	github.com/jmeaster30/vore/libvore/files v0.0.0
	github.com/jmeaster30/vore/libvore/testutils v0.0.0
	github.com/jmeaster30/vore/libvore/algo v0.0.0
)

replace github.com/jmeaster30/vore/libvore/algo => ./algo

replace github.com/jmeaster30/vore/libvore/ds => ./ds

replace github.com/jmeaster30/vore/libvore/files => ./files

replace github.com/jmeaster30/vore/libvore/testutils => ./testutils
