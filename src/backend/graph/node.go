package graph

type NodeMeta struct {
	Depth int
	Up    []*Node
}

type Node struct {
	ID         uint64
	TagName    string
	Attributes map[string]string
	Texts      []string
	Parent     *Node
	Children   []*Node
	Meta       NodeMeta
}

func NewNode(tagName string, attributes map[string]string) *Node {
	return &Node{
		TagName:    tagName,
		Attributes: cloneAttributes(attributes),
		Texts:      make([]string, 0),
		Children:   make([]*Node, 0),
	}
}

func AssignStableIDs(root *Node) {
	if root == nil {
		return
	}

	var nextID uint64 = 1
	stack := []*Node{root}

	for len(stack) > 0 {
		last := len(stack) - 1
		current := stack[last]
		stack = stack[:last]
		if current == nil {
			continue
		}

		current.ID = nextID
		nextID++

		for i := len(current.Children) - 1; i >= 0; i-- {
			stack = append(stack, current.Children[i])
		}
	}
}

func (n *Node) AddChild(child *Node) {
	if n == nil || child == nil {
		return
	}

	child.Parent = n
	n.Children = append(n.Children, child)
}

func (n *Node) AddText(text string) {
	if n == nil {
		return
	}

	if text != "" {
		n.Texts = append(n.Texts, text)
	}
}

func (n *Node) SetDepth(depth int) {
	if n == nil {
		return
	}

	n.Meta.Depth = depth
}

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

func MaxDepth(root *Node) int {
	if root == nil {
		return 0
	}

	maxDepth := 0
	stack := []nodeDepth{{node: root, depth: 0}}

	for len(stack) > 0 {
		last := len(stack) - 1
		current := stack[last]
		stack = stack[:last]

		if current.depth > maxDepth {
			maxDepth = current.depth
		}

		for _, child := range current.node.Children {
			if child != nil {
				stack = append(stack, nodeDepth{node: child, depth: current.depth + 1})
			}
		}
	}

	return maxDepth
}

func CountNodes(root *Node) int {
	if root == nil {
		return 0
	}

	count := 0
	stack := []*Node{root}

	for len(stack) > 0 {
		last := len(stack) - 1
		current := stack[last]
		stack = stack[:last]

		if current == nil {
			continue
		}

		count++
		stack = append(stack, current.Children...)
	}

	return count
}

type nodeDepth struct {
	node  *Node
	depth int
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
