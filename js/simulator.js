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
  clearChildren(svg);

  for (const t of pda.transitions) {
    const node = newTransitionSvg(t);
    makeDraggable(node);
    svg.appendChild(node);
  }

  for (const s of pda.states) {
    const node = newStateSvg(s, STATE_RADIUS);
    makeDraggable(node);
    svg.appendChild(node);
  }
}

function newTransitionSvg(t) {
  const subSvg = document.createElementNS(SVGNS, "svg");
  const path = document.createElementNS(SVGNS, "path");
  const arrow = document.createElementNS(SVGNS, "polygon");
  const svgText = document.createElementNS(SVGNS, "text");
  const name = document.createTextNode(transitionToString(t));

  subSvg.setAttribute("height", 100);
  subSvg.setAttribute("width", 100);
  subSvg.classList.add("draggable");
  subSvg.classList.add("transition");

  if (t.to === t.from) {
    path.setAttribute("d", "M 40,80 Q 50,-20 60,75");
    arrow.setAttribute("points", "55,75 65,75 60,85");
    svgText.setAttribute("x", "50%");
    svgText.setAttribute("y", "20%");
  } else {
    path.setAttribute("d", "M 0,50 Q 50,25 90,50");
    arrow.setAttribute("points", "90,55 90,45 100,50");
    svgText.setAttribute("x", "50%");
    svgText.setAttribute("y", "50%");
  }

  svgText.appendChild(name);

  subSvg.appendChild(path);
  subSvg.appendChild(arrow);
  // Last so text is selectable
  subSvg.appendChild(svgText);

  return subSvg;
}

function clearChildren(element) {
  while (element.firstChild) {
    element.firstChild.remove();
  }
}

const SVGNS = "http://www.w3.org/2000/svg";
const STATE_RADIUS = 40;
function newStateSvg(s) {
  const subSvg = document.createElementNS(SVGNS, "svg");
  const circle = document.createElementNS(SVGNS, "circle");
  const acceptCircle = document.createElementNS(SVGNS, "circle");
  const svgText = document.createElementNS(SVGNS, "text");
  const name = document.createTextNode(s.name);

  subSvg.setAttribute("height", Math.ceil(2 * (STATE_RADIUS + 2)));
  subSvg.setAttribute("width", Math.ceil(2 * (STATE_RADIUS + 2)));
  subSvg.classList.add("draggable");
  subSvg.id = s.name;
  circle.setAttribute("cx", STATE_RADIUS + 2);
  circle.setAttribute("cy", STATE_RADIUS + 2);
  circle.setAttribute("r", STATE_RADIUS);
  acceptCircle.setAttribute("cx", STATE_RADIUS + 2);
  acceptCircle.setAttribute("cy", STATE_RADIUS + 2);
  acceptCircle.setAttribute("r", STATE_RADIUS - 2);
  svgText.setAttribute("x", "50%");
  svgText.setAttribute("y", "50%");
  svgText.appendChild(name);

  if (s.accept) {
    subSvg.appendChild(acceptCircle);
  }
  subSvg.appendChild(circle);
  // Last so text is selectable
  subSvg.appendChild(svgText);

  return subSvg;
}

function transitionToString(t) {
  const push = t.push.join("");
  const pop = t.pop.join("");
  return `${t.input},${push ? push : "ε"},${pop ? pop : "ε"}`;
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
        return;
      }
      selected = selected.parentNode;
    }
    offset = getMousePosition(event);
    offset.x -= selected.x.baseVal.value;
    offset.y -= selected.y.baseVal.value;
  }

  function drag(event) {
    if (selected) {
      event.preventDefault();
      const coord = getMousePosition(event);
      selected.x.baseVal.value = coord.x - offset.x;
      selected.y.baseVal.value = coord.y - offset.y;
    }
  }

  function endDrag(_event) {
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
