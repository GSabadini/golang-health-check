package check

import (
	"github.com/go-redis/redis"
	"log"
)

type Redis struct {
	DNS string
}

func (r Redis) Check() error {
	rdb := redis.NewClient(&redis.Options{
		Addr: r.DNS,
	})
	defer rdb.Close()

	pong, err := rdb.Ping().Result()
	if err != nil {
		log.Printf("redis ping failed: %w", err)
		return err
	}

	if pong != "PONG" {
		log.Printf("unexpected response for redis ping: %q", pong)
		return err
	}

	return nil
}
