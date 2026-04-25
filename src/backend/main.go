package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
	"tubes2_izzat/src/backend/graph"
	"tubes2_izzat/src/backend/parser"
	"tubes2_izzat/src/backend/scraper"
	"tubes2_izzat/src/backend/selector"
)

const defaultPort = "5175"

type traversalRequest struct {
	URL       string `json:"url"`
	HTML      string `json:"html"`
	Selector  string `json:"selector"`
	Algorithm string `json:"algorithm"`
	Limit     int    `json:"limit"`
	LCA       bool   `json:"lca"`
	FirstNode string `json:"firstNodeId"`
	SecondNode string `json:"secondNodeId"`
}

type treeResponse struct {
	ID         string            `json:"id"`
	Value      string            `json:"value"`
	Tag        string            `json:"tag"`
	Attributes map[string]string `json:"attributes,omitempty"`
	Text       string            `json:"text,omitempty"`
	Texts      []string          `json:"texts,omitempty"`
	Depth      int               `json:"depth"`
	Children   []treeResponse    `json:"children"`
}

type logEntryResponse struct {
	ID      string `json:"id"`
	Tag     string `json:"tag"`
	Depth   int    `json:"depth"`
	Matched bool   `json:"matched"`
}

type traversalResponse struct {
	Tree         *treeResponse      `json:"tree"`
	Visited      []string           `json:"visited"`
	Matched      []string           `json:"matched"`
	TraversalLog []logEntryResponse `json:"traversalLog"`
	LCA          *lcaResponse       `json:"lca,omitempty"`
	Time         float64            `json:"time"`
	VisitedCount int                `json:"visitedCount"`
	MatchedCount int                `json:"matchedCount"`
	MaxDepth     int                `json:"maxDepth"`
	NodeCount    int                `json:"nodeCount"`
	Algorithm    string             `json:"algorithm"`
	Selector     string             `json:"selector"`
	Source       string             `json:"source"`
}

type errorResponse struct {
	Error string `json:"error"`
}

type lcaResponse struct {
	Enabled        bool   `json:"enabled"`
	FirstNodeID    string `json:"firstNodeId"`
	FirstNodeLabel string `json:"firstNodeLabel"`
	SecondNodeID   string `json:"secondNodeId"`
	SecondNodeLabel string `json:"secondNodeLabel"`
	AncestorID     string `json:"ancestorId"`
	AncestorLabel  string `json:"ancestorLabel"`
	AncestorDepth  int    `json:"ancestorDepth"`
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/health", handleHealth)
	mux.HandleFunc("/api/traversal", handleTraversal)

	addr := ":" + defaultPort
	log.Printf("backend listening on http://localhost%s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	withCORS(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func handleTraversal(w http.ResponseWriter, r *http.Request) {
	withCORS(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method tidak didukung")
		return
	}

	var request traversalRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "payload JSON tidak valid")
		return
	}

	response, err := runTraversal(request)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, response)
}

func runTraversal(request traversalRequest) (*traversalResponse, error) {
	request.Selector = strings.TrimSpace(request.Selector)
	request.Algorithm = strings.ToLower(strings.TrimSpace(request.Algorithm))
	request.URL = strings.TrimSpace(request.URL)
	request.HTML = strings.TrimSpace(request.HTML)

	if request.Selector == "" {
		return nil, fmt.Errorf("CSS selector wajib diisi")
	}

	if request.Algorithm == "" {
		request.Algorithm = "bfs"
	}
	if request.Algorithm != "bfs" && request.Algorithm != "dfs" {
		return nil, fmt.Errorf("algoritma harus bfs atau dfs")
	}

	if request.Limit < 0 {
		return nil, fmt.Errorf("jumlah hasil tidak boleh negatif")
	}

	htmlInput, source, err := resolveHTMLInput(request)
	if err != nil {
		return nil, err
	}

	root, err := parser.ParseHTML(htmlInput)
	if err != nil {
		return nil, err
	}

	matcher := selector.ParseSelector(request.Selector)
	if matcher == nil {
		return nil, fmt.Errorf("CSS selector tidak valid")
	}

	var (
		matches      []*graph.Node
		visitedCount int
		traversalLog []*graph.Node
		elapsed      time.Duration
	)

	if request.Algorithm == "dfs" {
		matches, visitedCount, traversalLog, elapsed = graph.SearchDFS(root, matcher, request.Limit)
	} else {
		matches, visitedCount, traversalLog, elapsed = graph.SearchBFS(root, matcher, request.Limit)
	}

	matchedIDs := nodeIDs(matches)
	visitedIDs := nodeIDs(traversalLog)
	matchedSet := make(map[uint64]bool, len(matches))
	for _, node := range matches {
		if node != nil {
			matchedSet[node.ID] = true
		}
	}

	tree := serializeTree(root)
	lcaResult, err := resolveLCA(root, request)
	if err != nil {
		return nil, err
	}
	return &traversalResponse{
		Tree:         tree,
		Visited:      visitedIDs,
		Matched:      matchedIDs,
		TraversalLog: serializeTraversalLog(traversalLog, matchedSet),
		LCA:          lcaResult,
		Time:         float64(elapsed.Microseconds()) / 1000.0,
		VisitedCount: visitedCount,
		MatchedCount: len(matches),
		MaxDepth:     graph.MaxDepth(root),
		NodeCount:    graph.CountNodes(root),
		Algorithm:    request.Algorithm,
		Selector:     request.Selector,
		Source:       source,
	}, nil
}

func resolveLCA(root *graph.Node, request traversalRequest) (*lcaResponse, error) {
	if !request.LCA {
		return nil, nil
	}

	nodesByID := collectNodesByID(root)
	if len(nodesByID) < 2 {
		return nil, fmt.Errorf("mode LCA membutuhkan minimal dua node pada pohon DOM")
	}

	firstNode, secondNode, err := pickLCANodes(nodesByID, request.FirstNode, request.SecondNode)
	if err != nil {
		return nil, err
	}

	ancestor := graph.GetLCA(firstNode, secondNode)
	if ancestor == nil {
		return nil, fmt.Errorf("LCA tidak ditemukan untuk dua node yang dipilih")
	}

	return &lcaResponse{
		Enabled:         true,
		FirstNodeID:     formatNodeID(firstNode),
		FirstNodeLabel:  nodeLabel(firstNode),
		SecondNodeID:    formatNodeID(secondNode),
		SecondNodeLabel: nodeLabel(secondNode),
		AncestorID:      formatNodeID(ancestor),
		AncestorLabel:   nodeLabel(ancestor),
		AncestorDepth:   ancestor.Meta.Depth,
	}, nil
}

func collectNodesByID(root *graph.Node) map[string]*graph.Node {
	nodesByID := make(map[string]*graph.Node)
	if root == nil {
		return nodesByID
	}

	stack := []*graph.Node{root}
	for len(stack) > 0 {
		last := len(stack) - 1
		current := stack[last]
		stack = stack[:last]
		if current == nil {
			continue
		}

		nodesByID[formatNodeID(current)] = current
		stack = append(stack, current.Children...)
	}

	return nodesByID
}

func pickLCANodes(nodesByID map[string]*graph.Node, firstID, secondID string) (*graph.Node, *graph.Node, error) {
	firstID = strings.TrimSpace(firstID)
	secondID = strings.TrimSpace(secondID)

	if firstID != "" && secondID != "" {
		firstNode := nodesByID[firstID]
		if firstNode == nil {
			return nil, nil, fmt.Errorf("node pertama untuk LCA tidak ditemukan")
		}

		secondNode := nodesByID[secondID]
		if secondNode == nil {
			return nil, nil, fmt.Errorf("node kedua untuk LCA tidak ditemukan")
		}

		return firstNode, secondNode, nil
	}

	orderedKeys := make([]string, 0, len(nodesByID))
	for id := range nodesByID {
		orderedKeys = append(orderedKeys, id)
	}

	sort.Slice(orderedKeys, func(i, j int) bool {
		parsedA, errA := strconv.ParseUint(orderedKeys[i], 10, 64)
		parsedB, errB := strconv.ParseUint(orderedKeys[j], 10, 64)
		if errA != nil || errB != nil {
			return orderedKeys[i] < orderedKeys[j]
		}
		return parsedA < parsedB
	})

	if len(orderedKeys) < 2 {
		return nil, nil, fmt.Errorf("mode LCA membutuhkan minimal dua node pada pohon DOM")
	}

	return nodesByID[orderedKeys[0]], nodesByID[orderedKeys[1]], nil
}

func resolveHTMLInput(request traversalRequest) (string, string, error) {
	if request.URL != "" {
		htmlInput, err := scraper.FetchHTML(request.URL)
		if err != nil {
			return "", "", err
		}
		return htmlInput, "url", nil
	}

	if request.HTML == "" {
		return "", "", fmt.Errorf("input HTML atau URL wajib diisi")
	}

	if err := parser.ValidateHTMLStructure(request.HTML); err != nil {
		return "", "", err
	}

	return request.HTML, "html", nil
}

func serializeTree(node *graph.Node) *treeResponse {
	if node == nil {
		return nil
	}

	children := make([]treeResponse, 0, len(node.Children))
	for _, child := range node.Children {
		serialized := serializeTree(child)
		if serialized != nil {
			children = append(children, *serialized)
		}
	}

	return &treeResponse{
		ID:         formatNodeID(node),
		Value:      nodeLabel(node),
		Tag:        node.TagName,
		Attributes: node.Attributes,
		Text:       strings.Join(node.Texts, " "),
		Texts:      node.Texts,
		Depth:      node.Meta.Depth,
		Children:   children,
	}
}

func serializeTraversalLog(nodes []*graph.Node, matchedSet map[uint64]bool) []logEntryResponse {
	logs := make([]logEntryResponse, 0, len(nodes))
	for _, node := range nodes {
		if node == nil {
			continue
		}
		logs = append(logs, logEntryResponse{
			ID:      formatNodeID(node),
			Tag:     nodeLabel(node),
			Depth:   node.Meta.Depth,
			Matched: matchedSet[node.ID],
		})
	}
	return logs
}

func nodeIDs(nodes []*graph.Node) []string {
	ids := make([]string, 0, len(nodes))
	for _, node := range nodes {
		if node != nil {
			ids = append(ids, formatNodeID(node))
		}
	}
	return ids
}

func nodeLabel(node *graph.Node) string {
	if node == nil {
		return "unknown"
	}

	parts := []string{node.TagName}
	if id := node.Attributes["id"]; id != "" {
		parts = append(parts, "#"+id)
	}
	if className := node.Attributes["class"]; className != "" {
		for _, classPart := range strings.Fields(className) {
			parts = append(parts, "."+classPart)
		}
	}
	return strings.Join(parts, "")
}

func formatNodeID(node *graph.Node) string {
	return strconv.FormatUint(node.ID, 10)
}

func withCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("failed writing response: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, errorResponse{Error: message})
}