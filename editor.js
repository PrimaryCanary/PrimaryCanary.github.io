const searchInput = document.getElementById("search");
const resultsDiv = document.getElementById("results");
const editorDiv = document.getElementById("editor");
const currentNameEl = document.getElementById("current-name");

const fString = document.getElementById("field-string");
const fFlagged = document.getElementById("field-flagged");
const fLink = document.getElementById("field-link");
const fNotes = document.getElementById("field-notes");

let currentIndex = null;

function renderResults(query) {
    resultsDiv.innerHTML = "";
    if (!query) return;

    data
        .map((obj, i) => ({ ...obj, i }))
        .filter((obj) => obj.string.toLowerCase().includes(query.toLowerCase()))
        .slice(0, 5)
        .forEach((obj) => {
            const div = document.createElement("div");
            div.textContent = obj.string;
            if (obj.i === currentIndex) div.classList.add("selected");
            div.onclick = () => loadObject(obj.i);
            div.addEventListener(
                "click",
                (event) => selectResult(obj.i),
            );
            resultsDiv.appendChild(div);
        });
}

function selectResult(result) {
    resultsDiv.innerHTML = "";
    // TODO
    // searchInput = result.innerHTML;
}

function loadObject(index) {
    saveObject();
    currentIndex = index;
    const obj = data[index];
    editorDiv.style.display = "block";

    currentNameEl.textContent = obj.string || "(unnamed)";
    fString.value = obj.string;
    fFlagged.checked = obj.flagged;
    fLink.value = obj.link;
    fNotes.innerHTML = obj.notes;
}

function saveObject() {
    if (currentIndex === null) return;

    data[currentIndex] = {
        string: fString.value,
        flagged: fFlagged.checked,
        link: fLink.value,
        notes: fNotes.innerHTML,
    };
}

function addNew() {
    saveObject();
    data.push({ string: "", flagged: false, link: "", notes: "" });
    searchInput.value = "";
    resultsDiv.innerHTML = "";
    loadObject(data.length - 1);
    fString.focus();
}

function downloadData() {
    saveObject();
    const lines = data.map((obj) => JSON.stringify(obj, null, 2));
    const dataStr = "const data = [\n" + lines.join(",\n") + "\n];";
    const blob = new Blob([dataStr], { type: "text/javascript" });
    const link = document.getElementById("downloadLink");
    link.href = URL.createObjectURL(blob);
    link.click();
}

// WYSIWYG helpers

function toggleWrap(tagName) {
    const sel = window.getSelection();
    if (!sel.rangeCount) return;
    const range = sel.getRangeAt(0);

    // Walk up from the selection anchor to see if we're already inside this tag
    let node = range.commonAncestorContainer;
    if (node.nodeType === Node.TEXT_NODE) node = node.parentNode;
    let existing = null;
    let cursor = node;
    while (cursor && cursor !== fNotes) {
        if (cursor.tagName && cursor.tagName.toLowerCase() === tagName) {
            existing = cursor;
            break;
        }
        cursor = cursor.parentNode;
    }

    if (existing) {
        // Unwrap: hoist children out, remove the wrapper
        const parent = existing.parentNode;
        while (existing.firstChild) {
            parent.insertBefore(existing.firstChild, existing);
        }
        parent.removeChild(existing);
    } else if (!range.collapsed) {
        // Wrap selected content in the tag
        const wrapper = document.createElement(tagName);
        wrapper.appendChild(range.extractContents());
        range.insertNode(wrapper);
        const newRange = document.createRange();
        newRange.selectNodeContents(wrapper);
        sel.removeAllRanges();
        sel.addRange(newRange);
    }
}

function insertList(type) {
    const sel = window.getSelection();
    if (!sel.rangeCount) return;
    const range = sel.getRangeAt(0);

    const list = document.createElement(type);
    const li = document.createElement("li");
    li.appendChild(range.extractContents());
    list.appendChild(li);
    range.insertNode(list);
}

function insertLink(url) {
    const sel = window.getSelection();
    if (!sel.rangeCount) return;
    const range = sel.getRangeAt(0);

    const a = document.createElement("a");
    a.href = url;
    a.target = "_blank";
    a.appendChild(range.extractContents());
    range.insertNode(a);
}

function insertImage(url) {
    const sel = window.getSelection();
    if (!sel.rangeCount) return;
    const range = sel.getRangeAt(0);

    const img = document.createElement("img");
    img.src = url;
    img.style.maxWidth = "100%";
    range.insertNode(img);
}

const toolbar = document.getElementById("toolbar");

// Prevent toolbar buttons from stealing focus from the notes field
toolbar.addEventListener("mousedown", (e) => {
    if (e.target.closest("button")) e.preventDefault();
});

toolbar.addEventListener("click", (e) => {
    const btn = e.target.closest("button");
    if (!btn || !btn.dataset.cmd) return;
    const cmd = btn.dataset.cmd;
    if (cmd === "bold") toggleWrap("b");
    else if (cmd === "italic") toggleWrap("i");
    else if (cmd === "underline") toggleWrap("u");
    else if (cmd === "insertUnorderedList") insertList("ul");
    else if (cmd === "insertOrderedList") insertList("ol");
    fNotes.focus();
});

document.getElementById("linkBtn").addEventListener(
    "mousedown",
    (e) => e.preventDefault(),
);
document.getElementById("linkBtn").addEventListener("click", () => {
    const url = prompt("Enter URL");
    if (url) insertLink(url);
    fNotes.focus();
});

document.getElementById("imgBtn").addEventListener(
    "mousedown",
    (e) => e.preventDefault(),
);
document.getElementById("imgBtn").addEventListener("click", () => {
    const url = prompt("Enter image URL");
    if (url) insertImage(url);
    fNotes.focus();
});

// Keyboard shortcuts in the notes field
fNotes.addEventListener("keydown", (e) => {
    if (!e.ctrlKey && !e.metaKey) return;
    const map = {
        b: () => toggleWrap("b"),
        i: () => toggleWrap("i"),
        u: () => toggleWrap("u"),
    };
    const fn = map[e.key.toLowerCase()];
    if (fn) {
        e.preventDefault();
        fn();
    }
});

searchInput.addEventListener("input", (e) => renderResults(e.target.value));
document.getElementById("downloadBtn").onclick = downloadData;
document.getElementById("addBtn").onclick = addNew;
