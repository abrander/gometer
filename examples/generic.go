package main

import (
	"fmt"

	"github.com/abrander/gometer/iec62056"
)

func main() {
	i, err := iec62056.NewIec62056Serial("/dev/ttyUSB0")
	if err != nil {
		panic(err.Error())
	}
	id, collection, err := i.Signin("")
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("ID: %s\n", id)

	for obis, value := range collection {
		fmt.Printf("%-10s %-60s \033[32m%s\033[0m\n", obis, obis.Description(), value.String())
	}
}
