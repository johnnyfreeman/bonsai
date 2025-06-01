package viewer

import "fmt"

// BuildTree creates a tree structure from JSON data
func BuildTree(data interface{}, key, path string) *Node {
	node := &Node{
		Key:  key,
		Path: path,
	}

	switch v := data.(type) {
	case map[string]interface{}:
		node.Type = ObjectNode
		node.Value = v
		for k, val := range v {
			child := BuildTree(val, k, path+"."+k)
			child.Parent = node
			node.Children = append(node.Children, child)
		}
	case []interface{}:
		node.Type = ArrayNode
		node.Value = v
		for i, val := range v {
			child := BuildTree(val, fmt.Sprintf("[%d]", i), fmt.Sprintf("%s[%d]", path, i))
			child.Parent = node
			node.Children = append(node.Children, child)
		}
	case string:
		node.Type = StringNode
		node.Value = v
	case float64:
		node.Type = NumberNode
		node.Value = v
	case bool:
		node.Type = BoolNode
		node.Value = v
	case nil:
		node.Type = NullNode
		node.Value = nil
	}

	return node
}

// CountNodes recursively counts the total number of nodes in the tree
func CountNodes(node *Node) int {
	count := 1
	for _, child := range node.Children {
		count += CountNodes(child)
	}
	return count
}

// GetDepth returns the depth of a node in the tree
func (n *Node) GetDepth() int {
	depth := 0
	current := n.Parent
	for current != nil {
		depth++
		current = current.Parent
	}
	return depth
}

// IsExpanded returns true if the node is expanded
func (n *Node) IsExpanded() bool {
	return n.Expanded
}

// HasChildren returns true if the node has children
func (n *Node) HasChildren() bool {
	return len(n.Children) > 0
}

// ToggleExpansion toggles the expansion state of the node
func (n *Node) ToggleExpansion() {
	n.Expanded = !n.Expanded
}

// Expand expands the node
func (n *Node) Expand() {
	n.Expanded = true
}

// Collapse collapses the node
func (n *Node) Collapse() {
	n.Expanded = false
}

// ExpandAll recursively expands all nodes in the subtree
func (n *Node) ExpandAll() {
	n.Expanded = true
	for _, child := range n.Children {
		child.ExpandAll()
	}
}

// CollapseAll recursively collapses all nodes in the subtree
func (n *Node) CollapseAll() {
	n.Expanded = false
	for _, child := range n.Children {
		child.CollapseAll()
	}
}

// FindPath finds a node by its path
func (n *Node) FindPath(path string) *Node {
	if n.Path == path {
		return n
	}
	
	for _, child := range n.Children {
		if found := child.FindPath(path); found != nil {
			return found
		}
	}
	
	return nil
}

// GetParentChain returns the chain of parent nodes up to the root
func (n *Node) GetParentChain() []*Node {
	var chain []*Node
	current := n.Parent
	for current != nil {
		chain = append([]*Node{current}, chain...)
		current = current.Parent
	}
	return chain
}