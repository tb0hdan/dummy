package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/gomodule/redigo/redis"
)

// ErrNil indicates that a reply value is nil.
var ErrNil = errors.New("(nil)")

// Redis describes connection to Redis server
type Redis struct {
	pool *redis.Pool
}

// New returns the initialized Redis object
func New(ctx context.Context, redisServer string) (*Redis, error) {
	log.Info("Redis init: host=", redisServer)
	c := new(Redis)
	c.initNewPool(redisServer)
	if err := c.Ping(); err != nil {
		return nil, err
	}
	go func() {
		<-ctx.Done()
		err := c.pool.Close()
		if err != nil {
			log.Error("close redis connection error:", err.Error())
			return
		}
		log.Info("close redis connection")
	}()

	return c, nil
}

func (c *Redis) initNewPool(addr string) {
	c.pool = &redis.Pool{
		//MaxIdle:     5,
		//IdleTimeout: 30 * time.Second,
		Dial: func() (redis.Conn, error) { return redis.Dial("tcp", addr) },
		//TestOnBorrow: func(c redis.Conn, t time.Time) error {
		//	_, err := c.Do("PING")
		//	return err
		//},
		//MaxActive: 5,
		//Wait:      true,
	}
}

// Ping checks if connection exists
func (c *Redis) Ping() error {
	r := c.pool.Get()
	defer r.Close()
	_, err := r.Do("PING")
	return err
}

func (c *Redis) SaveIDFA(xid int64, idfa string) error {
	const ttl30days = 2592000 // 30 days
	r := c.pool.Get()
	defer r.Close()
	_, err := r.Do("SETEX", xid, ttl30days, idfa)
	return err
}

func (c *Redis) ReadIDFA(xid int64) (string, error) {
	r := c.pool.Get()
	defer r.Close()
	idfa, err := redis.String(r.Do("GET", xid))
	if err == redis.ErrNil {
		err = ErrNil
	}
	return idfa, err
}

func (c *Redis) SaveRegKeys(xid int64, keys map[string]string) error {
	const ttl10min = 600
	js, err := json.Marshal(keys)
	if err != nil {
		return err
	}
	r := c.pool.Get()
	defer r.Close()
	_, err = r.Do("SETEX", fmt.Sprintf("%d:keys", xid), ttl10min, js)
	return err
}

func (c *Redis) ReadRegKeys(xid int64) (map[string]string, error) {
	r := c.pool.Get()
	defer r.Close()
	js, err := redis.Bytes(r.Do("GET", fmt.Sprintf("%d:keys", xid)))
	if err != nil {
		return nil, err
	}
	keys := make(map[string]string)
	err = json.Unmarshal(js, &keys)
	if err != nil {
		return nil, err
	}
	return keys, nil
}
