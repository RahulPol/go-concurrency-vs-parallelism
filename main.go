package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"runtime"
	"runtime/trace"
	"sync"
)

type RunType string

const (
	Sequential RunType = "Sequential"
	Concurrent RunType = "Concurrent"
	Parallel   RunType = "Parallel"
)

func worker(id int, wg *sync.WaitGroup) {

	defer wg.Done()

	fmt.Printf("Worker %d starting\n", id)

	for number := math.MinInt32; number < math.MaxInt32; number++ {
	}
	fmt.Printf("Worker %d done\n", id)
}

func runGoRoutine(rt RunType) {

	if rt == Sequential {
		for i := 1; i <= 5; i++ {
			for number := math.MinInt32; number < math.MaxInt32; number++ {
			}
		}
		return
	}

	if rt == Parallel {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	if rt == Concurrent {
		runtime.GOMAXPROCS(1)
	}

	var wg sync.WaitGroup

	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go worker(i, &wg)
	}

	fmt.Println("Waiting To Finish")
	wg.Wait()

	fmt.Println("\nTerminating Program")
}

func getFileName(rt RunType) string {
	switch rt {
	case Sequential:
		return "trace-sequential.out"
	case Concurrent:
		return "trace-concurrent.out"
	case Parallel:
		return "trace-parallel.out"
	}
	return "trace"
}

func main() {

	reader := bufio.NewReader(os.Stdin)
	var rt RunType

	fmt.Printf("Choose:\n1.Sequential\n2.Concurrent\n3.Parallel\n")

	// use following to read a character
	char, _, err := reader.ReadRune()
	if err != nil {
		fmt.Println("Error occurred: ", err)
	}
	//clear the input stream
	reader.Reset(os.Stdin)

	switch char {
	case '1':
		rt = Sequential
	case '2':
		rt = Concurrent
	case '3':
		rt = Parallel
	default:
		fmt.Println("Learn to read, stupid!!")
		return
	}

	fileName := getFileName(rt)

	f, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("failed to create trace output file: %v", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Fatalf("failed to close trace file: %v", err)
		}
	}()

	fmt.Printf("user input: %c\n", char)
	if err := trace.Start(f); err != nil {
		log.Fatalf("failed to start trace: %v", err)
	}
	defer trace.Stop()

	runGoRoutine(rt)

	// Basics of multithreading
	// Multithreading
	// startV1()
	// startV2()
	// startV3()
	// startV4()
	// startV5()
	// startV6()
	// startV7()
	// startV8()
	// startV9()
	// startV10()
	// tricky()
	// buffered()
	// trickyV2()
	// selectStatement()
	// workerPatternV2()
}
