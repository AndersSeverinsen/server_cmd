package main

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

func assertEqual(t *testing.T, actual, expected interface{}) {
	if expected == nil && actual == nil {
		return
	}

	if expected != actual {
		t.Errorf("Expected %v but got %v", expected, actual)
	}
}

func TestBooking(t *testing.T) {
	fmt.Println("TestBooking")
	// Call main function
	go main()
	// Delay 1 second to allow the server to start
	time.Sleep(1 * time.Second)

	// Create a new request
	request, err := http.NewRequest(http.MethodPost, "http://localhost:8080/book/1", nil)
	if err != nil {
		fmt.Println(err)
	}

	res, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
	}
	fmt.Println(res)

	assertEqual(t, res, "User 1 has no bookings. No lockers are available")
}
