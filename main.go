package main

import (
	"context"
	"log"
	"time"

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
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	err = client.Connect(ctx)
	db := client.Database(databaseName)

	if err != nil {
		return err
	}

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
