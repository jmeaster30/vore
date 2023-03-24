go build
./vore -profile profile.prof $@
go tool pprof vore profile.prof