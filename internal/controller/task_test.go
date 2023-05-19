package controller_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dmawardi/Go-Template/internal/db"
)

func TestTaskController_FindAll(t *testing.T) {
	// Test setup
	var createTask1 = db.Task{
		TaskName: "Broken light switches",
		Type:     "Maintenance",
		Notes:    "This is a note",
		Status:   "Pending",
	}
	var createTask2 = db.Task{
		TaskName: "Broken light switches",
		Type:     "Maintenance",
		Notes:    "This is a note",
		Status:   "Pending",
	}
	var createdTasks = []db.Task{createTask1, createTask2}
	// Create tasks in database
	seedErr := testConnection.dbClient.Create(createdTasks)
	if seedErr.Error != nil {
		t.Errorf("Error seeding database: %v", seedErr.Error)
	}

	// Create a new request
	req, err := http.NewRequest("GET", "/api/tasks?limit=10&offset=0&order=", nil)
	if err != nil {
		t.Fatal(err)
	}
	// Add auth token to header
	req.Header.Set("Authorization", fmt.Sprintf("bearer %v", testConnection.accounts.admin.token))
	// Create a response recorder
	rr := httptest.NewRecorder()

	// Use handler with recorder and created request
	testConnection.router.ServeHTTP(rr, req)

	// Check the response status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Convert response JSON to struct
	var body []db.Task
	json.Unmarshal(rr.Body.Bytes(), &body)

	// Check length of tasks array (should be two with seeded assets)
	if len(body) != len(createdTasks) {
		t.Errorf("Tasks array in findAll failed: expected %d, got %d", len(createdTasks), len(body))
	}

	// Iterate through tasks array received
	for _, actual := range body {
		// Iterate through created tasks to determine a match
		for _, created := range createdTasks {
			// If match found
			if actual.ID == created.ID {
				// Check the details of the feature
				checkTaskDetails(&actual, &created, t, false)
			}
		}
	}

	// Test parameter inputs
	//
	var failParameterTests = []struct {
		limit                  string
		offset                 string
		order                  string
		expectedResponseStatus int
	}{
		// Bad order by
		{limit: "10", offset: "", order: "none", expectedResponseStatus: http.StatusBadRequest},
		// No limit should result in bad request
		{limit: "", offset: "", order: "", expectedResponseStatus: http.StatusBadRequest},
		// Check normal parameters functional with order by
		{limit: "20", offset: "1", order: "ID ASC", expectedResponseStatus: http.StatusOK},
	}
	// Iterate through URL parameter tests
	for _, v := range failParameterTests {
		request := fmt.Sprintf("/api/tasks?limit=%v&offset=%v&order=%v", v.limit, v.offset, v.order)
		// Create a new request
		req, err := http.NewRequest("GET", request, nil)
		if err != nil {
			t.Fatal(err)
		}
		// Add auth token to header
		req.Header.Set("Authorization", fmt.Sprintf("bearer %v", testConnection.accounts.admin.token))
		// Create a response recorder
		rr := httptest.NewRecorder()

		// Use handler with recorder and created request
		testConnection.router.ServeHTTP(rr, req)

		// Check the response status code
		if status := rr.Code; status != v.expectedResponseStatus {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, v.expectedResponseStatus)
		}
	}

	// Clean up created fixtures
	deleteResult := testConnection.dbClient.Delete(createdTasks)
	if deleteResult.Error != nil {
		t.Fatalf("Couldn't clean up seeded features for find all: %v", deleteResult.Error)
	}
}

func TestTaskController_Find(t *testing.T) {
	// Test setup
	var createTask = db.Task{
		TaskName: "Broken light switches",
		Type:     "Maintenance",
		Notes:    "This is a note",
	}

	var createdTasks = []db.Task{createTask}
	// Create tasks in database
	seedErr := testConnection.dbClient.Create(createdTasks)
	if seedErr.Error != nil {
		t.Errorf("Error seeding database: %v", seedErr.Error)
	}

	// Create a request with an "id" URL parameter
	requestUrl := fmt.Sprintf("/api/tasks/%v", createdTasks[0].ID)
	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		t.Fatal(err)
	}

	// // Add auth token to header
	req.Header.Set("Authorization", fmt.Sprintf("bearer %v", testConnection.accounts.admin.token))
	// // Create a response recorder
	rr := httptest.NewRecorder()

	// // Use handler with recorder and created request
	testConnection.router.ServeHTTP(rr, req)

	// Check the response status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v: %v\n Created: %v",
			status, http.StatusOK, rr.Body.String(), createdTasks[0])
	}
	// Extract the response body
	var body db.Task
	json.Unmarshal(rr.Body.Bytes(), &body)

	checkTaskDetails(&body, &createdTasks[0], t, true)

	// Cleanup
	// Delete the created fixtures
	deleteResult := testConnection.dbClient.Delete(createdTasks)
	if deleteResult.Error != nil {
		t.Fatalf("Couldn't clean up seeded fixtures for find all: %v", deleteResult.Error)
	}
}

func TestTaskController_Delete(t *testing.T) {
	// Test setup
	var createTask = db.Task{
		TaskName: "Broken light switches",
		Type:     "Maintenance",
		Notes:    "This is a note",
		Status:   "Pending",
	}

	var createdTasks = []db.Task{createTask}
	// Create tasks in database
	seedErr := testConnection.dbClient.Create(createdTasks)
	if seedErr.Error != nil {
		t.Errorf("Error seeding database: %v", seedErr.Error)
	}

	// Create a request with an "id" URL parameter
	requestUrl := fmt.Sprintf("/api/tasks/%v", createdTasks[0].ID)
	req, err := http.NewRequest("DELETE", requestUrl, nil)
	if err != nil {
		t.Fatal(err)
	}
	// Array of tests
	var deleteTests = []struct {
		testName               string
		tokenToUse             string
		expectedResponseStatus int
	}{
		// Tests of deletion failure
		{testName: "Contacts basic user delete test", tokenToUse: testConnection.accounts.user.token, expectedResponseStatus: http.StatusForbidden},
		// Must be last
		// Tests of deletion success using admin token
		{testName: "Contacts admin delete test", tokenToUse: testConnection.accounts.admin.token, expectedResponseStatus: http.StatusOK},
	}

	// Iterate through tests
	for _, test := range deleteTests {
		// Create a response recorder
		rr := httptest.NewRecorder()

		// Add auth token to header
		req.Header.Set("Authorization", fmt.Sprintf("bearer %v", test.tokenToUse))
		// Send deletion requestion to mock server
		testConnection.router.ServeHTTP(rr, req)
		// Check response is as expected
		if status := rr.Code; status != test.expectedResponseStatus {
			t.Errorf("Task deletion test (%v): got %v want %v.", test.testName,
				status, test.expectedResponseStatus)
		}
	}
}

func TestTaskController_Update(t *testing.T) {
	// Test setup
	var createTask = db.Task{
		TaskName: "Broken light switches",
		Type:     "Maintenance",
		Notes:    "This is a note",
		Status:   "Pending",
	}

	var createdTasks = []db.Task{createTask}
	// Create tasks in database
	seedErr := testConnection.dbClient.Create(createdTasks)
	if seedErr.Error != nil {
		t.Errorf("Error seeding database: %v", seedErr.Error)
	}

	// Build test array
	var updateTests = []struct {
		data                   db.Task
		tokenToUse             string
		expectedResponseStatus int
		checkDetails           bool
		testName               string
	}{
		{db.Task{
			TaskName: "Broken light switches",
			Type:     "Maintenance",
			Notes:    "This is a note",
			Status:   "Pending",
		}, testConnection.accounts.user.token, http.StatusForbidden, false, "Contacts basic user update test"},
		{db.Task{
			TaskName: "Broken light switches",
			Type:     "Maintenance",
			Notes:    "This is a note",
			Status:   "Pending",
		}, testConnection.accounts.admin.token, http.StatusOK, true, "Contacts admin update test"},
		// Update should be disallowed due to being too short
		{db.Task{
			TaskName: "Br",
			Type:     "Maintenance",
			Notes:    "This",
		}, testConnection.accounts.admin.token, http.StatusBadRequest, false, "Contacts admin too short fail test"},
		// Update should be disallowed due to not being proper Type value
		{db.Task{
			TaskName: "Broken light switches",
			Type:     "eggos",
		}, testConnection.accounts.admin.token, http.StatusBadRequest, false, "Contacts admin bad email fail test"},
		// Update should be disallowed due to not being proper Status value
		{db.Task{
			TaskName: "Broken light switches",
			Type:     "Maintenance",
			Status:   "Been there",
		}, testConnection.accounts.admin.token, http.StatusBadRequest, false, "Contacts admin bad phone fail test"},
		// User should be forbidden before validating rather than Bad Request
		{db.Task{
			TaskName: "Br",
			Type:     "Inspection",
			Notes:    "This is a note",
			Status:   "Pending",
		}, testConnection.accounts.user.token, http.StatusForbidden, false, "Contacts basic user too short fail test"},
	}

	// Create a request url with an "id" URL parameter
	requestUrl := fmt.Sprintf("/api/tasks/%v", createdTasks[0].ID)

	// Iterate through update tests
	for _, v := range updateTests {
		// Make new request with task update in body
		req, err := http.NewRequest("PUT", requestUrl, buildReqBody(v.data))
		if err != nil {
			t.Fatal(err)
		}
		// Create a response recorder
		rr := httptest.NewRecorder()
		// Add auth token to header
		req.Header.Set("Authorization", fmt.Sprintf("bearer %v", v.tokenToUse))
		// Send update request to mock server
		testConnection.router.ServeHTTP(rr, req)

		// Convert response JSON to struct
		var body db.Task
		json.Unmarshal(rr.Body.Bytes(), &body)
		// Check response expected vs received
		if status := rr.Code; status != v.expectedResponseStatus {
			t.Errorf("Task update (%v): got %v want %v.", v.testName,
				status, v.expectedResponseStatus)
		}

		// If need to check details
		if v.checkDetails == true {
			// Get task details from database
			var expected db.Task
			findResult := testConnection.dbClient.Find(&expected, createdTasks[0].ID)
			if findResult.Error != nil {
				t.Errorf("Error finding updated task: %v", findResult.Error)
			}

			// Check task details using updated object
			checkTaskDetails(&body, &expected, t, false)
		}
	}

	// Check for failure if incorrect ID parameter detected
	//
	var failUpdateTests = []struct {
		urlExtension           string
		expectedResponseStatus int
	}{
		// alpha character instead
		{urlExtension: "x", expectedResponseStatus: http.StatusForbidden},
		// Index out of bounds
		{urlExtension: "9", expectedResponseStatus: http.StatusBadRequest},
	}
	for _, v := range failUpdateTests {
		// Make new request with feature update in body
		req, err := http.NewRequest("PUT", fmt.Sprint("/api/tasks/"+v.urlExtension), buildReqBody(&db.Contact{
			FirstName: "Scrappy Kid",
		}))
		if err != nil {
			t.Fatal(err)
		}
		// Create a response recorder
		rr := httptest.NewRecorder()
		// Add auth token to header
		req.Header.Set("Authorization", fmt.Sprintf("bearer %v", testConnection.accounts.admin.token))

		// Send update request to mock server
		testConnection.router.ServeHTTP(rr, req)
		// Check response is as expected
		if status := rr.Code; status != v.expectedResponseStatus {
			t.Errorf("Task update test: got %v want %v.",
				status, v.expectedResponseStatus)
		}
	}

	// Delete the created fixtures
	deleteResult := testConnection.dbClient.Delete(createdTasks)
	if deleteResult.Error != nil {
		t.Fatalf("Couldn't clean up seeded fixtures: %v", deleteResult.Error)
	}
}

func TestTaskController_Create(t *testing.T) {
	// Setup
	//
	var createTests = []struct {
		data                   db.Task
		expectedResponseStatus int
		tokenToUse             string
	}{
		// Should fail due to user role status of basic
		{db.Task{
			TaskName: "Broken light switches",
			Type:     "Maintenance",
			Notes:    "This is a note",
		}, http.StatusForbidden, testConnection.accounts.user.token},
		// Should pass as user is admin
		{db.Task{
			TaskName: "Broken light switches",
			Type:     "Maintenance",
			Notes:    "This is a note",
		}, http.StatusCreated, testConnection.accounts.admin.token},
		// Create should be disallowed due to being too short
		{db.Task{
			TaskName: "Br",
			Type:     "Maintenance",
			Notes:    "This is a note",
		}, http.StatusBadRequest, testConnection.accounts.admin.token},
		// Should be a bad request due to invalid type
		{db.Task{
			TaskName: "Broken wall sockets",
			Type:     "Regulate",
			Notes:    "This is a note",
			Status:   "Pending",
		}, http.StatusBadRequest, testConnection.accounts.admin.token},
		// Should be a bad request due to invalid status
		{db.Task{
			TaskName: "Sangsaka",
			Type:     "Maintenance",
			Notes:    "This is a note",
			Status:   "Broken",
		}, http.StatusBadRequest, testConnection.accounts.admin.token},
	}

	// Create a request url with an "id" URL parameter
	requestUrl := "/api/tasks"

	for _, v := range createTests {
		// Make new request with contact creation in body
		req, err := http.NewRequest("POST", requestUrl, buildReqBody(v.data))
		if err != nil {
			t.Fatal(err)
		}
		// Create a response recorder
		rr := httptest.NewRecorder()
		// Add auth token to header
		req.Header.Set("Authorization", fmt.Sprintf("bearer %v", v.tokenToUse))

		// Send update request to mock server
		testConnection.router.ServeHTTP(rr, req)
		// Check response is as expected
		if status := rr.Code; status != v.expectedResponseStatus {
			t.Errorf("Task create test (%v): got %v want %v.", v.data.TaskName,
				status, v.expectedResponseStatus)
		}

		// Init body for response extraction
		var body db.Task
		var foundTask db.Task
		// Grab ID from response body
		json.Unmarshal(rr.Body.Bytes(), &body)

		// Find the created task (to obtain full data with ID)
		testConnection.dbClient.Find(foundTask, uint(body.ID))

		// Compare found task details with those found in returned body
		checkTaskDetails(&body, &foundTask, t, true)

		// Delete the created fixtures
		delResult := testConnection.dbClient.Delete(&db.Task{}, uint(body.ID))
		if delResult.Error != nil {
			t.Fatalf("Issue encountered deleting seeded assets for task create test (%v): %v", v.data.TaskName, delResult.Error)
		}
	}
}

// Check if task details match expected
func checkTaskDetails(actual *db.Task, expected *db.Task, t *testing.T, checkId bool) {

	// Only check ID if parameter checkId is true
	if checkId == true {
		// Verify that the actual contact id matches the created contact
		if actual.ID != expected.ID {
			t.Errorf("found task has incorrect ID: expected %d, got %d", expected.ID, actual.ID)
		}
	}

	// Verify that the actual contact details match the expected contact details
	if actual.Completed != expected.Completed {
		t.Errorf("Task has incorrect Completed value: expected %v, got %v", expected.Completed, actual.Completed)
	}
	if actual.Notes != expected.Notes {
		t.Errorf("Task has incorrect Notes value: expected %s, got %s", expected.Notes, actual.Notes)
	}
	if actual.TaskName != expected.TaskName {
		t.Errorf("Task has incorrect TaskName value: expected %s, got %s", expected.TaskName, actual.TaskName)
	}
	if actual.Type != expected.Type {
		t.Errorf("Task has incorrect Type value: expected %s, got %s", expected.Type, actual.Type)
	}
	if actual.Snoozed != expected.Snoozed {
		t.Errorf("Task has incorrect Snoozed value: expected %v, got %v", expected.Snoozed, actual.Snoozed)
	}
}
