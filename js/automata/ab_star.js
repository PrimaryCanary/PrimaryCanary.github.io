export default {
    "description": "PDA for {ab}*",
    "states": [
        { "name": "a", "accept": true },
        { "name": "b", "accept": false },
    ],
    "transitions": [
        {
            "from": "a",
            "to": "b",
            "input": "a",
            "push": [],
            "pop": [],
        },
        {
            "from": "b",
            "to": "a",
            "input": "b",
            "push": [],
            "pop": [],
        },
    ],
    "start": "a",
};
