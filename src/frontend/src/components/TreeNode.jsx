import { useEffect, useState } from "react";

function formatAttributeValue(value) {
  if (value === "") {
    return '""';
  }

  return `"${value}"`;
}

function formatAttributes(attributes) {
  return attributes
    .map(([name, value]) => `${name}=${formatAttributeValue(value)}`)
    .join(" ");
}

export default function TreeNode({
  node,
  depth = 0,
  visited,
  matched,
  selectedA,
  selectedB,
  selectedLCA,
  animate = false,
}) {
  const [visible, setVisible] = useState(() => !animate);

  useEffect(() => {
    if (!animate) return undefined;

    const delay = (node?.animationOrder ?? 0) * 180;
    const timerId = window.setTimeout(() => {
      setVisible(true);
    }, delay);

    return () => window.clearTimeout(timerId);
  }, [animate, node?.animationOrder]);

  const nodeId = String(node?.id ?? node?.value ?? node?.tag ?? "unknown");
  const nodeValue = node?.value ?? node?.tag ?? node?.id ?? "unknown";
  const children = Array.isArray(node?.children) ? node.children : [];
  const attributes = Object.entries(node?.attributes ?? {});
  const texts = Array.isArray(node?.texts) && node.texts.length > 0
    ? node.texts.filter(Boolean)
    : [node?.text].filter(Boolean);
  const isVisited = visited.has(nodeId);
  const isMatched = matched.has(nodeId);
  const isSelectedA = selectedA.has(nodeId);
  const isSelectedB = selectedB.has(nodeId);
  const isSelectedLCA = selectedLCA.has(nodeId);

  let cls = "";
  if (isSelectedLCA) {
    cls = "selected-lca";
  } else if (isSelectedA) {
    cls = "selected-a";
  } else if (isSelectedB) {
    cls = "selected-b";
  } else if (isMatched) {
    cls = "matched";
  } else if (isVisited) {
    cls = "visited";
  }

  return (
    <div>
      <div className="tree-row" style={{ paddingLeft: depth * 28 }}>
        <div
          className={`tree-node ${cls} ${visible ? "is-visible" : "is-hidden"}`.trim()}
        >
          {"|- "}&lt;{nodeValue}&gt; <span className="tree-node-id">[id: {nodeId}]</span>
        </div>

        {attributes.length > 0 && (
          <div
            className={`attribute-node ${visible ? "is-visible" : "is-hidden"}`.trim()}
          >
            atribut: {formatAttributes(attributes)}
          </div>
        )}
      </div>

      {texts.map((text, index) => (
        <div
          className={`text-row ${visible ? "is-visible" : "is-hidden"}`.trim()}
          key={`${nodeId}-text-${index}`}
          style={{ marginLeft: depth * 28 + 34 }}
        >
          <div className="text-connector" aria-hidden="true" />
          <div className="text-node">text: {text}</div>
        </div>
      ))}

      {children.map((child, index) => (
        <TreeNode
          key={`${nodeId}-${index}-${child?.id ?? child?.animationOrder ?? index}`}
          node={child}
          depth={depth + 1}
          visited={visited}
          matched={matched}
          selectedA={selectedA}
          selectedB={selectedB}
          selectedLCA={selectedLCA}
          animate={animate}
        />
      ))}
    </div>
  );
}
