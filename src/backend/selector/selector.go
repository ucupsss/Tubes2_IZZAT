package selector

import (
	"strings"
	"tubes2_izzat/src/backend/graph"
)

// TagMatcher implements SelectorMatcher for HTML tag matching.
type TagMatcher struct {
	tagName string
}

// IsMatch checks if a node matches the tag criteria.
func (tm *TagMatcher) IsMatch(node *graph.Node) bool {
	if node == nil {
		return false
	}
	return strings.ToLower(node.TagName) == strings.ToLower(tm.tagName)
}

// ClassMatcher implements SelectorMatcher for HTML class attribute matching.
type ClassMatcher struct {
	className string
}

// IsMatch checks if a node matches the class criteria.
func (cm *ClassMatcher) IsMatch(node *graph.Node) bool {
	if node == nil {
		return false
	}

	classAttr, exists := node.Attributes["class"]
	if !exists {
		return false
	}

	classes := strings.Fields(classAttr)
	for _, cls := range classes {
		if cls == cm.className {
			return true
		}
	}
	return false
}

// IDMatcher implements SelectorMatcher for HTML id attribute matching.
type IDMatcher struct {
	id string
}

// IsMatch checks if a node matches the id criteria.
func (im *IDMatcher) IsMatch(node *graph.Node) bool {
	if node == nil {
		return false
	}

	idAttr, exists := node.Attributes["id"]
	if !exists {
		return false
	}

	return idAttr == im.id
}

// UniversalMatcher implements SelectorMatcher for universal selector matching (*).
type UniversalMatcher struct{}

// IsMatch checks if a node matches (always true for universal selector).
func (um *UniversalMatcher) IsMatch(node *graph.Node) bool {
	return node != nil
}

// AttributeMatcher implements SelectorMatcher for attribute matching (e.g., a[href=example.com]).
type AttributeMatcher struct {
	tagName   string
	attribute string
	value     string
}

// IsMatch checks if a node matches the attribute criteria.
func (am *AttributeMatcher) IsMatch(node *graph.Node) bool {
	if node == nil {
		return false
	}

	if am.tagName != "" && strings.ToLower(node.TagName) != strings.ToLower(am.tagName) {
		return false
	}

	attrValue, exists := node.Attributes[strings.ToLower(am.attribute)]
	if !exists {
		return false
	}

	return attrValue == am.value
}

// TagClassMatcher implements SelectorMatcher for tag+class matching (e.g., p.intro).
type TagClassMatcher struct {
	tagName   string
	className string
}

// IsMatch checks if a node matches the tag and class criteria.
func (tcm *TagClassMatcher) IsMatch(node *graph.Node) bool {
	if node == nil {
		return false
	}

	if strings.ToLower(node.TagName) != strings.ToLower(tcm.tagName) {
		return false
	}

	classAttr, exists := node.Attributes["class"]
	if !exists {
		return false
	}

	classes := strings.Fields(classAttr)
	for _, cls := range classes {
		if cls == tcm.className {
			return true
		}
	}
	return false
}

// MulticlassMatcher implements SelectorMatcher for multiclass matching (e.g., .btn.primary or p.btn.primary).
type MulticlassMatcher struct {
	tagName    string
	classNames []string
}

// IsMatch checks if a node has all the required classes.
func (mcm *MulticlassMatcher) IsMatch(node *graph.Node) bool {
	if node == nil || len(mcm.classNames) == 0 {
		return false
	}

	if mcm.tagName != "" && strings.ToLower(node.TagName) != strings.ToLower(mcm.tagName) {
		return false
	}

	classAttr, exists := node.Attributes["class"]
	if !exists {
		return false
	}

	classes := strings.Fields(classAttr)
	classMap := make(map[string]bool)
	for _, cls := range classes {
		classMap[cls] = true
	}

	// Check if all required classes are present
	for _, requiredClass := range mcm.classNames {
		if !classMap[requiredClass] {
			return false
		}
	}
	return true
}

// CombinatorMatcher implements SelectorMatcher for CSS combinators.
// Combinators: > (child), space (descendant), + (adjacent sibling), ~ (general sibling)
type CombinatorMatcher struct {
	left       graph.SelectorMatcher
	right      graph.SelectorMatcher
	combinator string // ">", " ", "+", "~"
}

// IsMatch checks if a node matches the combinator criteria.
func (cm *CombinatorMatcher) IsMatch(node *graph.Node) bool {
	if node == nil || cm.right == nil || cm.left == nil {
		return false
	}

	// First, check if the node matches the right selector
	if !cm.right.IsMatch(node) {
		return false
	}

	switch cm.combinator {
	case ">": // Child combinator: immediate parent must match left selector
		return node.Parent != nil && cm.left.IsMatch(node.Parent)

	case " ": // Descendant combinator: any ancestor must match left selector
		current := node.Parent
		for current != nil {
			if cm.left.IsMatch(current) {
				return true
			}
			current = current.Parent
		}
		return false

	case "+": // Adjacent sibling combinator: previous sibling must match left selector
		if node.Parent == nil {
			return false
		}
		for i, sibling := range node.Parent.Children {
			if sibling == node && i > 0 {
				return cm.left.IsMatch(node.Parent.Children[i-1])
			}
		}
		return false

	case "~": // General sibling combinator: any earlier sibling must match left selector
		if node.Parent == nil {
			return false
		}
		for i, sibling := range node.Parent.Children {
			if sibling == node {
				for j := 0; j < i; j++ {
					if cm.left.IsMatch(node.Parent.Children[j]) {
						return true
					}
				}
				return false
			}
		}
		return false

	default:
		return false
	}
}

// ParseSelector parses a CSS selector string and returns the corresponding SelectorMatcher implementation.
// Supported formats:
// - * (universal)
// - tag name (e.g., div, p)
// - .className (class selector)
// - .class1.class2 (multiclass selector)
// - #id (id selector)
// - tag.class (tag+class)
// - tag[attr=value] (attribute selector)
// - a[href=value] (tag+attribute)
// Supported combinators: > (child), space (descendant), + (adjacent sibling), ~ (general sibling)
func ParseSelector(input string) graph.SelectorMatcher {
	input = strings.TrimSpace(input)

	if input == "" {
		return nil
	}

	// Check for combinators (but not within brackets)
	combIdx, combType := findCombinatorOutsideBrackets(input)
	if combIdx != -1 {
		leftStr := input[:combIdx]
		rightStr := input[combIdx+len(combType):]

		left := ParseSelector(strings.TrimSpace(leftStr))
		right := ParseSelector(strings.TrimSpace(rightStr))

		comb := combType
		if comb != " " {
			comb = strings.TrimSpace(combType)
		}

		if left != nil && right != nil {
			return &CombinatorMatcher{
				left:       left,
				right:      right,
				combinator: comb,
			}
		}
	}

	// Check for attribute selector [attr=value]
	if strings.Contains(input, "[") && strings.Contains(input, "]") {
		return parseAttributeSelector(input)
	}

	// Check for tag+class or tag+multiclass (e.g., p.classname or p.class1.class2)
	if strings.Contains(input, ".") && !strings.HasPrefix(input, ".") {
		parts := strings.Split(input, ".")
		if len(parts) >= 2 && parts[0] != "" {
			if len(parts) == 2 {
				return &TagClassMatcher{tagName: parts[0], className: parts[1]}
			}
			return &MulticlassMatcher{tagName: parts[0], classNames: parts[1:]}
		}
	}

	// Check for multiclass selector (.class1.class2...)
	if strings.HasPrefix(input, ".") && strings.Count(input, ".") > 1 {
		classNames := strings.Split(strings.TrimPrefix(input, "."), ".")
		return &MulticlassMatcher{tagName: "", classNames: classNames}
	}

	// Simple selectors
	if input == "*" {
		return &UniversalMatcher{}
	} else if strings.HasPrefix(input, ".") {
		return &ClassMatcher{className: input[1:]}
	} else if strings.HasPrefix(input, "#") {
		return &IDMatcher{id: input[1:]}
	}

	return &TagMatcher{tagName: input}
}

// parseAttributeSelector parses attribute selectors like a[href=value] or [attr=value]
func parseAttributeSelector(input string) graph.SelectorMatcher {
	bracketIdx := strings.Index(input, "[")
	tagName := ""

	if bracketIdx > 0 {
		tagName = input[:bracketIdx]
	}

	start := strings.Index(input, "[")
	end := strings.Index(input, "]")

	if start == -1 || end == -1 || start >= end {
		return nil
	}

	attrContent := input[start+1 : end]

	parts := strings.Split(attrContent, "=")
	if len(parts) != 2 {
		return nil
	}

	attribute := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])

	value = strings.Trim(value, "\"'")

	return &AttributeMatcher{
		tagName:   tagName,
		attribute: attribute,
		value:     value,
	}
}

// findCombinatorOutsideBrackets finds the first combinator outside of brackets []
// Returns the index of the combinator and the string matched (e.g. " > ", " + ", "~", or " " for descendant)
func findCombinatorOutsideBrackets(input string) (int, string) {
	inBrackets := false
	inQuotes := false
	var quoteChar rune

	// look for combinators from right to left to build the tree correctly
	// e.g., "div > p span" -> combinator " ", left "div > p", right "span"
	for i := len(input) - 1; i >= 0; i-- {
		char := rune(input[i])

		if char == ']' {
			inBrackets = true
		} else if char == '[' {
			inBrackets = false
		} else if char == '"' || char == '\'' {
			if !inQuotes {
				inQuotes = true
				quoteChar = char
			} else if quoteChar == char {
				inQuotes = false
			}
		}

		if !inBrackets && !inQuotes {
			if char == '>' || char == '+' || char == '~' {
				return i, string(char)
			}
			// check for descendant space (must not be adjacent to another combinator, and not trailing space)
			if char == ' ' && i > 0 && i < len(input)-1 {
				prevChar := rune(input[i-1])
				nextChar := rune(input[i+1])
				if prevChar != '>' && prevChar != '+' && prevChar != '~' && prevChar != ' ' &&
					nextChar != '>' && nextChar != '+' && nextChar != '~' && nextChar != ' ' {
					return i, " "
				}
			}
		}
	}

	return -1, ""
}
