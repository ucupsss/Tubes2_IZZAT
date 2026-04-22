import { useEffect, useState } from "react";

export default function TreeNode({
  node,
  depth = 0,
  visited,
  matched,
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

  const nodeValue = node?.value ?? node?.tag ?? node?.id ?? "unknown";
  const children = Array.isArray(node?.children) ? node.children : [];
  const isMatched = matched.has(nodeValue);
  const isVisited = visited.has(nodeValue);
  const cls = isMatched ? "matched" : isVisited ? "visited" : "";

  return (
    <div>
      <div
        className={`tree-node ${cls} ${visible ? "is-visible" : "is-hidden"}`.trim()}
        style={{ paddingLeft: depth * 16 }}
      >
        {"|- "}&lt;{nodeValue}&gt;
      </div>
      {children.map((child, index) => (
        <TreeNode
          key={`${nodeValue}-${index}-${child?.animationOrder ?? index}`}
          node={child}
          depth={depth + 1}
          visited={visited}
          matched={matched}
          animate={animate}
        />
      ))}
    </div>
  );
}