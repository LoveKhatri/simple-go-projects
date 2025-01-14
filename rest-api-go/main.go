package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/LoveKhatri/rest-api-go/usecase"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client

func init() {
	// Load Env Variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error while loading env file", err)
	}

	mongoClient, err := mongo.Connect(context.Background(),
		options.Client().ApplyURI(os.Getenv("DATABASE_URL")))

	if err != nil {
		log.Fatal("Error while connecting to Mongodb", err)
	}

	err = mongoClient.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal("Error while pinging to Mongodb", err)
	}

	log.Println("Connected to Mongodb")
}

func main() {
	defer mongoClient.Disconnect(context.Background())

	coll := mongoClient.Database(os.Getenv("DATABASE_NAME")).Collection(os.Getenv("COLLECTION_NAME"))

	empService := usecase.EmployeeService{MongoCollection: coll}

	r := mux.NewRouter()

	r.HandleFunc("/health", healthHandler).Methods(http.MethodGet)
	r.HandleFunc("/employee", empService.CreateEmployee).Methods(http.MethodPost)
	r.HandleFunc("/employee", empService.GetAllEmployees).Methods(http.MethodGet)
	r.HandleFunc("/employee/{id}", empService.UpdateEmployeeById).Methods(http.MethodPut)
	r.HandleFunc("/employee/{id}", empService.DeleteEmployeeById).Methods(http.MethodDelete)
	r.HandleFunc("/employee", empService.DeleteAllEmployees).Methods(http.MethodDelete)

	log.Println("Server started at :8080")
	http.ListenAndServe(":8080", r)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Running..."))
}
