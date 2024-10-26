package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

type Locker struct {
	userid    string
	lockernum int
}

// Create a slice of pointers to integers to represent the lockers
var lockers = make([]*Locker, 10)

func hasLocker(id string) (bool, int) {
	for _, locker := range lockers {
		if locker.userid == id {
			return true, locker.lockernum
		}
	}
	return false, -1
}

func initLockers() {
	for i := range lockers {
		lockers[i] = &Locker{userid: "", lockernum: i}
	}
}

func book(userid string) string {
	fmt.Println("Lockers:", lockers)
	// Check if the user has any bookings
	hasLocker, lockernum := hasLocker(userid)
	if hasLocker {
		// If the user has a booking, return the booking
		return fmt.Sprintf("User %s already has locker %v", userid, lockernum)
	} else {
		// If the user has no bookings, check if a locker is available
		for _, locker := range lockers {
			if locker.userid == "" {
				// If a locker is available, book it for the user
				locker.userid = userid
				return fmt.Sprintf("User %s has no booking, so booking locker %d", userid, locker.lockernum)
			}
		}
		// If no lockers are available, return a message
		return fmt.Sprintf("User %s has no booking, and no lockers are available", userid)
	}
}

func bookHandler(w http.ResponseWriter, r *http.Request) {
	// Only accept POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
		return
	}
	// Parse URL path for "book/{userid}"
	userid := strings.TrimPrefix(r.URL.Path, "/book/")

	// Call the book function
	response := book(userid)

	// Write response to the client
	w.Write([]byte(response))
}

func cancel(userid string) string {
	// Check if the user has any bookings
	hasLocker, lockernum := hasLocker(userid)
	if hasLocker {
		// If the user has a booking, cancel it
		lockers[lockernum].userid = ""
		return fmt.Sprintf("User %s has cancelled the booking for locker %d", userid, lockernum)
	} else {
		// If the user has no bookings, return a message
		return fmt.Sprintf("User %s has no booking", userid)
	}
}

func cancelHandler(w http.ResponseWriter, r *http.Request) {
	// Only accept POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
		return
	}
	// Parse URL path for "cancelBooking/{userid}"
	userid := strings.TrimPrefix(r.URL.Path, "/cancelBooking/")

	// Call the cancel function
	response := cancel(userid)

	// Write response to the client
	w.Write([]byte(response))
}

func keep(userid string) string {
	// Check if the user has any bookings
	hasLocker, lockernum := hasLocker(userid)
	if hasLocker {
		// If the user has a booking, return the booking
		return fmt.Sprintf("User %s has kept the booking for locker %v", userid, lockernum)
	} else {
		// If the user has no bookings, return a message
		return fmt.Sprintf("User %s has no booking", userid)
	}
}

func keepHandler(w http.ResponseWriter, r *http.Request) {
	// Only accept POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
		return
	}
	// Parse URL path for "keepBooking/{userid}"
	userid := strings.TrimPrefix(r.URL.Path, "/keepBooking/")

	// Call the keep function
	response := keep(userid)

	// Write response to the client
	w.Write([]byte(response))
}

func main() {
	initLockers()
	// Start the server using mux as the root handler
	http.HandleFunc("/book/", bookHandler)

	http.HandleFunc("/cancelBooking/", cancelHandler)

	http.HandleFunc("/keepBooking/", keepHandler)

	// Start the server on port 8080
	//fmt.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
