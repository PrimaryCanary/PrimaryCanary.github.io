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
    subSvg.classList.add("draggable");
    circle.setAttribute("cx", cx + 2);
    circle.setAttribute("cy", cy + 2);
    circle.setAttribute("r", r);
    acceptCircle.setAttribute("cx", cx + 2);
    acceptCircle.setAttribute("cy", cy + 2);
    acceptCircle.setAttribute("r", r - 2);
    svgText.setAttribute("x", "50%");
    svgText.setAttribute("y", "50%");
    svgText.appendChild(name);

    if (state.accept) {
        subSvg.appendChild(acceptCircle);
    }
    subSvg.appendChild(circle);
    // Last so text is selectable
    subSvg.appendChild(svgText);

    makeDraggable(subSvg);
    return subSvg;
}

function makeDraggable(element) {
    element.addEventListener("mousedown", startDrag);
    element.addEventListener("mousemove", drag);
    element.addEventListener("mouseup", endDrag);
    element.addEventListener("mouseleave", endDrag);

    let selected = undefined;
    let offset = undefined;
    function startDrag(event) {
        selected = event.target;
        while (!selected.classList.contains("draggable")) {
            if (selected === svg) {
                //console.log("canvas clicked");
                return;
            }
            selected = selected.parentNode;
        }
        offset = getMousePosition(event);
        offset.x -= selected.x.baseVal.value;
        offset.y -= selected.y.baseVal.value;
        // console.log(selected);
        // console.log(offset);
    }

    function drag(event) {
        if (selected) {
            event.preventDefault();
            const coord = getMousePosition(event);
            // console.log(selected);
            selected.x.baseVal.value = coord.x - offset.x;
            selected.y.baseVal.value = coord.y - offset.y;
        }
    }

    function endDrag(event) {
        selected = undefined;
    }

    function getMousePosition(event) {
        const ctm = selected.getScreenCTM();
        return {
            x: (event.clientX - ctm.e) / ctm.a,
            y: (event.clientY - ctm.f) / ctm.d,
        };
    }
}
