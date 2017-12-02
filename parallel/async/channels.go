package main

import (
	"fmt"
	"time"
)

//
// Demo 1
//

func d1() {
	messages := make(chan string)
	go func() { messages <- "ping" }()
	msg := <-messages
	fmt.Println(msg)
}

//
// Demo 2: Multiple messages
//

func d2() {
	messages := make(chan string, 2)
	messages <- "one"
	messages <- "two"
	fmt.Println(<-messages)
	fmt.Println(<-messages)
}

//
// Demo 3: Channel directions
//

func ping(pings chan<- string, msg string) {
	pings <- msg
}

func pong(pings <-chan string, pongs chan<- string) {
	msg := <-pings
	pongs <- msg
}

func d3() {
	pings := make(chan string, 1)
	pongs := make(chan string, 1)
	ping(pings, "ping-pong message")
	pong(pings, pongs)
	fmt.Println(<-pongs)
}

//
// Demo 4: Select
//

func d4() {
	c1 := make(chan string)
	c2 := make(chan string)

	go func() {
		time.Sleep(time.Second * 1)
		c1 <- "one sec"
	}()
	go func() {
		time.Sleep(time.Second * 2)
		c1 <- "two sec"
	}()

	fmt.Printf("[%s] Before select\n", time.Now().String())
	for i := 0; i < 2; i++ {
		select {
		case msg1 := <-c1:
			fmt.Printf("[%s] received=%s\n", time.Now(), msg1)
		case msg2 := <-c2:
			fmt.Printf("[%s] received=%s\n", time.Now(), msg2)
		}
	}
}

//
// Demo 5: Timeouts
//

func d5() {
	c1 := make(chan string, 1)
	go func() {
		time.Sleep(time.Second * 2)
		c1 <- "result #1"
	}()

	select {
	case res := <-c1:
		fmt.Println("SHOULD NOT HAPPEN", res)
	case <-time.After(time.Second * 1):
		fmt.Println("timeout 1")
	}

	c2 := make(chan string, 1)
	go func() {
		time.Sleep(time.Second * 2)
		c2 <- "result #2"
	}()

	select {
	case res := <-c2:
		fmt.Println(res)
	case <-time.After(time.Second * 3):
		fmt.Println("SHOULD NOT HAPPEN: timeout 2")
	}
}

//
// Demo 6: Non-Blocking Channel Operations
// https://gobyexample.com/non-blocking-channel-operations

func d6() {
	messages := make(chan string)
	signals := make(chan bool)

	// non-blocking receive
	select {
	case msg := <-messages:
		fmt.Println("[1] received message", msg)
	default:
		fmt.Println("[1] no message received")
	}

	// non-blocking send
	msg := "hi"
	select {
	case messages <- msg:
		fmt.Println("[2] sent message", msg)
	default:
		fmt.Println("[2] no message sent")
	}

	// non-blocking multi-receive cases
	select {
	case msg := <-messages:
		fmt.Println("[3] received message", msg)
	case sig := <-signals:
		fmt.Println("[3] received signel", sig)
	default:
		fmt.Println("[3] no activity")
	}

	// try to receive msg
	/*select {
	case messages <- msg:
		fmt.Println("[4] blocking (timeout-based) - sent message", msg)
	case <-time.After(time.Second * 1):
		fmt.Println("[4] SHOULD NOT HAPPEN: no msg sent")
	}
	select {
	case msg := <-messages:
		fmt.Println("[5] blocking (timeout-based) receive", msg)
	case <-time.After(time.Second * 1):
		fmt.Println("[5] SHOULD NOT HAPPEN: no msg received")
	}*/
}

//
// Demo 7: Closing Channels
// https://gobyexample.com/range-over-channels

func d7() {
	jobs := make(chan int, 5)
	done := make(chan bool)

	go func() {
		for {
			j, more := <-jobs
			if more {
				fmt.Println("received job", j)
			} else {
				fmt.Println("received all jobs")
				done <- true
				return
			}
		}
	}()

	for j := 0; j < 3; j++ {
		jobs <- (j + 1)
		fmt.Println("sent job", j)
	}
	close(jobs)
	fmt.Println("sent all jobs")

	<-done
}

func d8() {
	queue := make(chan string, 2)
	queue <- "one"
	queue <- "two"
	close(queue)

	for e := range queue {
		fmt.Println("enumerated element:", e)
	}
}

//
// Entry Point
//

type demofunc func()

func main() {
	demos := []demofunc{
		d1,
		d2,
		d3,
		d4,
		d5,
		d6,
		d7,
		d8,
	}
	for index, df := range demos {
		fmt.Printf("Demo #%d:\n\n", index+1)
		df()
		fmt.Println("---")
	}
}
