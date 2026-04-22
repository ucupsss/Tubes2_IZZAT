package parser

import (
	"fmt"
	"strings"
	"tubes2_izzat/src/backend/graph"

	"golang.org/x/net/html"
)

// ParseHTML converts an HTML document string into a graph.Node tree.
func ParseHTML(input string) (*graph.Node, error) {
	if strings.TrimSpace(input) == "" {
		return nil, fmt.Errorf("HTML kosong")
	}

	document, err := html.Parse(strings.NewReader(input))
	if err != nil {
		return nil, fmt.Errorf("gagal parsing HTML: %w", err)
	}

	nodes := convertChildren(document)
	if len(nodes) == 0 {
		return nil, fmt.Errorf("tidak ditemukan elemen HTML")
	}

	if len(nodes) == 1 {
		prepareTree(nodes[0])
		return nodes[0], nil
	}

	root := graph.NewNode("document", nil)
	for _, node := range nodes {
		root.AddChild(node)
	}
	prepareTree(root)

	return root, nil
}

func convertChildren(htmlNode *html.Node) []*graph.Node {
	nodes := make([]*graph.Node, 0)

	for child := htmlNode.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.ElementNode {
			nodes = append(nodes, convertElement(child))
			continue
		}

		nodes = append(nodes, convertChildren(child)...)
	}

	return nodes
}

func convertElement(htmlNode *html.Node) *graph.Node {
	node := graph.NewNode(strings.ToLower(htmlNode.Data), convertAttributes(htmlNode.Attr))

	for child := htmlNode.FirstChild; child != nil; child = child.NextSibling {
		switch child.Type {
		case html.TextNode:
			text := strings.Join(strings.Fields(child.Data), " ")
			node.AddText(text)
		case html.ElementNode:
			node.AddChild(convertElement(child))
		default:
			for _, nestedChild := range convertChildren(child) {
				node.AddChild(nestedChild)
			}
		}
	}

	return node
}

func convertAttributes(attributes []html.Attribute) map[string]string {
	result := make(map[string]string, len(attributes))

	for _, attr := range attributes {
		result[strings.ToLower(attr.Key)] = attr.Val
	}

	return result
}

func prepareTree(root *graph.Node) {
	graph.ComputeLCA(root, graph.CountNodes(root))
}
