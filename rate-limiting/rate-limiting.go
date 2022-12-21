package main

import (
	"fmt"
	"time"
)

/*
Rate limiting is an important mechanism for controlling
resource utilization and maintaining quality of service.
Go elegantly supports rate limiting with goroutines,
channels, and tickers.
*/
func main() {

	/*
		First we’ll look at basic rate limiting. Suppose we
		want to limit our handling of incoming requests.
		We’ll serve these requests off a channel of the same
		name.
	*/
	requests := make(chan int, 5)
	for i := 1; i <= 5; i++ {
		requests <- i
	}
	close(requests)

	/*
		This limiter channel will receive a value every 200
		milliseconds. This is the regulator in our rate
		limiting scheme.
	*/
	limiter := time.Tick(200 * time.Millisecond)

	/*
		By blocking on a receive from the limiter channel before
		serving each request, we limit ourselves to 1 request
		every 200 milliseconds.
	*/
	for req := range requests {
		<-limiter
		fmt.Println("request", req, time.Now())
	}

	/*
		We may want to allow short bursts of requests in our rate
		limiting scheme while preserving the overall rate limit.
		We can accomplish this by buffering our limiter channel.
		This burstyLimiter channel will allow bursts of up to 3
		events.
	*/
	burstyLimiter := make(chan time.Time, 3)

	/*
		Fill up the channel to represent allowed bursting.
	*/
	for i := 0; i < 3; i++ {
		burstyLimiter <- time.Now()
	}

	/*
		Every 200 milliseconds we’ll try to add a new
		value to burstyLimiter, up to its limit of 3.
	*/
	go func() {
		for t := range time.Tick(200 * time.Millisecond) {
			burstyLimiter <- t
		}
	}()

	/*
		Now simulate 5 more incoming requests. The first 3
		of these will benefit from the burst capability of
		burstyLimiter.
	*/
	burstyRequests := make(chan int, 5)
	for i := 1; i <= 5; i++ {
		burstyRequests <- i
	}
	close(burstyRequests)
	for req := range burstyRequests {
		<-burstyLimiter
		fmt.Println("request", req, time.Now())
	}
}

/*
Running our program we see the first batch of requests handled
once every ~200 milliseconds as desired.

➜  gru git:(main) ✗ go run rate-limiting/rate-limiting.go
request 1 2022-12-21 14:07:30.063037 -0500 EST m=+0.201340417
request 2 2022-12-21 14:07:30.263033 -0500 EST m=+0.401337376
request 3 2022-12-21 14:07:30.463017 -0500 EST m=+0.601322501
request 4 2022-12-21 14:07:30.663012 -0500 EST m=+0.801318042
request 5 2022-12-21 14:07:30.863023 -0500 EST m=+1.001330376

For the second batch of requests we serve the first 3 immediately
because of the burstable rate limiting, then serve the remaining
2 with ~200ms delays each.

request 1 2022-12-21 14:07:30.863263 -0500 EST m=+1.001570667
request 2 2022-12-21 14:07:30.863277 -0500 EST m=+1.001584834
request 3 2022-12-21 14:07:30.863285 -0500 EST m=+1.001592376
request 4 2022-12-21 14:07:31.063349 -0500 EST m=+1.201657459
request 5 2022-12-21 14:07:31.26334 -0500 EST m=+1.401649459
*/
