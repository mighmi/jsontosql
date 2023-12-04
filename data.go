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
	Results []struct { // only 1 per request
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

//	https://dummyjson.com/users is one single call for many people

type DummyJsonPerson struct { // well people...
	Users []struct { // many people/results in one request
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Email     string `json:"email"`
		Username  string `json:"username"`
		Password  string `json:"password"`
		BirthDate string `json:"birthDate"`
		Address   struct {
			Coordinates struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"coordinates"`
		} `json:"address"`
	} `json:"users"`
}
