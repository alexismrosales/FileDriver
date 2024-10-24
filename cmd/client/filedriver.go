package main

import (
	"log"
	"os"
)

func main() {
	app := RunApp()
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
