package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Locker struct {
	userid    string
	lockernum int
	lockerip  string
}

type BookingResponse struct {
	Locker          int  `json:"locker,omitempty"`
	ExistingBooking bool `json:"existingBooking"`
	FreeLocker      bool `json:"freeLocker"`
}

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
		lockers[i] = &Locker{userid: "", lockernum: i, lockerip: ""}
	}
}

func book(userid string) string {
	hasBooking, lockernum := hasLocker(userid)
	var response BookingResponse

	if hasBooking {
		response = BookingResponse{
			ExistingBooking: true,
			Locker:          lockernum,
			FreeLocker:      true,
		}
	} else {
		for i := range lockers {
			if lockers[i].userid == "" {
				lockers[i].userid = userid
				response = BookingResponse{
					ExistingBooking: false,
					Locker:          lockers[i].lockernum,
					FreeLocker:      true,
				}
				break
			}
		}
		if !response.FreeLocker {
			response = BookingResponse{
				ExistingBooking: false,
				FreeLocker:      false,
			}
		}
	}

	jsonResponse, _ := json.Marshal(response)
	return string(jsonResponse)
}

func bookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "Only POST requests are allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	userid := strings.TrimPrefix(r.URL.Path, "/book/")
	response := book(userid)
	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write([]byte(response))
	if err != nil {
		return
	}
}

func cancel(userid string) string {
	hasBooking, lockernum := hasLocker(userid)
	if hasBooking {
		lockers[lockernum].userid = ""
		return fmt.Sprintf(`{"message": "User %s has cancelled the booking for locker %d"}`, userid, lockernum)
	}
	return fmt.Sprintf(`{"error": "User %s has no booking"}`, userid)
}

func cancelHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "Only POST requests are allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	userid := strings.TrimPrefix(r.URL.Path, "/cancelBooking/")
	response := cancel(userid)
	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write([]byte(response))
	if err != nil {
		return
	}
}

func keep(userid string) string {
	hasBooking, lockernum := hasLocker(userid)
	if hasBooking {
		return fmt.Sprintf(`{"message": "User %s has kept the booking for locker %d"}`, userid, lockernum)
	}
	return fmt.Sprintf(`{"error": "User %s has no booking"}`, userid)
}

func keepHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "Only POST requests are allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	userid := strings.TrimPrefix(r.URL.Path, "/keepBooking/")
	response := keep(userid)
	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write([]byte(response))
	if err != nil {
		return
	}
}

func unlockHandler(w http.ResponseWriter, r *http.Request) {
	// Only accept POST requests
	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "Only POST requests are allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	// Parse URL path for "unlock/{lockerindex}"
	lockerindex := strings.TrimPrefix(r.URL.Path, "/unlock/")

	// Call the cancel function
	response := cancel(lockerindex)
	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write([]byte(response))
	if err != nil {
		return
	}
}

func unlock(lockerindex string) string {
	// Convert lockerindex to an integer
	lockernum, _ := strconv.Atoi(lockerindex)

	// Check if the locker is available
	for _, locker := range lockers {
		if locker.lockernum == lockernum {
			// If the locker is available, unlock it
			ip := locker.lockerip

			// Send a POST request to the locker to unlock it
			_, err := http.Post("http://"+ip+":8080/unlock", "application/json", nil)
			if err != nil {
				return fmt.Sprintf("Error unlocking locker %d", lockernum)
			}
			return fmt.Sprintf("Unlocked locker %d", lockernum)
		}
	}
	return fmt.Sprintf("Locker %d is not initialized", lockernum)
}

func lockerStatus(w http.ResponseWriter, r *http.Request) {
	// Only accept POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
		return
	}
	// Parse URL path for "unlock/{userid}"
	lockerindex := strings.TrimPrefix(r.URL.Path, "/lockerStatus/")

	if lockerindex == "" {
		bytes, _ := json.Marshal(lockers)
		w.Write(bytes)
	}

	// Call the unlock function
	response := unlock(lockerindex)

	// Write response to the client
	w.Write([]byte(response))
}

func main() {
	initLockers()
	http.HandleFunc("/book/", bookHandler)
	http.HandleFunc("/cancelBooking/", cancelHandler)
	http.HandleFunc("/keepBooking/", keepHandler)
	http.HandleFunc("/lockerStatus/", lockerStatus)
	http.HandleFunc("/unlock/", unlockHandler)

	// Start the server on port 8080
	fmt.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
