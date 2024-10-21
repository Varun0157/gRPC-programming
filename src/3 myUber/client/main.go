package main

import (
	"fmt"
	"log"
)

func getDriverDetails() (name string) {
	fmt.Println("Enter your name: ")
	fmt.Scan(&name)

	return name
}

func getRiderDetails() (name string, source string, dest string) {
	fmt.Println("Enter your name: ")
	fmt.Scan(&name)
	fmt.Println("Enter your source loc: ")
	fmt.Scan(&source)
	fmt.Println("Enter your destination: ")
	fmt.Scan(&dest)

	return name, source, dest
}

func main() {
	var choice string
	fmt.Println("rider or driver (r/d)?")
	fmt.Scan(&choice)

	if choice == "d" {
		name := getDriverDetails()
		for {
			err := connectDriver(name)
			if err != nil {
				log.Fatalf("error creating driver client %v", err)
			}

			var choice string
			fmt.Print("try another server? (y/n)")
			fmt.Scan(&choice)

			if choice == "n" {
				break
			}
		}

	} else {
		name, source, dest := getRiderDetails()
		// choose a random port
		err := connectRider(name, source, dest)
		if err != nil {
			log.Fatalf("error creating rider client: %v", err)
		}
	}
}
