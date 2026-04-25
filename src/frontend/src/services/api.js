const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || "http://localhost:5175";

export async function runTraversal(data) {
  const response = await fetch(`${API_BASE_URL}/api/traversal`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(data),
  });

  let payload = null;

  try {
    payload = await response.json();
  } catch {
    payload = null;
  }

  if (!response.ok) {
    throw new Error(payload?.error || `Request failed with status ${response.status}`);
  }

  return payload || { tree: null, visited: [], matched: [], time: 0 };
}

export async function runLCA(data) {
  const response = await fetch(`${API_BASE_URL}/api/lca`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(data),
  });

  let payload = null;

  try {
    payload = await response.json();
  } catch {
    payload = null;
  }

  if (!response.ok) {
    throw new Error(payload?.error || `Request failed with status ${response.status}`);
  }

  return payload || { nodeA: null, nodeB: null, lca: null };
}
