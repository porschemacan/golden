package libs

import (
	"bou.ke/monkey"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	//"time"
)

type conn struct {
}

func (this *conn) Close() error {
	return nil
}

func (this *conn) Err() error {
	return nil
}

func (this *conn) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	return nil, nil
}

func (this *conn) Send(commandName string, args ...interface{}) error {
	return nil
}

func (this *conn) Flush() error {
	return nil
}

func (this *conn) Receive() (reply interface{}, err error) {
	return nil, nil
}

/*
func TestDial(t *testing.T) {
	client, err := NewRedisClient(&RedisConf{
		ConnConf: ConnConf{
			ConnectTimeout: 3000,
			ReadTimeout:    1000,
			WriteTimeout:   2000,
			RetryTimes:     4000,
		},
		Host: []string{
			"1.1.1.2",
			"2.2.2.3",
		},
		InitConnSize: 1000,
		MaxConnSize:  1000,
		MaxIdleSec:   1000,
		ExpireSec:    1000,
		Password:     "xxx",
	})
	assert.Nil(t, err)

	guard1 := monkey.Patch(redis.Dial, func(network, address string, options ...redis.DialOption) (redis.Conn, error) {

		assert.True(t, network == "tcp")
		assert.True(t, address == "1.1.1.2" || address == "2.2.2.3")

		if address == "2.2.2.3" {
			return &conn{}, nil
		}

		return nil, fmt.Errorf("-")
	})
	defer guard1.Unpatch()

	cnt := 0
	guard2 := monkey.Patch(redis.DialPassword, func(password string) redis.DialOption {
		assert.Equal(t, password, "xxx")
		cnt = cnt + 1
		return redis.DialOption{}
	})
	defer guard2.Unpatch()

	guard3 := monkey.Patch(redis.DialReadTimeout, func(d time.Duration) redis.DialOption {
		assert.Equal(t, int64(d), int64(1000*1000000))
		cnt = cnt + 2
		return redis.DialOption{}
	})
	defer guard3.Unpatch()

	guard4 := monkey.Patch(redis.DialWriteTimeout, func(d time.Duration) redis.DialOption {
		assert.Equal(t, int64(d), int64(2000*1000000))
		cnt = cnt + 4
		return redis.DialOption{}
	})
	defer guard4.Unpatch()

	guard5 := monkey.Patch(redis.DialConnectTimeout, func(d time.Duration) redis.DialOption {
		assert.Equal(t, int64(d), int64(3000*1000000))
		cnt = cnt + 8
		return redis.DialOption{}
	})
	defer guard5.Unpatch()

	client.redisPool.Dial()

}

func TestExist(t *testing.T) {
	client, err := NewRedisClient(&RedisConf{
		ConnConf: ConnConf{
			ConnectTimeout: 1000,
			ReadTimeout:    1000,
			WriteTimeout:   2000,
			RetryTimes:     4000,
		},
		Host: []string{
			"1.1.1.2",
			"2.2.2.3",
		},
		InitConnSize: 1000,
		MaxConnSize:  1000,
		MaxIdleSec:   1000,
		ExpireSec:    1000,
		Password:     "xxx",
	})
	assert.Nil(t, err)

	c := &conn{}
	guard1 := monkey.Patch(redis.Dial, func(network, address string, options ...redis.DialOption) (redis.Conn, error) {

		assert.True(t, network == "tcp")
		assert.True(t, address == "1.1.1.2" || address == "2.2.2.3")

		if address == "2.2.2.3" {
			return c, nil
		}

		return nil, fmt.Errorf("-")
	})
	defer guard1.Unpatch()

	guard2 := monkey.PatchInstanceMethod(reflect.TypeOf(c), "Do", func(_ *conn, commandName string, args ...interface{}) (reply interface{}, err error) {

		assert.Equal(t, commandName, "EXISTS")
		assert.Equal(t, string(args[0].(string)), "ABC")
		return nil, nil
	})
	defer guard2.Unpatch()

	ret, errExist := client.Exist("ABC")
	assert.Equal(t, ret, 0)
	assert.Nil(t, errExist)

}

func TestHMSetWithExpire(t *testing.T) {
	client, err := NewRedisClient(&RedisConf{
		ConnConf: ConnConf{
			ConnectTimeout: 1000,
			ReadTimeout:    1000,
			WriteTimeout:   2000,
			RetryTimes:     4000,
		},
		Host: []string{
			"1.1.1.2",
			"2.2.2.3",
		},
		InitConnSize: 1000,
		MaxConnSize:  1000,
		MaxIdleSec:   1000,
		ExpireSec:    1000,
		Password:     "xxx",
	})
	assert.Nil(t, err)

	c := &conn{}
	guard1 := monkey.Patch(redis.Dial, func(network, address string, options ...redis.DialOption) (redis.Conn, error) {

		assert.True(t, network == "tcp")
		assert.True(t, address == "1.1.1.2" || address == "2.2.2.3")

		if address == "2.2.2.3" {
			return c, nil
		}

		return nil, fmt.Errorf("-")
	})
	defer guard1.Unpatch()

	guard2 := monkey.PatchInstanceMethod(reflect.TypeOf(c), "Do", func(commandName string, args ...interface{}) (reply interface{}, err error) {

		assert.True(t, commandName == "HMSET" || commandName == "EXPIRE")
		if commandName == "HMSET" {
			assert.Equal(t, string(args[0].(string)), "a b c d")
		}
		assert.Equal(t, string(args[0].(string)), "44")

		return nil, nil
	})
	defer guard2.Unpatch()

	errExist := client.HMSetWithExpire("ABC", map[string]string{
		"a": "b",
		"c": "d",
	}, 44)
	assert.Nil(t, errExist)

}
*/

func TestHGetAll(t *testing.T) {
	client, err := NewRedisClient(&RedisConf{
		ConnConf: ConnConf{
			ConnectTimeout: 1000,
			ReadTimeout:    1000,
			WriteTimeout:   2000,
			RetryTimes:     1,
		},
		Host: []string{
			"1.1.1.2",
			"2.2.2.3",
		},
		InitConnSize: 10,
		MaxConnSize:  10,
		MaxIdleSec:   1000,
		ExpireSec:    1000,
		Password:     "xxx",
	})
	assert.Nil(t, err)

	c := &conn{}
	guard1 := monkey.Patch(redis.Dial, func(network, address string, options ...redis.DialOption) (redis.Conn, error) {
		assert.True(t, network == "tcp")
		assert.True(t, address == "1.1.1.2" || address == "2.2.2.3")

		if address == "2.2.2.3" {
			return c, nil
		}
		return nil, fmt.Errorf("-")
	})
	defer guard1.Unpatch()

	guard2 := monkey.PatchInstanceMethod(reflect.TypeOf(c), "Do", func(_ *conn, commandName string, args ...interface{}) (reply interface{}, err error) {
		assert.True(t, commandName == "HGETALL" || commandName == "")
		if commandName == "HGETALL" {
			assert.Equal(t, string(args[0].(string)), "ABC")
		}

		return []interface{}{[]byte("a"), []byte("b"), []byte("c"), []byte("d")}, nil
	})
	defer guard2.Unpatch()

	ret, errExist := client.HGetAll("ABC")
	assert.Equal(t, ret["a"], "b")
	assert.Equal(t, ret["c"], "d")
	assert.Nil(t, errExist)

}
