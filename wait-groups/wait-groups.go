package main

import (
	"fmt"
	"sync"
	"time"
)

/*
To wait for multiple goroutines to finish,
we can use a wait group.
*/

/*
This is the function we’ll run in every
goroutine.
*/
func worker(id int) {
	fmt.Printf("Worker %d starting\n", id)

	/*
		Sleep to simulate an expensive task.
	*/
	time.Sleep(time.Second)
	fmt.Printf("Worker %d done\n", id)
}

func main() {

	/*
		This WaitGroup is used to wait for all the
		goroutines launched here to finish. Note: if
		a WaitGroup is explicitly passed into
		functions, it should be done by pointer.
	*/
	var wg sync.WaitGroup

	/*
		Launch several goroutines and increment the
		WaitGroup counter for each.
	*/
	for i := 1; i <= 5; i++ {
		wg.Add(1)

		/*
			Avoid re-use of the same i value in each
			goroutine closure.
		*/
		i := i

		/*
			Wrap the worker call in a closure that makes
			sure to tell the WaitGroup that this worker
			is done. This way the worker itself does not
			have to be aware of the concurrency primitives
			involved in its execution.
		*/
		go func() {
			defer wg.Done()
			worker(i)
		}()
	}

	/*
		Block until the WaitGroup counter goes back to 0;
		all the workers notified they’re done.
	*/
	wg.Wait()

	/*
		Note that this approach has no straightforward
		way to propagate errors from workers.
	*/
}

/*
➜  gru git:(main) ✗ go run wait-groups/wait-groups.go
Worker 2 starting
Worker 5 starting
Worker 3 starting
Worker 1 starting
Worker 4 starting
Worker 1 done
Worker 3 done
Worker 2 done
Worker 4 done
Worker 5 done
*/
