package controller_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dmawardi/Go-Template/internal/db"
	"github.com/dmawardi/Go-Template/internal/models"
)

func TestTaskLogController_FindAll(t *testing.T) {
	// Test setup
	// Create task
	taskToCreate := &db.Task{
		TaskName: "Test Task",
		Type:     "Maintenance",
		Notes:    "Yohoo",
	}
	createdTasks := []db.Task{*taskToCreate}
	createResult := testConnection.dbClient.Create(createdTasks)
	if createResult.Error != nil {
		t.Fatal("Failed to create property for task log find all test: ", createResult.Error)
	}
	// Create property logs
	taskLogToCreate1 := &db.TaskLog{
		TaskID:     createdTasks[0].ID,
		LogMessage: "This is a note",
	}
	taskLogToCreate2 := &db.TaskLog{
		TaskID:     createdTasks[0].ID,
		LogMessage: "This is a second note",
	}
	createdTaskLogs := []db.TaskLog{*taskLogToCreate1, *taskLogToCreate2}
	// Create task logs in db
	createResult = testConnection.dbClient.Create(createdTaskLogs)
	if createResult.Error != nil {
		t.Fatal("Failed to create task logs for task log find all test: ", createResult.Error)
	}

	// Create a new request
	req, err := http.NewRequest("GET", "/api/task-logs?limit=10&offset=0&order=", nil)
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
	var body []db.TaskLog
	json.Unmarshal(rr.Body.Bytes(), &body)

	// Check length of array (should be two with seeded assets)
	if len(body) != len(createdTaskLogs) {
		t.Errorf("Features array in findAll failed: expected %d, got %d", len(createdTaskLogs), len(body))
	}

	// Iterate through feature array received
	for _, actualTaskLog := range body {
		// Iterate through created task logs to determine a match
		for _, createdTaskLog := range createdTaskLogs {
			// If match found
			if actualTaskLog.ID == createdTaskLog.ID {
				// Check the details of the task log message
				checkTaskLogDetails(&actualTaskLog, &createdTaskLog, t, true)
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
		request := fmt.Sprintf("/api/task-logs?limit=%v&offset=%v&order=%v", v.limit, v.offset, v.order)
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
	deleteResult := testConnection.dbClient.Delete(createdTaskLogs)
	if deleteResult.Error != nil {
		t.Fatalf("Couldn't clean up seeded prop logs: %v", deleteResult.Error)
	}
	deleteResult = testConnection.dbClient.Delete(createdTasks)
	if deleteResult.Error != nil {
		t.Fatalf("Couldn't clean up seeded props: %v", deleteResult.Error)
	}
}

func TestTaskLogController_Find(t *testing.T) {
	// Test setup
	// Create task
	taskToCreate := &db.Task{
		TaskName: "Swell Task",
		Type:     "Maintenance",
		Notes:    "Yohoo",
	}
	createdTasks := []db.Task{*taskToCreate}
	createResult := testConnection.dbClient.Create(createdTasks)
	if createResult.Error != nil {
		t.Fatal("Failed to create task for task logs find all test: ", createResult.Error)
	}
	// Create task logs
	taskLogToCreate := &db.TaskLog{
		TaskID:     createdTasks[0].ID,
		LogMessage: "This is a note",
	}
	createdTaskLogs := []db.TaskLog{*taskLogToCreate}
	// Create task logs in db
	createResult = testConnection.dbClient.Create(createdTaskLogs)
	if createResult.Error != nil {
		t.Fatal("Failed to create task logs for task log find all test: ", createResult.Error)
	}

	// Create a request with an "id" URL parameter
	requestUrl := fmt.Sprintf("/api/task-logs/%v", createdTaskLogs[0].ID)
	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		t.Fatal(err)
	}

	// // Add auth token to header
	req.Header.Set("Authorization", fmt.Sprintf("bearer %v", testConnection.accounts.admin.token))
	// // Create a response recorder
	rr := httptest.NewRecorder()

	// // Serve request using recorder and created request
	testConnection.router.ServeHTTP(rr, req)

	// Check the response status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v. Error: %v",
			status, http.StatusOK, rr.Body.String())
	}

	// Extract the response body
	var body db.TaskLog
	json.Unmarshal(rr.Body.Bytes(), &body)
	checkTaskLogDetails(&body, &createdTaskLogs[0], t, true)

	// Cleanup
	// Delete the created task logs
	deleteResult := testConnection.dbClient.Delete(createdTaskLogs)
	if deleteResult.Error != nil {
		t.Fatalf("Couldn't clean up seeded task logs: %v", deleteResult.Error)
	}
	// Delete the created tasks
	deleteResult = testConnection.dbClient.Delete(createdTasks)
	if deleteResult.Error != nil {
		t.Fatalf("Couldn't clean up seeded props: %v", deleteResult.Error)
	}
}

func TestTaskLogController_Delete(t *testing.T) {
	// Test setup
	// Create property
	taskToCreate := &db.Task{
		TaskName: "Swell Task",
		Type:     "Maintenance",
		Notes:    "Yohoo",
	}
	createdTasks := []db.Task{*taskToCreate}
	createResult := testConnection.dbClient.Create(createdTasks)
	if createResult.Error != nil {
		t.Fatal("Failed to create seeded task: ", createResult.Error)
	}

	// Create task log
	taskLogToCreate := &db.TaskLog{
		TaskID:     createdTasks[0].ID,
		LogMessage: "This is a log message 1",
	}
	createdTaskLogs := []db.TaskLog{*taskLogToCreate}
	// Create property logs in db
	createResult = testConnection.dbClient.Create(createdTaskLogs)
	if createResult.Error != nil {
		t.Fatal("Failed to create property logs for property log find all test: ", createResult.Error)
	}

	// Create a request with an "id" URL parameter
	requestUrl := fmt.Sprintf("/api/task-logs/%v", testConnection.propertyLogs.created[0].ID)
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
		{testName: "Task log basic user delete test", tokenToUse: testConnection.accounts.user.token, expectedResponseStatus: http.StatusForbidden},
		// Must be last
		// Tests of deletion success using admin privileges
		{testName: "Task log admin delete test", tokenToUse: testConnection.accounts.admin.token, expectedResponseStatus: http.StatusOK},
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
			t.Errorf("Task log deletion test (%v): got %v want %v.", test.testName,
				status, test.expectedResponseStatus)
		}
	}
	// Cleanup
	// Delete the created task logs
	deleteResult := testConnection.dbClient.Delete(createdTaskLogs)
	if deleteResult.Error != nil {
		t.Fatalf("Couldn't clean up seeded task logs: %v", deleteResult.Error)
	}
	// Delete the created tasks
	deleteResult = testConnection.dbClient.Delete(createdTasks)
	if deleteResult.Error != nil {
		t.Fatalf("Couldn't clean up seeded tasks: %v", deleteResult.Error)
	}
}

func TestTaskLogController_Update(t *testing.T) {
	// Test setup
	// Create task
	taskToCreate := &db.Task{
		TaskName: "Swell Task",
		Type:     "Maintenance",
		Notes:    "Yohoo",
	}
	createdTasks := []db.Task{*taskToCreate}
	createResult := testConnection.dbClient.Create(createdTasks)
	if createResult.Error != nil {
		t.Fatal("Failed to create task for task log: ", createResult.Error)
	}
	// Create task log
	taskLogToCreate := &db.TaskLog{
		TaskID:     createdTasks[0].ID,
		LogMessage: "This is a note",
	}
	createdTaskLogs := []db.TaskLog{*taskLogToCreate}
	// Create task logs in db
	createResult = testConnection.dbClient.Create(createdTaskLogs)
	if createResult.Error != nil {
		t.Fatal("Failed to create task logs for test: ", createResult.Error)
	}

	// Build test array
	var updateTests = []struct {
		data                   models.RecvTaskLog
		tokenToUse             string
		expectedResponseStatus int
		checkDetails           bool
	}{
		// Test of update failure: basic user
		{models.RecvTaskLog{
			LogMessage: "This is a note",
			Task:       createdTasks[0],
		}, testConnection.accounts.user.token, http.StatusForbidden, false},
		// Update should be allowed: admin
		{models.RecvTaskLog{
			LogMessage: "This is a note",
			Task:       createdTasks[0],
		}, testConnection.accounts.admin.token, http.StatusOK, true},
		// Update should be disallowed due to being too short
		{models.RecvTaskLog{
			LogMessage: "Th",
			Task:       createdTasks[0],
		}, testConnection.accounts.admin.token, http.StatusBadRequest, false},
		// User should be forbidden before validating rather than Bad Request
		{models.RecvTaskLog{
			LogMessage: "Tx",
			Task:       createdTasks[0],
		}, testConnection.accounts.user.token, http.StatusForbidden, false},
	}

	// Create a request url with an "id" URL parameter
	requestUrl := fmt.Sprintf("/api/task-logs/%v", createdTaskLogs[0].ID)

	// Iterate through update tests
	for _, v := range updateTests {
		// Make new request with feature update in body
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
		var body db.TaskLog
		json.Unmarshal(rr.Body.Bytes(), &body)
		// Check response expected vs received
		if status := rr.Code; status != v.expectedResponseStatus {
			t.Errorf("Update test: got %v want %v. \nBody: %v", status, v.expectedResponseStatus, body)

		}

		// If need to check details
		if v.checkDetails == true {
			// Get task details from database
			var expected db.TaskLog
			findResult := testConnection.dbClient.Find(&expected, createdTaskLogs[0].ID)
			if findResult.Error != nil {
				t.Errorf("Error finding updated task log: %v", findResult.Error)
			}

			// Check task log details using updated object
			checkTaskLogDetails(&body, &expected, t, false)
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
		// Make new request with task log update in body
		req, err := http.NewRequest("PUT", fmt.Sprint("/api/property-logs/"+v.urlExtension), buildReqBody(&db.PropertyLog{
			LogMessage: "Gustav",
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
			t.Errorf("Fail update test: got %v want %v.",
				status, v.expectedResponseStatus)
		}
	}

	// Cleanup
	// Delete the created task logs
	deleteResult := testConnection.dbClient.Delete(createdTaskLogs)
	if deleteResult.Error != nil {
		t.Fatalf("Couldn't clean up seeded task logs: %v", deleteResult.Error)
	}
	// Delete the created tasks
	deleteResult = testConnection.dbClient.Delete(createdTasks)
	if deleteResult.Error != nil {
		t.Fatalf("Couldn't clean up seeded tasks: %v", deleteResult.Error)
	}
}

func TestTaskLogController_Create(t *testing.T) {
	// Setup
	//
	// Create task
	taskToCreate := &db.Task{
		TaskName: "Swell Task",
		Type:     "Inspection",
		Notes:    "The bylaws are unclear on this",
	}
	createdTasks := []db.Task{*taskToCreate}
	createResult := testConnection.dbClient.Create(createdTasks)
	if createResult.Error != nil {
		t.Fatal("Failed to create task for task log find all test: ", createResult.Error)
	}

	var createTests = []struct {
		data                   models.RecvTaskLog
		expectedResponseStatus int
		tokenToUse             string
	}{
		// Should fail due to user role status of basic
		{models.RecvTaskLog{
			LogMessage: "This is a note",
			Task:       createdTasks[0],
		}, http.StatusForbidden, testConnection.accounts.user.token},
		// Should pass as user is admin
		{models.RecvTaskLog{
			LogMessage: "This is a note",
			Task:       createdTasks[0],
		}, http.StatusCreated, testConnection.accounts.admin.token},
		// Create should be disallowed due to being too short
		{models.RecvTaskLog{
			LogMessage: "Tx",
			Task:       createdTasks[0],
		}, http.StatusBadRequest, testConnection.accounts.admin.token},
		// Create should be disallowed due to being too long
		{models.RecvTaskLog{
			LogMessage: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Suspendisse vulputate, nunc sit amet efficitur bibendum, sapien odio auctor nisi, a interdum magna nisl ac purus. Fusce condimentum malesuada mi at eleifend. Sed laoreet varius risus, id mattis libero tristique nec. Sed eget malesuada magna. Morbi feugiat sapien euismod neque commodo suscipit. Vivamus vehicula euismod dui, id imperdiet elit lacinia non. Integer hendrerit, enim ac gravida malesuada, dolor leo dictum purus, nec bibendum velit est vel nulla. Nulla sagittis nulla non elit imperdiet convallis. Sed bibendum sollicitudin nunc, vel facilisis nulla convallis a. Nunc id ex feugiat, finibus magna sit amet, ultricies lacus.",
			Task:       createdTasks[0],
		}, http.StatusBadRequest, testConnection.accounts.admin.token},
	}

	// Create a request url with an "id" URL parameter
	requestUrl := "/api/task-logs"

	for _, v := range createTests {
		// Make new request with task log creation in body
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
			t.Errorf("Task log create test (%v): got %v want %v. \nBody: %v\n", v.data.LogMessage,
				status, v.expectedResponseStatus, rr.Body.String())
		}

		// Init body for response extraction
		var body db.TaskLog
		var foundTaskLog db.TaskLog
		// Grab ID from response body
		json.Unmarshal(rr.Body.Bytes(), &body)

		// Find the created task log (to obtain full data with ID)
		testConnection.dbClient.Find(foundTaskLog, uint(body.ID))

		// Compare found details with those found in returned body
		checkTaskLogDetails(&body, &foundTaskLog, t, true)

		// If the task log was created successfully, check that it's deleted after test
		if v.expectedResponseStatus == http.StatusCreated {
			// Cleanup
			//
			// Delete the created task logs
			deleteResult := testConnection.dbClient.Delete(&db.TaskLog{}, uint(body.ID))
			if deleteResult.Error != nil {
				t.Fatalf("Couldn't clean up created task logs: %v", deleteResult.Error)
			}
		}
	}

	// Overall test cleanup
	// Delete the created properties
	deleteResult := testConnection.dbClient.Delete(createdTasks)
	if deleteResult.Error != nil {
		t.Fatalf("Couldn't clean up seeded tasks: %v", deleteResult.Error)
	}
}

// Check the task log message details
func checkTaskLogDetails(actual *db.TaskLog, expected *db.TaskLog, t *testing.T, checkId bool) {

	// Only check ID if parameter checkId is true
	if checkId == true {
		// Verify that the actual task log message id matches the created task log message id
		if actual.ID != expected.ID {
			t.Errorf("found feature has incorrect ID: expected %d, got %d", expected.ID, actual.ID)
		}
	}

	// Verify that the actual details matches the expected details
	if actual.LogMessage != expected.LogMessage {
		t.Errorf("found task log message has incorrect log message: expected %s, got %s", expected.LogMessage, actual.LogMessage)
	}
	if actual.TaskID != expected.TaskID {
		t.Errorf("found task log message has incorrect Task ID: expected %d, got %d", expected.TaskID, actual.TaskID)
	}
	if actual.UserID != expected.UserID {
		t.Errorf("found task log message has incorrect user ID: expected %d, got %d", expected.UserID, actual.UserID)
	}
}
