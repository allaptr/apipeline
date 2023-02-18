package main

import (
	"fmt"
	"log"
)

func main() {
	var err error
	fmt.Println("----Running Sky Pipeline----	")
	pl := NewSkyCPIXPipeline()
	for pl != nil {
		err = pl.Handle()
		if err != nil {
			log.Fatal(err)
		}
		pl = pl.Next()
	}
	fmt.Println("----Running Fire Pipeline----	")
	fpl := NewFireCPIXPipeline()
	for fpl != nil {
		err = fpl.Handle()
		if err != nil {
			log.Fatal(err)
		}
		fpl = fpl.Next()
	}
}
