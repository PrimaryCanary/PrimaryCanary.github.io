export class Stack {
    #buffer;
    constructor() {
        this.#buffer = [];
    }

    peek() {
        if (!this.empty()) {
            return this.#buffer[this.#buffer.length - 1];
        } else {
            return undefined;
        }
    }

    push(data) {
        this.#buffer.push(data);
    }
    pop() {
        if (!this.empty()) {
            return this.#buffer.pop();
        } else {
            return undefined;
        }
    }
    empty() {
        return this.#buffer.length === 0;
    }
}
