package zwf

import (
	"net/http"
	"testing"
)

func TestDefault(t *testing.T) {
	r := Default()

	r.GET("/", func(c *Context) {
		c.String(http.StatusOK, "hello world\n")
	})

	// index out of range for testing Recovery()
	r.GET("/panic", func(c *Context) {
		names := []string{"hello world"}
		c.String(http.StatusOK, names[100])
	})

	v1 := r.Group("/v1")
	{
		v1.GET("/hello/:name", func(c *Context) {
			c.String(http.StatusOK, "hello %s", c.Param("name"))
		})

		v1.GET("/hello", func(c *Context) {
			c.String(http.StatusOK, "hello")
		})
	}

	v2 := r.Group("/v2")
	{
		v2.GET("/a/c", func(c *Context) {
			c.String(http.StatusOK, "hello")
		})

		v2.GET("/a/:b", func(c *Context) {
			c.String(http.StatusOK, "hello %s", c.Param("name"))
		})
	}

	r.Run(":9999")
}
