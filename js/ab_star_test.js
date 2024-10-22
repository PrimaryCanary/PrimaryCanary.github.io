import { DFA } from "./dfa.js";
import { assert } from "jsr:@std/assert";
import abs from "./automata/ab_star.js";

Deno.test("stringify => fromJson", () => {
    const jabs = JSON.stringify(abs);
    const dabs = DFA.fromJson(jabs);
    assert(dabs);
});

const jabs = JSON.stringify(abs);
const dabs = DFA.fromJson(jabs);
const success = ["ab", "abab", ""];
const fail = ["aba", "abc", "bba", "abb"];
for (const s of success) {
    Deno.test(s + " accepted", () => {
        assert(dabs.simulate(s));
    });
}
for (const s of fail) {
    Deno.test(s + " rejected", () => {
        assert(!dabs.simulate(s));
    });
}
