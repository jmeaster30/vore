mkdir static
cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" ./static/wasm_exec.js
GOOS=js GOARCH=wasm go build -o ./static/libvore.wasm ./wasm