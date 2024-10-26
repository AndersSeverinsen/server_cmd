package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type Locker struct {
	userid    string
	lockernum int
}

type BookingResponse struct {
	Locker          int  `json:"locker,omitempty"`
	ExistingBooking bool `json:"existingBooking"`
	FreeLocker      bool `json:"freeLocker"`
}

var lockers = make([]*Locker, 10)

func initLockers() {
	for i := range lockers {
		lockers[i] = &Locker{userid: "", lockernum: i + 1} // lockernum starts from 1
	}
}

func hasLocker(id string) (bool, int) {
	for _, locker := range lockers {
		if locker.userid == id {
			return true, locker.lockernum
		}
	}
	return false, -1
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

func main() {
	initLockers()
	http.HandleFunc("/book/", bookHandler)
	http.HandleFunc("/cancelBooking/", cancelHandler)
	http.HandleFunc("/keepBooking/", keepHandler)
	fmt.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
