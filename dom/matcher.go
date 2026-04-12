package dom

// SelectorMatcher abstracts CSS selector matching so traversal code can stay
// decoupled from selector parsing and evaluation.
type SelectorMatcher interface {
	IsMatch(node *Node) bool
}
