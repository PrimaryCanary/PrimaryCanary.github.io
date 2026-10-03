import { Chess } from "https://cdn.jsdelivr.net/npm/chess.js@1.4.0/+esm";

const API = "https://explorer.lichess.ovh/lichess";
const OAUTH_AUTHORIZE = "https://lichess.org/oauth";
const OAUTH_TOKEN = "https://lichess.org/api/token";
const OAUTH_REVOKE = "https://lichess.org/api/token";
const CLIENT_ID = "repertoire-popularity";
const TOKEN_KEY = "repertoire_popularity_lichess_token";
const VERIFIER_KEY = "repertoire_popularity_pkce_verifier";
const STATE_KEY = "repertoire_popularity_oauth_state";

const authStatus = document.querySelector("#authStatus");
const loginButton = document.querySelector("#login");
const logoutButton = document.querySelector("#logout");

let accessToken = localStorage.getItem(TOKEN_KEY);

function base64url(bytes) {
  let s = "";
  for (const b of bytes) s += String.fromCharCode(b);
  return btoa(s).replace(/\+/g, "-").replace(/\//g, "_").replace(/=+$/, "");
}

function randomString(length = 64) {
  const bytes = new Uint8Array(length);
  crypto.getRandomValues(bytes);
  return base64url(bytes);
}

async function pkceChallenge(verifier) {
  const digest = await crypto.subtle.digest(
    "SHA-256",
    new TextEncoder().encode(verifier),
  );
  return base64url(new Uint8Array(digest));
}

function redirectUri() {
  if (location.protocol === "file:") {
    throw new Error(
      "OAuth login requires the app to be served over HTTP. Run: python -m http.server 8000",
    );
  }
  return location.origin + location.pathname;
}

async function login() {
  const verifier = randomString();
  const challenge = await pkceChallenge(verifier);
  const state = randomString(32);
  sessionStorage.setItem(VERIFIER_KEY, verifier);
  sessionStorage.setItem(STATE_KEY, state);

  const params = new URLSearchParams({
    response_type: "code",
    client_id: CLIENT_ID,
    redirect_uri: redirectUri(),
    code_challenge_method: "S256",
    code_challenge: challenge,
    state,
  });
  location.assign(OAUTH_AUTHORIZE + "?" + params);
}

async function finishOAuth() {
  const params = new URLSearchParams(location.search);
  const error = params.get("error");
  const code = params.get("code");
  if (!code && !error) return;
  if (error) throw new Error(`Lichess OAuth error: ${error}`);

  const state = params.get("state");
  const savedState = sessionStorage.getItem(STATE_KEY);
  const verifier = sessionStorage.getItem(VERIFIER_KEY);
  if (!state || state !== savedState) throw new Error("OAuth state mismatch.");
  if (!verifier) throw new Error("Missing PKCE verifier.");

  const body = new URLSearchParams({
    grant_type: "authorization_code",
    code,
    redirect_uri: redirectUri(),
    client_id: CLIENT_ID,
    code_verifier: verifier,
  });
  const response = await fetch(OAUTH_TOKEN, {
    method: "POST",
    headers: { "Content-Type": "application/x-www-form-urlencoded" },
    body,
  });
  if (!response.ok) {
    throw new Error(`Lichess token request failed: HTTP ${response.status}`);
  }
  const token = await response.json();
  if (!token.access_token) {
    throw new Error("Lichess did not return an access token.");
  }

  accessToken = token.access_token;
  localStorage.setItem(TOKEN_KEY, accessToken);
  sessionStorage.removeItem(VERIFIER_KEY);
  sessionStorage.removeItem(STATE_KEY);
  history.replaceState({}, document.title, redirectUri());
}

async function updateAuthStatus() {
  if (!accessToken) {
    authStatus.textContent = "Not logged in";
    loginButton.hidden = false;
    logoutButton.hidden = true;
    return;
  }
  const response = await fetch("https://lichess.org/api/account", {
    headers: { Authorization: `Bearer ${accessToken}` },
  });
  if (!response.ok) {
    accessToken = null;
    localStorage.removeItem(TOKEN_KEY);
    authStatus.textContent = "Not logged in";
    loginButton.hidden = false;
    logoutButton.hidden = true;
    return;
  }
  const account = await response.json();
  authStatus.textContent = `Logged in as ${account.username}`;
  loginButton.hidden = true;
  logoutButton.hidden = false;
}

async function logout() {
  if (accessToken) {
    try {
      await fetch(OAUTH_REVOKE, {
        method: "DELETE",
        headers: { Authorization: `Bearer ${accessToken}` },
      });
    } catch {}
  }
  accessToken = null;
  localStorage.removeItem(TOKEN_KEY);
  await updateAuthStatus();
}

const fileInput = document.querySelector("#file");
const runButton = document.querySelector("#run");
const statusEl = document.querySelector("#status");
const resultsEl = document.querySelector("#results");
const bar = document.querySelector("#bar");

let pgnText = "";

loginButton.addEventListener("click", () =>
  login().catch((e) => {
    statusEl.innerHTML = `<span class="error">${escapeHtml(e.message)}</span>`;
  }));
logoutButton.addEventListener("click", () =>
  logout().catch((e) => {
    statusEl.innerHTML = `<span class="error">${escapeHtml(e.message)}</span>`;
  }));

(async () => {
  try {
    await finishOAuth();
    await updateAuthStatus();
  } catch (e) {
    statusEl.innerHTML = `<span class="error">${escapeHtml(e.message)}</span>`;
    accessToken = null;
    localStorage.removeItem(TOKEN_KEY);
    await updateAuthStatus();
  }
})();

fileInput.addEventListener("change", async () => {
  const f = fileInput.files[0];
  if (!f) return;
  pgnText = await f.text();
  document.querySelector("#filename").textContent =
    `${f.name} (${pgnText.length.toLocaleString()} characters)`;
  resultsEl.innerHTML = "";
  statusEl.textContent = "PGN loaded.";
});

const sleep = (ms) => new Promise((r) => setTimeout(r, ms));

function cloneChess(c) {
  return new Chess(c.fen());
}

function tokenizeMovetext(s) {
  s = s.replace(/\{[^}]*\}/gs, " ")
    .replace(/;[^\n]*/g, " ")
    .replace(/\$\d+/g, " ");
  return s.match(/\(|\)|1-0|0-1|1\/2-1\/2|\*|(?:\d+\.(?:\.\.)?)|[^\s()]+/g) ||
    [];
}

function parseGameMovetext(movetext) {
  let board = new Chess();
  let path = [];
  let sans = [];
  let lastBefore = null;
  const branches = new Map();
  const leaves = new Map();
  const stack = [];

  const addChild = (parentPath, uci) => {
    const k = parentPath.join(" ");
    if (!branches.has(k)) branches.set(k, new Set());
    branches.get(k).add(uci);
  };

  const addLeaf = () => {
    const key = path.join(" ");
    if (key) leaves.set(key, { uci: [...path], san: [...sans] });
  };

  for (const tok of tokenizeMovetext(movetext)) {
    if (tok === "(") {
      if (!lastBefore) continue;

      // A PGN variation starts from the position immediately before the
      // preceding move. Save both the position we will resume at and the
      // branch point. The latter is needed so that sibling variations like
      // (A...) (B...) both start from the same position.
      stack.push({
        resumeBoard: cloneChess(board),
        resumePath: [...path],
        resumeSans: [...sans],
        resumeLastBefore: lastBefore,
      });
      board = cloneChess(lastBefore.board);
      path = [...lastBefore.path];
      sans = [...lastBefore.sans];
      lastBefore = null;
      continue;
    }

    if (tok === ")") {
      // Everything parsed since the matching '(' is one terminal variation
      // unless another nested variation has already continued this path.
      addLeaf();
      const x = stack.pop();
      if (x) {
        board = x.resumeBoard;
        path = x.resumePath;
        sans = x.resumeSans;
        lastBefore = x.resumeLastBefore;
      }
      continue;
    }

    if (
      /^(?:\d+\.(?:\.\.)?)$/.test(tok) ||
      /^(?:1-0|0-1|1\/2-1\/2|\*)$/.test(tok)
    ) continue;

    try {
      const before = cloneChess(board);
      const beforePath = [...path];
      const beforeSans = [...sans];
      const move = board.move(tok, { strict: false });
      if (!move) continue;

      const uci = move.from + move.to + (move.promotion || "");
      addChild(beforePath, uci);
      path.push(uci);
      sans.push(move.san);
      lastBefore = { board: before, path: beforePath, sans: beforeSans };
    } catch (_) {
      // Ignore non-move tokens.
    }
  }

  // The main line (or a variation that reaches EOF) is also a leaf.
  addLeaf();
  return { branches, leaves };
}

function splitGames(pgn) {
  const starts = [];
  const re = /^\[Event\s+/gm;
  let m;
  while ((m = re.exec(pgn))) starts.push(m.index);
  if (!starts.length) return [pgn];
  return starts.map((start, i) =>
    pgn.slice(start, starts[i + 1] ?? pgn.length)
  );
}

function getMovetext(game) {
  return game.replace(/^\s*\[[^\n]*\]\s*$/gm, " ").trim();
}

function parseRepertoire(pgn, depth) {
  const allLeaves = new Map();

  for (const game of splitGames(pgn)) {
    const { leaves } = parseGameMovetext(getMovetext(game));
    for (const [k, line] of leaves) allLeaves.set(k, line);
  }

  // Keep one repertoire line for each distinct prefix at the requested depth.
  const prefixes = new Map();
  for (const line of allLeaves.values()) {
    const d = (line.uci.length < depth) ? line.uci.length : depth;
    const uci = line.uci.slice(0, d);
    const key = uci.join(" ");
    if (!prefixes.has(key)) prefixes.set(key, uci);
  }

  const lines = [];
  for (const uci of prefixes.values()) {
    const board = new Chess();
    const san = [];

    for (const moveUci of uci) {
      const move = board.move({
        from: moveUci.slice(0, 2),
        to: moveUci.slice(2, 4),
        ...(moveUci.length === 5 ? { promotion: moveUci[4] } : {}),
      });
      san.push(move.san);
    }

    lines.push({
      uci,
      san,
      scoreFen: board.fen(),
    });
  }

  return lines;
}

function fmt(n) {
  return n == null ? "—" : n.toLocaleString();
}

function pct(a, b) {
  return b ? (100 * a / b).toFixed(1) + "%" : "—";
}

async function explorer(fen, ratings, speeds) {
  const params = new URLSearchParams({
    variant: "standard",
    ratings: ratings.join(","),
    speeds: speeds.join(","),
    fen,
  });

  for (let attempt = 0; attempt < 3; attempt++) {
    const response = await fetch(API + "?" + params.toString(), {
      headers: accessToken ? { Authorization: `Bearer ${accessToken}` } : {},
    });

    if (response.status === 429) {
      statusEl.textContent = "Lichess rate limit hit; waiting 60 seconds…";
      await sleep(60000);
      continue;
    }

    if (!response.ok) {
      throw new Error(`Lichess returned HTTP ${response.status}`);
    }

    return response.json();
  }

  throw new Error("Lichess rate limit persisted after retries.");
}

function escapeHtml(s) {
  return s.replace(/[&<>"']/g, (c) => ({
    "&": "&amp;",
    "<": "&lt;",
    ">": "&gt;",
    '"': "&quot;",
    "'": "&#39;",
  }[c]));
}

function render(rows, depth) {
  rows.sort((a, b) => b.games - a.games || b.uci.length - a.uci.length);

  let limit = Number(document.querySelector("#limit").value);
  if (limit <= 0 || limit > rows.length) {
    limit = rows.length;
  }
  const shown = rows.slice(0, limit);

  let html = `<table>
    <thead><tr>
      <th>#</th>
      <th>Repertoire line</th>
      <th class="num">Games after ${depth} ply</th>
      <th class="num">% of most common</th>
    </tr></thead><tbody>`;

  const maxGames = Math.max(...shown.map((elem) => elem.games));
  shown.forEach((r, i) => {
    html += `<tr>
      <td>${i + 1}</td>
      <td><code>${escapeHtml(r.san.join(" "))}</code></td>
      <td class="num"><b>${fmt(r.games)}</b></td>
      <td class="num">${fmt(r.games / maxGames * 100)}%</td>
    </tr>`;
  });

  html += `</tbody></table>
    <p class="muted">Showing ${shown.length} of ${rows.length} repertoire lines.</p>`;

  resultsEl.innerHTML = html;
}

runButton.addEventListener("click", async () => {
  if (!accessToken) {
    statusEl.innerHTML =
      "<span class='error'>Log in with Lichess before analyzing.</span>";
    return;
  }
  if (!pgnText) {
    statusEl.innerHTML = "<span class='error'>Choose a PGN first.</span>";
    return;
  }

  const ratings = [...document.querySelectorAll(".rating:checked")].map((x) =>
    x.value
  );
  const speeds = [...document.querySelectorAll(".speed:checked")].map((x) =>
    x.value
  );
  const depth = Math.max(
    1,
    Math.min(200, Number(document.querySelector("#depth").value) || 10),
  );

  if (!ratings.length || !speeds.length) {
    statusEl.innerHTML =
      "<span class='error'>Select at least one rating bracket and one time control.</span>";
    return;
  }

  runButton.disabled = true;
  resultsEl.innerHTML = "";
  bar.style.width = "0%";

  try {
    statusEl.textContent = "Parsing repertoire…";
    const lines = parseRepertoire(pgnText, depth);

    if (!lines.length) {
      throw new Error(
        `No repertoire lines reach move ${depth}. Try a smaller depth.`,
      );
    }

    const uniqueFens = [...new Set(lines.map((x) => x.scoreFen))];
    const data = new Map();

    for (let i = 0; i < uniqueFens.length; i++) {
      statusEl.textContent = `Querying Lichess: ${
        i + 1
      }/${uniqueFens.length} positions…`;

      data.set(
        uniqueFens[i],
        await explorer(uniqueFens[i], ratings, speeds),
      );

      bar.style.width = `${((i + 1) / uniqueFens.length) * 100}%`;
      await sleep(120);
    }

    statusEl.textContent = "Calculating rankings…";
    const rows = lines.map((line) => {
      const d = data.get(line.scoreFen);
      return {
        ...line,
        games: (d.white || 0) + (d.draws || 0) + (d.black || 0),
      };
    });

    render(rows, depth);

    statusEl.textContent =
      `Done. ${lines.length} repertoire lines after move ${depth}; ` +
      `${uniqueFens.length} unique positions queried.`;
  } catch (e) {
    console.error(e);
    statusEl.innerHTML = `<span class="error">${
      escapeHtml(e.message || String(e))
    }</span>`;
  } finally {
    runButton.disabled = false;
  }
});
