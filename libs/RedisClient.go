package libs

import (
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"time"
)

// redis的api按需打开，现在不用的都注释了

var ErrKeyExist = errors.New("key already exist")
var ErrKeyNotExist = errors.New("key not exist")
var BadKeyValue = errors.New("bad key value")

var RedisErrExhausted = redis.ErrPoolExhausted

type RedisClient struct {
	serverList     []string
	maxIdleConn    int
	maxActiveConn  int
	maxIdleTime    time.Duration
	maxOpRetry     int
	expireSec      int
	serverIndex    int
	redisPool      *redis.Pool
	password       string
	connectTimeout int
	readTimeout    int
	writeTimeout   int
}

/*
func StringMap(result interface{}, columns []string, err error) (map[string]string, error) {
	values, err := redis.Values(result, err)
	if err != nil {
		return nil, err
	}
	m := make(map[string]string, len(values))
	for i := 0; i < len(columns); i++ {
		key := columns[i]
		value, okValue := values[i].([]byte)
		if !okValue {
			m[key] = ""
		}
		m[key] = string(value)
	}
	return m, nil
}
*/

func (r *RedisClient) GetConnectionStatus() string {
	if nil == r {
		return "active:0 max:0"
	}

	return fmt.Sprintf("active:%d max:%d", r.redisPool.ActiveCount(), r.redisPool.MaxActive)
}

func (r *RedisClient) Init(config *RedisConf) (rc *RedisClient, err error) {
	defer func() {
		if fatal := recover(); fatal != nil {
			err = fmt.Errorf("Fatal error: %v", fatal)
		}
	}()
	if len(config.Host) == 0 {
		return nil, fmt.Errorf("Empty redis server list.")
	}
	r.serverList = config.Host

	if config.InitConnSize < 1 || config.MaxConnSize < 1 {
		return nil, fmt.Errorf("InitConnSize and MaxConnSize shouldn't be < 1")
	}
	r.maxIdleConn = config.InitConnSize
	r.maxActiveConn = config.MaxConnSize
	//r.maxIdleTime, err = time.ParseDuration(config.Redis.MaxIdleSec)
	//if err != nil {
	//	return nil, err
	//}
	r.maxIdleTime = time.Second * time.Duration(config.MaxIdleSec)
	r.maxOpRetry = config.RetryTimes
	if config.RetryTimes < 0 {
		r.maxOpRetry = 0
	}
	r.expireSec = config.ExpireSec
	r.serverIndex = 0
	r.password = config.Password
	r.connectTimeout = config.ConnectTimeout
	r.readTimeout = config.ReadTimeout
	r.writeTimeout = config.WriteTimeout
	serverCnt := len(r.serverList)
	// Init redis pool
	r.redisPool = &redis.Pool{
		MaxIdle:     r.maxIdleConn,
		MaxActive:   r.maxActiveConn,
		IdleTimeout: r.maxIdleTime,
		Dial: func() (redis.Conn, error) {
			fmt.Println(">> PPP")
			serverIdx := r.serverIndex
			r.serverIndex += 1
			var lastError error
			// If connect failed, we try next
			for try := 0; try < serverCnt; try++ {
				server := r.serverList[(serverIdx+try)%serverCnt]
				c, err := redis.Dial("tcp", server,
					redis.DialPassword(r.password),
					redis.DialConnectTimeout(time.Duration(r.connectTimeout)*time.Millisecond),
					redis.DialReadTimeout(time.Duration(r.readTimeout)*time.Millisecond),
					redis.DialWriteTimeout(time.Duration(r.writeTimeout)*time.Millisecond))
				if err != nil {
					lastError = fmt.Errorf("Dial %s Fail", server)
					continue
				}

				return c, err
			}
			return nil, lastError
		},
	}

	return r, nil
}

/*
func (r *RedisClient) Set(key string, value string) (err error) {
	defer func() {
		if fatal := recover(); fatal != nil {
			err = fmt.Errorf("Fatal error: %v", fatal)
		}
	}()

	var lastError error

	for try := 0; try <= r.maxOpRetry; try++ {
		conn := r.redisPool.Get()
		defer conn.Close()
		_, err = redis.String(conn.Do("SET", key, value))
		if err != nil {
			lastError = err
			continue
		}
		return nil
	}

	return fmt.Errorf("Set key %s faild %d times: %v", key, r.maxOpRetry, lastError)
}


func (r *RedisClient) SetWithExpireTime(key string, value string, expireSec int) (err error) {
	defer func() {
		if fatal := recover(); fatal != nil {
			err = fmt.Errorf("Fatal error: %v", fatal)
		}
	}()

	var lastError error

	for try := 0; try <= r.maxOpRetry; try++ {
		conn := r.redisPool.Get()
		defer conn.Close()
		_, err := redis.String(conn.Do("SET", key, value, "EX", expireSec))
		if err != nil {
			lastError = err
			continue
		}
		return nil
	}

	return fmt.Errorf("Set key %s faild %d times: %v", key, r.maxOpRetry, lastError)
}


func (r *RedisClient) SetIfNotExist(key string, value string, expireSec int) (err error) {
	defer func() {
		if fatal := recover(); fatal != nil {
			err = fmt.Errorf("Fatal error: %v", fatal)
		}
	}()

	var lastError error

	for try := 0; try <= r.maxOpRetry; try++ {
		conn := r.redisPool.Get()
		defer conn.Close()
		_, err := redis.String(conn.Do("SET", key, value, "NX", "EX", expireSec))
		if err != nil {
			if err == redis.ErrNil {
				return ErrKeyExist
			}
			lastError = err
			continue
		}
		return nil
	}
	return fmt.Errorf("Set key %s faild %d times: %v", key, r.maxOpRetry, lastError)
}
*/

func (r *RedisClient) Exist(key string) (ret int64, err error) {
	defer func() {
		if fatal := recover(); fatal != nil {
			err = fmt.Errorf("Fatal error: %v", fatal)
		}
	}()

	var lastError error

	for try := 0; try <= r.maxOpRetry; try++ {
		conn := r.redisPool.Get()
		defer conn.Close()
		ret, err = redis.Int64(conn.Do("EXISTS", key))
		if err != nil {
			lastError = err
			continue
		}
		return
	}
	err = lastError
	return
}

/*
func (r *RedisClient) Expire(key string) (err error) {
	defer func() {
		if fatal := recover(); fatal != nil {
			err = fmt.Errorf("Fatal error: %v", fatal)
		}
	}()

	var lastError error

	for try := 0; try <= r.maxOpRetry; try++ {
		conn := r.redisPool.Get()
		defer conn.Close()
		_, err := conn.Do("EXPIRE", key, r.expireSec)
		if err != nil {
			lastError = err
			continue
		}
		return nil
	}
	err = lastError
	return
}

func (r *RedisClient) ExpireSec(key string, sec int) (err error) {
	defer func() {
		if fatal := recover(); fatal != nil {
			err = fmt.Errorf("Fatal error: %v", fatal)
		}
	}()

	var lastError error

	for try := 0; try <= r.maxOpRetry; try++ {
		conn := r.redisPool.Get()
		defer conn.Close()
		_, err := conn.Do("EXPIRE", key, sec)
		if err != nil {
			lastError = err
			continue
		}
		return nil
	}
	err = lastError
	return
}

func (r *RedisClient) HMSet(key string, kvs map[string]string) (err error) {
	if len(kvs) < 1 {
		err = BadKeyValue
		return
	}
	var lastError error
	args := make([]interface{}, 0, len(kvs)+1)
	args = append(args, key)
	for k, v := range kvs {
		args = append(args, k)
		args = append(args, v)
	}

	for try := 0; try <= r.maxOpRetry; try++ {
		conn := r.redisPool.Get()
		defer conn.Close()
		_, err := conn.Do("HMSET", args...)
		if err != nil {
			lastError = err
			continue
		}
		//logger.V(0).Tracef("Codis HSET, key:%s value:%s", key, kvs)
		return nil
	}
	err = lastError
	return
}
*/

func (r *RedisClient) HMSetWithExpire(key string, kvs map[string]string, sec int) (err error) {
	if len(kvs) < 1 {
		err = BadKeyValue
		return
	}
	if sec < 1 {
		sec = r.expireSec
	}
	var lastError error
	args := make([]interface{}, 0, 2*len(kvs)+1)
	args = append(args, key)
	for k, v := range kvs {
		args = append(args, k)
		args = append(args, v)
	}

	for try := 0; try <= r.maxOpRetry; try++ {
		conn := r.redisPool.Get()
		defer conn.Close()
		_, err0 := conn.Do("HMSET", args...)
		if err0 != nil {
			lastError = err0
			continue
		}
		_, err0 = conn.Do("EXPIRE", key, sec)
		if err0 != nil {
			lastError = err0
			continue
		}
		return nil
	}
	err = lastError
	return err
}

func (r *RedisClient) HGetAll(key string) (ret map[string]string, err error) {
	defer func() {
		if fatal := recover(); fatal != nil {
			err = fmt.Errorf("Fatal error: %v", fatal)
		}
	}()

	var lastError error

	for try := 0; try <= r.maxOpRetry; try++ {
		conn := r.redisPool.Get()
		defer conn.Close()
		ret, err = redis.StringMap(conn.Do("HGETALL", key))
		if err != nil {
			if err == redis.ErrNil {
				err = ErrKeyNotExist
				return
			}
			lastError = err
			continue
		}
		return
	}
	err = lastError
	return
}

/*
func (r *RedisClient) HMGet(key string, columns []string) (ret map[string]string, err error) {
	//defer func() {
	//	if fatal := recover(); fatal != nil {
	//		err = fmt.Errorf("Fatal error: %v", fatal)
	//	}
	//}()

	var lastError error

	args := make([]interface{}, 0, len(columns)+1)
	args = append(args, key)
	for _, k := range columns {
		args = append(args, k)
	}

	for try := 0; try <= r.maxOpRetry; try++ {
		conn := r.redisPool.Get()
		defer conn.Close()
		//fmt.Println("args:", args)
		res, e := conn.Do("HMGET", args...)
		//fmt.Println("res:", res)
		if e != nil {
			if e == redis.ErrNil {
				err = ErrKeyNotExist
				return
			}
			lastError = err
			continue
		}
		ret, err = StringMap(res, columns, e)
		if err != nil {
			return
		}
		//ret, err = redis.StringMap(res, e)
		//logger.V(0).Tracef("codis HMGET key:%s columns:%s res:%d", key, columns, len(res))
		//logger.V(0).Tracef("codis HMGET key:%s columns:%s ret:%s", key, columns, ret)
		return
	}
	err = lastError
	return
}
*/

/*

func (r *RedisClient) Get(key string) (ret string, err error) {
	defer func() {
		if fatal := recover(); fatal != nil {
			err = fmt.Errorf("Fatal error: %v", fatal)
		}
	}()

	var lastError error

	for try := 0; try <= r.maxOpRetry; try++ {
		conn := r.redisPool.Get()
		defer conn.Close()
		value, err := redis.String(conn.Do("GET", key))
		if err != nil {
			if err == redis.ErrNil {
				return "", nil
			}
			lastError = err
			continue
		}
		return value, nil
	}
	return "", lastError
}


func (r *RedisClient) Del2(key1 string, key2 string) error {
	var lastError error

	for try := 0; try <= r.maxOpRetry; try++ {
		conn := r.redisPool.Get()
		defer conn.Close()
		_, err := redis.Int(conn.Do("DEL", key1, key2))
		if err != nil {
			lastError = err
			continue
		}
		return nil
	}

	return lastError

}

func (r *RedisClient) Del(key string) error {
	var lastError error

	for try := 0; try <= r.maxOpRetry; try++ {
		conn := r.redisPool.Get()
		defer conn.Close()
		_, err := redis.Int(conn.Do("DEL", key))
		if err != nil {
			lastError = err
			continue
		}
		return nil
	}

	return lastError
}



// LPushWithExpire 对Redis做LPUSH操作
func (r *RedisClient) LPushWithExpire(key string, values [][]byte, sec int) (listSize int, err error) {
	var lastError error
	args := make([]interface{}, 0, len(values)+1)
	args = append(args, key)
	for _, v := range values {
		args = append(args, v)
	}

	for try := 0; try <= r.maxOpRetry; try++ {
		conn := r.redisPool.Get()
		defer conn.Close()
		listSize, err = redis.Int(conn.Do("LPUSH", args...))
		if err != nil {
			lastError = err
			continue
		}
		_, err = conn.Do("EXPIRE", key, sec)
		if err != nil {
			lastError = err
			continue
		}
		return listSize, nil
	}

	return listSize, lastError
}



// RPushWithExpire 对Redis做RPUSH操作
func (r *RedisClient) RPushWithExpire(key string, values [][]byte, sec int) (listSize int, err error) {
	var lastError error
	args := make([]interface{}, 0, len(values)+1)
	args = append(args, key)
	for _, v := range values {
		args = append(args, v)
	}

	for try := 0; try <= r.maxOpRetry; try++ {
		conn := r.redisPool.Get()
		defer conn.Close()
		listSize, err = redis.Int(conn.Do("RPUSH", args...))
		if err != nil {
			lastError = err
			continue
		}
		_, err = conn.Do("EXPIRE", key, sec)
		if err != nil {
			lastError = err
			continue
		}
		return listSize, nil
	}

	return listSize, lastError
}

func (r *RedisClient) LPop(key string) (value string, err error) {

	for try := 0; try <= r.maxOpRetry; try++ {
		conn := r.redisPool.Get()
		defer conn.Close()
		value, err = redis.String(conn.Do("LPOP", key))
		if err != nil {
			continue
		}
		return
	}
	return
}


// LRange 对 Redis 做LRANGE操作
func (r *RedisClient) LRange(key string, s int, e int) (values [][]byte, err error) {
	for try := 0; try <= r.maxOpRetry; try++ {
		conn := r.redisPool.Get()
		defer conn.Close()
		values, err = redis.ByteSlices(conn.Do("LRANGE", key, s, e))
		if err != nil {
			continue
		}
		return
	}
	return

}

func (r *RedisClient) LIndex(key string, idx int) (value []byte, err error) {
	for try := 0; try <= r.maxOpRetry; try++ {
		conn := r.redisPool.Get()
		defer conn.Close()
		value, err = redis.Bytes(conn.Do("LINDEX", key, idx))
		if err != nil {
			continue
		}
		return
	}
	return
}
*/

/*
func (r *RedisClient) ExpireMultiKeys(keys []string, sec int) error {
	conn := r.redisPool.Get()
	defer conn.Close()
	var lastError error
	for _, key := range keys {
		_, err := conn.Do("EXPIRE", key, sec)
		if err != nil {
			lastError = err
		}
	}
	return lastError
}
*/

/*
func (r *RedisClient) Publish(key string, val interface{}) (err error) {
	defer func() {
		if fatal := recover(); fatal != nil {
			err = fmt.Errorf("Fatal error: %v", fatal)
		}
	}()

	var lastError error

	for try := 0; try <= r.maxOpRetry; try++ {
		conn := r.redisPool.Get()
		defer conn.Close()
		_, err = redis.Int(conn.Do("PUBLISH", key, val))
		if err != nil {
			lastError = err
			continue
		}
		return nil
	}

	return fmt.Errorf("Set key %s faild %d times: %v", key, r.maxOpRetry, lastError)
}
*/

func NewRedisClient(config *RedisConf) (*RedisClient, error) {
	rclient := new(RedisClient)
	return rclient.Init(config)
}
