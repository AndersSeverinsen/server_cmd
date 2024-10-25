package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Create a slice of pointers to integers to represent the lockers
var lockers = make([]*int, 10)

// Create a map to store the user's booking
var users = make(map[int]int)

func book(userid int) string {
	// Check if the user has any bookings
	if users[userid] != 0 {
		// If the user has a booking, return the booking
		return fmt.Sprintf("User %d has the following booking: %v", userid, users[userid])
	} else {
		// If the user has no bookings, check if a locker is available
		for i, locker := range lockers {
			if locker == nil {
				// If a locker is available, book it for the user
				lockers[i] = &userid
				users[userid] = i
				return fmt.Sprintf("User %d has booked locker %d", userid, i)
			}
		}
		// If no lockers are available, return a message
		return fmt.Sprintf("User %d has no bookings. No lockers are available", userid)
	}
}

func bookHandler(w http.ResponseWriter, r *http.Request) {
	// Only accept POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
		return
	}
	// Parse URL path for "book/{id}/{value}"
	path := strings.TrimPrefix(r.URL.Path, "/book/")
	parts := strings.Split(path, "/")
	if len(parts) != 2 {
		http.Error(w, "Invalid URL format. Expected /book/{userid}", http.StatusBadRequest)
		return
	}
	// Convert the ID and value from strings to integers
	userid, err1 := strconv.Atoi(parts[0])
	if err1 != nil {
		http.Error(w, "UserID should be an integer", http.StatusBadRequest)
		return
	}
	// Call the book function
	response := book(userid)

	// Write response to the client
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}

func main() {
	// Start the server using mux as the root handler
	http.HandleFunc("/book/", bookHandler)

	// Start the server on port 8080
	fmt.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
