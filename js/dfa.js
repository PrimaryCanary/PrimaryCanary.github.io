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

// import abs from "./ab-star.js";
// let jabs = JSON.stringify(abs);
// console.log(DFA.fromJson(jabs));
