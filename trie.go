package zwf

import (
	"fmt"
	"strings"
)

// 使用前缀树记录，实现动态路由
type node struct {
	pattern  string  // 路由 仅子节点有值 /v1/part/:lang
	part     string  // 路由存储在节点的部分 例如 :lang
	children []*node // 子节点
	isWild   bool    // 是否模糊匹配
}

func (n *node) String() string {
	return fmt.Sprintf("node{pattern=%s, part=%s, isWild=%t}", n.pattern, n.part, n.isWild)
}

// matchChild 匹配第一个节点 用于添加
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// matchChildren 匹配所有节点 用于查找
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

//func (n *node) insert(pattern string, parts []string, height int) {
//	if len(parts) == height {
//		n.pattern = pattern
//		return
//	}
//
//	part := parts[height]
//	child := n.matchChild(part)
//	if child == nil {
//		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
//		n.children = append(n.children, child)
//	}
//
//	child.insert(pattern, parts, height+1)
//}

// insert 添加路由到前缀树 逐层匹配 如果不存在则添加节点 直到最终节点 添加pattern到这个节点 pattern不为空字符串也用于判断是否在最终节点
// todo fix 路由冲突问题 /a/:b /a/c 冲突
func (n *node) insert(pattern string, parts []string, height int) {
	currentNode := n
	for height < len(parts) {
		part := parts[height]
		child := currentNode.matchChild(part)
		if child == nil {
			child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
			currentNode.children = append(currentNode.children, child)
		}

		currentNode = child
		height++
	}

	// 只有最终节点有pattern
	currentNode.pattern = pattern
}

//func (n *node) search(parts []string, height int) *node {
//	if len(parts) == height || strings.HasPrefix(n.part, "*") {
//		if n.pattern == "" {
//			return nil
//		}
//		return n
//	}
//
//	part := parts[height]
//	for _, child := range n.matchChildren(part) {
//		result := child.search(parts, height+1)
//		if result != nil {
//			return result
//		}
//	}
//
//	return nil
//}

// search 查找节点 匹配传入路由到已注册到前缀树 要注意节点模糊匹配到情况
func (n *node) search(parts []string, height int) *node {
	currentNode := n
	for height < len(parts) {
		if currentNode.pattern != "" || strings.HasPrefix(currentNode.part, "*") {
			return currentNode
		}

		part := parts[height]
		match := false
		for _, child := range currentNode.matchChildren(part) {
			currentNode = child
			match = true
			break
		}

		if !match {
			return nil
		}
		height++
	}

	if currentNode.pattern != "" {
		return currentNode
	}
	return nil
}

//func (n *node) travel(list *[]*node) {
//	if n.pattern != "" {
//		*list = append(*list, n)
//	}
//	for _, child := range n.children {
//		child.travel(list)
//	}
//	return
//}

// travel 获取所有最终节点 DFS
func (n *node) travel(list *[]*node) {
	stack := []*node{n} // 初始化栈，将根节点推入栈中

	for len(stack) > 0 {
		node := stack[len(stack)-1]  // 获取栈顶元素
		stack = stack[:len(stack)-1] // 弹出栈顶元素

		if node.pattern != "" {
			*list = append(*list, node) // 如果节点的 pattern 不为空，将其添加到列表中
		}

		for _, child := range node.children {
			stack = append(stack, child) // 将子节点推入栈中
		}
	}
}
