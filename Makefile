GO = go.exe
GOFLAGS = 
LOX_BINARY = cmd/loxogon/loxogon.exe
PWSH = pwsh.exe -Command

.PHONY: all
all: loxogon-wasm test run-loxogon

.PHONY: run-loxogon
run-loxogon: loxogon
	cd loxogon && "./$(LOX_BINARY)"

.PHONY: loxogon
loxogon:
	$(GO) build -C loxogon/cmd/loxogon $(GOFLAGS) 

.PHONY: loxogon-wasm
loxogon-wasm:
	$(PWSH) '$$env:GOOS="js"; $$env:GOARCH="wasm"; go build -C loxogon/cmd/wasm -o loxogon.wasm'
	cp loxogon/cmd/wasm/loxogon.wasm playground/

.PHONY: test
test:
	cd loxogon && $(GO) test ./...

.PHONY: build-dll
build-dll:
	cmake -DCMAKE_C_COMPILER=clang -S doubly-linked-list/ -B doubly-linked-list/build
	cmake --build doubly-linked-list/build

.PHONY: run-dll
run-dll: build-dll
	./doubly-linked-list/build/doubly-linked-list