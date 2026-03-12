const go = new Go();

let wasmInstance = null;

async function initWasm() {
    const result = await WebAssembly.instantiateStreaming(
        fetch("loxogon.wasm"),
        go.importObject,
    );

    wasmInstance = result.instance;

    go.run(wasmInstance);
}

async function runCode() {
    const src = document.getElementById("source").value;
    const output = document.getElementById("output");

    if (!window.runLox) {
        output.textContent = "Interpreter not loaded";
        return;
    }

    try {
        output.textContent = "Running...\n";
        // ????????, wait for DOM update I think
        await sleep(10);
        const result = await window.runLox(src);
        output.textContent = "";
        if (result.stdout) {
            output.textContent = result.stdout;
        }
        if (result.exitCode || result.error) {
            // console.log(result.exitCode);
            output.textContent = result.error;
        } else {
            output.textContent += "last expression: " + result.lastExpr;
        }
    } catch (err) {
        output.textContent = err;
    }
}
async function sleep(ms) {
    return new Promise((resolve) => {
        setTimeout(() => {
            resolve();
        }, ms);
    });
}

function share() {
    const code = encodeURIComponent(
        document.getElementById("source").value,
    );

    const url = `${location.origin}${location.pathname}?code=${code}`;

    navigator.clipboard.writeText(url);
    alert("Link copied to clipboard");
}

function loadSharedCode() {
    const params = new URLSearchParams(location.search);
    const code = params.get("code");

    if (code) {
        document.getElementById("source").value = decodeURIComponent(code);
    }
}

const examples = {
    hello: `print "Hello, world!";`,

    math: `print 1 + 2 * 3;
print (10 - 3) / 5;`,

    fib: `fun fib(n) {
  if (n <= 1) return n;
  return fib(n-1) + fib(n-2);
}

print fib(10);`,

    loop: `var i = 0;
while (i < 5) {
  print i;
  i = i + 1;
}`,
};

function loadExample() {
    const select = document.getElementById("examples");
    const key = select.value;

    if (!key) return;

    document.getElementById("source").value = examples[key];
}

document.getElementById("examples").addEventListener("input", loadExample);
document.getElementById("run").onclick = runCode;
document.getElementById("share").onclick = share;

loadExample();
loadSharedCode();
initWasm();
