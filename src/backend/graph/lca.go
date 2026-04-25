package graph

const defaultLCALog = 20

func ComputeLCA(root *Node, maxNodes int) {
	if root == nil {
		return
	}

	log := computeLCALog(maxNodes)
	stack := make([]*Node, 1, 16)
	stack[0] = root

	root.Parent = nil
	root.Meta.Depth = 0
	root.InitUpTable(log)

	for len(stack) > 0 {
		last := len(stack) - 1
		current := stack[last]
		stack[last] = nil
		stack = stack[:last]

		for _, child := range current.Children {
			if child == nil {
				continue
			}

			child.Parent = current
			child.Meta.Depth = current.Meta.Depth + 1
			child.InitUpTable(log)
			child.Meta.Up[0] = current

			for i := 1; i < log; i++ {
				ancestor := child.Meta.Up[i-1]
				if ancestor == nil {
					break
				}
				child.Meta.Up[i] = ancestor.Meta.Up[i-1]
			}

			stack = append(stack, child)
		}
	}
}

func GetLCA(u, v *Node) *Node {
	if u == nil || v == nil {
		return nil
	}

	if len(u.Meta.Up) == 0 || len(v.Meta.Up) == 0 {
		if u == v {
			return u
		}
		return nil
	}

	if u.Meta.Depth < v.Meta.Depth {
		u, v = v, u
	}

	u = liftNode(u, u.Meta.Depth-v.Meta.Depth)
	if u == nil {
		return nil
	}

	if u == v {
		return u
	}

	log := minInt(len(u.Meta.Up), len(v.Meta.Up))
	for i := log - 1; i >= 0; i-- {
		if u.Meta.Up[i] != v.Meta.Up[i] {
			u = u.Meta.Up[i]
			v = v.Meta.Up[i]
		}
	}

	return u.Parent
}

func liftNode(node *Node, distance int) *Node {
	if node == nil || distance < 0 {
		return nil
	}

	for i := 0; distance > 0 && node != nil; i++ {
		if distance&1 == 1 {
			if i >= len(node.Meta.Up) {
				return nil
			}
			node = node.Meta.Up[i]
		}
		distance >>= 1
	}

	return node
}

func computeLCALog(maxNodes int) int {
	if maxNodes <= 1 {
		return defaultLCALog
	}

	log := 1
	for (1 << log) <= maxNodes {
		log++
	}

	if log < defaultLCALog {
		return defaultLCALog
	}

	return log
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
