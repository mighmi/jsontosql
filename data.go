package main

type Person struct {
	FirstName   string
	LastName    string
	Latitude    float64
	Longitude   float64
	Email       string
	Username    string
	Password    string
	DateOfBirth string
}

// https://randomuser.me/api/ is a new call per person

type RandomUserPerson struct {
	Results []struct {
		Name struct {
			First string `json:"first"`
			Last  string `json:"last"`
		} `json:"name"`
		Location struct {
			Coordinates struct {
				Latitude  string `json:"latitude"`
				Longitude string `json:"longitude"`
			} `json:"coordinates"`
		} `json:"location"`
		Email string `json:"email"`
		Login struct {
			Username string `json:"username"`
			Password string `json:"password"`
		} `json:"login"`
		DOB struct {
			Date string `json:"date"`
		} `json:"dob"`
	} `json:"results"`
}

// name
// 	first
// 	last

// location
// 	country
// 	postcode
// 	coordinates
// 		latitude
// 		longtitude
// email
// login
// 	username
// 	password
// dob
// 	date

//	https://dummyjson.com/users is one single call for many people

type DummyJsonPerson struct {
	FirstName   string  `json:"firstName"`
	LastName    string  `json:""`
	Latitude    float64 `json:""`
	Longtitude  float64 `json:""`
	Email       string  `json:""`
	Username    string  `json:""`
	Password    string  `json:""`
	DateOfBirth string  `json:""`
}

// 	firstName
// 	lastName
// username
// password
// email
// address
// 	coordinates:
// 		lat:
// 		lng:
// birthDate
