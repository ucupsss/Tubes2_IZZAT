package graph

import "time"

func SearchBFS(root *Node, matcher SelectorMatcher, limit int) ([]*Node, int, []*Node, time.Duration) {
	start := time.Now()
	if root == nil || matcher == nil || limit < 0 {
		return nil, 0, nil, time.Since(start)
	}

	results := make([]*Node, 0)
	traversalLog := make([]*Node, 0)
	queue := make([]*Node, 1, 16)
	queue[0] = root
	visitedCount := 0
	head := 0

	for head < len(queue) {
		current := queue[head]
		queue[head] = nil
		head++

		visitedCount++
		traversalLog = append(traversalLog, current)

		if matcher.IsMatch(current) {
			results = append(results, current)
			if limit > 0 && len(results) >= limit {
				break
			}
		}

		queue = append(queue, current.Children...)
	}

	return results, visitedCount, traversalLog, time.Since(start)
}

func SearchDFS(root *Node, matcher SelectorMatcher, limit int) ([]*Node, int, []*Node, time.Duration) {
	start := time.Now()
	if root == nil || matcher == nil || limit < 0 {
		return nil, 0, nil, time.Since(start)
	}

	results := make([]*Node, 0)
	traversalLog := make([]*Node, 0)
	stack := make([]*Node, 1, 16)
	stack[0] = root
	visitedCount := 0

	for len(stack) > 0 {
		last := len(stack) - 1
		current := stack[last]
		stack[last] = nil
		stack = stack[:last]

		visitedCount++
		traversalLog = append(traversalLog, current)

		if matcher.IsMatch(current) {
			results = append(results, current)
			if limit > 0 && len(results) >= limit {
				break
			}
		}

		for i := len(current.Children) - 1; i >= 0; i-- {
			stack = append(stack, current.Children[i])
		}
	}

	return results, visitedCount, traversalLog, time.Since(start)
}
