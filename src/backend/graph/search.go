package graph

import (
	"runtime"
	"sync"
	"time"
)

type bfsVisitResult struct {
	node     *Node
	matched  bool
	children []*Node
}

type dfsVisitResult struct {
	traversal []*Node
	matches   []*Node
}

func SearchBFS(root *Node, matcher SelectorMatcher, limit int) ([]*Node, int, []*Node, time.Duration) {
	start := time.Now()
	if root == nil || matcher == nil || limit < 0 {
		return nil, 0, nil, time.Since(start)
	}

	results := make([]*Node, 0)
	traversalLog := make([]*Node, 0)
	frontier := []*Node{root}

	for len(frontier) > 0 {
		levelResults := processBFSLevel(frontier, matcher)
		nextFrontier := make([]*Node, 0)
		stop := false

		for _, item := range levelResults {
			if item.node == nil {
				continue
			}

			traversalLog = append(traversalLog, item.node)
			if item.matched {
				results = append(results, item.node)
				if limit > 0 && len(results) >= limit {
					stop = true
				}
			}

			if stop {
				break
			}

			nextFrontier = append(nextFrontier, item.children...)
		}

		if stop {
			break
		}

		frontier = nextFrontier
	}

	return results, len(traversalLog), traversalLog, time.Since(start)
}

func SearchDFS(root *Node, matcher SelectorMatcher, limit int) ([]*Node, int, []*Node, time.Duration) {
	start := time.Now()
	if root == nil || matcher == nil || limit < 0 {
		return nil, 0, nil, time.Since(start)
	}

	sem := make(chan struct{}, traversalWorkerCount())
	result := processDFSNode(root, matcher, sem)

	traversalLog := result.traversal
	results := result.matches

	if limit > 0 && len(results) >= limit {
		trimmedMatches := results[:limit]
		allowed := make(map[uint64]struct{}, len(trimmedMatches))
		for _, node := range trimmedMatches {
			if node != nil {
				allowed[node.ID] = struct{}{}
			}
		}

		matchCount := 0
		cutIndex := len(traversalLog)
		for index, node := range traversalLog {
			if node == nil {
				continue
			}
			if _, ok := allowed[node.ID]; ok {
				matchCount++
				if matchCount == limit {
					cutIndex = index + 1
					break
				}
			}
		}

		results = trimmedMatches
		traversalLog = traversalLog[:cutIndex]
	}

	return results, len(traversalLog), traversalLog, time.Since(start)
}

func processBFSLevel(frontier []*Node, matcher SelectorMatcher) []bfsVisitResult {
	results := make([]bfsVisitResult, len(frontier))
	if len(frontier) == 0 {
		return results
	}

	workerCount := minInt(traversalWorkerCount(), len(frontier))
	jobs := make(chan int, len(frontier))
	var waitGroup sync.WaitGroup

	for worker := 0; worker < workerCount; worker++ {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			for index := range jobs {
				node := frontier[index]
				if node == nil {
					continue
				}

				children := make([]*Node, 0, len(node.Children))
				for _, child := range node.Children {
					if child != nil {
						children = append(children, child)
					}
				}

				results[index] = bfsVisitResult{
					node:     node,
					matched:  matcher.IsMatch(node),
					children: children,
				}
			}
		}()
	}

	for index := range frontier {
		jobs <- index
	}
	close(jobs)
	waitGroup.Wait()

	return results
}

func processDFSNode(node *Node, matcher SelectorMatcher, sem chan struct{}) dfsVisitResult {
	if node == nil {
		return dfsVisitResult{
			traversal: []*Node{},
			matches:   []*Node{},
		}
	}

	result := dfsVisitResult{
		traversal: []*Node{node},
		matches:   []*Node{},
	}
	if matcher.IsMatch(node) {
		result.matches = append(result.matches, node)
	}

	childResults := make([]dfsVisitResult, len(node.Children))
	var waitGroup sync.WaitGroup

	for index, child := range node.Children {
		if child == nil {
			continue
		}

		if tryAcquireWorkerSlot(sem) {
			waitGroup.Add(1)
			go func(position int, currentChild *Node) {
				defer waitGroup.Done()
				defer releaseWorkerSlot(sem)
				childResults[position] = processDFSNode(currentChild, matcher, sem)
			}(index, child)
			continue
		}

		childResults[index] = processDFSNode(child, matcher, sem)
	}

	waitGroup.Wait()

	for _, childResult := range childResults {
		result.traversal = append(result.traversal, childResult.traversal...)
		result.matches = append(result.matches, childResult.matches...)
	}

	return result
}

func traversalWorkerCount() int {
	workerCount := runtime.NumCPU()
	if workerCount < 2 {
		return 2
	}
	return workerCount
}

func tryAcquireWorkerSlot(sem chan struct{}) bool {
	select {
	case sem <- struct{}{}:
		return true
	default:
		return false
	}
}

func releaseWorkerSlot(sem chan struct{}) {
	select {
	case <-sem:
	default:
	}
}