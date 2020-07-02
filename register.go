package golden

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type HandlerHttpFunc func(*HttpContext)

func (golden *Golden) Get(path string, handlers ...HandlerHttpFunc) {
	size := len(handlers)
	handlerFuncs := make([]gin.HandlerFunc, size)
	for i := 0; i < size; i++ {
		index := i
		handlerFuncs[index] = func(c *gin.Context) {
			handlers[index](&HttpContext{
				context: c,
				options: golden.options,
				client:  golden.client,
			})
			// 不允许隐式的遍历后续的handler
			c.Abort()
		}
	}
	golden.router.GET(path, handlerFuncs...)
}

func (golden *Golden) Post(path string, handlers ...HandlerHttpFunc) {
	size := len(handlers)
	handlerFuncs := make([]gin.HandlerFunc, size)
	for i := 0; i < size; i++ {
		index := i
		handlerFuncs[index] = func(c *gin.Context) {
			handlers[index](&HttpContext{
				context: c,
				options: golden.options,
				client:  golden.client,
			})
			// 不允许隐式的遍历后续的handler
			c.Abort()
		}
	}
	golden.router.POST(path, handlerFuncs...)
}

func (golden *Golden) Delete(path string, handlers ...HandlerHttpFunc) {
	size := len(handlers)
	handlerFuncs := make([]gin.HandlerFunc, size)
	for i := 0; i < size; i++ {
		index := i
		handlerFuncs[index] = func(c *gin.Context) {
			handlers[index](&HttpContext{
				context: c,
				options: golden.options,
				client:  golden.client,
			})
			// 不允许隐式的遍历后续的handler
			c.Abort()
		}
	}
	golden.router.DELETE(path, handlerFuncs...)
}

func (golden *Golden) Patch(path string, handlers ...HandlerHttpFunc) {
	size := len(handlers)
	handlerFuncs := make([]gin.HandlerFunc, size)
	for i := 0; i < size; i++ {
		index := i
		handlerFuncs[index] = func(c *gin.Context) {
			handlers[index](&HttpContext{
				context: c,
				options: golden.options,
				client:  golden.client,
			})
			// 不允许隐式的遍历后续的handler
			c.Abort()
		}
	}
	golden.router.PATCH(path, handlerFuncs...)
}

func (golden *Golden) Put(path string, handlers ...HandlerHttpFunc) {
	size := len(handlers)
	handlerFuncs := make([]gin.HandlerFunc, size)
	for i := 0; i < size; i++ {
		index := i
		handlerFuncs[index] = func(c *gin.Context) {
			handlers[index](&HttpContext{
				context: c,
				options: golden.options,
				client:  golden.client,
			})
			// 不允许隐式的遍历后续的handler
			c.Abort()
		}
	}
	golden.router.PUT(path, handlerFuncs...)
}

// GET/PUT/DELETE/POST/PATCH都能处理
func (golden *Golden) AnyMethod(path string, handlers ...HandlerHttpFunc) {
	size := len(handlers)
	handlerFuncs := make([]gin.HandlerFunc, size)
	for i := 0; i < size; i++ {
		index := i
		handlerFuncs[index] = func(c *gin.Context) {
			handlers[index](&HttpContext{
				context: c,
				options: golden.options,
				client:  golden.client,
			})
			// 不允许隐式的遍历后续的handler
			c.Abort()
		}
	}
	golden.router.Any(path, handlerFuncs...)
}

func (golden *Golden) AllRequest(handlers ...HandlerHttpFunc) {
	size := len(handlers)
	handlerFuncs := make([]gin.HandlerFunc, size)
	for i := 0; i < size; i++ {
		index := i
		handlerFuncs[index] = func(c *gin.Context) {
			handlers[index](&HttpContext{
				context: c,
				options: golden.options,
				client:  golden.client,
			})
			// 不允许隐式的遍历后续的handler
			c.Abort()
		}
	}
	golden.router.Use(handlerFuncs...)
}

func (golden *Golden) NotFound(handlers ...HandlerHttpFunc) {
	size := len(handlers)
	handlerFuncs := make([]gin.HandlerFunc, size)
	for i := 0; i < size; i++ {
		index := i
		handlerFuncs[index] = func(c *gin.Context) {
			handlers[index](&HttpContext{
				context: c,
				options: golden.options,
				client:  golden.client,
			})
			// 不允许隐式的遍历后续的handler
			c.Abort()
		}
	}
	golden.router.NoRoute(handlerFuncs...)
}

/**

RegisterXX 是一套支持原生http回调函数的接口，用于轻松兼容代码中的老接口。

**/
func (golden *Golden) RegisterGet(path string, handler http.Handler) {
	wrapper := func(c *gin.Context) {
		request := c.Request
		response := c.Writer
		handler.ServeHTTP(response, request)
	}
	if strings.HasSuffix(path, "/") {
		golden.router.Group(path).GET("/*subpath", wrapper)
	} else {
		golden.router.GET(path, wrapper)
	}
}

func (golden *Golden) RegisterPost(path string, handler http.Handler) {
	wrapper := func(c *gin.Context) {
		request := c.Request
		response := c.Writer
		handler.ServeHTTP(response, request)
	}
	if strings.HasSuffix(path, "/") {
		golden.router.Group(path).POST("/*subpath", wrapper)
	} else {
		golden.router.POST(path, wrapper)
	}
}

func (golden *Golden) RegisterDelete(path string, handler http.Handler) {
	wrapper := func(c *gin.Context) {
		request := c.Request
		response := c.Writer
		handler.ServeHTTP(response, request)
	}
	if strings.HasSuffix(path, "/") {
		golden.router.Group(path).DELETE("/*subpath", wrapper)
	} else {
		golden.router.DELETE(path, wrapper)
	}
}

func (golden *Golden) RegisterPatch(path string, handler http.Handler) {
	wrapper := func(c *gin.Context) {
		request := c.Request
		response := c.Writer
		handler.ServeHTTP(response, request)
	}
	if strings.HasSuffix(path, "/") {
		golden.router.Group(path).PATCH("/*subpath", wrapper)
	} else {
		golden.router.PATCH(path, wrapper)
	}
}

func (golden *Golden) RegisterPut(path string, handler http.Handler) {
	wrapper := func(c *gin.Context) {
		request := c.Request
		response := c.Writer
		handler.ServeHTTP(response, request)
	}
	if strings.HasSuffix(path, "/") {
		golden.router.Group(path).PUT("/*subpath", wrapper)
	} else {
		golden.router.PUT(path, wrapper)
	}
}

// 指定精确路径，不区别对待包含/结尾的
func (golden *Golden) RegisterExactPathHandler(path string, handler http.Handler) {
	wrapper := func(c *gin.Context) {
		request := c.Request
		response := c.Writer
		handler.ServeHTTP(response, request)
	}
	golden.router.Any(path, wrapper)
}

// 指定精确路径，不区别对待包含/结尾的
func (golden *Golden) RegisterExactPathHandleFunc(path string, handleFunc func(http.ResponseWriter, *http.Request)) {
	wrapper := func(c *gin.Context) {
		request := c.Request
		response := c.Writer
		handleFunc(response, request)
	}
	golden.router.Any(path, wrapper)
}

func (golden *Golden) RegisterHandler(path string, handler http.Handler) {
	wrapper := func(c *gin.Context) {
		request := c.Request
		response := c.Writer
		handler.ServeHTTP(response, request)
	}
	if strings.HasSuffix(path, "/") {
		golden.router.Group(path).Any("/*subpath", wrapper)
	} else {
		golden.router.Any(path, wrapper)
	}
}

func (golden *Golden) RegisterHandleFunc(path string, handleFunc func(http.ResponseWriter, *http.Request)) {
	wrapper := func(c *gin.Context) {
		request := c.Request
		response := c.Writer
		handleFunc(response, request)
	}
	if strings.HasSuffix(path, "/") {
		golden.router.Group(path).Any("/*subpath", wrapper)
	} else {
		golden.router.Any(path, wrapper)
	}
}

func (golden *Golden) Register404Tip(tip string) {
	golden.router.NoRoute(func(c *gin.Context) {
		c.String(404, tip)
	})
}
