package libs

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestDLogf_1(t *testing.T) {
	InitLog(&LogConfig{
		Path:  "/tmp/",
		Level: 7,
	})
	DLogf("AAAA=%d=", 5)
	dat, err := ioutil.ReadFile("/tmp/didi_log.libs.test.log")
	assert.Nil(t, err)
	assert.Contains(t, string(dat), "AAAA=5=")

	DTagf("GGG", "AAAA=%d=", 6)
	dat, err = ioutil.ReadFile("/tmp/didi_log.libs.test.log")
	assert.Nil(t, err)
	assert.Contains(t, string(dat), "GGG")
	assert.Contains(t, string(dat), "AAAA=6=")
}
