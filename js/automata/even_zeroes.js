export default {
    "description":
        "DFA for (1* 0 1* 0 1*)* i.e. binary strings with even number of zeroes",
    "states": [
        { "name": "even", "accept": true },
        { "name": "odd", "accept": false },
    ],
    "transitions": [
        { "from": "even", "to": "odd", "input": "0" },
        { "from": "even", "to": "even", "input": "1" },
        { "from": "odd", "to": "even", "input": "0" },
        { "from": "odd", "to": "odd", "input": "1" },
    ],
    "start": "even",
};
