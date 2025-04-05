.PHONY: build
build:
	go build

.PHONY: run
run: build
	./main

.PHONY: build-dll
build-dll:
	cmake -DCMAKE_C_COMPILER=clang -S doubly-linked-list/ -B doubly-linked-list/build
	cmake --build doubly-linked-list/build

.PHONY: run-dll
run-dll: build-dll
	./doubly-linked-list/build/doubly-linked-list