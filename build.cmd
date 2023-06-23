mkdir build
mkdir build\web
set GOOS=js
set GOARCH=wasm
go build -o ./build/web/herder-legacy.wasm .
copy wasm_exec.js build\web\wasm_exec.js
copy index.html build\web\index.html