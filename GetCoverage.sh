go test -v -coverprofile="coverage.text" ./libvore
go tool cover -html="coverage.text" -o "coverage.html"
firefox coverage.html