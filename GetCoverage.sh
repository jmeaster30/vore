go test -v -coverprofile="coverage.text" ./libvore
go tool cover -html="coverage.text" -o "coverage.html"
open coverage.html