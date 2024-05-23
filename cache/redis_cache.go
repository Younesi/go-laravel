package cache

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
)

type RedisCache struct {
	Conn   *redis.Pool
	Prefix string
}

func (c *RedisCache) Has(str string) (bool, error) {
	key := fmt.Sprintf("%s:%s", c.Prefix, str)
	conn := c.Conn.Get()
	defer conn.Close()

	ok, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return false, err
	}

	return ok, err
}

func (c *RedisCache) Get(str string) (interface{}, error) {
	key := fmt.Sprintf("%s:%s", c.Prefix, str)
	conn := c.Conn.Get()
	defer conn.Close()

	cacheEntry, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return false, err
	}

	decoded, err := decode(string(cacheEntry))
	if err != nil {
		return false, err
	}
	item := decoded[key]

	return item, err
}

func (c *RedisCache) Set(str string, value interface{}, ttl ...int) error {
	key := fmt.Sprintf("%s:%s", c.Prefix, str)
	conn := c.Conn.Get()
	defer conn.Close()

	entry := Entry{}
	entry[key] = value
	encoded, err := encode(entry)
	if err != nil {
		return err
	}

	if len(ttl) > 0 {
		_, err := conn.Do("SETEX", key, ttl[0], string(encoded))
		if err != nil {
			return err
		}
	} else {
		// refactor to use always one of them if possible
		_, err := conn.Do("SET", key, string(encoded))
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *RedisCache) Forget(str string) error {
	key := fmt.Sprintf("%s:%s", c.Prefix, str)
	conn := c.Conn.Get()
	defer conn.Close()

	_, err := conn.Do("DEL", key)
	if err != nil {
		return err
	}

	return nil
}

func (c *RedisCache) EmptyByMatch(str string) error {
	// implment later
	return nil
}

func (c *RedisCache) Empty() error {
	// implment later
	return nil
}
