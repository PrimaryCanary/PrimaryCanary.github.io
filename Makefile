GO = go.exe
LOX_BINARY = loxogon.exe
GOFLAGS = -gcflags="all=-N -l"

.PHONY: all
all: run-loxogon

.PHONY: run-loxogon
run-loxogon: loxogon
	cd loxogon && "./$(LOX_BINARY)"

.PHONY: loxogon
loxogon:
	$(GO) build -C loxogon $(GOFLAGS) 

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