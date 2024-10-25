import { PDA } from "./pda.js";
import { assert } from "jsr:@std/assert";
import e0 from "./automata/even_zeroes.js";

Deno.test("stringify => fromJson", () => {
    const je0 = JSON.stringify(e0);
    const de0 = PDA.fromJson(je0);
    assert(de0);
});

const je0 = JSON.stringify(e0);
const de0 = PDA.fromJson(je0);
const success = ["1110010001011", "0000", "11111", ""];
const fail = ["110111", "100110", "1234"];
for (const s of success) {
    Deno.test(s + " accepted", () => {
        assert(de0.simulate(s));
    });
}
for (const s of fail) {
    Deno.test(s + " rejected", () => {
        assert(!de0.simulate(s));
    });
}
