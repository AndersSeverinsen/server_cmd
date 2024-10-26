package main

import (
	"encoding/json"
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
	// Send the request
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
	assertEqual(t, bodystr, "{\"existingBooking\":false,\"freeLocker\":true}")

	// Check locker status
	assertEqual(t, lockers[0].Userid, "1")
}

func TestBookingWhenUserAlreadyHasBooking(t *testing.T) {
	fmt.Println("TestBookingWhenUserAlreadyHasBooking")
	// Call main function
	go main()
	// Delay 1 millisecond to allow the server to start
	time.Sleep(1 * time.Millisecond)

	// Create a new request
	request, err := http.NewRequest(http.MethodPost, "http://localhost:8080/book/2", nil)
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
	assertEqual(t, bodystr, "{\"existingBooking\":true,\"freeLocker\":true}")
}

func TestBookingWhenUserHasNoBookingAndNoLockerAvailable(t *testing.T) {
	fmt.Println("TestBookingWhenUserHasNoBookingAndNoLockerAvailable")
	// Call main function
	go main()

	// Delay 1 millisecond to allow the server to start
	time.Sleep(1 * time.Millisecond)

	// Set all lockers to be occupied by user 123
	for i := range lockers {
		lockers[i] = &Locker{Userid: "123", Lockernum: i}
	}

	// Create a new request
	request, err := http.NewRequest(http.MethodPost, "http://localhost:8080/book/3", nil)
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
	assertEqual(t, bodystr, "{\"existingBooking\":false,\"freeLocker\":false}")
}

func TestBookingWhenUserHasBookingAndCancels(t *testing.T) {
	fmt.Println("TestBookingWhenUserHasBookingAndCancels")
	// Call main function
	go main()

	// Delay 1 millisecond to allow the server to start
	time.Sleep(1 * time.Millisecond)

	// Create a new request
	request, err := http.NewRequest(http.MethodPost, "http://localhost:8080/book/4", nil)
	if err != nil {
		fmt.Println(err)
	}
	// Send the request
	_, err1 := http.DefaultClient.Do(request)
	if err1 != nil {
		fmt.Printf("client: error making http request: %s\n", err)
	}
	time.Sleep(1 * time.Millisecond)
	// Create a new request
	request, err = http.NewRequest(http.MethodPost, "http://localhost:8080/cancelBooking/4", nil)
	if err != nil {
		fmt.Println(err)
	}
	// Send the request
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
	assertEqual(t, bodystr, "{\"message\": \"User 4 has cancelled the booking for locker 0\"}")
}

func TestBookingWhenUserHasBookingAndKeeps(t *testing.T) {
	fmt.Println("TestBookingWhenUserHasBookingAndKeeps")
	// Call main function
	go main()

	// Delay 1 millisecond to allow the server to start
	time.Sleep(1 * time.Millisecond)

	// Create a new request
	request, err := http.NewRequest(http.MethodPost, "http://localhost:8080/book/5", nil)
	if err != nil {
		fmt.Println(err)
	}
	// Send the request
	_, err1 := http.DefaultClient.Do(request)
	if err1 != nil {
		fmt.Printf("client: error making http request: %s\n", err)
	}
	time.Sleep(1 * time.Millisecond)
	// Create a new request
	request, err = http.NewRequest(http.MethodPost, "http://localhost:8080/keepBooking/5", nil)
	if err != nil {
		fmt.Println(err)
	}
	// Send the request
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
	assertEqual(t, bodystr, "{\"message\": \"User 5 has kept the booking for locker 0\"}")
}

func TestAddLocker(t *testing.T) {
	fmt.Println("TestAddLocker")
	// Call main function
	go main()

	// Delay 1 millisecond to allow the server to start
	time.Sleep(1 * time.Millisecond)

	addLocker("192.123.1.100")

	assertEqual(t, len(lockers), 11)
}

func TestLockerStatus(t *testing.T) {
	fmt.Println("TestLockerStatus")
	// Call main function
	go main()

	// Delay 1 millisecond to allow the server to start
	time.Sleep(10 * time.Millisecond)

	// Create a new request
	request, err := http.NewRequest(http.MethodGet, "http://localhost:8080/statusLocker/", nil)
	if err != nil {
		fmt.Println(err)
	}
	// Send the request
	res, err1 := http.DefaultClient.Do(request)
	if err1 != nil {
		fmt.Printf("client: error making http request: %s\n", err)
	}
	time.Sleep(1 * time.Millisecond)
	// Create a new request
	lockers := make([]Locker, 10)
	body, err := io.ReadAll(res.Body)
	json.Unmarshal(body, &lockers)
	fmt.Println(lockers)

	res.Body.Close()

}
