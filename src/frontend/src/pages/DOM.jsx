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
    traversalLog.map((entry, index) => [String(entry?.id ?? ""), index])
  );
  return annotateTree(node, nextOrderRef, traversalOrder);
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

  async function handleStart() {
    setLoading(true);
    setError("");
    const parsedLimit = Number.parseInt(resultLimit, 10);

    try {
      const response = await runTraversal({
        html: inputType === "html" ? html : undefined,
        url: inputType === "url" ? url : undefined,
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
                placeholder="https://example.com"
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
