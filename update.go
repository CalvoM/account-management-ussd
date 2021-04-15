//+build ignore

package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	filename := os.Args[1]
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	vars := os.Args[2:]
	for _, i := range vars {
		fmt.Println(i)
		_, e := f.Write([]byte(i + "=\n"))
		if e != nil {
			f.Close()
			log.Fatal(e)
		}
	}
	f.Close()
}
