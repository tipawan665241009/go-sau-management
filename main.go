package main

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
)

func main() {
	db, error := ConnectDatabase()
	if error != nil {
		fmt.Println(error, db)
	}

	app := fiber.New()

	studentRoute(app, db)

	if error := app.Listen(":5599"); error != nil {
		fmt.Println(error)
	}
}
