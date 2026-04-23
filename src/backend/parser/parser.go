package parser

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"tubes2_izzat/src/backend/graph"

	"golang.org/x/net/html"
)

var voidElements = map[string]bool{
	"area": true, "base": true, "br": true, "col": true, "embed": true,
	"hr": true, "img": true, "input": true, "link": true, "meta": true,
	"param": true, "source": true, "track": true, "wbr": true,
}

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

// ValidateHTMLStructure performs a lightweight structural validation for manually entered HTML.
func ValidateHTMLStructure(input string) error {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return fmt.Errorf("HTML kosong")
	}

	if !strings.Contains(trimmed, "<") || !strings.Contains(trimmed, ">") {
		return fmt.Errorf("struktur HTML tidak valid: input tidak mengandung tag HTML")
	}

	tokenizer := html.NewTokenizer(strings.NewReader(input))
	stack := make([]string, 0)
	foundElement := false

	for {
		switch tokenizer.Next() {
		case html.ErrorToken:
			if err := tokenizer.Err(); err != nil && !errors.Is(err, io.EOF) {
				return fmt.Errorf("struktur HTML tidak valid: %v", err)
			}

			if !foundElement {
				return fmt.Errorf("struktur HTML tidak valid: tidak ditemukan elemen HTML")
			}

			if len(stack) > 0 {
				return fmt.Errorf("struktur HTML tidak valid: tag <%s> belum ditutup", stack[len(stack)-1])
			}

			return nil

		case html.StartTagToken:
			token := tokenizer.Token()
			tagName := strings.ToLower(token.Data)
			if tagName == "" {
				continue
			}

			foundElement = true
			if !voidElements[tagName] {
				stack = append(stack, tagName)
			}

		case html.SelfClosingTagToken:
			token := tokenizer.Token()
			if token.Data != "" {
				foundElement = true
			}

		case html.EndTagToken:
			token := tokenizer.Token()
			tagName := strings.ToLower(token.Data)
			if tagName == "" || voidElements[tagName] {
				continue
			}

			if len(stack) == 0 {
				return fmt.Errorf("struktur HTML tidak valid: ditemukan tag penutup </%s> tanpa tag pembuka", tagName)
			}

			expected := stack[len(stack)-1]
			if expected != tagName {
				return fmt.Errorf("struktur HTML tidak valid: tag penutup </%s> tidak cocok, seharusnya </%s>", tagName, expected)
			}

			stack = stack[:len(stack)-1]
		}
	}
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
