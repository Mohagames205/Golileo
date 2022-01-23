package main

import (
	"github.com/Mohagames205/Golileo/database"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"log"
)

func main() {
	app := fiber.New(fiber.Config{
		AppName: "Golileo",
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Dag wereld 👋!")
	})

	app.Post("/api/skin", func(ctx *fiber.Ctx) error {
		database, context := database.Database()
		skinCollection := database.Collection("skins")
		_, err := skinCollection.InsertOne(context, bson.D{
			{Key: "username", Value: "test"},
			{Key: "skinstring", Value: "aabbcc"},
		})

		if err != nil{
			log.Println(err)
		}

		return ctx.SendString("Skin has been inserted into the database")
	})

	app.Listen(":3000")
}
