package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"io"
	"net/http"
	"sync"
)

// RequestResult consists the address and MD5 hash of the response
type RequestResult struct {
	Address string
	Hash    string
}

// Worker performs the HTTP request and computes the MD5 hash of the response
func Worker(address string, results chan<- RequestResult, wg *sync.WaitGroup) {
	defer wg.Done()

	response, err := http.Get(address)
	if err != nil {
		fmt.Printf("Error requesting %s: %s\n", address, err)
		return
	}
	defer response.Body.Close()

	hash := md5.New()
	_, err = io.Copy(hash, response.Body)
	if err != nil {
		fmt.Printf("Error computing hash for %s: %s\n", address, err)
		return
	}

	result := RequestResult{
		Address: address,
		Hash:    fmt.Sprintf("%x", hash.Sum(nil)),
	}
	results <- result
}

func main() {
	parallelLimit := flag.Int("parallel", 10, "number of parallel requests")
	flag.Parse()

	addresses := flag.Args()
	if len(addresses) == 0 {
		fmt.Println("No addresses provided.")
		return
	}

	// Create a buffered channel for the results
	results := make(chan RequestResult, len(addresses))

	// WaitGroup to wait for all workers to finish
	var wg sync.WaitGroup

	// Start workers
	workerSemaphore := make(chan struct{}, *parallelLimit)
	for _, address := range addresses {
		wg.Add(1)
		workerSemaphore <- struct{}{}
		go func(addr string) {
			defer func() { <-workerSemaphore }()
			Worker(addr, results, &wg)
		}(address)
	}

	wg.Wait()
	close(results)

	for result := range results {
		fmt.Printf("Address: %s, Hash: %s\n", result.Address, result.Hash)
	}
}
