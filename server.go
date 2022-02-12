package main

import (
	"context"
	"github.com/Mohagames205/Golileo/skin"
	"github.com/Mohagames205/Golileo/util"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

func main() {
	app := fiber.New(fiber.Config{
		AppName: "Golileo",
	})

	app.Use(cors.New())

	_ = godotenv.Load()

	if os.Getenv("AUTH") == "true" {
		app.Use(basicauth.New(basicauth.Config{
			Users: map[string]string{
				os.Getenv("USERNAME"): os.Getenv("PASSWORD"),
			},
			Unauthorized: func(c *fiber.Ctx) error {
				return c.SendFile("./static/unauthorized.html")
			},
		}))
	}

	err := util.InitFs()
	if err != nil {
		log.Fatal(err)
	}
	util.InitDatabase()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendFile("./static/skindex.html")
	})

	app.Post("/api/skin", func(ctx *fiber.Ctx) error {

		payload := skin.Skin{
			Username: `json:"username"`,
			Skin:     `json:"skin"`,
		}

		if err := ctx.BodyParser(&payload); err != nil {
			return err
		}

		skinCollection := util.Database().Collection("skins")

		opts := options.Update().SetUpsert(true)
		filter := bson.D{
			{"username", payload.Username},
		}
		update := bson.D{{"$set", bson.D{
			{Key: "username", Value: payload.Username},
			{Key: "skinstring", Value: payload.Skin}},
		}}

		_, err := skinCollection.UpdateOne(context.TODO(), filter, update, opts)

		if err != nil {
			log.Println(err)
		}

		return ctx.SendString("Skin has been inserted into the util")
	})

	app.Get("/api/skin/:username/raw", func(ctx *fiber.Ctx) error {

		skinCollection := util.Database().Collection("skins")

		var skinResult bson.M
		err := skinCollection.FindOne(context.TODO(), bson.M{"username": ctx.Params("username")}).Decode(&skinResult)
		if err != nil {
			return fiber.NewError(404, "Username not found")
		}

		return ctx.JSON(fiber.Map{"username": ctx.Params("username"), "skinstring": skinResult["skinstring"]})
	})

	app.Get("/api/skin/:username/img/head", func(ctx *fiber.Ctx) error {

		skinCollection := util.Database().Collection("skins")

		var skinResult bson.M
		err := skinCollection.FindOne(context.TODO(), bson.M{"username": ctx.Params("username")}).Decode(&skinResult)
		if err != nil {
			return fiber.NewError(404, "Username not found")
		}

		skinStruct := skin.S(skinResult["username"].(string), skinResult["skinstring"].(string))

		uuid, err := skinStruct.SaveHeadImage()

		return ctx.JSON(fiber.Map{"url": "/cdn/skinImage/" + uuid})
	})

	app.Get("/api/skin/:username/img/full", func(ctx *fiber.Ctx) error {

		skinCollection := util.Database().Collection("skins")

		var skinResult bson.M
		err := skinCollection.FindOne(context.TODO(), bson.M{"username": ctx.Params("username")}).Decode(&skinResult)
		if err != nil {
			return fiber.NewError(404, "Username not found")
		}

		skinStruct := skin.S(skinResult["username"].(string), skinResult["skinstring"].(string))

		uuid, err := skinStruct.SaveFullImage()

		return ctx.JSON(fiber.Map{"url": "/cdn/skinImage/" + uuid})
	})

	app.Get("/cdn/skinImage/:uuid", func(ctx *fiber.Ctx) error {
		workingDir, _ := os.Getwd()
		uuid := ctx.Params("uuid")
		return ctx.SendFile(workingDir + "/images/" + uuid + ".png")
	})

	log.Fatal(app.Listen(":3000"))
}
