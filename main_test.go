package main

import (
	"fmt"
	"testing"
)

func BenchmarkGetRandomUser(b *testing.B) {

	var choice int = 100

	var people []Person
	for i := 0; i < choice; i++ {
		person := getRandomUser()
		people = append(people, person)
	}
	fmt.Print(people)
}

// Results:
// with 500:

// | encoding/json | segmentio | goccy/go-json |
// | :----------- | :------: | ------------: |
// | 32238224 B/op | 47847576 B/op | 31763104 B/op |
// | 88445 allocs/op | 81512 allocs/op | 79894 allocs/op |

// with 100:
// | encoding/json | segmentio | goccy/go-json |
// | :----------- | :------: | ------------: |
// | 9800360 B/op | 12924096 B/op | 9721896 B/op |
// | 48080 allocs/op | 46584 allocs/op | 46316 allocs/op |
