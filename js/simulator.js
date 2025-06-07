import { PDA } from "./pda.js";

// Global elements
const load = document.getElementById("load");
const automText = document.getElementById("automata-text");
const canvas = document.getElementById("simulator-svg");
const simulate = document.getElementById("simulate");
const acceptBox = document.getElementById("accept");
const rejectBox = document.getElementById("reject");
const results = document.getElementById("result-box");

// Simulator state
let pda = undefined;
let pdaSvgs = new Map();

// Button handlers
load.addEventListener("click", loadPda);
simulate.addEventListener("click", runSimulation);

// Drag handlers
canvas.addEventListener("mousemove", drag);
canvas.addEventListener("mouseup", endDrag);
canvas.addEventListener("mouseleave", endDrag);
let dragTarget = null;
let dragOffset = { x: 0, y: 0 };

// SVG drawing constants
const SVGNS = "http://www.w3.org/2000/svg";
// Arbitrarily chosen; they happen to look good together
const STATE_RADIUS = 40;
const TRIANGLE = { size: 10, rotate: 10, offset: 3, rotY: STATE_RADIUS / 2 };
const SELF_CURVE = "-22 -60, 28 -60, 11 -7";
const TEXT_HEIGHT = 60;
const EDGE_PUSH = 2; // Compensates for lines drawn on the edge of the bounding box

function runSimulation() {
  // TODO ignore empty accept/reject box
  if (!pda) {
    results.value = "Error: no automata loaded";
  }

  const acc = parseStringPerLine(acceptBox.value).map((s) => {
    return new Result(s, true, pda.simulate(s));
  });
  const rej = parseStringPerLine(rejectBox.value).map((s) => {
    return new Result(s, false, pda.simulate(s));
  });
  const resultStrings = acc.concat(rej).map((res) => {
    return res.toString();
  });
  results.value = resultStrings.join("\n");
}

function loadPda() {
  const json = automText.value;
  results.value = "";
  try {
    pda = PDA.fromJson(json);
  } catch (e) {
    results.value = "Error: failed to create PDA: " + e;
  }
  clearChildren(canvas);
  renderPda(canvas, pda, pdaSvgs);
}

class Result {
  constructor(str, shouldAccept, didAccept) {
    this.str = str;
    this.shouldAccept = shouldAccept;
    this.didAccept = didAccept;
  }

  toString() {
    const verb = this.shouldAccept ? "accepted" : "rejected";
    const str = this.str || '""';
    return `${str} was ${verb}: ${this.didAccept === this.shouldAccept}`;
  }
}

function parseStringPerLine(str) {
  return str.trimEnd().split("\n");
}

function renderPda(canvas, pda, pdaSvgs) {
  pdaSvgs.clear();
  for (const s of pda.states) {
    const node = newStateSvg(s, canvas);
    makeDraggable(node);
    canvas.appendChild(node);
    pdaSvgs.set(s.name, node);
  }

  for (const t of pda.transitions) {
    const fromSvg = pdaSvgs.get(t.from.name);
    const toSvg = pdaSvgs.get(t.to.name);
    const node = newTransitionSvg(t, fromSvg, toSvg);
    canvas.appendChild(node);
    pdaSvgs.set(t.toId(), node);
  }
}

function newTransitionSvg(t, fromSvg, toSvg) {
  const subSvg = document.createElementNS(SVGNS, "svg");
  const path = document.createElementNS(SVGNS, "path");
  const arrow = document.createElementNS(SVGNS, "polygon");
  const svgText = document.createElementNS(SVGNS, "text");
  const name = document.createTextNode(t.toString());

  subSvg.classList.add("transition");

  const fromX = parseFloat(fromSvg.getAttribute("x")) + STATE_RADIUS +
    EDGE_PUSH;
  const fromY = parseFloat(fromSvg.getAttribute("y")) + STATE_RADIUS +
    EDGE_PUSH;
  const toX = parseFloat(toSvg.getAttribute("x")) + STATE_RADIUS + EDGE_PUSH;
  const toY = parseFloat(toSvg.getAttribute("y")) + STATE_RADIUS + EDGE_PUSH;

  if (t.to === t.from) {
    const loopY = fromY - STATE_RADIUS; // Alias for top of circle
    path.setAttribute(
      "d",
      `M ${fromX - TRIANGLE.offset},${loopY} c ${SELF_CURVE}`,
    );

    arrow.setAttribute(
      "points",
      `${fromX - TRIANGLE.size / 2 + TRIANGLE.offset},${
        loopY - TRIANGLE.size
      } ${fromX + TRIANGLE.size / 2 + TRIANGLE.offset},${
        loopY - TRIANGLE.size
      } ${fromX + TRIANGLE.offset},${loopY}`,
    );
    arrow.setAttribute(
      "transform",
      `rotate(${TRIANGLE.rotate} ${fromX} ${loopY + TRIANGLE.rotY})`,
    );
    svgText.setAttribute("x", fromX);
    svgText.setAttribute("y", loopY - TEXT_HEIGHT);
  } else {
    // Adjust points to be on the circle's edge
    const angle = Math.atan2(toY - fromY, toX - fromX);
    const startX = fromX + STATE_RADIUS * Math.cos(angle);
    const startY = fromY + STATE_RADIUS * Math.sin(angle);
    const endX = toX - STATE_RADIUS * Math.cos(angle);
    const endY = toY - STATE_RADIUS * Math.sin(angle);

    // Control point for bezier curve
    const midX = (startX + endX) / 2;
    const midY = (startY + endY) / 2;
    const controlPointOffset = 30;
    const controlX = midX + controlPointOffset * Math.sin(angle);
    const controlY = midY - controlPointOffset * Math.cos(angle);
    path.setAttribute(
      "d",
      `M ${startX},${startY} Q ${controlX},${controlY} ${endX},${endY}`,
    );

    // Calculate arrow points
    const arrowAngle = Math.atan2(endY - controlY, endX - controlX);
    const p1x = endX - TRIANGLE.size * Math.cos(arrowAngle - Math.PI / 6);
    const p1y = endY - TRIANGLE.size * Math.sin(arrowAngle - Math.PI / 6);
    const p2x = endX - TRIANGLE.size * Math.cos(arrowAngle + Math.PI / 6);
    const p2y = endY - TRIANGLE.size * Math.sin(arrowAngle + Math.PI / 6);
    arrow.setAttribute("points", `${endX},${endY} ${p1x},${p1y} ${p2x},${p2y}`);

    // TODO transitions between two states lay on top of each other, makes t.toString() unreadable
    svgText.setAttribute("x", midX);
    svgText.setAttribute("y", midY);
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

// TODO x, y arguments
function newStateSvg(s, canvas) {
  const subSvg = document.createElementNS(SVGNS, "svg");
  const circle = document.createElementNS(SVGNS, "circle");
  const acceptCircle = document.createElementNS(SVGNS, "circle");
  const svgText = document.createElementNS(SVGNS, "text");
  const name = document.createTextNode(s.name);

  // Set initial random position for states
  // TODO this is garbage
  const x = Math.random() * (canvas.clientWidth - 2 * STATE_RADIUS);
  const y = Math.random() * (canvas.clientHeight - 2 * STATE_RADIUS);
  subSvg.setAttribute("x", x);
  subSvg.setAttribute("y", y);

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

function makeDraggable(element) {
  element.classList.add("draggable");
  element.addEventListener("mousedown", startDrag);
}

function startDrag(event) {
  // Move entire draggable element, not just the clicked element
  dragTarget = event.target;
  while (!dragTarget.classList.contains("draggable")) {
    dragTarget = dragTarget.parentNode;
  }

  const mousePos = getMousePosition(event);
  const elementX = parseFloat(dragTarget.getAttribute("x")) || 0;
  const elementY = parseFloat(dragTarget.getAttribute("y")) || 0;

  dragOffset.x = mousePos.x - elementX;
  dragOffset.y = mousePos.y - elementY;
}

function drag(event) {
  if (dragTarget) {
    event.preventDefault();
    const coord = getMousePosition(event);
    const newX = coord.x - dragOffset.x;
    const newY = coord.y - dragOffset.y;

    dragTarget.setAttribute("x", newX);
    dragTarget.setAttribute("y", newY);
    updateTransitions(pda, pdaSvgs, dragTarget.id);
  }
}

function endDrag() {
  dragTarget = null;
}

function getMousePosition(event) {
  const ctm = canvas.getScreenCTM();
  return {
    x: (event.clientX - ctm.e) / ctm.a,
    y: (event.clientY - ctm.f) / ctm.d,
  };
}

// Function to update transitions connected to a specific state
function updateTransitions(pda, pdaSvgs, stateName) {
  const transitions = pda.transitions.filter(
    (t) => t.from.name === stateName || t.to.name === stateName,
  );
  for (const t of transitions) {
    const svg = pdaSvgs.get(t.toId());
    svg.remove();
    const fromSvg = pdaSvgs.get(t.from.name);
    const toSvg = pdaSvgs.get(t.to.name);
    const node = newTransitionSvg(t, fromSvg, toSvg);
    canvas.appendChild(node);
    pdaSvgs.set(t.toId(), node);
  }
}
