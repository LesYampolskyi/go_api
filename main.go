package main

import (
	"api/api"
	"api/db"
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const dbUri = "mongodb://localhost:27017"
const dbName = "hotel-reservation"
const userColl = "users"

var config = fiber.Config{
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return c.JSON(map[string]string{"error": err.Error()})
	},
}

func main() {

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dbUri))
	if err != nil {
		log.Fatal(err)
	}

	// handlers initialization
	userHandler := api.NewUserHandler(db.NewMongoUserStore(client, db.DBNAME))

	fmt.Println(client)
	listenAddr := flag.String("listenAddr", ":5000", "Listen address of API server")
	flag.Parse()
	app := fiber.New(config)

	// user := types.User{
	// 	FirstName: "ME",
	// 	LastName:  "MARIO",
	// }
	// ctx := context.Background()
	// res, err := client.Database("hotel-reservation").Collection("users").InsertOne(ctx, user)
	if err != nil {
		panic(err)
	}
	// fmt.Println(res)
	fmt.Println("TTTT")

	apiv1 := app.Group("/api/v1")

	// app.Get("/foo", handleFoo)

	apiv1.Post("/user", userHandler.HandlePostUser)
	apiv1.Get("/user", userHandler.HandleGetUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiv1.Put("/user/:id", userHandler.HandlePutUser)

	app.Listen(*listenAddr)
}

func handleFoo(c *fiber.Ctx) error {
	return c.JSON(map[string]string{"msg": "working fine"})
}
