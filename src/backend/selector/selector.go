package selector

import (
	"strings"
	"tubes2_izzat/src/backend/graph"
)

type TagMatcher struct {
	tagName string
}

func (tm *TagMatcher) IsMatch(node *graph.Node) bool {
	if node == nil {
		return false
	}
	return strings.ToLower(node.TagName) == strings.ToLower(tm.tagName)
}

type ClassMatcher struct {
	className string
}

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

type IDMatcher struct {
	id string
}

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

type UniversalMatcher struct{}

func (um *UniversalMatcher) IsMatch(node *graph.Node) bool {
	return node != nil
}

type AttributeMatcher struct {
	tagName   string
	attribute string
	value     string
}

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

type TagClassMatcher struct {
	tagName   string
	className string
}

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

type MulticlassMatcher struct {
	tagName    string
	classNames []string
}

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

	for _, requiredClass := range mcm.classNames {
		if !classMap[requiredClass] {
			return false
		}
	}
	return true
}

type CombinatorMatcher struct {
	left       graph.SelectorMatcher
	right      graph.SelectorMatcher
	combinator string
}

func (cm *CombinatorMatcher) IsMatch(node *graph.Node) bool {
	if node == nil || cm.right == nil || cm.left == nil {
		return false
	}

	if !cm.right.IsMatch(node) {
		return false
	}

	switch cm.combinator {
	case ">":
		return node.Parent != nil && cm.left.IsMatch(node.Parent)

	case " ":
		current := node.Parent
		for current != nil {
			if cm.left.IsMatch(current) {
				return true
			}
			current = current.Parent
		}
		return false

	case "+":
		if node.Parent == nil {
			return false
		}
		for i, sibling := range node.Parent.Children {
			if sibling == node && i > 0 {
				return cm.left.IsMatch(node.Parent.Children[i-1])
			}
		}
		return false

	case "~":
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

func ParseSelector(input string) graph.SelectorMatcher {
	input = strings.TrimSpace(input)

	if input == "" {
		return nil
	}

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

	if strings.Contains(input, "[") && strings.Contains(input, "]") {
		return parseAttributeSelector(input)
	}

	if strings.Contains(input, ".") && !strings.HasPrefix(input, ".") {
		parts := strings.Split(input, ".")
		if len(parts) >= 2 && parts[0] != "" {
			if len(parts) == 2 {
				return &TagClassMatcher{tagName: parts[0], className: parts[1]}
			}
			return &MulticlassMatcher{tagName: parts[0], classNames: parts[1:]}
		}
	}

	if strings.HasPrefix(input, ".") && strings.Count(input, ".") > 1 {
		classNames := strings.Split(strings.TrimPrefix(input, "."), ".")
		return &MulticlassMatcher{tagName: "", classNames: classNames}
	}

	if input == "*" {
		return &UniversalMatcher{}
	} else if strings.HasPrefix(input, ".") {
		return &ClassMatcher{className: input[1:]}
	} else if strings.HasPrefix(input, "#") {
		return &IDMatcher{id: input[1:]}
	}

	return &TagMatcher{tagName: input}
}

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

func findCombinatorOutsideBrackets(input string) (int, string) {
	inBrackets := false
	inQuotes := false
	var quoteChar rune

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
