package main

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoInstance struct {
	Client   *mongo.Client
	Database *mongo.Database
}

var mg MongoInstance

const databaseName = "go-hrms"
const mongoURI = "mongodb://localhost:27017/" + databaseName

type Employee struct {
	ID     string
	Name   string
	Salary float64
	Age    float64
}

func connect() error {
	clientOptions := options.Client().ApplyURI(mongoURI)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return err
	}

	db := client.Database(databaseName)

	mg = MongoInstance{
		Client:   client,
		Database: db,
	}

	log.Println("Conexão ao banco finalizada.")
	return nil
}

func main() {
	if err := connect(); err != nil {
		log.Fatal(err)
	}

	app := fiber.New()

	app.Get("/employee", func(c *fiber.Ctx) error {
		return nil
	})
}
