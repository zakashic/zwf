package zwf

import (
	"net/http"
	"strings"
)

type router struct {
	roots    map[string]*node       // method为key
	handlers map[string]HandlerFunc // method-pattern为key
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// parsePattern 解析pattern
func parsePattern(pattern string) []string {
	parts := make([]string, 0)
	for _, part := range strings.Split(pattern, "/") {
		if part != "" {
			parts = append(parts, part)
			// 如果是模糊匹配则停止 后续part无意义
			if part[0] == '*' {
				break
			}
		}
	}

	return parts
}

// addRoute 注册路由 添加到路由树节点和记录pattern对应的handler
func (r *router) addRoute(method, pattern string, handler HandlerFunc) {
	if _, ok := r.roots[method]; !ok {
		r.roots[method] = &node{}
	}

	parts := parsePattern(pattern)
	r.roots[method].insert(pattern, parts, 0)

	key := method + "-" + pattern
	r.handlers[key] = handler
}

// getRoute 获取路由 返回节点和参数
func (r *router) getRoute(method, pattern string) (*node, map[string]string) {
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}

	params := make(map[string]string)
	searchParts := parsePattern(pattern)
	resultNode := root.search(searchParts, 0)
	if resultNode != nil {
		parts := parsePattern(resultNode.pattern)
		for index, part := range parts {
			// 处理模糊匹配的情况
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return resultNode, params
	}

	return nil, nil
}

// getRoutes 获取所有路由
func (r *router) getRoutes(method string) []*node {
	root, ok := r.roots[method]
	if !ok {
		return nil
	}

	nodes := make([]*node, 0)
	root.travel(&nodes)
	return nodes
}

// handle 处理请求
func (r *router) handle(c *Context) {
	node, params := r.getRoute(c.Method, c.Path)

	if node != nil {
		key := c.Method + "-" + node.pattern
		c.Params = params
		// 获取到对应的handler 添加到中间件链
		c.handlers = append(c.handlers, r.handlers[key])
	} else {
		c.handlers = append(c.handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		})
	}
	c.Next()
}
