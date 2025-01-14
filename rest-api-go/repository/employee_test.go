package repository

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/LoveKhatri/rest-api-go/models"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func newMongoClient() *mongo.Client {
	// Load Env Variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error while loading env file", err)
	}

	mongoTestClient, err := mongo.Connect(context.Background(),
		options.Client().ApplyURI(os.Getenv("DATABASE_URL")))

	if err != nil {
		log.Fatal("Error while connecting to Mongodb", err)
	}

	log.Println("Connected to Mongodb")

	err = mongoTestClient.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatal("Error while pinging to Mongodb", err)
	}

	log.Println("Ping to Mongodb successful")

	return mongoTestClient
}

func TestMongoOperations(t *testing.T) {
	mongoTestClient := newMongoClient()
	defer mongoTestClient.Disconnect(context.Background())

	emp1 := uuid.New().String()
	emp2 := uuid.New().String()

	coll := mongoTestClient.Database("main").Collection("employee")

	empRepo := EmployeeRepository{MongoCollection: coll}

	// Insert Employee
	t.Run("Insert Employee 1", func(t *testing.T) {
		emp := models.Employee{
			EmployeeId: emp1,
			Name:       "John Doe",
			Department: "Avengers",
		}

		result, err := empRepo.InsertEmployee(&emp)

		if err != nil {
			t.Fatal("Error while inserting employee", err)
		}

		t.Log("Inserted Employee 1 with ID:", result)
	})
	t.Run("Insert Employee 2", func(t *testing.T) {
		emp := models.Employee{
			EmployeeId: emp2,
			Name:       "Jane Doe",
			Department: "Avengers Assemble",
		}

		result, err := empRepo.InsertEmployee(&emp)

		if err != nil {
			t.Fatal("Error while inserting employee", err)
		}

		t.Log("Inserted Employee 2 with ID:", result)
	})

	// Find Employee
	t.Run("Find Employee 1", func(t *testing.T) {
		emp, err := empRepo.FindEmployeeById(emp1)

		if err != nil {
			t.Fatal("Error while finding employee", err)
		}

		t.Log("Found Employee 1:", emp)
	})

	// Get All employees
	t.Run("Get All Employees", func(t *testing.T) {
		emps, err := empRepo.FindAllEmployees()

		if err != nil {
			t.Fatal("Error while getting all employees", err)
		}

		t.Log("Found Employees:", emps)
	})

	// Update Employee
	t.Run("Update Employee 1", func(t *testing.T) {
		emp := models.Employee{
			EmployeeId: emp1,
			Name:       "John Wick",
			Department: "Avengers",
		}

		result, err := empRepo.UpdateEmployee(emp1, &emp)

		if err != nil {
			t.Fatal("Error while updating employee", err)
		}

		t.Log("Updated Employee count:", result)
	})

	// Delete Employee
	t.Run("Delete Employee 1", func(t *testing.T) {
		result, err := empRepo.DeleteEmployee(emp1)

		if err != nil {
			t.Fatal("Error while deleting employee", err)
		}

		t.Log("Deleted Employee count:", result)
	})

	// Delete All employees
	t.Run("Delete All Employees", func(t *testing.T) {
		result, err := empRepo.DeleteAllEmployees()

		if err != nil {
			t.Fatal("Error while deleting all employees", err)
		}

		t.Log("Deleted Employee count:", result)
	})

	// Get All employees
	t.Run("Get All Employees after Cleanup", func(t *testing.T) {
		emps, err := empRepo.FindAllEmployees()

		if err != nil {
			t.Fatal("Error while getting all employees", err)
		}

		t.Log("Found Employees:", emps)
	})
}
