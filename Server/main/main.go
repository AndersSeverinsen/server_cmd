package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type Locker struct {
	Userid    string `json:"userid"`
	Lockernum int    `json:"lockernum"`
	Lockerip  string `json:"lockerip"`
}

type BookingResponse struct {
	Locker          int  `json:"locker,omitempty"`
	ExistingBooking bool `json:"existingBooking"`
	FreeLocker      bool `json:"freeLocker"`
}

var lockers = make([]*Locker, 0)

func addLocker(ip string) {
	lockers = append(lockers, &Locker{Userid: "", Lockernum: len(lockers), Lockerip: ip})
}

func hasLocker(id string) (bool, int) {
	for _, locker := range lockers {
		if locker.Userid == id {
			return true, locker.Lockernum
		}
	}
	return false, -1
}

func initLockers() {
	for i := range lockers {
		lockers[i] = &Locker{Userid: "", Lockernum: i, Lockerip: ""}
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
			if lockers[i].Userid == "" {
				unlock(i, "Green")
				lockers[i].Userid = userid
				response = BookingResponse{
					ExistingBooking: false,
					Locker:          lockers[i].Lockernum,
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
		unlock(lockernum, "Green")
		lockers[lockernum].Userid = ""
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
		unlock(lockernum, "Red")
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

func unlock(lockerindex int, color string) {
	// Check if the locker is available
	for _, locker := range lockers {
		if locker.Lockernum == lockerindex {
			// If the locker is available, unlock it
			ip := locker.Lockerip

			// Send a POST request to the locker to unlock it
			url := "http://" + ip + ":8080/unlock" + color
			_, _ = http.Post(url, "application/json", nil)
		}
	}
}

func lockerStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	bytes, err := json.Marshal(lockers)
	if err != nil {
		http.Error(w, "Unable to marshal locker data", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(bytes); err != nil {
		fmt.Printf("Failed to write response: %v\n", err)
	}
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
	//initLockers()
	// Prompt the user to add a locker in a while loop
	go addLockerPrompt()
	http.HandleFunc("/book/", bookHandler)
	http.HandleFunc("/cancelBooking/", cancelHandler)
	http.HandleFunc("/keepBooking/", keepHandler)
	http.HandleFunc("/lockerStatus/", lockerStatus)
	http.HandleFunc("/unlock/", unlockHandler)

	// Start the server on port 8080
	fmt.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
