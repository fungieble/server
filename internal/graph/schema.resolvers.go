package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/aliadotsh/alia/internal/graph/generated"
	"github.com/aliadotsh/alia/internal/graph/models"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
)

// CreateMessage is the resolver for the createMessage field.
func (r *mutationResolver) CreateMessage(ctx context.Context, message string) (*models.Message, error) {
	m := models.Message{
		Message: message,
	}

	r.RedisClient.XAdd(&redis.XAddArgs{
		Stream: "room",
		ID:     "*",
		Values: map[string]interface{}{
			"message": m.Message,
		},
	})

	return &m, nil
}

// EmailUserAuthChallenge is the resolver for the emailUserAuthChallenge field.
func (r *mutationResolver) EmailUserAuthChallenge(ctx context.Context, email string) (bool, error) {
	panic(fmt.Errorf("not implemented: EmailUserAuthChallenge - emailUserAuthChallenge"))
}

// EmailUserAuthTokenChallenge is the resolver for the emailUserAuthTokenChallenge field.
func (r *mutationResolver) EmailUserAuthTokenChallenge(ctx context.Context, email string, token string) (bool, error) {
	panic(fmt.Errorf("not implemented: EmailUserAuthTokenChallenge - emailUserAuthTokenChallenge"))
}

// Viewer is the resolver for the viewer field.
func (r *queryResolver) Viewer(ctx context.Context) (*models.User, error) {
	return user, nil
}

// User is the resolver for the user field.
func (r *queryResolver) User(ctx context.Context, input models.UserInput) (*models.User, error) {
	if input.ID != nil {
		if *input.ID != user.ID {
			return nil, nil
		}
	}

	if input.Username != nil {
		if *input.Username != user.Username {
			return nil, nil
		}
	}
	return user, nil
}

// Messages is the resolver for the messages field.
func (r *queryResolver) Messages(ctx context.Context) ([]*models.Message, error) {
	streams, err := r.RedisClient.XRead(&redis.XReadArgs{
		Streams: []string{"room", "0"},
	}).Result()
	if !errors.Is(err, nil) {
		panic(err)
	}

	stream := streams[0]

	ms := make([]*models.Message, len(stream.Messages))
	for i, m := range stream.Messages {
		ms[i] = &models.Message{
			ID:      m.ID,
			Message: m.Values["message"].(string),
		}
	}

	return ms, nil
}

// MessageCreated is the resolver for the messageCreated field.
func (r *subscriptionResolver) MessageCreated(ctx context.Context) (<-chan *models.Message, error) {
	t := uuid.New()
	mc := make(chan *models.Message)
	r.mutex.Lock()
	r.messageChannels[t.String()] = mc
	r.mutex.Unlock()

	go func() {
		<-ctx.Done()
		r.mutex.Lock()
		delete(r.messageChannels, t.String())
		r.mutex.Unlock()
		log.Println("Subscription closed")
	}()

	log.Println("Subscription: message created")

	return mc, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }

var user = &models.User{
	ID:          "d9070b51-d3b9-45eb-9c5f-5ec6187c5c1b",
	Email:       "teddy@alia.sh",
	FirstName:   "Theodore",
	LastName:    "Verhoeff",
	Username:    "teddy",
	Description: "I am a talented product designer and engineer with a passion for crafting beautiful, polished interfaces. My expertise in design and development allows me to create user-friendly products that meet customer needs and exceed expectations.",
}
