package main

import (
	"fmt"
	"os"

	"github.com/MattHauglustaine/ddleash"
)

func main() {
	ddleash, err := ddleash.New(ddleash.Config{
		Team:     "foo",
		Username: "bar",
		Password: "baz",
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := ddleash.Login(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	metrics, err := ddleash.FetchAllMetrics()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(len(metrics))
}
