package main

import (
	"fmt"
	"os"
)

func main() {
	// musicUrl := flag.String("url", "", "Adds a new music URL to the youtube playlist of the config file")

	// flag.Parse()

	err := StartBot()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return
}
