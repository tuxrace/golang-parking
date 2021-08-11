package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type FuelInfo struct {
	Capacity float64 `json:"capacity"`
	Level    float64 `json:"level"`
}

type Parking struct {
	LicencePlate string   `json:"licencePlate"`
	Size         string   `json:"size"`
	Fuel         FuelInfo `json:"fuel"`
}

type Assignment struct {
	LicencePlate string  `json:"licencePlate"`
	Employee     string  `json:"employee"`
	FuelAdded    float64 `json:"fuelAdded"`
	Price        float64 `json:"price"`
}

var SMALL = "small"
var PARK_RATE_SMALL float64 = 25
var PARK_RATE_LARGE float64 = 35
var FUEL_RATE = 1.75
var EMPLOYEE_A = "Employee A"
var EMPLOYEE_B = "Employee B"

func AssignEmployee(size string, fuelAdded float64) string {
	if size != SMALL && fuelAdded > 0 {
		return EMPLOYEE_B
	}
	return EMPLOYEE_A
}

func GetParkingRate(size string) float64 {
	if size == SMALL {
		return PARK_RATE_SMALL
	}
	return PARK_RATE_LARGE
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Home")
}

func ParkingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var assignments []Assignment
	var items []Parking
	_ = json.NewDecoder(r.Body).Decode(&items)

	for i := 0; i < len(items); i++ {
		var litersRefueled float64 = 0
		var price float64 = 0
		level := items[i].Fuel.Level
		capacity := items[i].Fuel.Capacity
		licencePlate := items[i].LicencePlate
		size := items[i].Size

		if level < 0.1 {
			litersLeft := ((level * 100) / 100) * capacity
			litersRefueled = capacity - litersLeft
			price = litersRefueled * FUEL_RATE
		}

		assignments = append(assignments, Assignment{
			LicencePlate: licencePlate,
			Employee:     AssignEmployee(size, price),
			FuelAdded:    litersRefueled,
			Price:        GetParkingRate(size) + price,
		})
	}

	json.NewEncoder(w).Encode(assignments)
}

func handleRequest() {
	router := mux.NewRouter()
	router.HandleFunc("/", HomeHandler)
	router.HandleFunc("/parking", ParkingHandler).Methods("POST")
	http.Handle("/", router)

	log.Fatal(http.ListenAndServe(":8000", router))
}

func main() {
	handleRequest()
}
