import { Stack } from "./stack.js";
import { assert } from "jsr:@std/assert";

Deno.test("Stack is empty", () => {
  let s = new Stack();
  assert(s.empty());
});
Deno.test("Empty peek", () => {
  let s = new Stack();
  assert(!s.peek());
});
Deno.test("Empty pop", () => {
  let s = new Stack();
  assert(!s.pop());
});

let s = new Stack();
s.push(1);
s.push(2);
s.push(3);
for (const i of [3, 2, 1]) {
  Deno.test("Peek " + i, () => {
    assert(s.peek() === i);
  });
  Deno.test("Pop " + i, () => {
    let r = s.pop();
    assert(r === i);
  });
}
