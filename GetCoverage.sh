go test -v -coverprofile="coverage.text" \
  ./libvore \
  ./libvore/algo \
  ./libvore/ast \
  ./libvore/bytecode \
  ./libvore/ds \
  ./libvore/engine \
  ./libvore/files
go tool cover -html="coverage.text" -o "coverage.html"
firefox coverage.html
