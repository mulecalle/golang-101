// https://medium.com/@gauravsingharoy/asynchronous-programming-with-go-546b96cd50c1

package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	// A slice of sample websites
	urls := []string{
		"https://www.easyjet.com/",
		"https://www.skyscanner.de/",
		"https://www.ryanair.com",
		"https://wizzair.com/",
		"https://www.swiss.com/",
	}

	c := make(chan urlStatus, len(urls))

	for _, url := range urls {
		go checkUrl(url, c)
	}

	result := make([]urlStatus, len(urls))

	fmt.Println("Starting")

	for i, _ := range result {
		result[i] = <-c
		if result[i].status {
			fmt.Println(result[i].url, "is up.")
		} else {
			fmt.Println(result[i].url, "is down !!")
		}
	}

	fmt.Println("Finished")

}

// checks and prints a message if a website is up or down
func checkUrl(url string, c chan urlStatus) {
	_, err := http.Get(url)

	min := 10
	max := 100
	asd := rand.Intn(rand.Intn(max-min) + min)

	time.Sleep(time.Duration(asd) * time.Second)

	if err != nil {
		// The website is down
		c <- urlStatus{url, false}
	} else {
		// The website is up
		c <- urlStatus{url, true}
	}
}

type urlStatus struct {
	url    string
	status bool
}
