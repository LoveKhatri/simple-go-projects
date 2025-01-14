package repository

import (
	"context"

	"github.com/LoveKhatri/rest-api-go/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type EmployeeRepository struct {
	MongoCollection *mongo.Collection
}

func (r *EmployeeRepository) InsertEmployee(emp *models.Employee) (interface{}, error) {
	result, err := r.MongoCollection.InsertOne(context.Background(), emp)

	if err != nil {
		return nil, err
	}

	return result.InsertedID, nil
}

func (r *EmployeeRepository) FindEmployeeById(empId string) (*models.Employee, error) {
	var emp models.Employee

	err := r.MongoCollection.FindOne(context.Background(),
		bson.D{{Key: "employee_id", Value: empId}}).Decode(&emp)

	if err != nil {
		return nil, err
	}

	return &emp, nil
}

func (r *EmployeeRepository) FindAllEmployees() ([]models.Employee, error) {
	result, err := r.MongoCollection.Find(context.Background(), bson.D{})

	if err != nil {
		return nil, err
	}

	var employees []models.Employee

	err = result.All(context.Background(), &employees)

	if err != nil {
		return nil, err
	}

	return employees, nil
}

func (r *EmployeeRepository) UpdateEmployee(empId string, emp *models.Employee) (int64, error) {
	result, err := r.MongoCollection.UpdateOne(context.Background(),
		bson.D{{Key: "employee_id", Value: empId}},
		bson.D{{Key: "$set", Value: emp}})

	if err != nil {
		return -1, err
	}

	return result.ModifiedCount, nil
}

func (r *EmployeeRepository) DeleteEmployee(empId string) (int64, error) {
	result, err := r.MongoCollection.DeleteOne(context.Background(),
		bson.D{{Key: "employee_id", Value: empId}})

	if err != nil {
		return -1, err
	}

	return result.DeletedCount, nil
}

func (r *EmployeeRepository) DeleteAllEmployees() (int64, error) {
	result, err := r.MongoCollection.DeleteMany(context.Background(), bson.D{})

	if err != nil {
		return -1, err
	}

	return result.DeletedCount, nil
}
