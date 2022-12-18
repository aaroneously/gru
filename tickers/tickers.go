package main

import (
	"fmt"
	"time"
)

/*
Timers are for when you want to do something once in
the future - tickers are for when you want to do
something repeatedly at regular intervals. Here’s an
example of a ticker that ticks periodically until we
stop it.
*/
func main() {
	/*
		Tickers use a similar mechanism to timers: a channel
		that is sent values. Here we’ll use the select builtin
		on the channel to await the values as they arrive
		every 500ms.
	*/
	ticker := time.NewTicker(500 * time.Millisecond)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				fmt.Println("Tick at", t)
			}
		}
	}()

	/*
		Tickers can be stopped like timers. Once a ticker is
		stopped it won’t receive any more values on its channel.
		We’ll stop ours after 1600ms.
	*/
	time.Sleep(1600 * time.Millisecond)
	ticker.Stop()
	done <- true
	fmt.Println("Ticker stopped")
}

/*
When we run this program the ticker should tick 3 times before
we stop it.
*/
/*
➜  gru git:(main) ✗ go run tickers/tickers.go
Tick at 2022-12-18 14:41:39.418795 -0500 EST m=+0.501255376
Tick at 2022-12-18 14:41:39.918602 -0500 EST m=+1.001048751
Tick at 2022-12-18 14:41:40.418788 -0500 EST m=+1.501220709
Ticker stopped

*/
