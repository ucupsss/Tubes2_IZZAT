import { useState } from "react";
import TreeNode from "../components/TreeNode";
import { runTraversal } from "../services/api";

const fallbackHtml =
  "<html><body><div><p>Hello</p><div class='box'>World</div></div></body></html>";

function annotateTree(node, nextOrderRef, traversalOrder) {
  if (!node || typeof node !== "object") {
    const order = nextOrderRef.current;
    nextOrderRef.current += 1;
    return { id: "root", value: "root", children: [], animationOrder: order };
  }

  const id = String(node.id ?? node.value ?? node.tag ?? nextOrderRef.current);
  const value = node.value ?? node.tag ?? node.id ?? "root";
  const order = traversalOrder.get(id) ?? nextOrderRef.current;
  nextOrderRef.current += 1;
  const children = Array.isArray(node.children)
    ? node.children.map((child) => annotateTree(child, nextOrderRef, traversalOrder))
    : [];

  return {
    id,
    value,
    tag: node.tag,
    depth: node.depth,
    attributes: node.attributes ?? {},
    text: node.text ?? "",
    texts: Array.isArray(node.texts) ? node.texts : [],
    children,
    animationOrder: order,
  };
}

function normalizeTree(node, traversalLog = []) {
  const nextOrderRef = { current: 0 };
  const traversalOrder = new Map(
    traversalLog.map((entry, index) => [String(entry?.id ?? ""), index]),
  );
  return annotateTree(node, nextOrderRef, traversalOrder);
}

function flattenNodeOptions(node) {
  if (!node) {
    return [];
  }

  const options = [];
  const stack = [node];

  while (stack.length > 0) {
    const current = stack.pop();
    if (!current) {
      continue;
    }

    options.push({
      id: String(current.id ?? ""),
      label: current.value ?? current.tag ?? current.id ?? "unknown",
      depth: Number.isFinite(current.depth) ? current.depth : 0,
    });

    const children = Array.isArray(current.children) ? [...current.children].reverse() : [];
    stack.push(...children);
  }

  return options;
}

function normalizeURLInput(value) {
  const trimmedValue = value.trim();
  if (!trimmedValue) {
    return "";
  }

  if (!trimmedValue.includes("://")) {
    return `http://${trimmedValue}`;
  }

  return trimmedValue;
}

function validateTraversalInput({
  inputType,
  url,
  html,
  selector,
  resultMode,
  resultLimit,
  lcaMode,
  firstNodeId,
  secondNodeId,
}) {
  if (!selector.trim()) {
    return "CSS selector wajib diisi.";
  }

  if (inputType === "url") {
    const trimmedURL = url.trim();
    if (!trimmedURL) {
      return "URL wajib diisi.";
    }

    try {
      const parsedURL = new URL(normalizeURLInput(trimmedURL));
      if (!["http:", "https:"].includes(parsedURL.protocol)) {
        return "URL harus menggunakan http:// atau https://.";
      }
    } catch {
      return "Format URL tidak valid.";
    }
  }

  if (inputType === "html") {
    const trimmedHTML = html.trim();
    if (!trimmedHTML) {
      return "HTML wajib diisi.";
    }

    if (!trimmedHTML.includes("<") || !trimmedHTML.includes(">")) {
      return "Input HTML tidak tampak valid.";
    }
  }

  if (resultMode === "top") {
    const parsedLimit = Number.parseInt(resultLimit, 10);
    if (!Number.isFinite(parsedLimit) || parsedLimit < 1) {
      return "Jumlah hasil Top n harus berupa angka minimal 1.";
    }
  }

  if (lcaMode && firstNodeId && secondNodeId && firstNodeId === secondNodeId) {
    return "Node pertama dan node kedua untuk LCA harus berbeda.";
  }

  return "";
}

export default function DOM() {
  const [inputType, setInputType] = useState("html");
  const [url, setUrl] = useState("");
  const [html, setHtml] = useState(fallbackHtml);
  const [algorithm, setAlgorithm] = useState("bfs");
  const [selector, setSelector] = useState("div");
  const [resultMode, setResultMode] = useState("all");
  const [resultLimit, setResultLimit] = useState(10);
  const [animateNodes, setAnimateNodes] = useState(false);
  const [lcaMode, setLcaMode] = useState(false);
  const [firstNodeId, setFirstNodeId] = useState("");
  const [secondNodeId, setSecondNodeId] = useState("");
  const [animationKey, setAnimationKey] = useState(0);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [result, setResult] = useState(null);

  const nodeOptions = flattenNodeOptions(result?.tree);

  function resetLCASelection() {
    setFirstNodeId("");
    setSecondNodeId("");
  }

  function resetDOMSourceState() {
    resetLCASelection();
    setResult(null);
    setError("");
  }

  function handleInputTypeChange(nextInputType) {
    setInputType(nextInputType);
    resetDOMSourceState();
  }

  function handleURLChange(event) {
    setUrl(event.target.value);
    resetDOMSourceState();
  }

  function handleHTMLChange(event) {
    setHtml(event.target.value);
    resetDOMSourceState();
  }

  async function handleStart() {
    setError("");
    const inputError = validateTraversalInput({
      inputType,
      url,
      html,
      selector,
      resultMode,
      resultLimit,
      lcaMode,
      firstNodeId,
      secondNodeId,
    });
    if (inputError) {
      setResult(null);
      setError(inputError);
      return;
    }

    setLoading(true);
    const parsedLimit = Number.parseInt(resultLimit, 10);
    const availableNodeIDs = new Set(nodeOptions.map((option) => option.id));
    const nextFirstNodeId = availableNodeIDs.has(firstNodeId) ? firstNodeId : "";
    const nextSecondNodeId = availableNodeIDs.has(secondNodeId) ? secondNodeId : "";

    try {
      const response = await runTraversal({
        html: inputType === "html" ? html : undefined,
        url: inputType === "url" ? normalizeURLInput(url) : undefined,
        selector,
        algorithm,
        limit: resultMode === "all" ? 0 : Math.max(1, Number.isFinite(parsedLimit) ? parsedLimit : 1),
        lca: lcaMode,
        firstNodeId: lcaMode ? nextFirstNodeId : "",
        secondNodeId: lcaMode ? nextSecondNodeId : "",
      });

      if (response?.error) {
        setError(response.error);
      }

      const nextResult = {
        tree: normalizeTree(response?.tree, response?.traversalLog),
        visited: Array.isArray(response?.visited) ? response.visited.map(String) : [],
        matched: Array.isArray(response?.matched) ? response.matched.map(String) : [],
        traversalLog: Array.isArray(response?.traversalLog) ? response.traversalLog : [],
        lca: response?.lca ?? null,
        time: Number.isFinite(response?.time) ? response.time : 0,
        maxDepth: Number.isFinite(response?.maxDepth) ? response.maxDepth : 0,
        nodeCount: Number.isFinite(response?.nodeCount) ? response.nodeCount : 0,
        visitedCount: Number.isFinite(response?.visitedCount)
          ? response.visitedCount
          : response?.visited?.length ?? 0,
        matchedCount: Number.isFinite(response?.matchedCount)
          ? response.matchedCount
          : response?.matched?.length ?? 0,
      };

      setResult(nextResult);

      if (response?.lca?.enabled) {
        setFirstNodeId(String(response.lca.firstNodeId ?? ""));
        setSecondNodeId(String(response.lca.secondNodeId ?? ""));
      }

      setAnimationKey((current) => current + 1);
    } catch (err) {
      setResult(null);
      setError(err instanceof Error ? err.message : "Request failed");
    } finally {
      setLoading(false);
    }
  }

  function handleToggleAnimation(event) {
    setAnimateNodes(event.target.checked);
    setAnimationKey((current) => current + 1);
  }

  function handleToggleLCA(event) {
    const checked = event.target.checked;
    setLcaMode(checked);

    if (!checked) {
      resetLCASelection();
      return;
    }

    if (nodeOptions.length >= 2) {
      setFirstNodeId((current) => current || nodeOptions[0].id);
      setSecondNodeId((current) => current || nodeOptions[1].id);
    }
  }

  const visitedSet = new Set(result?.visited ?? []);
  const matchedSet = new Set(result?.matched ?? []);

  return (
    <>
      <h1 className="page-title">Penelusuran DOM</h1>
      <div className="DOM-grid">
        <div className="card">
          <h2>Input</h2>

          <div className="form-group">
            <label>Input Type</label>
            <div className="radio-row">
              <label>
                <input
                  type="radio"
                  checked={inputType === "url"}
                  onChange={() => handleInputTypeChange("url")}
                />
                URL
              </label>
              <label>
                <input
                  type="radio"
                  checked={inputType === "html"}
                  onChange={() => handleInputTypeChange("html")}
                />
                HTML
              </label>
            </div>
          </div>

          {inputType === "url" ? (
            <div className="form-group">
              <label>URL</label>
              <input
                type="url"
                value={url}
                onChange={handleURLChange}
                placeholder="example.com atau https://example.com"
              />
            </div>
          ) : (
            <div className="form-group">
              <label>HTML</label>
              <textarea value={html} onChange={handleHTMLChange} />
            </div>
          )}

          <div className="form-group">
            <label>Algorithm</label>
            <select value={algorithm} onChange={(event) => setAlgorithm(event.target.value)}>
              <option value="bfs">BFS</option>
              <option value="dfs">DFS</option>
            </select>
          </div>

          <div className="form-group">
            <label>CSS Selector</label>
            <input
              type="text"
              value={selector}
              onChange={(event) => setSelector(event.target.value)}
              placeholder="div, .class, #id"
            />
          </div>

          <div className="form-group">
            <label>Jumlah Hasil</label>
            <div className="radio-row">
              <label>
                <input
                  type="radio"
                  checked={resultMode === "all"}
                  onChange={() => setResultMode("all")}
                />
                Semua
              </label>
              <label>
                <input
                  type="radio"
                  checked={resultMode === "top"}
                  onChange={() => setResultMode("top")}
                />
                Top n
              </label>
            </div>
            {resultMode === "top" && (
              <input
                className="number-input"
                type="number"
                min="1"
                value={resultLimit}
                onChange={(event) => setResultLimit(event.target.value)}
              />
            )}
          </div>

          <div className="form-group">
            <label className="checkbox-row">
              <span>Aktifkan Mode Animasi</span>
              <input type="checkbox" checked={animateNodes} onChange={handleToggleAnimation} />
            </label>
          </div>

          <div className="form-group">
            <label className="checkbox-row">
              <span>Aktifkan Mode LCA</span>
              <input type="checkbox" checked={lcaMode} onChange={handleToggleLCA} />
            </label>
          </div>

          {lcaMode && (
            <div className="lca-config">
              <p className="lca-description">
                Lowest Common Ancestor mencari parent bersama terdekat dari dua node pada pohon
                DOM dengan teknik Binary Lifting.
              </p>

              <div className="form-group">
                <label>Node Pertama</label>
                <select
                  value={firstNodeId}
                  onChange={(event) => setFirstNodeId(event.target.value)}
                  disabled={nodeOptions.length < 2}
                >
                  <option value="">
                    {nodeOptions.length < 2 ? "Jalankan traversal dulu" : "Pilih node pertama"}
                  </option>
                  {nodeOptions.map((option) => (
                    <option key={`first-${option.id}`} value={option.id}>
                      {`${" ".repeat(option.depth * 2)}${option.label} [${option.id}]`}
                    </option>
                  ))}
                </select>
              </div>

              <div className="form-group">
                <label>Node Kedua</label>
                <select
                  value={secondNodeId}
                  onChange={(event) => setSecondNodeId(event.target.value)}
                  disabled={nodeOptions.length < 2}
                >
                  <option value="">
                    {nodeOptions.length < 2 ? "Jalankan traversal dulu" : "Pilih node kedua"}
                  </option>
                  {nodeOptions.map((option) => (
                    <option key={`second-${option.id}`} value={option.id}>
                      {`${" ".repeat(option.depth * 2)}${option.label} [${option.id}]`}
                    </option>
                  ))}
                </select>
              </div>
            </div>
          )}

          <button className="primary" onClick={handleStart} disabled={loading}>
            {loading ? "Running..." : "Start"}
          </button>
        </div>

        <div className="card">
          <h2>Result</h2>

          {error && <div className="error">{error}</div>}
          {!result && !error && <p style={{ color: "#666" }}>Run a traversal to see results.</p>}

          {result && (
            <>
              <div className="stats">
                <span>
                  Visited: <b>{result.visitedCount}</b>
                </span>
                <span>
                  Matched: <b>{result.matchedCount}</b>
                </span>
                <span>
                  Max Depth: <b>{result.maxDepth}</b>
                </span>
                <span>
                  Nodes: <b>{result.nodeCount}</b>
                </span>
                <span>
                  Time: <b>{result.time.toFixed(3)} ms</b>
                </span>
              </div>

              {result.lca?.enabled && (
                <div className="lca-result">
                  <h3 className="section-subtitle">LCA Binary Lifting</h3>
                  <p>
                    Node 1: <b>{result.lca.firstNodeLabel}</b> [{result.lca.firstNodeId}]
                  </p>
                  <p>
                    Node 2: <b>{result.lca.secondNodeLabel}</b> [{result.lca.secondNodeId}]
                  </p>
                  <p>
                    Common Ancestor: <b>{result.lca.ancestorLabel}</b> [{result.lca.ancestorId}]
                  </p>
                  <p>
                    Depth: <b>{result.lca.ancestorDepth}</b>
                  </p>
                </div>
              )}

              <h3 className="section-subtitle">Tree</h3>
              <div className="tree">
                <TreeNode
                  key={`${animationKey}-${animateNodes ? "animated" : "static"}`}
                  node={result.tree}
                  visited={visitedSet}
                  matched={matchedSet}
                  animate={animateNodes}
                />
              </div>

              <h3 className="section-subtitle">Traversal Order</h3>
              <div className="order-list">
                {result.traversalLog.map((entry, index) => (
                  <span
                    className={entry.matched ? "log-entry matched" : "log-entry"}
                    key={`${entry.id}-${index}`}
                  >
                    {index + 1}. {entry.tag}
                  </span>
                ))}
              </div>
            </>
          )}
        </div>
      </div>
    </>
  );
}
