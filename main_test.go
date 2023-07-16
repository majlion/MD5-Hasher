package main

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
)

func TestWorker(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, world!")
	}))
	defer server.Close()

	results := make(chan RequestResult, 1)
	var wg sync.WaitGroup
	wg.Add(1)
	go Worker(server.URL, results, &wg)
	wg.Wait()
	close(results)

	result := <-results
	expectedHash := fmt.Sprintf("%x", md5.Sum([]byte("Hello, world!")))
	if result.Address != server.URL || result.Hash != expectedHash {
		t.Errorf("Worker did not produce the expected result.")
	}
}

func TestMainFunction(t *testing.T) {
	testAddresses := []string{
		"http://example.com",
		"http://google.com",
		"http://github.com",
	}

	mockedResponse := "Mocked response"
	mockedHash := fmt.Sprintf("%x", md5.Sum([]byte(mockedResponse)))
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, mockedResponse)
	}))
	defer mockServer.Close()

	os.Args = []string{"./myapp", "-parallel", "10"}
	os.Args = append(os.Args, testAddresses...)

	output := captureOutput(func() {
		main()
	})

	for _, address := range testAddresses {
		expectedResult := fmt.Sprintf("Address: %s, Hash: %s\n", address, mockedHash)
		fmt.Println("address :", address, " expected:", expectedResult)

		if !strings.Contains(output, expectedResult) {
			t.Errorf("Main function did not produce the expected result for address: %s", address)
		}
	}
}

// Helper function
func captureOutput(fn func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	fn()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)

	return buf.String()
}
