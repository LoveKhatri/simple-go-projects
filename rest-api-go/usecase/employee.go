package usecase

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/LoveKhatri/rest-api-go/models"
	"github.com/LoveKhatri/rest-api-go/repository"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

type EmployeeService struct {
	MongoCollection *mongo.Collection
}

type Response struct {
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}

func (svc *EmployeeService) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	res := &Response{}
	defer json.NewEncoder(w).Encode(res)

	var emp models.Employee
	err := json.NewDecoder(r.Body).Decode(&emp)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Fatal("Error while decoding request body", err)
		res.Error = err.Error()
		return
	}

	emp.EmployeeId = uuid.New().String()

	repo := repository.EmployeeRepository{MongoCollection: svc.MongoCollection}

	insertedId, err := repo.InsertEmployee(&emp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal("Error while inserting employee", err)
		res.Error = err.Error()
		return
	}

	res.Data = emp.EmployeeId
	w.WriteHeader(http.StatusCreated)

	log.Println("Employee created with ID:", insertedId)

	return
}

func (svc *EmployeeService) GetEmployeeById(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	res := &Response{}
	defer json.NewEncoder(w).Encode(res)

	empId := mux.Vars(r)["id"]

	repo := repository.EmployeeRepository{MongoCollection: svc.MongoCollection}

	emp, err := repo.FindEmployeeById(empId)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal("Error while finding employee", err)
		res.Error = err.Error()
		return
	}

	res.Data = emp
	w.WriteHeader(http.StatusOK)

	log.Println("Employee found with ID:", empId)

	return
}

func (svc *EmployeeService) GetAllEmployees(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	res := &Response{}
	defer json.NewEncoder(w).Encode(res)

	repo := repository.EmployeeRepository{MongoCollection: svc.MongoCollection}

	employees, err := repo.FindAllEmployees()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal("Error while finding employees", err)
		res.Error = err.Error()
		return
	}

	res.Data = employees
	w.WriteHeader(http.StatusOK)

	log.Println("Employees found")

	return
}

func (svc *EmployeeService) UpdateEmployeeById(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	res := &Response{}
	defer json.NewEncoder(w).Encode(res)

	var emp models.Employee
	err := json.NewDecoder(r.Body).Decode(&emp)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Fatal("Error while decoding request body", err)
		res.Error = err.Error()
		return
	}

	emp.EmployeeId = mux.Vars(r)["id"]

	repo := repository.EmployeeRepository{MongoCollection: svc.MongoCollection}

	_, err = repo.UpdateEmployee(emp.EmployeeId, &emp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal("Error while updating employee", err)
		res.Error = err.Error()
		return
	}

	res.Data = emp.EmployeeId
	w.WriteHeader(http.StatusOK)

	log.Println("Employee updated with ID:", emp.EmployeeId)

	return
}

func (svc *EmployeeService) DeleteEmployeeById(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	res := &Response{}
	defer json.NewEncoder(w).Encode(res)

	empId := mux.Vars(r)["id"]

	repo := repository.EmployeeRepository{MongoCollection: svc.MongoCollection}

	_, err := repo.DeleteEmployee(empId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal("Error while deleting employee", err)
		res.Error = err.Error()
		return
	}

	res.Data = empId
	w.WriteHeader(http.StatusOK)

	log.Println("Employee deleted with ID:", empId)

	return
}

func (svc *EmployeeService) DeleteAllEmployees(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	res := &Response{}
	defer json.NewEncoder(w).Encode(res)

	repo := repository.EmployeeRepository{MongoCollection: svc.MongoCollection}

	_, err := repo.DeleteAllEmployees()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal("Error while deleting all employees", err)
		res.Error = err.Error()
		return
	}

	w.WriteHeader(http.StatusOK)

	log.Println("All employees deleted")

	return
}
