package main

import "log"

func main() {
	log.Println("Consuming as much CPU as possible across all cores...")
	const concurrency = 100
	for i := 0; i < concurrency; i++ {
		go func() {
			i := 0
			for {
				i++
				useCPU(i)
				if i > 100000000 {
					i = 0
				}
			}
		}()
	}
	forever := make(chan struct{})
	<-forever
}

func useCPU(i int) int {
	return i * i
}
