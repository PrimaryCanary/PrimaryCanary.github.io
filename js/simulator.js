import { PDA } from "./pda.js";

const load = document.getElementById("load");
const aut_text = document.getElementById("automata-text");
const svg = document.getElementById("simulator-svg");

load.addEventListener("click", loadPda);

function loadPda() {
    const json = aut_text.value;
    const pda = PDA.fromJson(json);
    renderPda(pda);
}

function renderPda(pda) {
    // console.log(pda);
    clearSimulator();

    for (const s of pda.states) {
        const node = newStateSvg(s, 40, 40, 40);
        svg.appendChild(node);
    }
}

function clearSimulator() {
    while (svg.firstChild) {
        svg.firstChild.remove();
    }
}

const svgNS = "http://www.w3.org/2000/svg";
function newStateSvg(state, cx, cy, r) {
    let subSvg = document.createElementNS(svgNS, "svg");
    let circle = document.createElementNS(svgNS, "circle");
    let acceptCircle = document.createElementNS(svgNS, "circle");
    let svgText = document.createElementNS(svgNS, "text");
    let name = document.createTextNode(state.name);

    subSvg.setAttribute("height", Math.ceil(2 * (r + 2)));
    subSvg.setAttribute("width", Math.ceil(2 * (r + 2)));
    circle.setAttribute("cx", cx + 2);
    circle.setAttribute("cy", cy + 2);
    circle.setAttribute("r", r);
    acceptCircle.setAttribute("cx", cx + 2);
    acceptCircle.setAttribute("cy", cy + 2);
    acceptCircle.setAttribute("r", r - 2);
    svgText.setAttribute("x", "50%");
    svgText.setAttribute("y", "50%");
    svgText.appendChild(name);

    subSvg.appendChild(circle);
    subSvg.appendChild(svgText);
    if (state.accept) {
        subSvg.appendChild(acceptCircle);
    }

    return subSvg;
}
