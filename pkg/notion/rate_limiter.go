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
	debug    bool
	key      string
	limit    uint64
	Interval time.Duration
	ctx      context.Context
	store    limiter.Store
}

func NewRateLimiter(ctx context.Context, key string, tokens uint64, interval time.Duration, debug bool) (*RateLimiter, error) {
	store, err := memorystore.New(&memorystore.Config{
		Tokens:   tokens,
		Interval: interval,
	})
	if err != nil {
		return nil, err
	}

	return &RateLimiter{
		ctx:      ctx,
		key:      key,
		limit:    tokens,
		store:    store,
		debug:    debug,
		Interval: interval,
	}, nil
}

func (rl *RateLimiter) Take() (ok bool, remaining uint64, err error) {
	limit, remaining, reset, ok, err := rl.store.Take(rl.ctx, rl.key)
	if err != nil {
		log.Fatal(err)
		return ok, remaining, err
	}

	if rl.debug {
		fmt.Println("fetching: ", rl.key)
		fmt.Println("limit: ", limit)
		fmt.Println("remaining: ", remaining)
		fmt.Println("reset: ", reset)
		fmt.Println("ok: ", ok)
		fmt.Println("+---------------+")
	}

	return ok, remaining, nil
}

func (rl *RateLimiter) Idle() {
	if rl.debug {
		fmt.Println("idling...")
	}
	time.Sleep(time.Second)
}

func (rl *RateLimiter) Close() error {
	if err := rl.store.Close(rl.ctx); err != nil {
		return err
	}
	return nil
}
