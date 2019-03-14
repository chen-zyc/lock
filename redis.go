package lock

import (
	"time"

	"github.com/garyburd/redigo/redis"
)

var unlockScript = redis.NewScript(1, `
if redis.call("get",KEYS[1]) == ARGV[1]
then
    return redis.call("del",KEYS[1])
else
    return 0
end
`)

type RedisLock struct{}

func (rl RedisLock) TryLock(conn redis.Conn, key, token string, expire time.Duration) (ok bool, err error) {
	reply, err := redis.String(conn.Do("SET", key, token, "PX", int64(expire/time.Millisecond), "NX"))
	if err != nil {
		if err == redis.ErrNil {
			return false, nil
		}
		return
	}
	return reply == "OK", nil
}

func (rl RedisLock) Unlock(conn redis.Conn, key, token string) (ok bool, err error) {
	reply, err := redis.Int(unlockScript.Do(conn, key, token))
	if err != nil {
		return
	}
	return reply == 1, nil
}
