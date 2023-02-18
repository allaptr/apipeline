package main

import (
	"fmt"
	"log"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		var err error
		defer wg.Done()
		fmt.Println("----Running Sky Pipeline----	")
		pl := NewSkyCPIXPipeline()
		for pl != nil {
			err = pl.Handle()
			if err != nil {
				log.Fatal(err)
			}
			pl = pl.Next()
		}
	}()

	go func() {
		var err error
		defer wg.Done()
		fmt.Println("----Running Fire Pipeline----	")
		fpl := NewFireCPIXPipeline()
		for fpl != nil {
			err = fpl.Handle()
			if err != nil {
				log.Fatal(err)
			}
			fpl = fpl.Next()
		}
	}()

	wg.Wait()
}
