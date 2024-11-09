export default {
  "description": "PDA for a^n b^2n i.e. n a's followed by 2n b's",
  "states": [
    { "name": "empty", "accept": true },
    { "name": "eating a's", "accept": false },
    { "name": "eating b's", "accept": true },
  ],
  "transitions": [
    {
      "from": "empty",
      "to": "eating a's",
      "input": "a",
      "push": ["a", "a"],
      "pop": [],
    },
    {
      "from": "eating a's",
      "to": "eating a's",
      "input": "a",
      "push": ["a", "a"],
      "pop": [],
    },
    {
      "from": "eating a's",
      "to": "eating b's",
      "input": "b",
      "push": [],
      "pop": ["a"],
    },
    {
      "from": "eating b's",
      "to": "eating b's",
      "input": "b",
      "push": [],
      "pop": ["a"],
    },
  ],
  "start": "empty",
};
