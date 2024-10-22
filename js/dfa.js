class DFA {
    constructor(states, transitions, startState) {
        this.states = states;
        this.transitions = transitions;
        this.start = startState;
    }

    static fromJson(string) {
        // TODO handle invalid JSON
        const json = JSON.parse(string);
        // TODO ensure unique names
        const states = json["states"].map((s) => {
            return new State(s.name, s.accept);
        });
        const stateTable = new Map(states.map((s) => {
            return [s.name, s];
        }));
        // TODO fallible transitions
        const transitions = json["transitions"].map((s) => {
            return new Transition(
                stateTable.get(s.from),
                stateTable.get(s.to),
                s.input,
            );
        });
        // TODO ensure start is an existing state
        const start = stateTable.get(json["start"]);

        return new DFA(states, transitions, start);
    }

    simulate(input) {
        let currentState = this.start;
        for (const i of input) {
            let t = this.#findViableTransitions(currentState, i);
            if (t) {
                currentState = t[0].to;
            } else {
                return false;
            }
        }
        return currentState.accept;
    }

    #findViableTransitions(currentState, input) {
        let vts = this.transitions.filter((t) => {
            if (t.from.name === currentState.name && t.input === input) {
                return t;
            }
        });
        if (vts.length) {
            return vts;
        } else {
            return undefined;
        }
    }
}

class State {
    constructor(name, accept) {
        this.name = name;
        this.accept = accept;
    }
}

class Transition {
    constructor(from, to, input) {
        this.to = to;
        this.from = from;
        this.input = input;
    }
}

import abs from "./ab-star.js";
let jabs = JSON.stringify(abs);
let dabs = DFA.fromJson(jabs);
// console.log(dabs);
console.assert(dabs.simulate("ab"), "ab");
console.assert(dabs.simulate("abab"), "abab");
console.assert(dabs.simulate(""), "empty string");
console.assert(!dabs.simulate("aba"), "aba");
console.assert(!dabs.simulate("abc"), "abc");
console.assert(!dabs.simulate("bba"), "bba");
console.assert(!dabs.simulate("abb"), "abb");
import e0 from "./even-zeroes.js";
let je0 = JSON.stringify(e0);
let de0 = DFA.fromJson(je0);
// console.log(de0);
console.assert(de0.simulate("1110010001011"), "1110010001011");
console.assert(de0.simulate("0000"), "0000");
console.assert(!de0.simulate("110111"), "110111");
console.assert(!de0.simulate("100110"), "100110");
console.assert(!de0.simulate("1234"), "1234");
console.assert(de0.simulate("11111"), "11111");
console.assert(de0.simulate(""), "empty string");
