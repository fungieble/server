//go:generate go run github.com/99designs/gqlgen generate

package graph

import (
	"errors"
	"log"
	"sync"

	"github.com/fungiedotsh/fungie/internal/graph/models"
	"github.com/go-redis/redis"
)

type Resolver struct {
	RedisClient     *redis.Client
	messageChannels map[string]chan *models.Message
	mutex           sync.Mutex
}

func NewResolver(client *redis.Client) *Resolver {
	return &Resolver{
		RedisClient:     client,
		messageChannels: map[string]chan *models.Message{},
		mutex:           sync.Mutex{},
	}
}

func (r *Resolver) SubscribeRedis() {
	log.Println("Initiating redis stream…")

	go func() {
		for {
			log.Println("Stream starting…")

			streams, err := r.RedisClient.XRead(&redis.XReadArgs{
				Streams: []string{"room", "$"},
				Block:   0,
			}).Result()
			if !errors.Is(err, nil) {
				panic(err)
			}

			stream := streams[0]
			m := &models.Message{
				ID:      stream.Messages[0].ID,
				Message: stream.Messages[0].Values["message"].(string),
			}
			r.mutex.Lock()
			for _, ch := range r.messageChannels {
				// Wait for a response from the channel
				ch <- m
			}
			r.mutex.Unlock()

			log.Println("Stream finished…")
		}
	}()
}
