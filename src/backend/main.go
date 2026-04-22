package main

import (
	"fmt"
	"strings"
	"tubes2_izzat/dom"
	"tubes2_izzat/src/backend/scraper"
	"tubes2_izzat/src/backend/selector"
)

func main() {
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println("TESTING SELECTOR MATCHER")
	fmt.Println(strings.Repeat("=", 70))

	// ========== TEST 1: SIMPLE SELECTORS ==========
	fmt.Println("\n--- TEST 1: Simple Selectors ---\n")

	// Test Tag Selector
	fmt.Println("1.1: Tag Selector (div)")
	divNode := dom.NewNode("div", map[string]string{})
	tagMatcher := selector.ParseSelector("div")
	testResult("div selector on <div>", tagMatcher.IsMatch(divNode), true)

	// Test Class Selector
	fmt.Println("\n1.2: Class Selector (.container)")
	node1 := dom.NewNode("div", map[string]string{"class": "container"})
	classMatcher := selector.ParseSelector(".container")
	testResult(".container selector", classMatcher.IsMatch(node1), true)

	// Test ID Selector
	fmt.Println("\n1.3: ID Selector (#main)")
	node2 := dom.NewNode("div", map[string]string{"id": "main"})
	idMatcher := selector.ParseSelector("#main")
	testResult("#main selector", idMatcher.IsMatch(node2), true)

	// Test Universal Selector
	fmt.Println("\n1.4: Universal Selector (*)")
	universalMatcher := selector.ParseSelector("*")
	testResult("* selector on any node", universalMatcher.IsMatch(node1), true)

	// ========== TEST 2: TAG + CLASS ==========
	fmt.Println("\n--- TEST 2: Tag + Class Selector (p.intro) ---\n")

	p1 := dom.NewNode("p", map[string]string{"class": "intro"})
	p2 := dom.NewNode("p", map[string]string{"class": "outro"})
	div1 := dom.NewNode("div", map[string]string{"class": "intro"})

	tagClassMatcher := selector.ParseSelector("p.intro")
	testResult("p.intro on <p class='intro'>", tagClassMatcher.IsMatch(p1), true)
	testResult("p.intro on <p class='outro'>", tagClassMatcher.IsMatch(p2), false)
	testResult("p.intro on <div class='intro'>", tagClassMatcher.IsMatch(div1), false)

	// ========== TEST 3: MULTICLASS ==========
	fmt.Println("\n--- TEST 3: Multiclass Selector (.btn.primary) ---\n")

	btn1 := dom.NewNode("button", map[string]string{"class": "btn primary"})
	btn2 := dom.NewNode("button", map[string]string{"class": "btn"})
	btn3 := dom.NewNode("button", map[string]string{"class": "btn primary secondary"})

	multiclassMatcher := selector.ParseSelector(".btn.primary")
	testResult(".btn.primary on <button class='btn primary'>", multiclassMatcher.IsMatch(btn1), true)
	testResult(".btn.primary on <button class='btn'>", multiclassMatcher.IsMatch(btn2), false)
	testResult(".btn.primary on <button class='btn primary secondary'>", multiclassMatcher.IsMatch(btn3), true)

	// ========== TEST 4: ATTRIBUTE SELECTOR ==========
	fmt.Println("\n--- TEST 4: Attribute Selector (input[type=text]) ---\n")

	input1 := dom.NewNode("input", map[string]string{"type": "text"})
	input2 := dom.NewNode("input", map[string]string{"type": "password"})
	input3 := dom.NewNode("input", map[string]string{"type": "text", "name": "username"})

	attrMatcher := selector.ParseSelector("input[type=text]")
	testResult("input[type=text] on <input type='text'>", attrMatcher.IsMatch(input1), true)
	testResult("input[type=text] on <input type='password'>", attrMatcher.IsMatch(input2), false)
	testResult("input[type=text] on <input type='text' name='username'>", attrMatcher.IsMatch(input3), true)

	// Test attribute without tag
	fmt.Println("\n4.2: Attribute Selector without tag ([type=text])")
	attrMatcher2 := selector.ParseSelector("[type=text]")
	testResult("[type=text] on <input type='text'>", attrMatcher2.IsMatch(input1), true)

	// ========== TEST 5: COMBINATORS ==========
	fmt.Println("\n--- TEST 5: Combinator Selectors ---\n")

	// Build DOM structure:
	// <div id="container">
	//   <ul class="list">
	//     <li>Item 1</li>
	//     <li class="active">Item 2</li>
	//     <li>Item 3</li>
	//   </ul>
	//   <p>Paragraph</p>
	// </div>

	container := dom.NewNode("div", map[string]string{"id": "container"})
	ulList := dom.NewNode("ul", map[string]string{"class": "list"})
	li1 := dom.NewNode("li", map[string]string{})
	li2 := dom.NewNode("li", map[string]string{"class": "active"})
	li3 := dom.NewNode("li", map[string]string{})
	pTag := dom.NewNode("p", map[string]string{})

	container.AddChild(ulList)
	container.AddChild(pTag)
	ulList.AddChild(li1)
	ulList.AddChild(li2)
	ulList.AddChild(li3)

	fmt.Println("DOM Structure:")
	fmt.Println(`<div id="container">
  <ul class="list">
    <li>Item 1</li>
    <li class="active">Item 2</li>
    <li>Item 3</li>
  </ul>
  <p>Paragraph</p>
</div>`)

	// 5.1: Child Combinator (>)
	fmt.Println("\n5.1: Child Combinator (div > ul)")
	childMatcher1 := selector.ParseSelector("div > ul")
	testResult("div > ul on <ul> child of <div>", childMatcher1.IsMatch(ulList), true)
	testResult("div > ul on <li> child of <ul>", childMatcher1.IsMatch(li1), false)

	// 5.2: Descendant Combinator (space)
	fmt.Println("\n5.2: Descendant Combinator (div li)")
	descMatcher := selector.ParseSelector("div li")
	testResult("div li on <li> descendant of <div>", descMatcher.IsMatch(li1), true)
	testResult("div li on <li> descendant of <div>", descMatcher.IsMatch(li2), true)
	testResult("div li on <li> descendant of <div>", descMatcher.IsMatch(li3), true)

	// 5.3: Adjacent Sibling Combinator (+)
	fmt.Println("\n5.3: Adjacent Sibling Combinator (li + li)")
	adjMatcher := selector.ParseSelector("li + li")
	testResult("li + li on Item 2 (adjacent to Item 1)", adjMatcher.IsMatch(li2), true)
	testResult("li + li on Item 3 (adjacent to Item 2)", adjMatcher.IsMatch(li3), true)
	testResult("li + li on Item 1 (no previous sibling)", adjMatcher.IsMatch(li1), false)

	// 5.4: General Sibling Combinator (~)
	fmt.Println("\n5.4: General Sibling Combinator (li ~ li)")
	genMatcher := selector.ParseSelector("li ~ li")
	testResult("li ~ li on Item 2 (has sibling before)", genMatcher.IsMatch(li2), true)
	testResult("li ~ li on Item 3 (has sibling before)", genMatcher.IsMatch(li3), true)
	testResult("li ~ li on Item 1 (no sibling before)", genMatcher.IsMatch(li1), false)

	// 5.5: Class with Descendant Combinator
	fmt.Println("\n5.5: Class with Combinator (ul .active)")
	classDescMatcher := selector.ParseSelector("ul .active")
	testResult("ul .active on <li class='active'>", classDescMatcher.IsMatch(li2), true)
	testResult("ul .active on <li> without class", classDescMatcher.IsMatch(li1), false)

	fmt.Println("\n" + strings.Repeat("=", 70))
	fmt.Println("TESTING SCRAPER")
	fmt.Println(strings.Repeat("=", 70))

	fmt.Println("\n--- TEST 6: Scraper FetchHTML ---\n")
	
	// Test on a known fast website
	url := "http://example.com"
	fmt.Printf("Fetching HTML from: %s\n", url)
	htmlContent, err := scraper.FetchHTML(url)
	
	if err != nil {
		fmt.Printf("  Fetch failed: %v\n", err)
	} else {
		fmt.Printf("  Fetch successful! Received %d bytes.\n", len(htmlContent))
		hasTitle := strings.Contains(htmlContent, "<title>Example Domain</title>")
		testResult("Contains <title>Example Domain</title>", hasTitle, true)
	}

	fmt.Println("\n" + strings.Repeat("=", 70))
	fmt.Println("ALL TESTS COMPLETED")
	fmt.Println(strings.Repeat("=", 70))
}

// testResult prints test result in a clear format
func testResult(testName string, actual, expected bool) {
	status := "✓"
	if actual != expected {
		status = "✗"
	}
	result := "PASS"
	if actual != expected {
		result = "FAIL"
	}
	fmt.Printf("  %s %s: %v\n", status, testName, result)
}