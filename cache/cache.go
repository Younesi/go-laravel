package cache

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/gomodule/redigo/redis"
)

type Cache interface {
	Has(string) (bool, error)
	Get(string) (interface{}, error)
	Set(string, interface{}, ...int) error
	Forget(string) error
	EmptyByMatch(string) error
	Empty() error
}

type RedisCache struct {
	Conn   *redis.Pool
	Prefix string
}

type entry map[string]interface{}

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

	entry := entry{}
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

func encode(item entry) ([]byte, error) {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	if err := e.Encode(item); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func decode(str string) (entry, error) {
	item := entry{}
	b := bytes.Buffer{}
	b.Write([]byte(str))
	d := gob.NewDecoder(&b)
	if err := d.Decode(&item); err != nil {
		return nil, err
	}

	return item, nil
}
