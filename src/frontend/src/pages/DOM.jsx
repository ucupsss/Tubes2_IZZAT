import { useState } from "react";
import TreeNode from "../components/TreeNode";
import { runLCA, runTraversal } from "../services/api";

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
    traversalLog.map((entry, index) => [String(entry?.id ?? ""), index])
  );
  return annotateTree(node, nextOrderRef, traversalOrder);
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

function validateTraversalInput({ inputType, url, html, selector, resultMode, resultLimit }) {
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
  const [animationKey, setAnimationKey] = useState(0);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [result, setResult] = useState(null);
  const [activeDOMInput, setActiveDOMInput] = useState(null);
  const [lcaNodeIdA, setLcaNodeIdA] = useState("");
  const [lcaNodeIdB, setLcaNodeIdB] = useState("");
  const [lcaLoading, setLcaLoading] = useState(false);
  const [lcaError, setLcaError] = useState("");
  const [lcaResult, setLcaResult] = useState(null);

  async function handleStart() {
    setError("");
    const inputError = validateTraversalInput({
      inputType,
      url,
      html,
      selector,
      resultMode,
      resultLimit,
    });
    if (inputError) {
      setResult(null);
      setError(inputError);
      return;
    }

    setLoading(true);
    const parsedLimit = Number.parseInt(resultLimit, 10);

    try {
      const normalizedURL = inputType === "url" ? normalizeURLInput(url) : "";
      const response = await runTraversal({
        html: inputType === "html" ? html : undefined,
        url: inputType === "url" ? normalizedURL : undefined,
        selector,
        algorithm,
        limit: resultMode === "all" ? 0 : Math.max(1, Number.isFinite(parsedLimit) ? parsedLimit : 1),
      });

      if (response?.error) {
        setError(response.error);
      }

      setResult({
        tree: normalizeTree(response?.tree, response?.traversalLog),
        visited: Array.isArray(response?.visited) ? response.visited.map(String) : [],
        matched: Array.isArray(response?.matched) ? response.matched.map(String) : [],
        traversalLog: Array.isArray(response?.traversalLog) ? response.traversalLog : [],
        time: Number.isFinite(response?.time) ? response.time : 0,
        maxDepth: Number.isFinite(response?.maxDepth) ? response.maxDepth : 0,
        nodeCount: Number.isFinite(response?.nodeCount) ? response.nodeCount : 0,
        visitedCount: Number.isFinite(response?.visitedCount)
          ? response.visitedCount
          : response?.visited?.length ?? 0,
        matchedCount: Number.isFinite(response?.matchedCount)
          ? response.matchedCount
          : response?.matched?.length ?? 0,
      });
      setActiveDOMInput({
        html: inputType === "html" ? html : undefined,
        url: inputType === "url" ? normalizedURL : undefined,
      });
      setLcaNodeIdA("");
      setLcaNodeIdB("");
      setLcaError("");
      setLcaResult(null);
      setAnimationKey((current) => current + 1);
    } catch (err) {
      setResult(null);
      setActiveDOMInput(null);
      setLcaResult(null);
      setError(err instanceof Error ? err.message : "Request failed");
    } finally {
      setLoading(false);
    }
  }

  function handleToggleAnimation(event) {
    setAnimateNodes(event.target.checked);
    setAnimationKey((current) => current + 1);
  }

  async function handleFindLCA() {
    setLcaError("");

    if (!result?.tree || !activeDOMInput) {
      setLcaResult(null);
      setLcaError("Jalankan traversal atau muat DOM terlebih dahulu.");
      return;
    }

    if (!lcaNodeIdA.trim() || !lcaNodeIdB.trim()) {
      setLcaResult(null);
      setLcaError("Node ID A dan Node ID B wajib diisi.");
      return;
    }

    setLcaLoading(true);

    try {
      const response = await runLCA({
        html: activeDOMInput.html,
        url: activeDOMInput.url,
        nodeIdA: lcaNodeIdA.trim(),
        nodeIdB: lcaNodeIdB.trim(),
      });

      setLcaResult(response);
    } catch (err) {
      setLcaResult(null);
      setLcaError(err instanceof Error ? err.message : "Request failed");
    } finally {
      setLcaLoading(false);
    }
  }

  function handleResetLCA() {
    setLcaNodeIdA("");
    setLcaNodeIdB("");
    setLcaError("");
    setLcaResult(null);
  }

  const visitedSet = new Set(result?.visited ?? []);
  const matchedSet = new Set(result?.matched ?? []);
  const selectedASet = new Set(lcaResult?.nodeA?.id ? [String(lcaResult.nodeA.id)] : []);
  const selectedBSet = new Set(lcaResult?.nodeB?.id ? [String(lcaResult.nodeB.id)] : []);
  const selectedLCASet = new Set(lcaResult?.lca?.id ? [String(lcaResult.lca.id)] : []);

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
                  onChange={() => setInputType("url")}
                />
                URL
              </label>
              <label>
                <input
                  type="radio"
                  checked={inputType === "html"}
                  onChange={() => setInputType("html")}
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
                onChange={(event) => setUrl(event.target.value)}
                placeholder="example.com atau https://example.com"
              />
            </div>
          ) : (
            <div className="form-group">
              <label>HTML</label>
              <textarea value={html} onChange={(event) => setHtml(event.target.value)} />
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

              <h3 className="section-subtitle">Tree</h3>
              <div className="tree">
                <TreeNode
                  key={`${animationKey}-${animateNodes ? "animated" : "static"}`}
                  node={result.tree}
                  visited={visitedSet}
                  matched={matchedSet}
                  selectedA={selectedASet}
                  selectedB={selectedBSet}
                  selectedLCA={selectedLCASet}
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

      <div className="card">
        <h2>Pencarian LCA</h2>

        {!result?.tree && (
          <p className="muted-text">Jalankan traversal atau muat DOM terlebih dahulu.</p>
        )}

        {result?.tree && (
          <>
            <div className="lca-grid">
              <div className="form-group">
                <label>Node ID A</label>
                <input
                  type="text"
                  value={lcaNodeIdA}
                  onChange={(event) => setLcaNodeIdA(event.target.value)}
                  placeholder="contoh: 4"
                />
              </div>

              <div className="form-group">
                <label>Node ID B</label>
                <input
                  type="text"
                  value={lcaNodeIdB}
                  onChange={(event) => setLcaNodeIdB(event.target.value)}
                  placeholder="contoh: 7"
                />
              </div>
            </div>

            <div className="lca-actions">
              <button className="primary" onClick={handleFindLCA} disabled={lcaLoading}>
                {lcaLoading ? "Mencari..." : "Cari LCA"}
              </button>
              <button className="secondary-button" onClick={handleResetLCA} disabled={lcaLoading}>
                Reset LCA
              </button>
            </div>

            {lcaError && <div className="error">{lcaError}</div>}

            {lcaResult && (
              <div className="lca-result">
                <div className="lca-result-item">
                  <span className="lca-label">Node A</span>
                  <span>{lcaResult.nodeA.value} [id: {lcaResult.nodeA.id}]</span>
                </div>
                <div className="lca-result-item">
                  <span className="lca-label">Node B</span>
                  <span>{lcaResult.nodeB.value} [id: {lcaResult.nodeB.id}]</span>
                </div>
                <div className="lca-result-item lca-highlight">
                  <span className="lca-label">LCA</span>
                  <span>
                    {lcaResult.lca.value} [id: {lcaResult.lca.id}] depth: {lcaResult.lca.depth}
                  </span>
                </div>
              </div>
            )}
          </>
        )}
      </div>
    </>
  );
}
