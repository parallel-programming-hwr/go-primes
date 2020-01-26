package main

import (
	"fmt"
	"os"
	"bufio"
	"strconv"
	"runtime"
)

var ch chan uint64

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func main() {
	ch = make(chan uint64)
	numThreads := runtime.NumCPU()
	fmt.Printf("Starting to calculate primes\n")
	f, err := os.OpenFile("primes.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	check(err)
	defer f.Close()
	w := bufio.NewWriter(f)
	defer w.Flush()
	start := uint64(1)
	if len(os.Args) >= 2 {
		start, err = strconv.ParseUint(os.Args[1], 10, 64)
		check(err)
	}
	if start % 2 == 0 {
		start += 1
	}
	for i := 0; i < numThreads; i++ {
		go getPrimes(start + uint64(i * 2), uint64(numThreads * 2))
	}
	for {
		prime := <- ch
		fmt.Printf("%d\n", prime)
		w.WriteString(strconv.FormatUint(prime, 10) + "\n")
	}
}

// coroutine that calculates primes and adds it to the channel
func getPrimes(start, incr uint64) {
	num := uint64(start)
	for {
		isPrime := true
		if (num < 3 || num % 2 == 0) {
			isPrime = false
			num += incr
			continue
		}
		var i uint64
		for i = 3; i < num/2; i+=2 {
			if num % i == 0 {
				isPrime = false
				break
			}
		}
		if isPrime {
			ch <- num
		}
		num += incr
	}
}