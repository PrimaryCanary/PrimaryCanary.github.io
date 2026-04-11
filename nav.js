const searchInput = document.getElementById("search");
const resultsDiv = document.getElementById("results");
const outputDiv = document.getElementById("output");
const notesDiv = document.getElementById("notes");

let currentResults = [];
let highlightedIndex = 0;

function renderResults(query) {
    resultsDiv.innerHTML = "";
    highlightedIndex = 0;
    if (!query) return;

    currentResults = data.filter((obj) =>
        obj.string.toLowerCase().includes(query.toLowerCase())
    ).slice(0, 5);

    currentResults.forEach((obj, idx) => {
        const div = document.createElement("div");
        div.textContent = obj.string;
        div.addEventListener("click", () => selectResult(idx));
        resultsDiv.appendChild(div);
    });
}

function selectResult(index) {
    const obj = currentResults[index];
    if (!obj) return;

    outputDiv.innerHTML =
        `<a href="${obj.link}" target="_blank">${obj.link}</a>`;

    if (obj.flagged) {
        outputDiv.innerHTML +=
            `<div class="flagged">This key is flagged. Read the notes for special instructions and double check the route.</div>`;
    }

    notesDiv.innerHTML = obj.notes;

    resultsDiv.innerHTML = "";
    searchInput.value = obj.string;
}

function updateHighlight() {
    Array.from(resultsDiv.children).forEach((child, idx) => {
        child.classList.toggle("highlight", idx === highlightedIndex);
    });
}

searchInput.addEventListener("input", (e) => {
    renderResults(e.target.value);
    updateHighlight();
});

searchInput.addEventListener("keydown", (e) => {
    if (e.key === "ArrowDown") {
        highlightedIndex = Math.min(
            highlightedIndex + 1,
            currentResults.length - 1,
        );
        updateHighlight();
        e.preventDefault();
    } else if (e.key === "ArrowUp") {
        highlightedIndex = Math.max(highlightedIndex - 1, 0);
        updateHighlight();
        e.preventDefault();
    } else if (e.key === "Enter") {
        if (highlightedIndex >= 0) {
            selectResult(highlightedIndex);
        }
    }
});
