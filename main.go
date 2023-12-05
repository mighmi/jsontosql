package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

// This program will save json from 2 sources in postgres:
// - https://dummyjson.com/
// - https://randomuser.me/api/
// Aditionally, it can query/print DB information for the user, and allow the user to post new information.

func getUsers() []Person {
	var choice int
	fmt.Print("To add a user:\n1. randomuser.me\n2. dummyuser.com\n3. create your own\n")
	fmt.Scan(&choice)

	switch choice {
	case 1:
		fmt.Println("How many users?")
		fmt.Scan(&choice)
		return getRandomUsers(choice)
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

func getRandomUsers(count int) []Person {
	var people []Person
	ch := make(chan Person, count) // buffer with capacity of count
	var wg sync.WaitGroup
	done := make(chan struct{})

	// token bucket to avoid 429 error
	var tokenpersecond int = 10 // rate limited at 11
	rateLimiter := make(chan time.Time, tokenpersecond)
	go func() {
		ticker := time.NewTicker(time.Second / time.Duration(tokenpersecond))
		defer close(rateLimiter)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				rateLimiter <- time.Now()

			case <-done:
				return
			}
		}
	}()

	for i := 0; i < count; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-rateLimiter // wait for token
			getRandomUser(ch)
		}()
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for person := range ch {
		people = append(people, person)
	}
	return people
}

func getRandomUser(ch chan<- Person) {

	client := &http.Client{
		Timeout: 5 * time.Minute,
	}
	resp, err := client.Get("https://randomuser.me/api/")
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Print("status not ok\n")
		panic(resp.StatusCode)
	}
	decoder := json.NewDecoder(resp.Body)
	var RUP RandomUserPerson
	err = decoder.Decode(&RUP)
	if err != nil {
		fmt.Printf("%v", resp.Body)
		panic(err)
	}
	person := convertToPerson(RUP)
	ch <- person
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
	resp, err := http.Get("https://dummyjson.com/users")
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
	host     = "json_to_sql_db_1" // "127.0.0.1" // use docker name if from docker, ip if not in container!
	port     = 5432
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

func GetRowsFromPostgres() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo) // first arg is driver name (pq!)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var count int
	err = db.QueryRow(`SELECT COUNT(*) FROM people`).Scan(&count)
	if err != nil {
		panic(err)
	}

	fmt.Printf("There are %d rows in the people table\n", count)

}

func WhatDoesUserWantToDo() {
	var choice int
	fmt.Print("What do you want to do?\n1. Add a user\n2. check DB\n3. Exit\n")
	fmt.Scan(&choice)

	switch choice {
	case 1:
		AddToPostgres(getUsers())
	case 2:
		GetRowsFromPostgres()
	case 3:
		fmt.Println("Exiting program.")
		os.Exit(0)
	default:
		fmt.Println("Input not allowed")
	}
}

func main() {

	for {
		WhatDoesUserWantToDo()
	}
}
