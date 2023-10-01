package main

import (
	"log"

	"github.com/arrow2nd/nekomata/app"
)

func main() {
	app := app.New()

	if err := app.Init(); err != nil {
		log.Fatalln(err)
	}

	if err := app.Run(); err != nil {
		log.Fatalln(err)
	}
}
