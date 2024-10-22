export default
{
    "description": "DFA for {ab}*",
    "states": [
        { "name": "a", "accept": true },
        { "name": "b", "accept": false }
    ],
    "transitions": [
        { "from": "a", "to": "b", "input": "a" },
        { "from": "b", "to": "a", "input": "b" }
    ],
    "start": "a"
}
