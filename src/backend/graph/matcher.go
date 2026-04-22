package graph

type SelectorMatcher interface {
	IsMatch(node *Node) bool
}
