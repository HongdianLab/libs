package cache

import (
	"github.com/garyburd/redigo/redis"
	. "github.com/smartystreets/goconvey/convey"

	"fmt"
	"testing"
	"time"
)

func newRedisPool(server, password string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 120 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if len(password) > 0 {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

type RedisLoader struct {
	pool *redis.Pool
}

func (this *RedisLoader) do(commandName string, args ...interface{}) (reply interface{}, err error) {
	c := this.pool.Get()
	defer c.Close()
	reply, err = c.Do(commandName, args...)
	return
}

func (this *RedisLoader) Load(key string) interface{} {
	value, err := redis.String(this.do("GET", key))
	if err != nil {
		fmt.Println(err)
	}
	return value
}

func (this *RedisLoader) Modify(key string, value interface{}) error {
	_, err := this.do("SET", key, value)
	if err != nil {
		fmt.Println(err)
	}
	return err
}

var (
	key    = "foo"
	value  = "bar"
	value2 = "bar2"
	loader = &RedisLoader{
		pool: newRedisPool("127.0.0.1:6379", ""),
	}
	c, err = NewCache(loader, 10)
)

func TestLoadingGet(t *testing.T) {
	Convey("GET", t, func() {
		loader.do("SET", key, value)
		v := c.Get(key)
		So(v, ShouldEqual, value)
	})
}

func TestModify(t *testing.T) {
	Convey("Modify", t, func() {
		err = c.Modify(key, value2)
		So(err, ShouldBeNil)
		v := c.Get(key)
		So(v, ShouldEqual, value2)
	})
}

func TestInvalid(t *testing.T) {
	Convey("Invalid", t, func() {
		c.Stop()

		err = c.Modify(key, value)
		v := c.Get(key)
		So(v, ShouldEqual, value)

		loader.do("SET", key, value2)
		v = c.Get(key)
		So(v, ShouldEqual, value)

		err = c.Invalid(key)
		So(err, ShouldBeNil)
		v = c.Get(key)
		So(v, ShouldEqual, value2)
	})
}
