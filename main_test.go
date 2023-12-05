package main

import (
	"fmt"
	"testing"
)

func BenchmarkGetRandomUsers(b *testing.B) {

	var count int = 100

	var people []Person = getRandomUsers(count)
	fmt.Print(people)
}

// Non concurrent results:
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

// tokenbucket 5 takes 20s (encoding/json)
// 20310413868 ns/op	10218400 B/op	   50439 allocs/op

// token bucket 10 took 10.369s
// 10363534365 ns/op	10223976 B/op	   50450 allocs/op

// rate limited at 11

// | Max Speed | 10 tokens/second | https://randomuser.me/api |
// | :----------- | :------: | ------------: |
// | 0363534365 ns/op | 10223976 B/op | 50450 allocs/op |
