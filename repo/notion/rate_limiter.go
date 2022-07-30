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

func NewRateLimiter(ctx context.Context, key string, tokens uint64, debug bool) (*RateLimiter, error) {
	store, err := memorystore.New(&memorystore.Config{
		Tokens:   tokens,
		Interval: time.Second / 3,
	})
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &RateLimiter{
		ctx:      ctx,
		key:      key,
		limit:    tokens,
		Interval: time.Second / 3,
		store:    store,
		debug:    debug,
	}, nil
}

func (rl *RateLimiter) Take() (ok bool, remaining uint64, err error) {
	fmt.Printf("--- fetching: %s\n", rl.key)
	limit, remaining, reset, ok, err := rl.store.Take(rl.ctx, rl.key)
	if err != nil {
		log.Fatal(err)
		return ok, remaining, err
	}

	if rl.debug {
		fmt.Println("limit: ", limit)
		fmt.Println("remaining: ", remaining)
		fmt.Println("reset: ", reset)
		fmt.Println("ok: ", ok)
		fmt.Println("+---------------+")
	}

	return ok, remaining, nil
}

func (rl *RateLimiter) Burst(requests uint64) {
	//count := make()
}

func (rl *RateLimiter) Idle() {
	fmt.Println("idling...")
	time.Sleep(time.Second)
}

func (rl *RateLimiter) Close() error {
	if err := rl.store.Close(rl.ctx); err != nil {
		return err
	}
	return nil
}
