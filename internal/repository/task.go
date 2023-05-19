package repository

import (
	"fmt"

	"github.com/dmawardi/Go-Template/internal/db"
	"gorm.io/gorm"
)

type TaskRepository interface {
	FindAll(int, int, string) (*[]db.Task, error)
	FindById(int) (*db.Task, error)
	Create(*db.Task) (*db.Task, error)
	Update(int, *db.Task) (*db.Task, error)
	Delete(int) error
}

type taskRepository struct {
	DB *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db}
}

// Creates a task in the database
func (r *taskRepository) Create(task *db.Task) (*db.Task, error) {
	// Create parameter task in database
	result := r.DB.Create(&task)
	if result.Error != nil {
		return nil, fmt.Errorf("failed creating task: %w", result.Error)
	}

	var assResult error
	// Build associations
	if len(task.Assignment) > 0 {
		assResult = r.DB.Model(&task).Association("Assignment").Replace(task.Assignment)
	}
	// Check if association update failed
	if assResult != nil {
		fmt.Println("Task association update failed: ", assResult)
		return nil, assResult
	}

	return task, nil
}

// Find a list of tasks in the database
func (r *taskRepository) FindAll(limit int, offset int, order string) (*[]db.Task, error) {
	// Query all tasks based on the received parameters
	tasks, err := QueryAllTasksBasedOnParams(limit, offset, order, r.DB)
	if err != nil {
		fmt.Printf("Error querying db for list of tasks: %s", err)
		return nil, err
	}

	return &tasks, nil
}

// Find task in database by ID
func (r *taskRepository) FindById(id int) (*db.Task, error) {
	// Create an empty ref object of type task
	task := db.Task{}
	// Check if task exists in db
	result := r.DB.Preload("Assignment").Preload("Log.User").First(&task, id)

	// Extract error result
	err := result.Error
	// If error detected
	if err != nil {
		return nil, err
	}
	// else
	return &task, nil
}

// Delete task in database
func (r *taskRepository) Delete(id int) error {
	// Create an empty ref object of type task
	task := db.Task{}
	// Check if task exists in db
	result := r.DB.Delete(&task, id)

	// If error detected
	if result.Error != nil {
		fmt.Println("error in deleting task: ", result.Error)
		return result.Error
	}
	// else
	return nil
}

// Updates task in database
func (r *taskRepository) Update(id int, task *db.Task) (*db.Task, error) {
	// Init
	var err error
	// Find task by id
	foundTask, err := r.FindById(id)
	if err != nil {
		fmt.Println("Property to update not found: ", err)
		return nil, err
	}

	// Update found task using new details
	updateResult := r.DB.Model(&foundTask).Updates(task)

	// Extract error
	err = updateResult.Error
	// Check if update failed
	if err != nil {
		fmt.Println("Task update failed: ", err)
		return nil, err
	}

	// Init associate error
	var assResult error
	// Update assigments if available in struct
	if len(task.Assignment) > 0 {
		// Replace
		assResult = r.DB.Model(&foundTask).Association("Assignment").Replace(task.Assignment)
	}
	// Check if association update failed
	if assResult != nil {
		fmt.Println("Task association update failed: ", assResult)
		return nil, assResult
	}

	// Retrieve updated task by id
	updatedTask, err := r.FindById(id)
	if err != nil {
		fmt.Println("Updated task not found: ", err)
		return nil, assResult
	}
	return updatedTask, nil
}

// Takes limit, offset, and order parameters, builds a query and executes returning a list of tasks
func QueryAllTasksBasedOnParams(limit int, offset int, order string, dbClient *gorm.DB) ([]db.Task, error) {
	// Build model to query database
	tasks := []db.Task{}
	// Build base query for tasks table
	query := dbClient.Model(&tasks).Preload("Assignment")

	// Add parameters into query as needed
	if limit != 0 {
		query.Limit(limit)
	}
	if offset != 0 {
		query.Offset(offset)
	}
	// order format should be "column_name ASC/DESC" eg. "created_at ASC"
	if order != "" {
		query.Order(order)
	}
	// Query database
	result := query.Find(&tasks)
	if result.Error != nil {
		return nil, result.Error
	}
	// Return if no errors with result
	return tasks, nil
}
