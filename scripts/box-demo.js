let box = document.getElementById("color-box");
let button = document.querySelector("button");

function changeColor() {
  const r = randomIntRange(0, 256);
  const g = randomIntRange(0, 256);
  const b = randomIntRange(0, 256);
  box.style.background = rgb(r, g, b);
}

// Non-inclusive range
function randomIntRange(min, max) {
  return Math.floor(Math.random() * (max - min) + min);
}

function rgb(r, g, b) {
  return `rgb(${r}, ${g}, ${b})`;
}
