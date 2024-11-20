export default {
  "description": "PDA for a^n b^2n i.e. n a's followed by 2n b's",
  "states": [
    { "name": "empty", "accept": true },
    { "name": "eat a's", "accept": false },
    { "name": "eat b's", "accept": true },
  ],
  "transitions": [
    {
      "from": "empty",
      "to": "eat a's",
      "input": "a",
      "push": ["a", "a"],
      "pop": [],
    },
    {
      "from": "eat a's",
      "to": "eat a's",
      "input": "a",
      "push": ["a", "a"],
      "pop": [],
    },
    {
      "from": "eat a's",
      "to": "eat b's",
      "input": "b",
      "push": [],
      "pop": ["a"],
    },
    {
      "from": "eat b's",
      "to": "eat b's",
      "input": "b",
      "push": [],
      "pop": ["a"],
    },
  ],
  "start": "empty",
};
