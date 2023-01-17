go test -v -coverprofile="coverage.text" ./libvore
go tool cover -html="coverage.text" -o "coverage.html"
Invoke-Item coverage.html