package dom

import "sync/atomic"

var nextNodeID uint64

// NodeMeta stores traversal and ancestor data that higher-level algorithms
// can populate without changing the structural DOM representation.
type NodeMeta struct {
	Depth int
	Up    []*Node
}

// Node represents a single HTML element in the DOM tree.
type Node struct {
	ID         uint64
	TagName    string
	Attributes map[string]string
	Parent     *Node
	Children   []*Node
	Meta       NodeMeta
}

// NewNode creates a DOM node with a unique ID and a defensive copy of attributes.
func NewNode(tagName string, attributes map[string]string) *Node {
	return &Node{
		ID:         atomic.AddUint64(&nextNodeID, 1),
		TagName:    tagName,
		Attributes: cloneAttributes(attributes),
		Children:   make([]*Node, 0),
	}
}

// AddChild links a child to the node and updates the child's parent pointer.
func (n *Node) AddChild(child *Node) {
	if n == nil || child == nil {
		return
	}

	child.Parent = n
	n.Children = append(n.Children, child)
}

// SetDepth stores the node depth for traversal and binary lifting preparation.
func (n *Node) SetDepth(depth int) {
	if n == nil {
		return
	}

	n.Meta.Depth = depth
}

// InitUpTable allocates the binary lifting table with the requested size.
func (n *Node) InitUpTable(levels int) {
	if n == nil {
		return
	}

	if levels <= 0 {
		n.Meta.Up = nil
		return
	}

	n.Meta.Up = make([]*Node, levels)
}

func cloneAttributes(attributes map[string]string) map[string]string {
	if len(attributes) == 0 {
		return make(map[string]string)
	}

	cloned := make(map[string]string, len(attributes))
	for key, value := range attributes {
		cloned[key] = value
	}

	return cloned
}
