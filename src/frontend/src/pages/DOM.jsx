import { useState } from "react";
import TreeNode from "../components/TreeNode";
import { runTraversal } from "../services/api";

const fallbackHtml =
  "<html><body><div><p>Hello</p><div class='box'>World</div></div></body></html>";

function annotateTree(node, nextOrderRef) {
  if (!node || typeof node !== "object") {
    const order = nextOrderRef.current;
    nextOrderRef.current += 1;
    return { value: "root", children: [], animationOrder: order };
  }

  const value = node.value ?? node.tag ?? node.id ?? "root";
  const order = nextOrderRef.current;
  nextOrderRef.current += 1;
  const children = Array.isArray(node.children)
    ? node.children.map((child) => annotateTree(child, nextOrderRef))
    : [];

  return { value, children, animationOrder: order };
}

function normalizeTree(node) {
  const nextOrderRef = { current: 0 };
  return annotateTree(node, nextOrderRef);
}

export default function DOM() {
  const [inputType, setInputType] = useState("html");
  const [url, setUrl] = useState("");
  const [html, setHtml] = useState(fallbackHtml);
  const [algorithm, setAlgorithm] = useState("bfs");
  const [selector, setSelector] = useState("div");
  const [animateNodes, setAnimateNodes] = useState(false);
  const [animationKey, setAnimationKey] = useState(0);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [result, setResult] = useState(null);

  async function handleStart() {
    setLoading(true);
    setError("");

    try {
      const response = await runTraversal({
        html: inputType === "html" ? html : undefined,
        url: inputType === "url" ? url : undefined,
        selector,
        algorithm,
      });

      if (response?.error) {
        setError(response.error);
      }

      setResult({
        tree: normalizeTree(response?.tree),
        visited: Array.isArray(response?.visited) ? response.visited : [],
        matched: Array.isArray(response?.matched) ? response.matched : [],
        time: Number.isFinite(response?.time) ? response.time : 0,
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
                  Visited: <b>{result.visited.length}</b>
                </span>
                <span>
                  Matched: <b>{result.matched.length}</b>
                </span>
                <span>
                  Time: <b>{result.time} ms</b>
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
              <div className="order-list">[{result.visited.join(", ")}]</div>
            </>
          )}
        </div>
      </div>
    </>
  );
}