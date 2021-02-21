package cacherepo

import (
	"sync"
	"term-frequency/repository"
	"time"

	"github.com/gomodule/redigo/redis"
)

var wg sync.WaitGroup

func NewCacheRepository(redisPool *redis.Pool) repository.ICacheRepository {
	return &RedisPool{
		redisPool: redisPool,
	}
}

type RedisPool struct {
	redisPool *redis.Pool
}

func (redisPool RedisPool) InsertTokens(tokens []string) {
	now := time.Now()
	keySet := now.Format("20060102_15")

	for _, token := range tokens {
		wg.Add(1)
		go cacheTokensInRedis(token, keySet, redisPool)
	}

	wg.Wait()
}

func (redisPool RedisPool) StoreKeyUnionOfTokens(keyTop string, t string, keys []string) error {
	conn := redisPool.redisPool.Get()
	defer conn.Close()
	var args []interface{}
	args = append(args, keyTop)
	args = append(args, t)
	for _, k := range keys {
		args = append(args, k)
	}
	_, err := conn.Do("ZUNIONSTORE", args...)
	if err != nil {
		return err
	}
	return nil
}

func (redisPool RedisPool) GetCountOfTokensInSortedSet(key string) (int, error) {
	conn := redisPool.redisPool.Get()
	defer conn.Close()
	totalTokenCount, err := redis.Int(conn.Do("ZCOUNT", key, "-inf", "+inf"))
	return totalTokenCount, err
}

func (redisPool RedisPool) ExpireKey(key string, t int) error {
	conn := redisPool.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("EXPIRE", key, t)
	return err
}

func (redisPool RedisPool) GetTopValuesOfSortedSet(key string, n string) ([]string, error) {
	conn := redisPool.redisPool.Get()
	defer conn.Close()
	values, err := redis.Strings(conn.Do("ZREVRANGEBYSCORE", key, "+inf", "-inf", "LIMIT", "0", n, "withscores"))
	return values, err
}

func cacheTokensInRedis(token string, keySet string, redisPool RedisPool) {
	defer wg.Done()
	c := redisPool.redisPool.Get()
	c.Do("ZINCRBY", keySet, 1, token)
	// expire key after 168 hours or 1 week
	c.Do("EXPIRE", keySet, 604800)
	c.Close()
}
