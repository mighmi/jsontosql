package main

import (
	"crypto/tls"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

// This program will save json from 2 sources in postgres:
// - https://dummyjson.com/
// - https://randomuser.me/api/
// Aditionally, it can query/print DB information for the user, and allow the user to post new information.

// Create a custom HTTP client to bypass:
// panic: Get "https://dummyjson.com/users": tls: failed to verify certificate: x509: certificate signed by unknown authority
// in Dockerfile there are ideas to overcome this
var client = &http.Client{
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	},
}

func getUsers() []Person {
	var choice int
	fmt.Print("To add a user:\n1. randomuser.me\n2. dummyuser.com\n3. create your own\n")
	fmt.Scan(&choice)

	switch choice {
	case 1:
		fmt.Println("How many users?")
		fmt.Scan(&choice)
		var people []Person
		for i := 0; i < choice; i++ {
			person := getRandomUser()
			people = append(people, person)
		}
		return people
	case 2:
		return getDummyJsonUsers()
	case 3:
		fmt.Print("Human can't make users yet")
		return []Person{}
	default:
		fmt.Print("Only 1-3")
		return []Person{}
	}
}

func getRandomUser() Person {
	resp, err := client.Get("https://randomuser.me/api/")
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

func getDummyJsonUsers() []Person {
	resp, err := client.Get("https://dummyjson.com/users")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	var DJP DummyJsonPerson
	err = decoder.Decode(&DJP) // output is this RUP struct
	if err != nil {
		panic(err)
	}

	return convertToPeople(DJP)
}

func convertToPeople(rUP DummyJsonPerson) []Person {
	var people []Person
	for _, result := range rUP.Users {

		// latitude, _ := strconv.ParseFloat(result.Address.Coordinates.Lat, 64)
		// longitude, _ := strconv.ParseFloat(result.Address.Coordinates.Lng, 64)

		person := Person{
			FirstName:   result.FirstName,
			LastName:    result.LastName,
			Latitude:    result.Address.Coordinates.Lat,
			Longitude:   result.Address.Coordinates.Lng,
			Email:       result.Email,
			Username:    result.Username,
			Password:    result.Password,
			DateOfBirth: result.BirthDate,
		}
		people = append(people, person)
	}

	return people
}

// DB stuff
const (
	host     = "127.0.0.1" //"sql_db_1" // "127.0.0.1" // use docker name if from docker, ip if not in container!
	port     = 5435
	user     = "postgres"
	password = "password2"
	dbname   = "humans"
)

func AddToPostgres(people []Person) {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo) // first arg is driver name (pq!)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	for _, person := range people {
		_, err = db.Exec(`INSERT INTO people (FirstName, LastName, Latitude, Longitude, Username, passwd, Email, DateOfBirth)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
			person.FirstName, person.LastName, person.Latitude, person.Longitude,
			person.Email, person.Username, person.Password, person.DateOfBirth)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate") {
				fmt.Println("Duplicate value not saved")
			} else {
				panic(err)
			}
		}
	}
}

func main() {

	result := getUsers()
	fmt.Printf("%+v\n", result)
	AddToPostgres(result)

}
