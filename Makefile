.PHONY: all
all: run-loxogon

.PHONY: run
run-loxogon: loxogon
	cd loxogon && "./loxogon"

.PHONY: loxogon
loxogon:
	go build -C loxogon

.PHONY: build-dll
build-dll:
	cmake -DCMAKE_C_COMPILER=clang -S doubly-linked-list/ -B doubly-linked-list/build
	cmake --build doubly-linked-list/build

.PHONY: run-dll
run-dll: build-dll
	./doubly-linked-list/build/doubly-linked-list