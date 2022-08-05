package notion

import "time"

// Rate request per second to notion API
const Rate uint64 = 3
const Limit = time.Second / 3

type Burst struct {
	rateLimiter *RateLimiter
}

func NewBurst() *Burst {
	//rl := NewRateLimiter()
	return &Burst{}
}
