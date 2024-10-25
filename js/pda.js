import { Stack } from "./stack.js";

export class PDA {
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
                s.push,
                s.pop,
            );
        });
        // TODO ensure start is an existing state
        const start = stateTable.get(json["start"]);

        return new PDA(states, transitions, start);
    }

    simulate(input) {
        let currentState = this.start;
        let s = new Stack();
        for (const i of input) {
            let vts = this.#findViableTransitions(currentState, i);
            if (vts) {
                let t = vts[0];
                currentState = t.to;
                for (const _ of t.pop) {
                    if (s.empty()) {
                        return false;
                    }
                    s.pop();
                }
                for (const p of t.push) {
                    s.push(p);
                }
            } else {
                return false;
            }
        }
        return currentState.accept && s.empty();
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

// const EPSILON = [];
class Transition {
    constructor(from, to, input, push, pop) {
        this.to = to;
        this.from = from;
        this.input = input;
        this.push = push;
        this.pop = pop;
    }
}
