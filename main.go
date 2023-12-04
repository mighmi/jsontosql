package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// This program will save json from 2 sources in postgres:
// - https://dummyjson.com/
// - https://randomuser.me/api/
// Aditionally, it can query/print DB information for the user, and allow the user to post new information.

func getRandomUser() Person {

	resp, err := http.Get("https://randomuser.me/api/")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	var RUP RandomUserPerson
	err = decoder.Decode(&RUP) // output is this RUP struct
	if err != nil {
		panic(err)
	}

	return convertToPerson(RUP)

}

func convertToPerson(rUP RandomUserPerson) Person {
	result := rUP.Results[0]

	latitude, _ := strconv.ParseFloat(result.Location.Coordinates.Latitude, 64)
	longitude, _ := strconv.ParseFloat(result.Location.Coordinates.Longitude, 64)

	person := Person{
		FirstName:   result.Name.First,
		LastName:    result.Name.Last,
		Latitude:    latitude,
		Longitude:   longitude,
		Email:       result.Email,
		Username:    result.Login.Username,
		Password:    result.Login.Password,
		DateOfBirth: result.DOB.Date,
	}

	return person
}

func main() {

	result := getRandomUser()
	fmt.Printf("%+v\n", result)

}
