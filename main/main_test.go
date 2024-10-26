package main

import (
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"
)

func assertEqual(t *testing.T, actual, expected interface{}) {
	if expected == nil && actual == nil {
		return
	}

	if expected != actual {
		t.Errorf("\nExpected: %v\nBut got: %v\n", expected, actual)
	}
}

func TestBooking(t *testing.T) {
	fmt.Println("TestBooking")
	// Call main function
	go main()
	// Delay 1 millisecond to allow the server to start
	time.Sleep(1 * time.Millisecond)

	// Create a new request
	request, err := http.NewRequest(http.MethodPost, "http://localhost:8080/book/1", nil)
	if err != nil {
		fmt.Println(err)
	}

	res, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		fmt.Printf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}
	if err != nil {
		fmt.Println(err)
	}
	bodystr := string(body[:])
	assertEqual(t, bodystr, "User 1 has no booking, so booking locker 0")
}

func TestBookingWhenUserAlreadyHasBooking(t *testing.T) {
	fmt.Println("TestBookingWhenUserAlreadyHasBooking")
	// Call main function
	go main()
	// Delay 1 millisecond to allow the server to start
	time.Sleep(1 * time.Millisecond)

	// Create a new request
	request, err := http.NewRequest(http.MethodPost, "http://localhost:8080/book/1", nil)
	if err != nil {
		fmt.Println(err)
	}
	// Send the request
	_, err1 := http.DefaultClient.Do(request)
	if err1 != nil {
		fmt.Printf("client: error making http request: %s\n", err)
	}
	time.Sleep(1 * time.Millisecond)
	// Send the request again
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
	}
	// Read the response and handle errors
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		fmt.Printf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}
	if err != nil {
		fmt.Println(err)
	}
	bodystr := string(body[:])
	assertEqual(t, bodystr, "User 1 already has the following locker: 0")
}

func TestBookingWhenUserHasNoBookingAndNoLockerAvailable(t *testing.T){
	fmt.Println("TestBookingWhenUserHasNoBookingAndNoLockerAvailable")
	// Call main function
	go main()

	// Delay 1 millisecond to allow the server to start
	time.Sleep(1 * time.Millisecond)
	
	
	// Create a new request
	request, err := http.NewRequest(http.MethodPost, "http://localhost:8080/book/1", nil)
	if err != nil {
		fmt.Println(err)
	}
	// Send the request
	res, err1 := http.DefaultClient.Do(request)
	if err1 != nil {
		fmt.Printf("client: error making http request: %s\n", err)
	}
	time.Sleep(1 * time.Millisecond)
	// Read the response and handle errors

	body, err := io.ReadAll(res.Body)
	res.Body.Close()

	if res.StatusCode > 299 {
		fmt.Printf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}
	if err != nil {
		fmt.Println(err)
	}
	bodystr := string(body[:])
	assertEqual(t, bodystr, "User 1 has no booking, and no lockers are available")


}
