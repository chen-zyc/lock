package lock

import (
	"testing"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/stretchr/testify/require"
)

func TestRedisLock(t *testing.T) {
	conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	require.NoError(t, err)
	defer func() { _ = conn.Close() }()

	rl := RedisLock{}
	ok, err := rl.TryLock(conn, "key1", "123", 10*time.Second)
	require.NoError(t, err)
	require.True(t, ok)

	ok, err = rl.TryLock(conn, "key1", "456", 10*time.Second)
	require.NoError(t, err)
	require.False(t, ok)

	ok, err = rl.Unlock(conn, "key1", "456")
	require.NoError(t, err)
	require.False(t, ok)

	ok, err = rl.Unlock(conn, "key1", "123")
	require.NoError(t, err)
	require.True(t, ok)

	ok, err = rl.TryLock(conn, "key1", "456", 10*time.Second)
	require.NoError(t, err)
	require.True(t, ok)
}
