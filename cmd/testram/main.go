package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	megabytesEnv := os.Getenv("MEGABYTES")
	if megabytesEnv == "" {
		log.Fatal("MEGABYTES environment variable is not set.")
	}
	megabytes, err := strconv.Atoi(megabytesEnv)
	if err != nil {
		log.Fatal("Error converting MEGABYTES to integer")
	}
	fmt.Printf("Allocating %d megabytes of memory....\n", megabytes)

	// Allocate memory
	bytes := make([]byte, megabytes*1024*1024)

	// Fill the allocated memory to ensure it's not optimized away
	for i := range bytes {
		bytes[i] = byte(i % 255)
	}

	// Hold onto the memory indefinitely
	fmt.Printf("Allocated %d megabytes of memory.\n", megabytes)
	for {
		time.Sleep(time.Hour)
	}
}
