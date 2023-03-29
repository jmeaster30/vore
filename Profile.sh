go build
./vore -profile profile.prof -no-output $@
go tool pprof vore profile.prof