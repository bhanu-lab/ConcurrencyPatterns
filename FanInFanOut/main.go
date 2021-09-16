package main

import (
	"fmt"
	"sync"
)

// FanInFanOut concurrency pattern
func main() {
	// take input of slice which are of any data type for example int
	nums := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	// feed all the data into a single channel
	fout := genNums(nums)
	//					   						___  channel 1
	// Fanning out channel __ /___  channel 2
	//					 					  \____ channel 3
	// Pass that single channel as input param and spin a go routine and try reading values from that routine
	// once reading values is done close the resulting channel where you are sending in
	// it can be any number of workers. This is FANOUT pattern.
	c1 := multiplexer(fout)
	c2 := multiplexer(fout)

	//	Chan1 ___
	//	Chan2 ___	\ ____ Chan
	//	Chan3 ___	/
	// Take multiple channels and combine all results into one channel.
	// Make fanIn as variadic function with chan int
	// take all channels which has to be combined.
	// range over all channels for each channel spin a go routine and range over all values fed to channel
	// write it to same channel across all go routines this is important.
	// create a wait group, wait for all routines and close the output channel.
	res := demultiplexer(c1, c2)

	for x := range res {
		fmt.Println("value of x is", x)
	}
}

func demultiplexer(channels ...chan int) chan int {
	var wg sync.WaitGroup
	out := make(chan int)
	for _, c := range channels {
		// spin a go routine for every channel so that all channels can
		// send data to same channel
		wg.Add(1)
		go func(m chan int) {
			defer wg.Done()
			for k := range m {
				out <- k
			}
		}(c)
	}
	// wg.Wait has to be in separate routine becuase since we are sending to out channel
	// and no one will be reading out from channel if wg.Wait is not inside go routine
	// as wg.Wait will block main routine until all go routines are completed wg.wait
	// has to be in separate go routine so that it is read from out channel other wise all routines will get blocked
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

// multiplexer
// invoke multiple times to have multiple workers to read from
// from cnum channel and send to a new channel causing Fanning Out
func multiplexer(cnum chan int) chan int {
	out := make(chan int)
	// spin a go routine and try reading from that channel
	go func(chNum chan int) {
		for i := range chNum {
			out <- i
		}
		close(out)
	}(cnum)
	return out
}

//creates a channel with all input data fed into it
func genNums(nums []int) chan int {
	out := make(chan int)
	go func() {
		for _, i := range nums {
			out <- i
		}
		close(out)
	}()
	return out
}
