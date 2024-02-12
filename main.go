package main

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoInstance struct {
	Client   *mongo.Client
	Database *mongo.Database
}

type Employee struct {
	ID     string  `json:"id,omitempty" bson:"_id,omitempty"`
	Name   string  `json:"nome"`
	Salary float64 `json:"salario"`
	Age    float64 `json:"idade"`
}

var mg MongoInstance

const databaseName = "go-hrms"
const mongoURI = "mongodb://localhost:27017/" + databaseName

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

	app.Get("/employees", func(c *fiber.Ctx) error {

		// Buscando por todos os registros na coleção.
		query := bson.D{{}} // Query vazia para buscar por tudo.

		cursor, err := mg.Database.Collection("employees").Find(c.Context(), query)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		// Passando todos os registros para uma lista de funcionários
		var employees []Employee = make([]Employee, 0)

		if err := cursor.All(c.Context(), &employees); err != nil {
			return c.Status(500).SendString(err.Error())
		}

		return c.JSON(employees)
	})

	app.Post("/new/employee", func(c *fiber.Ctx) error {
		collection := mg.Database.Collection("employees")

		// Coletando dados da requisição
		employee := new(Employee)

		if err := c.BodyParser(employee); err != nil {
			return c.Status(400).SendString(err.Error())
		}

		employee.ID = "" // Tirando id para deixar o Mongo criar automaticamente.

		// Inserindo o novo funcionário
		result, err := collection.InsertOne(c.Context(), employee)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		// Buscando pelo funcionário inserido para retornar.
		filter := bson.D{{Key: "_id", Value: result.InsertedID}}
		createdRecord := collection.FindOne(c.Context(), filter)

		createdEmployee := &Employee{}
		createdRecord.Decode(createdEmployee)

		return c.Status(201).JSON(createdEmployee)
	})
}
