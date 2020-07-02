package golden

import (
	//"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type header struct {
	Key   string
	Value string
}

func performRequest(r http.Handler, method, path string, headers ...header) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, nil)
	for _, h := range headers {
		req.Header.Add(h.Key, h.Value)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestRouteNotOK(t *testing.T) {
	passed := false
	g := New(Address(":9090"), Timeout(1, 1))
	g.Post("/test1", func(c *HttpContext) {
		passed = true
	})
	g.Register404Tip("whoru")

	w := performRequest(g, "GET", "/test1")

	assert.False(t, passed)
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "whoru", w.Body.String())
}

func TestMiddlewareNoRoute(t *testing.T) {
	signature := ""
	g := New(Address(":9090"), Timeout(1, 1))
	g.AllRequest(func(c *HttpContext) {
		signature += "A"
		c.Next()
		signature += "B"
	})
	g.AllRequest(func(c *HttpContext) {
		signature += "C"
		c.Next()
		c.Next()
		c.Next()
		c.Next()
		signature += "D"
	})
	g.NotFound(func(c *HttpContext) {
		signature += "E"
		c.Next()
		signature += "F"
	}, func(c *HttpContext) {
		signature += "G"
		c.Next()
		signature += "H"
	})
	// RUN
	w := performRequest(g, "GET", "/")

	// TEST
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "ACEGHFDB", signature)
}

func TestGet(t *testing.T) {
	signature := ""
	g := New(Address(":9090"), Timeout(1, 1))
	g.Get("/", func(c *HttpContext) {
		signature += "A"
		c.Next()
		signature += "B"
	}, func(c *HttpContext) {
		signature += "C"
		c.Next()
		c.Next()
		c.Next()
		c.Next()
		signature += "D"
	})
	g.NotFound(func(c *HttpContext) {
		signature += "E"
		c.Next()
		signature += "F"
	}, func(c *HttpContext) {
		signature += "G"
		c.Next()
		signature += "H"
	})
	// RUN
	w := performRequest(g, "GET", "/")

	// TEST
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "ACDB", signature)
}

func TestPost(t *testing.T) {
	signature := ""
	g := New(Address(":9090"), Timeout(1, 1))
	g.Post("/", func(c *HttpContext) {
		signature += "A"
		c.Next()
		signature += "B"
	}, func(c *HttpContext) {
		signature += "C"
		c.Next()
		c.Next()
		c.Next()
		c.Next()
		signature += "D"
	})
	g.NotFound(func(c *HttpContext) {
		signature += "E"
		c.Next()
		signature += "F"
	}, func(c *HttpContext) {
		signature += "G"
		c.Next()
		signature += "H"
	})
	// RUN
	w := performRequest(g, "POST", "/")

	// TEST
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "ACDB", signature)
}

func TestPost404(t *testing.T) {
	signature := ""
	g := New(Address(":9090"), Timeout(1, 1))
	g.Post("/", func(c *HttpContext) {
		signature += "A"
		c.Next()
		signature += "B"
	}, func(c *HttpContext) {
		signature += "C"
		c.Next()
		c.Next()
		c.Next()
		c.Next()
		signature += "D"
	})
	g.NotFound(func(c *HttpContext) {
		signature += "E"
		c.Next()
		signature += "F"
	}, func(c *HttpContext) {
		signature += "G"
		c.Next()
		signature += "H"
	})
	// RUN
	w := performRequest(g, "POST", "/x")

	// TEST
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "EGHF", signature)
}

func TestAnyMethod(t *testing.T) {
	signature := ""
	g := New(Address(":9090"), Timeout(1, 1))
	g.AnyMethod("/", func(c *HttpContext) {
		signature += "A"
		c.Next()
		signature += "B"
	}, func(c *HttpContext) {
		signature += "C"
		c.Next()
		c.Next()
		c.Next()
		c.Next()
		signature += "D"
	})
	g.NotFound(func(c *HttpContext) {
		signature += "E"
		c.Next()
		signature += "F"
	}, func(c *HttpContext) {
		signature += "G"
		c.Next()
		signature += "H"
	})
	// RUN
	w := performRequest(g, "PUT", "/")

	// TEST
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "ACDB", signature)
}

func TestDelete(t *testing.T) {
	signature := ""
	g := New(Address(":9090"), Timeout(1, 1))
	g.Delete("/", func(c *HttpContext) {
		signature += "A"
		c.Next()
		signature += "B"
	}, func(c *HttpContext) {
		signature += "C"
		c.Next()
		c.Next()
		c.Next()
		c.Next()
		signature += "D"
	})
	g.NotFound(func(c *HttpContext) {
		signature += "E"
		c.Next()
		signature += "F"
	}, func(c *HttpContext) {
		signature += "G"
		c.Next()
		signature += "H"
	})
	// RUN
	w := performRequest(g, "DELETE", "/")

	// TEST
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "ACDB", signature)
}

func TestPut(t *testing.T) {
	signature := ""
	g := New(Address(":9090"), Timeout(1, 1))
	g.Put("/", func(c *HttpContext) {
		signature += "A"
		c.Next()
		signature += "B"
	}, func(c *HttpContext) {
		signature += "C"
		c.Next()
		c.Next()
		c.Next()
		c.Next()
		signature += "D"
	})
	g.NotFound(func(c *HttpContext) {
		signature += "E"
		c.Next()
		signature += "F"
	}, func(c *HttpContext) {
		signature += "G"
		c.Next()
		signature += "H"
	})
	// RUN
	w := performRequest(g, "PUT", "/")

	// TEST
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "ACDB", signature)
}

func TestPatch(t *testing.T) {
	signature := ""
	g := New(Address(":9090"), Timeout(1, 1))
	g.Patch("/", func(c *HttpContext) {
		signature += "A"
		c.Next()
		signature += "B"
	}, func(c *HttpContext) {
		signature += "C"
		c.Next()
		c.Next()
		c.Next()
		c.Next()
		signature += "D"
	})
	g.NotFound(func(c *HttpContext) {
		signature += "E"
		c.Next()
		signature += "F"
	}, func(c *HttpContext) {
		signature += "G"
		c.Next()
		signature += "H"
	})
	// RUN
	w := performRequest(g, "PATCH", "/")

	// TEST
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "ACDB", signature)
}

type MyController struct {
}

func (c *MyController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	httpStatus := http.StatusOK
	w.WriteHeader(httpStatus)
}

func TestRegisterGet(t *testing.T) {
	{
		g := New(Address(":9090"), Timeout(1, 1))
		var c MyController
		g.RegisterGet("/aaa", &c)
		w := performRequest(g, "GET", "/aaa")
		assert.Equal(t, http.StatusOK, w.Code)
		w = performRequest(g, "GET", "/aaa/")
		assert.Equal(t, 301, w.Code)
	}
	{
		g := New(Address(":9090"), Timeout(1, 1))
		var c MyController
		g.RegisterGet("/aaa/", &c)
		w := performRequest(g, "GET", "/aaa")
		assert.Equal(t, 301, w.Code)
		w = performRequest(g, "GET", "/aaa/")
		assert.Equal(t, http.StatusOK, w.Code)
	}

}

func TestRegisterPost(t *testing.T) {
	{
		g := New(Address(":9090"), Timeout(1, 1))
		var c MyController
		g.RegisterPost("/aaa", &c)
		w := performRequest(g, "POST", "/aaa")
		assert.Equal(t, http.StatusOK, w.Code)
		w = performRequest(g, "POST", "/aaa/")
		assert.Equal(t, 307, w.Code)
	}
	{
		g := New(Address(":9090"), Timeout(1, 1))
		var c MyController
		g.RegisterPost("/aaa/", &c)
		w := performRequest(g, "POST", "/aaa")
		assert.Equal(t, 307, w.Code)
		w = performRequest(g, "POST", "/aaa/1")
		assert.Equal(t, http.StatusOK, w.Code)
		w = performRequest(g, "POST", "/aaa/")
		assert.Equal(t, http.StatusOK, w.Code)
	}

}

func TestRegisterPut(t *testing.T) {
	{
		g := New(Address(":9090"), Timeout(1, 1))
		var c MyController
		g.RegisterPut("/aaa", &c)
		w := performRequest(g, "PUT", "/aaa")
		assert.Equal(t, http.StatusOK, w.Code)
		w = performRequest(g, "PUT", "/aaa/")
		assert.Equal(t, 307, w.Code)
	}
	{
		g := New(Address(":9090"), Timeout(1, 1))
		var c MyController
		g.RegisterPut("/aaa/", &c)
		w := performRequest(g, "PUT", "/aaa")
		assert.Equal(t, 307, w.Code)
		w = performRequest(g, "PUT", "/aaa/1")
		assert.Equal(t, http.StatusOK, w.Code)
		w = performRequest(g, "PUT", "/aaa/")
		assert.Equal(t, http.StatusOK, w.Code)
	}

}

func TestRegisterPatch(t *testing.T) {
	{
		g := New(Address(":9090"), Timeout(1, 1))
		var c MyController
		g.RegisterPatch("/aaa", &c)
		w := performRequest(g, "PATCH", "/aaa")
		assert.Equal(t, http.StatusOK, w.Code)
		w = performRequest(g, "PATCH", "/aaa/")
		assert.Equal(t, 307, w.Code)
	}
	{
		g := New(Address(":9090"), Timeout(1, 1))
		var c MyController
		g.RegisterPatch("/aaa/", &c)
		w := performRequest(g, "PATCH", "/aaa")
		assert.Equal(t, 307, w.Code)
		w = performRequest(g, "PATCH", "/aaa/1")
		assert.Equal(t, http.StatusOK, w.Code)
		w = performRequest(g, "PATCH", "/aaa/")
		assert.Equal(t, http.StatusOK, w.Code)
	}

}

func TestRegisterDelete(t *testing.T) {
	{
		g := New(Address(":9090"), Timeout(1, 1))
		var c MyController
		g.RegisterDelete("/aaa", &c)
		w := performRequest(g, "DELETE", "/aaa")
		assert.Equal(t, http.StatusOK, w.Code)
		w = performRequest(g, "DELETE", "/aaa/")
		assert.Equal(t, 307, w.Code)
	}
	{
		g := New(Address(":9090"), Timeout(1, 1))
		var c MyController
		g.RegisterDelete("/aaa/", &c)
		w := performRequest(g, "DELETE", "/aaa")
		assert.Equal(t, 307, w.Code)
		w = performRequest(g, "DELETE", "/aaa/1")
		assert.Equal(t, http.StatusOK, w.Code)
		w = performRequest(g, "DELETE", "/aaa/")
		assert.Equal(t, http.StatusOK, w.Code)
	}

}

func TestRegisterAny(t *testing.T) {
	{
		g := New(Address(":9090"), Timeout(1, 1))
		var c MyController
		g.RegisterHandler("/aaa", &c)
		w := performRequest(g, "DELETE", "/aaa")
		assert.Equal(t, http.StatusOK, w.Code)
		w = performRequest(g, "DELETE", "/aaa/")
		assert.Equal(t, 307, w.Code)
	}
	{
		g := New(Address(":9090"), Timeout(1, 1))
		var c MyController
		g.RegisterHandler("/aaa/", &c)
		w := performRequest(g, "DELETE", "/aaa")
		assert.Equal(t, 307, w.Code)
		w = performRequest(g, "DELETE", "/aaa/1")
		assert.Equal(t, http.StatusOK, w.Code)
		w = performRequest(g, "DELETE", "/aaa/")
		assert.Equal(t, http.StatusOK, w.Code)
	}
	//
	{
		g := New(Address(":9090"), Timeout(1, 1))
		var c MyController
		g.RegisterHandler("/aaa", &c)
		w := performRequest(g, "POST", "/aaa")
		assert.Equal(t, http.StatusOK, w.Code)
		w = performRequest(g, "POST", "/aaa/")
		assert.Equal(t, 307, w.Code)
	}
	{
		g := New(Address(":9090"), Timeout(1, 1))
		var c MyController
		g.RegisterHandler("/aaa/", &c)
		w := performRequest(g, "POST", "/aaa")
		assert.Equal(t, 307, w.Code)
		w = performRequest(g, "POST", "/aaa/1")
		assert.Equal(t, http.StatusOK, w.Code)
		w = performRequest(g, "POST", "/aaa/")
		assert.Equal(t, http.StatusOK, w.Code)
	}
	//
	{
		g := New(Address(":9090"), Timeout(1, 1))
		var c MyController
		g.RegisterHandler("/aaa", &c)
		w := performRequest(g, "GET", "/aaa")
		assert.Equal(t, http.StatusOK, w.Code)
		w = performRequest(g, "GET", "/aaa/")
		assert.Equal(t, 301, w.Code)
	}
	{
		g := New(Address(":9090"), Timeout(1, 1))
		var c MyController
		g.RegisterHandler("/aaa/", &c)
		w := performRequest(g, "GET", "/aaa")
		assert.Equal(t, 301, w.Code)
		w = performRequest(g, "GET", "/aaa/1")
		assert.Equal(t, http.StatusOK, w.Code)
		w = performRequest(g, "GET", "/aaa/")
		assert.Equal(t, http.StatusOK, w.Code)
	}
	//
	{
		g := New(Address(":9090"), Timeout(1, 1))
		var c MyController
		g.RegisterHandler("/aaa", &c)
		w := performRequest(g, "PUT", "/aaa")
		assert.Equal(t, http.StatusOK, w.Code)
		w = performRequest(g, "PUT", "/aaa/")
		assert.Equal(t, 307, w.Code)
	}
	{
		g := New(Address(":9090"), Timeout(1, 1))
		var c MyController
		g.RegisterHandler("/aaa/", &c)
		w := performRequest(g, "PUT", "/aaa")
		assert.Equal(t, 307, w.Code)
		w = performRequest(g, "PUT", "/aaa/")
		assert.Equal(t, http.StatusOK, w.Code)
		w = performRequest(g, "PUT", "/aaa/1")
		assert.Equal(t, http.StatusOK, w.Code)
	}
	//
	{
		g := New(Address(":9090"), Timeout(1, 1))
		var c MyController
		g.RegisterHandler("/aaa", &c)
		w := performRequest(g, "PATCH", "/aaa")
		assert.Equal(t, http.StatusOK, w.Code)
		w = performRequest(g, "PATCH", "/aaa/")
		assert.Equal(t, 307, w.Code)
	}
	{
		g := New(Address(":9090"), Timeout(1, 1))
		var c MyController
		g.RegisterHandler("/aaa/", &c)
		w := performRequest(g, "PATCH", "/aaa")
		assert.Equal(t, 307, w.Code)
		w = performRequest(g, "PATCH", "/aaa/")
		assert.Equal(t, http.StatusOK, w.Code)
		w = performRequest(g, "PATCH", "/aaa/1")
		assert.Equal(t, http.StatusOK, w.Code)
	}
}

func TestRegisterHandlerFunc(t *testing.T) {
	{
		g := New(Address(":9090"), Timeout(1, 1))
		var c MyController
		g.RegisterHandleFunc("/aaa", c.ServeHTTP)
		w := performRequest(g, "DELETE", "/aaa")
		assert.Equal(t, http.StatusOK, w.Code)
		w = performRequest(g, "DELETE", "/aaa/")
		assert.Equal(t, 307, w.Code)
	}
	{
		g := New(Address(":9090"), Timeout(1, 1))
		var c MyController
		g.RegisterHandleFunc("/aaa/", c.ServeHTTP)
		w := performRequest(g, "DELETE", "/aaa")
		assert.Equal(t, 307, w.Code)
		w = performRequest(g, "DELETE", "/aaa/1")
		assert.Equal(t, http.StatusOK, w.Code)
		w = performRequest(g, "DELETE", "/aaa/")
		assert.Equal(t, http.StatusOK, w.Code)
	}
	//
	{
		g := New(Address(":9090"), Timeout(1, 1))
		var c MyController
		g.RegisterHandleFunc("/aaa", c.ServeHTTP)
		w := performRequest(g, "POST", "/aaa")
		assert.Equal(t, http.StatusOK, w.Code)
		w = performRequest(g, "POST", "/aaa/")
		assert.Equal(t, 307, w.Code)
	}
	{
		g := New(Address(":9090"), Timeout(1, 1))
		var c MyController
		g.RegisterHandleFunc("/aaa/", c.ServeHTTP)
		w := performRequest(g, "POST", "/aaa")
		assert.Equal(t, 307, w.Code)
		w = performRequest(g, "POST", "/aaa/1")
		assert.Equal(t, http.StatusOK, w.Code)
		w = performRequest(g, "POST", "/aaa/")
		assert.Equal(t, http.StatusOK, w.Code)
	}
	//
	{
		g := New(Address(":9090"), Timeout(1, 1))
		var c MyController
		g.RegisterHandleFunc("/aaa", c.ServeHTTP)
		w := performRequest(g, "GET", "/aaa")
		assert.Equal(t, http.StatusOK, w.Code)
		w = performRequest(g, "GET", "/aaa/")
		assert.Equal(t, 301, w.Code)
	}
	{
		g := New(Address(":9090"), Timeout(1, 1))
		var c MyController
		g.RegisterHandleFunc("/aaa/", c.ServeHTTP)
		w := performRequest(g, "GET", "/aaa")
		assert.Equal(t, 301, w.Code)
		w = performRequest(g, "GET", "/aaa/")
		assert.Equal(t, http.StatusOK, w.Code)
		w = performRequest(g, "GET", "/aaa/1")
		assert.Equal(t, http.StatusOK, w.Code)
	}
	//
	{
		g := New(Address(":9090"), Timeout(1, 1))
		var c MyController
		g.RegisterHandleFunc("/aaa", c.ServeHTTP)
		w := performRequest(g, "PUT", "/aaa")
		assert.Equal(t, http.StatusOK, w.Code)
		w = performRequest(g, "PUT", "/aaa/")
		assert.Equal(t, 307, w.Code)
	}
	{
		g := New(Address(":9090"), Timeout(1, 1))
		var c MyController
		g.RegisterHandleFunc("/aaa/", c.ServeHTTP)
		w := performRequest(g, "PUT", "/aaa")
		assert.Equal(t, 307, w.Code)
		w = performRequest(g, "PUT", "/aaa/")
		assert.Equal(t, http.StatusOK, w.Code)
		w = performRequest(g, "PUT", "/aaa/1")
		assert.Equal(t, http.StatusOK, w.Code)
	}
	//
	{
		g := New(Address(":9090"), Timeout(1, 1))
		var c MyController
		g.RegisterHandleFunc("/aaa", c.ServeHTTP)
		w := performRequest(g, "PATCH", "/aaa")
		assert.Equal(t, http.StatusOK, w.Code)
		w = performRequest(g, "PATCH", "/aaa/")
		assert.Equal(t, 307, w.Code)
	}
	{
		g := New(Address(":9090"), Timeout(1, 1))
		var c MyController
		g.RegisterHandleFunc("/aaa/", c.ServeHTTP)
		w := performRequest(g, "PATCH", "/aaa")
		assert.Equal(t, 307, w.Code)
		w = performRequest(g, "PATCH", "/aaa/")
		assert.Equal(t, http.StatusOK, w.Code)
		w = performRequest(g, "PATCH", "/aaa/1")
		assert.Equal(t, http.StatusOK, w.Code)
	}
}

func TestRegisterExactPathHandler(t *testing.T) {
	{
		g := New(Address(":9090"), Timeout(1, 1))
		var c MyController
		g.RegisterExactPathHandler("/aaa", &c)
		w := performRequest(g, "DELETE", "/aaa")
		assert.Equal(t, http.StatusOK, w.Code)
		w = performRequest(g, "DELETE", "/aaa/")
		assert.Equal(t, 307, w.Code)
	}
	{
		g := New(Address(":9090"), Timeout(1, 1))
		var c MyController
		g.RegisterExactPathHandler("/aaa/", &c)
		w := performRequest(g, "PUT", "/aaa")
		assert.Equal(t, 307, w.Code)
		w = performRequest(g, "PUT", "/aaa/")
		assert.Equal(t, http.StatusOK, w.Code)
		w = performRequest(g, "PUT", "/aaa/1")
		assert.Equal(t, 404, w.Code)
	}

}

func TestRegisterExactPathHandleFunc(t *testing.T) {
	{
		g := New(Address(":9090"), Timeout(1, 1))
		var c MyController
		g.RegisterExactPathHandleFunc("/aaa", c.ServeHTTP)
		w := performRequest(g, "POST", "/aaa")
		assert.Equal(t, http.StatusOK, w.Code)
		w = performRequest(g, "POST", "/aaa/")
		assert.Equal(t, 307, w.Code)
	}
	{
		g := New(Address(":9090"), Timeout(1, 1))
		var c MyController
		g.RegisterExactPathHandleFunc("/aaa/", c.ServeHTTP)
		w := performRequest(g, "PATCH", "/aaa")
		assert.Equal(t, 307, w.Code)
		w = performRequest(g, "PATCH", "/aaa/")
		assert.Equal(t, http.StatusOK, w.Code)
		w = performRequest(g, "PATCH", "/aaa/1")
		assert.Equal(t, 404, w.Code)
	}

}

func TestUrlParam(t *testing.T) {
	signature := ""
	g := New(Address(":9090"), Timeout(1, 1))

	g.Get("/abc:def", func(c *HttpContext) {
		signature += "A"
		c.Next()
		signature += "B"
	}, func(c *HttpContext) {
		signature += "C"
		c.Next()
		c.Next()
		c.Next()
		c.Next()
		signature += "D"
	})
	g.NotFound(func(c *HttpContext) {
		signature += "E"
		c.Next()
		signature += "F"
	}, func(c *HttpContext) {
		signature += "G"
		c.Next()
		signature += "H"
	})
	// RUN
	w := performRequest(g, "GET", "/abcd")

	// TEST
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "ACDB", signature)

	w = performRequest(g, "GET", "/abc")

	// TEST
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "ACDBEGHF", signature)
}
