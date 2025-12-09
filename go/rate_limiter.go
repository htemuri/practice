/*
Build a Rate Limiter (Leaky Bucket)

Core skills: timers, tickers, select, non-blocking channel patterns.

What to implement:
	•	NewRateLimiter(rate int, per time.Duration)
	•	Allow() bool returns true only if a token is available
	•	Background goroutine drips tokens
	•	Fixed max bucket size

This teaches event scheduling + concurrency correctness.
*/

package main

import (
	"log"
	"math/rand"
	"time"
)

type RateLimiter struct {
	counter   int
	threshold int
	rate      int
	quit      chan int
	ticker    <-chan time.Time
}

func NewRateLimiter(rate int, per time.Duration) *RateLimiter {
	rateLimiter := RateLimiter{
		counter:   0,
		threshold: 100,
		rate:      rate,
		ticker:    time.Tick(per),
		quit:      make(chan int),
	}

	go func() {
		for {
			select {
			case <-rateLimiter.ticker:
				rateLimiter.counter -= rate
				if rateLimiter.counter < 0 {
					rateLimiter.counter = 0
				}
				log.Println("drip...")
			case <-rateLimiter.quit:
				log.Println("Quitting...")
				return
			}
		}
	}()
	return &rateLimiter
}

func (r *RateLimiter) Allow(random int) bool {
	return r.counter+random < r.threshold
}

func (r *RateLimiter) Stop() {
	r.quit <- 1
}

func main() {
	start := time.Now()
	rateLimiter := NewRateLimiter(10, time.Duration(time.Second*2))
	go func() {
		for {
			select {
			case <-rateLimiter.quit:
				return
			default:
				random := rand.Intn(50)
				if rateLimiter.Allow(random) {
					log.Println("Sending token size: ", random)
					rateLimiter.counter += random
				} else {
					log.Printf("Token size %d rejected. Bucket full!", random)
				}
				time.Sleep(time.Duration((rand.Intn(5000) + 1000) * int(time.Millisecond)))
			}
		}
	}()
	for {
		log.Println("Counter: ", rateLimiter.counter)
		time.Sleep(2000 * time.Millisecond)
		if time.Since(start).Round(time.Second) >= time.Second*20 {
			log.Println("sending stop signal")
			rateLimiter.Stop()
			return
		}
	}
}
