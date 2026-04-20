package main

import (
	"fmt"

	"github.com/zimlewis/shortened/application"
)

func main() {
	app := application.New()

	err := app.Start()
	if err != nil {
		fmt.Println(err)
	}
}
