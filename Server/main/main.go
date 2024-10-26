package main

import (
	"encoding/json"
	"fmt"
	"net/http"
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

func addLocker(ip string) {
	lockers = append(lockers, &Locker{userid: "", lockernum: len(lockers), lockerip: ip})
}

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
				unlock(i)
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
		unlock(lockernum)
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
		unlock(lockernum)
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

func unlock(lockerindex int) {
	// Check if the locker is available
	for _, locker := range lockers {
		if locker.lockernum == lockerindex {
			// If the locker is available, unlock it
			ip := locker.lockerip

			// Send a POST request to the locker to unlock it
			_, _ = http.Post("http://"+ip+":8080/unlock", "application/json", nil)
		}
	}
}

func lockerStatus(w http.ResponseWriter, r *http.Request) {
	bytes, _ := json.Marshal(lockers)
	w.Write(bytes)
}

func addLockerPrompt() {
	for {
		var ip string
		fmt.Print("Enter the IP address of the locker: ")
		fmt.Scanln(&ip)
		addLocker(ip)
		fmt.Println("Locker added successfully")
	}
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
	http.ListenAndServe(":8080", nil)

	// Prompt the user to add a locker in a while loop
	addLockerPrompt()
}
