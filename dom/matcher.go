package dom

type SelectorMatcher interface {
	IsMatch(node *Node) bool
}
