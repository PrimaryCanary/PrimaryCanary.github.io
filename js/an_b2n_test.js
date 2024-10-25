import { PDA } from "./pda.js";
import { assert } from "jsr:@std/assert";
import anb2n from "./automata/an_b2n.js";

Deno.test("stringify => fromJson", () => {
    const json = JSON.stringify(anb2n);
    const hooray = PDA.fromJson(json);
    assert(hooray);
});

const json = JSON.stringify(anb2n);
const pda = PDA.fromJson(json);
const success = ["", "aaabbbbbb", "abb", "aaaaabbbbbbbbbb"];
const fail = ["aabbbba", "aabbb", "abbb", "abc"];
for (const s of success) {
    Deno.test(s + " accepted", () => {
        assert(pda.simulate(s));
    });
}
for (const s of fail) {
    Deno.test(s + " rejected", () => {
        assert(!pda.simulate(s));
    });
}
