package main

import (
	"fmt"
	"sync"
	"time"
)

// synchronous, will just print sheep
func startV1() {
	count("sheep")
	count("fish")
}

// go keyword creates a go routine and run first count function in the new go routine
// the second count runs on main go routine
func startV2() {
	go count("sheep")
	count("fish")
}

// first go keyword creates and run first count function in one go routine
// next go keyword creates and run second count function in second go routine
// the program doesn't print anything as main go routine completes, as there is
// nothing stopping from main go routine to complete and the program exits
func startV3() {
	go count("sheep")
	go count("fish")
}

// lets fix V3.
// first go keyword creates and run first count function in one go routine
// next go keyword creates and run second count function in second go routine
// the main go routine will wait for 2 seconds and then main go routine will complete
func startV4() {
	go count("sheep")
	go count("fish")
	time.Sleep(time.Millisecond * 2000)
}

// another way to fix V3
func startV5() {
	go count("sheep")
	go count("fish")
	fmt.Scanln()
}

func count(thing string) {
	for i := 1; true; i++ {
		fmt.Println(i, thing)
		time.Sleep(time.Millisecond * 500)
	}
}

// both V4 and V5 are not standard practices instead we should use waitgroup
// wg.Add function adds a go routine and wg.Done function removes a go routine in the wait group
// wg.Wait function will wait until there are no go routine in the wait group
func startV6() {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		countV1("sheep")
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		countV1("fish")
	}()

	wg.Wait()
}

func countV1(thing string) {
	for i := 1; i <= 5; i++ {
		fmt.Println(i, thing)
		time.Sleep(time.Millisecond * 500)
	}
}

// till now we were spinning go routines and letting them execute independently
// lets see how we can communicate between new go routine and main go routine
// this is done through channel,
// think of channel as pipeline between two go routine through which messages flows
// this pipeline also has a type which specify what kind of messages can be passed through pipeline
func startV7() {
	// create a pipeline which flows string messages
	c := make(chan string)
	// check countV2 implementation
	go countV2("sheep", c)

	// note that sending and receiving are blocking operation.
	// thus msg := <-c statement blocks the main routine execution and
	// it will be blocked until some other go routine writes on this channel, which is done
	// in countV2 implementation
	msg := <-c
	fmt.Println(msg)
}

// in V7 the code will read only one message, prints that message and then
// main go routine completes and the program exits
// to read all values you can wrap the read statement in for loop
func startV8() {
	c := make(chan string)
	go countV2("sheep", c)

	for {
		msg := <-c
		fmt.Println(msg)
	}
}

func countV2(thing string, c chan string) {
	for i := 1; i <= 5; i++ {
		// similarly here c <- thing statement is blocking and
		// it will not execute next statement until some other go routine reads on this channel
		// which is done in startV7, startV8 implementation
		c <- thing
		time.Sleep(time.Millisecond * 500)
	}
}

// when you execute V8, you will get deadlock in the output
// this is because when countV2 has finished writing 5th thing
// there is nothing to read at V8 at msg := <-c and since it is
// blocking it will block main routine execution and the program exits with deadlock message
// to resolve this issue,
// at the sender (remember sender) we need to close the channel
// at the receiver we need to check if channel if open, if not then break the loop
func startV9() {
	c := make(chan string)
	// check countV3 implementation, where we close the channel
	go countV3("sheep", c)

	for {
		msg, open := <-c
		if !open {
			break
		}
		fmt.Println(msg)
	}
}

// In GO, there is a slightly nicer way to check if channel is open
// by the usage of range operator.
func startV10() {
	c := make(chan string)
	go countV3("sheep", c)

	for msg := range c {
		fmt.Println(msg)
	}
}

func countV3(thing string, c chan string) {
	defer close(c)
	for i := 1; i <= 5; i++ {
		c <- thing
		time.Sleep(time.Millisecond * 500)
	}

}

// as name suggest this is a tricky function
// in here we are reading from channel at msg := <-c which will block main routine
// as reading/writing from channel are blocking operation in go lang, thus it will block
// execution of main routine until some other routine writes on channel which never
// happens and thus the program will terminate with deadlock!
func tricky() {
	c := make(chan string)

	c <- "Hello"
	msg := <-c

	fmt.Println(msg)
}

// to resolve above issue we can use buffered channel instead of un-buffered channel
// buffered channel gives capacity to the channel that means the channel read and write
// will be unblocking until the capacity is full. Once the capacity is full the read / write
// becomes blocking. Another way of understanding this is thinking that un-buffered channel is
// nothing but buffered channel with capacity one
func buffered() {
	// we can write twice to the channel until it becomes blocking
	c := make(chan string, 2)

	c <- "hello"
	c <- "hi"
	// following statement will block the code as write will become blocking since capacity
	// of channel is full and some other go routine needs to read from the channel so that capacity
	// can be freed
	// c <- "blocked"

	msg := <-c
	fmt.Println(msg)

	msg = <-c
	fmt.Println(msg)
}

// This is another tricky function.
// In this function since there is infinite for loop on main go routine
// which is waiting on c1 and c2. This wait is causing faster go routine with sleep
// 500ms is being delayed as main go routine waits(blocks itself) to read from c2
// thus even though we expect four `Every 500ms` we get single `Every 500ms`
// and then single `Every two seconds`
func trickyV2() {
	c1 := make(chan string)
	c2 := make(chan string)

	go func() {
		for {
			c1 <- "Every 500ms"
			time.Sleep(time.Millisecond * 500)
		}
	}()

	go func() {
		for {
			c2 <- "Every two seconds"
			time.Sleep(time.Millisecond * 2000)
		}
	}()

	for {
		fmt.Println(<-c1)
		fmt.Println(<-c2)
	}

}

// in trickyV2 the issue is we are waiting on channel c2 even though c1 is free the output
// is synchronous, to avoid this issue we can use `select` statement in go. The select
// statement allows us to read from whichever channel is free
func selectStatement() {
	c1 := make(chan string)
	c2 := make(chan string)

	go func() {
		for {
			c1 <- "Every 500ms"
			time.Sleep(time.Millisecond * 500)
		}
	}()

	go func() {
		for {
			c2 <- "Every two seconds"
			time.Sleep(time.Millisecond * 2000)
		}
	}()

	for {
		select {
		case msg1 := <-c1:
			fmt.Println(msg1)
		case msg2 := <-c2:
			fmt.Println(msg2)
		}
	}
}

// Worker pool pattern
// This is a common pattern in go lang multithreading. In this pattern we have queue of work to
// be done and we can assign it to multiple concurrent workers, where each worker is pulling out
// an item off the queue

// In the following code we are representing the queue of work(jobs) as fib output
// of integers from 1 to 100 and in the first attempt we have only worker that is pulling
// off one number and calculating fib value at that number
func workerPatternV1() {
	jobs := make(chan int, 100)
	result := make(chan int, 100)

	// spin single worker
	go worker_(jobs, result)

	// on job channel send integers from 0 to 99 representing sequence of work
	for i := 0; i < 100; i++ {
		jobs <- i
	}
	// don't forget to close the jobs channel
	close(jobs)

	// on result channel receive nth fib
	for j := 0; j < 100; j++ {
		fmt.Println(<-result)
	}

	// ps: note that since jobs and result are buffered channel they are not blocking
}

// this is the worker function. In the worker pool pattern, the worker function
// takes two channels, one channel of jobs to do and one channel to send result on.
// instead of specifying bidirectional channel we can restrict the jobs as receiving only
// channel and result as sending only channel
func worker_(jobs <-chan int, result chan<- int) {
	// jobs is going to queue of numbers, we will receive on the jobs channel
	// using range keyword
	for n := range jobs {
		// calculate nth fib number and send it to result channel
		result <- fib(n)
	}
}

// simple fibonacci implementation
func fib(n int) int {
	if n <= 1 {
		return 1
	}

	return fib(n-1) + fib(n-2)
}

// In the workerPatternV1, we used only one worker but modern pcs has multiple cores so lets
// utilize their capacity and spin up multiple workers.
// instead of spinning only 4 workers you can check capacity of your pc using
// runtime.GOMAXPROCS(runtime.NumCPU())
func workerPatternV2() {
	jobs := make(chan int, 100)
	result := make(chan int, 100)

	// spin four workers
	// do check your task manager / activity monitor after this
	// we are taking advantage of multi core processors
	go worker_(jobs, result)
	go worker_(jobs, result)
	go worker_(jobs, result)
	go worker_(jobs, result)

	// on job channel send integers from 0 to 99 representing sequence of work
	for i := 0; i < 100; i++ {
		jobs <- i
	}
	// don't forget to close the jobs channel
	close(jobs)

	// on result channel receive nth fib
	for j := 0; j < 100; j++ {
		output := <-result
		fmt.Println(output)
	}

	// ps: note that since jobs and result are buffered channel they are not blocking
}
