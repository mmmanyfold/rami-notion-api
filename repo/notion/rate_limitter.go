package notion

import (
	"context"
	"fmt"
	"github.com/sethvargo/go-limiter"
	"github.com/sethvargo/go-limiter/memorystore"
	"log"
	"time"
)

type RateLimiter struct {
	ctx   context.Context
	store limiter.Store
}

func NewRateLimiter(ctx context.Context, tokens uint64) (*RateLimiter, error) {
	store, err := memorystore.New(&memorystore.Config{
		Tokens:   tokens,
		Interval: time.Second / 3,
	})
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &RateLimiter{
		ctx:   ctx,
		store: store,
	}, nil
}

func (rl *RateLimiter) Take() (bool, error) {
	limit, remaining, reset, ok, err := rl.store.Take(rl.ctx, "my-key")
	if err != nil {
		log.Fatal(err)
		return false, err
	}

	fmt.Println("limit: ", limit)
	fmt.Println("remaining: ", remaining)
	fmt.Println("reset: ", reset)
	fmt.Println("ok: ", ok)
	fmt.Println("+---------------+")

	return ok, nil
}

func (rl *RateLimiter) Close() error {
	if err := rl.store.Close(rl.ctx); err != nil {
		return err
	}
	return nil
}
