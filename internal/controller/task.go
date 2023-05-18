package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dmawardi/Go-Template/internal/helpers"
	"github.com/dmawardi/Go-Template/internal/models"
	"github.com/dmawardi/Go-Template/internal/service"
	"github.com/go-chi/chi"
)

type TaskController interface {
	FindAll(w http.ResponseWriter, r *http.Request)
	Find(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type taskController struct {
	service service.TaskService
}

func NewTaskController(service service.TaskService) TaskController {
	return &taskController{service}
}

// API/TASKS
// Find a list of tasks
// @Summary      Find a list of tasks
// @Description  Accepts limit, offset, and order params and returns list of tasks
// @Tags         Tasks
// @Accept       json
// @Produce      json
// @Param        limit   path      int  true  "limit"
// @Param        offset   path      int  true  "offset"
// @Param        order   path      int  true  "order by"
// @Success      200 {object} []db.Task
// @Failure      400 {string} string "Can't find tasks"
// @Failure      400 {string} string "Must include limit parameter with a max value of 50"
// @Router       /tasks [get]
// @Security BearerToken
func (c taskController) FindAll(w http.ResponseWriter, r *http.Request) {
	// Grab URL query parameters
	limitParam := r.URL.Query().Get("limit")
	offsetParam := r.URL.Query().Get("offset")
	orderBy := r.URL.Query().Get("order")

	// Convert to int
	limit, _ := strconv.Atoi(limitParam)
	offset, _ := strconv.Atoi(offsetParam)

	// Check that limit is present as requirement
	if (limit == 0) || (limit >= 50) {
		http.Error(w, "Must include limit parameter with a max value of 50", http.StatusBadRequest)
		return
	}

	// Query database for all tasks using query params
	foundTasks, err := c.service.FindAll(limit, offset, orderBy)
	if err != nil {
		http.Error(w, "Can't find tasks", http.StatusBadRequest)
		return
	}
	err = helpers.WriteAsJSON(w, foundTasks)
	if err != nil {
		http.Error(w, "Can't find tasks", http.StatusBadRequest)
		fmt.Println("error writing tasks to response: ", err)
		return
	}
}

// Find a created task
// @Summary      Find task
// @Description  Find a task by ID
// @Tags         Tasks
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Task ID"
// @Success      200 {object} db.Task
// @Failure      400 {string} string "Can't find task with ID: {id}"
// @Router       /tasks/{id} [get]
// @Security BearerToken
func (c taskController) Find(w http.ResponseWriter, r *http.Request) {
	// Grab URL parameter
	stringParameter := chi.URLParam(r, "id")
	// Convert to int
	idParameter, err := strconv.Atoi(stringParameter)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	foundTask, err := c.service.FindById(idParameter)
	if err != nil {
		http.Error(w, fmt.Sprintf("Can't find task with ID: %v\n", idParameter), http.StatusBadRequest)
		return
	}
	err = helpers.WriteAsJSON(w, foundTask)
	if err != nil {
		http.Error(w, fmt.Sprintf("Can't find task with ID: %v\n", idParameter), http.StatusBadRequest)
		return
	}
}

// Create a new task
// @Summary      Create a task
// @Description  Creates a new task
// @Tags         Tasks
// @Accept       json
// @Produce      json
// @Param        task body models.CreateTask true "New Task Json"
// @Success      201 {string} string "Task creation successful!"
// @Failure      400 {string} string "Task creation failed."
// @Router       /tasks [post]
func (c taskController) Create(w http.ResponseWriter, r *http.Request) {
	// Init
	var task models.CreateTask
	// Decode request body as JSON and store in login
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		fmt.Println("Decoding error: ", err)
	}

	// Validate the incoming DTO
	pass, valErrors := helpers.GoValidateStruct(&task)
	// If failure detected
	if !pass {
		// Write bad request header
		w.WriteHeader(http.StatusBadRequest)
		// Write validation errors to JSON
		helpers.WriteAsJSON(w, valErrors)
		return
	}
	// else, validation passes and allow through

	// Create property in db
	_, createErr := c.service.Create(&task)
	if createErr != nil {
		http.Error(w, "Task creation failed.", http.StatusBadRequest)
		return
	}

	// Set status to created
	w.WriteHeader(http.StatusCreated)
	// Send success message in body
	w.Write([]byte("Task creation successful!"))
}

// Update a task (using URL parameter id)
// @Summary      Update task
// @Description  Updates an existing task
// @Tags         Tasks
// @Accept       json
// @Produce      json
// @Param        task body models.UpdateTask true "Update Task Json"
// @Param        id   path      int  true  "Task ID"
// @Success      200 {object} db.Task
// @Failure      400 {string} string "Failed task update"
// @Failure      403 {string} string "Authentication Token not detected"
// @Router       /tasks/{id} [put]
// @Security BearerToken
func (c taskController) Update(w http.ResponseWriter, r *http.Request) {
	// grab id parameter
	var task models.UpdateTask
	// Decode request body as JSON and store in login
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		fmt.Println("Decoding error: ", err)
	}

	// Validate the incoming DTO
	pass, valErrors := helpers.GoValidateStruct(&task)
	// If failure detected
	if !pass {
		// Write bad request header
		w.WriteHeader(http.StatusBadRequest)
		// Write validation errors to JSON
		helpers.WriteAsJSON(w, valErrors)
		return
	}
	// else, validation passes and allow through

	// Grab URL parameter
	stringParameter := chi.URLParam(r, "id")
	// Convert to int
	idParameter, _ := strconv.Atoi(stringParameter)

	// Generate a log message frop task update in preparation for successful update
	// genTaskLogMessage := buildLogUpdate(task)
	// // Grab user id from token
	// userID, err := auth.GetUserIDFromToken(w, r)
	// if err != nil {
	// 	http.Error(w, "Authentication Token not detected", http.StatusForbidden)
	// 	return
	// }

	// Update task
	updatedTask, createErr := c.service.Update(idParameter, &task)
	if createErr != nil {
		http.Error(w, fmt.Sprintf("Failed task update: %s", createErr), http.StatusBadRequest)
		return
	}
	// Proceed to update the log with the update
	// c.log.Create(&models.CreatePropertyLog{
	// 	// From URL parameter
	// 	Property: db.Property{
	// 		ID: uint(idParameter),
	// 	},
	// 	// From JWT token
	// 	User: db.User{ID: uint(userID)},
	// 	// Generated message
	// 	LogMessage: genPropLogMessage,
	// 	Type:       "gen",
	// })

	// Write task to output
	err = helpers.WriteAsJSON(w, updatedTask)
	if err != nil {
		fmt.Printf("Error encountered when writing to JSON. Err: %s", err)
	}
}

// Delete task (using URL parameter id)
// @Summary      Delete task
// @Description  Deletes an existing task
// @Tags         Tasks
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Task ID"
// @Success      200 {string} string "Deletion successful!"
// @Failure      400 {string} string "Failed task deletion"
// @Router       /tasks/{id} [delete]
// @Security BearerToken
func (c taskController) Delete(w http.ResponseWriter, r *http.Request) {
	// Grab URL parameter
	stringParameter := chi.URLParam(r, "id")
	// Convert to int
	idParameter, _ := strconv.Atoi(stringParameter)

	// Attampt to delete task using id
	err := c.service.Delete(idParameter)

	// If error detected
	if err != nil {
		http.Error(w, "Failed task deletion", http.StatusBadRequest)
		return
	}
	// Else write success
	w.Write([]byte("Deletion successful!"))
}

// Build a log string for struct updates
// func buildLogUpdate(updateStruct interface{}) string {
// 	// Log update
// 	var updateString string = "UPDATE: "
// 	// Iterate through key value pairs within the struct
// 	// Get the type of the struct
// 	t := reflect.TypeOf(updateStruct)

// 	// Iterate through the fields of the struct
// 	for i := 0; i < t.NumField(); i++ {
// 		// Get the field
// 		field := t.Field(i)
// 		// Get the value of the field
// 		value := reflect.ValueOf(updateStruct).Field(i).Interface()
// 		// Get the type of value
// 		valueType := reflect.TypeOf(value)

// 		// If value type is string
// 		if valueType.String() == "string" {
// 			// and not empty
// 			if value != "" {
// 				fmt.Printf("\nString value found in Field name: %v", field.Name)
// 				updateString += fmt.Sprintf("%s (%v), ", field.Name, value.(string)[0:5]+"...")
// 			}
// 			// else if value type is numeric
// 		} else if valueType.String() == "int" || valueType.String() == "int64" || valueType.String() == "float64" || valueType.String() == "float32" {
// 			// and not empty
// 			if value != "0" && value != "0.0" {
// 				updateString += fmt.Sprintf("%s, ", field.Name)
// 			}
// 			// else if value type is struct
// 		} else if strings.Contains(valueType.String(), "[]") {
// 			updateString += fmt.Sprintf("[]%s, ", field.Name)
// 		}

// 	}
// 	// Get length of string
// 	stringLength := len(updateString)
// 	// Remove last two characters of string (comma and space)
// 	removeFromEnd := stringLength - 2
// 	croppedLogMessage := updateString[:removeFromEnd]

// 	return croppedLogMessage
// }
