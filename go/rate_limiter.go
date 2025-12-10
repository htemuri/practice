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
	"sync"
	"time"
)

type RateLimiter struct {
	counter   int
	threshold int
	dripRate  int
	quit      chan struct{}
	mutex     sync.Mutex
}

func NewRateLimiter(rate int, per time.Duration) *RateLimiter {
	rateLimiter := RateLimiter{
		counter:   0,
		threshold: 100,
		dripRate:  rate,
		quit:      make(chan struct{}),
	}
	log.Println(per / time.Duration(rate))
	ticker := time.Tick(per / time.Duration(rate))

	go func() {
		for {
			select {
			case <-ticker:
				rateLimiter.mutex.Lock()
				rateLimiter.counter--
				if rateLimiter.counter < 0 {
					rateLimiter.counter = 0
				}
				log.Println("drip...")
				log.Println("Current counter: ", rateLimiter.counter)
				rateLimiter.mutex.Unlock()
			case <-rateLimiter.quit:
				log.Println("Quitting...")
				return
			}
		}
	}()
	return &rateLimiter
}

func (r *RateLimiter) Allow() bool {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if r.counter < r.threshold {
		r.counter++
		log.Println("Counter after recieving request: ", r.counter)
		return true
	}
	return false
}

func (r *RateLimiter) Stop() {
	close(r.quit)
}

func main() {
	start := time.Now()
	rateLimiter := NewRateLimiter(5, time.Duration(time.Second*12))
	go func() {
		for {
			if rateLimiter.Allow() {
				log.Println("Request recieved")
			} else {
				log.Printf("Request rejected. Bucket full!")
			}
			time.Sleep(time.Duration((rand.Intn(500) + 1800) * int(time.Millisecond)))
		}
	}()
	for {
		if time.Since(start).Round(time.Second) >= time.Second*20 {
			log.Println("sending stop signal")
			rateLimiter.Stop()
			return
		}
	}
}
