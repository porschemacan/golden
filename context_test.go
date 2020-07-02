package golden

import (
	//"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGetSetHttpContext(t *testing.T) {
	c := HttpContext{
		context: &gin.Context{},
	}
	tmp := "1"
	c.Set("A", tmp)
	assert.Equal(t, c.GetString("A"), "1")
	c.Set("A'", &tmp)
	assert.Equal(t, c.GetString("A'"), "")

	c.Set("B", true)
	assert.Equal(t, c.GetBool("B"), true)

	c.Set("C", 3)
	assert.Equal(t, c.GetInt("C"), 3)

	c.Set("D", int64(4))
	assert.Equal(t, c.GetInt64("D"), int64(4))

	c.Set("E", float64(5))
	assert.Equal(t, c.GetFloat64("E"), float64(5))

	tt := time.Now()
	c.Set("F", tt)
	assert.Equal(t, c.GetTime("F"), tt)

	c.Set("G", time.Duration(3)*time.Second)
	assert.Equal(t, c.GetDuration("G"), time.Duration(3)*time.Second)

	c.Set("H", []string{"1", "2"})
	assert.Equal(t, c.GetStringSlice("H"), []string{"1", "2"})

	c.Set("I", map[string]interface{}{"a": "b"})
	assert.Equal(t, c.GetStringMap("I"), map[string]interface{}{"a": "b"})

	c.Set("J", map[string]string{"a": "b"})
	assert.Equal(t, c.GetStringMapString("J"), map[string]string{"a": "b"})

	c.Set("K", map[string][]string{"a": {"b"}})
	assert.Equal(t, c.GetStringMapStringSlice("K"), map[string][]string{"a": {"b"}})
}
