const searchInput = document.getElementById("search");
const resultsDiv = document.getElementById("results");
const editorDiv = document.getElementById("editor");

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
        .filter((o) => o.string.toLowerCase().includes(query.toLowerCase()))
        .slice(0, 5)
        .forEach((o) => {
            const div = document.createElement("div");
            div.textContent = o.string;
            div.onclick = () => loadObject(o.i);
            resultsDiv.appendChild(div);
        });
}

function loadObject(i) {
    saveObject();
    currentIndex = i;
    const obj = data[i];
    editorDiv.style.display = "block";

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
    const obj = {
        string: "",
        flagged: false,
        link: "",
        notes: "Notes go here.",
    };
    data.push(obj);
    loadObject(data.length - 1);
}

function downloadData() {
    saveObject();
    let dataStr = `const data = ${JSON.stringify(data)};`;
    let blob = new Blob([dataStr]);
    let link = document.getElementById("downloadLink");
    link.href = URL.createObjectURL(blob);
    link.click();
}

searchInput.addEventListener("input", (e) => renderResults(e.target.value));
document.getElementById("downloadBtn").onclick = downloadData;
document.getElementById("addBtn").onclick = addNew;

// WYSIWYG helpers
function wrapSelection(tagName) {
    const sel = window.getSelection();
    if (!sel.rangeCount) return;

    const range = sel.getRangeAt(0);
    if (range.collapsed) return;

    const wrapper = document.createElement(tagName);
    wrapper.appendChild(range.extractContents());
    range.insertNode(wrapper);

    // move cursor after wrapper
    range.setStartAfter(wrapper);
    range.collapse(true);
    sel.removeAllRanges();
    sel.addRange(range);
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

// toolbar wiring
const toolbar = document.getElementById("toolbar");

toolbar.querySelector('[data-cmd="bold"]').onclick = () => wrapSelection("b");
toolbar.querySelector('[data-cmd="italic"]').onclick = () => wrapSelection("i");
toolbar.querySelector('[data-cmd="underline"]').onclick = () =>
    wrapSelection("u");
toolbar.querySelector('[data-cmd="insertUnorderedList"]').onclick = () =>
    insertList("ul");
toolbar.querySelector('[data-cmd="insertOrderedList"]').onclick = () =>
    insertList("ol");

document.getElementById("linkBtn").onclick = () => {
    const url = prompt("Enter URL");
    if (url) insertLink(url);
};

document.getElementById("imgBtn").onclick = () => {
    const url = prompt("Enter image URL");
    if (url) insertImage(url);
};
