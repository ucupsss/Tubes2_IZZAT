package selector

import (
	"strings"
	"tubes2_izzat/dom"
)

// TagMatcher implements SelectorMatcher for HTML tag matching.
type TagMatcher struct {
	tagName string
}

// IsMatch checks if a node matches the tag criteria.
func (tm *TagMatcher) IsMatch(node *dom.Node) bool {
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
func (cm *ClassMatcher) IsMatch(node *dom.Node) bool {
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
func (im *IDMatcher) IsMatch(node *dom.Node) bool {
	if node == nil {
		return false
	}
	
	idAttr, exists := node.Attributes["id"]
	if !exists {
		return false
	}

	return idAttr == im.id
}

// ParseSelector parses a CSS selector string and returns the corresponding SelectorMatcher implementation.
// Supported formats: tag name, .className, #id
func ParseSelector(input string) dom.SelectorMatcher {
	input = strings.TrimSpace(input)

	if input == "" {
		return nil 
	}

	if strings.HasPrefix(input, ".") {
		return &ClassMatcher{className: input[1:]}
	} else if strings.HasPrefix(input, "#") {
		return &IDMatcher{id: input[1:]}
	}

	return &TagMatcher{tagName: input}
}