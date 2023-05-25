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

func TestMaintenanceController_FindAll(t *testing.T) {
	// Test setup
	// Create property
	propertyToCreate := &db.Property{
		Property_Name:    "maintenanceProperty1",
		Postcode:         80361,
		Suburb:           "Test Suburb",
		City:             "Test City",
		Street_Address_1: "Test Street Address 1",
		Bedrooms:         3,
		Bathrooms:        2,
		Description:      "Test Description",
		Managed:          true,
	}
	createdProperties := []db.Property{*propertyToCreate}
	createResult := testConnection.dbClient.Create(createdProperties)
	if createResult.Error != nil {
		t.Fatal("Failed to create properties for maintenance find all test: ", createResult.Error)
	}
	// Create task
	taskToCreate1 := &db.Task{
		TaskName: "Test Task",
		Type:     "Maintenance",
		Notes:    "Yohoo",
	}
	taskToCreate2 := &db.Task{
		TaskName: "Test Task 2",
		Type:     "Maintenance",
		Notes:    "Yohoo",
	}
	createdTasks := []db.Task{*taskToCreate1, *taskToCreate2}
	createResult = testConnection.dbClient.Create(createdTasks)
	if createResult.Error != nil {
		t.Fatal("Failed to create tasks for maintenance find all test: ", createResult.Error)
	}
	// Create maintenance requests
	requestToCreate1 := &db.MaintenanceRequest{
		Scale:          "Urgent",
		WorkDefinition: "Repair",
		Type:           "Electrical",
		Notes:          "Marketing team absolutely sucks",
		Property:       createdProperties[0],
		TaskID:         createdTasks[0].ID,
	}
	requestToCreate2 := &db.MaintenanceRequest{
		Scale:          "Urgent",
		WorkDefinition: "Repair",
		Type:           "Electrical",
		Notes:          "Marketing team absolutely sucks",
		Property:       createdProperties[0],
		TaskID:         createdTasks[1].ID,
	}
	createdRequests := []db.MaintenanceRequest{*requestToCreate1, *requestToCreate2}
	// Create task logs in db
	createResult = testConnection.dbClient.Create(&createdRequests)
	if createResult.Error != nil {
		t.Fatal("Failed to create maintenance requests for maintenance find all test: ", createResult.Error)
	}

	// Create a new request
	req, err := http.NewRequest("GET", "/api/maintenance?limit=10&offset=0&order=", nil)
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
	var body []db.MaintenanceRequest
	json.Unmarshal(rr.Body.Bytes(), &body)

	// Check length of array (should be two with seeded assets)
	if len(body) != len(createdRequests) {
		t.Errorf("Array length check in findAll failed: expected %d, got %d", len(createdRequests), len(body))
	}

	// Iterate through array received
	for _, actualRequest := range body {
		// Iterate through created transactions to determine a match
		for _, createdRequest := range createdRequests {
			// If match found
			if actualRequest.ID == createdRequest.ID {
				// Check the details of the transaction match
				// t.Fatal("Maintenance request details check failed:", createdRequests)
				checkMaintenanceRequestDetails(&actualRequest, &createdRequest, t, false)
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
		request := fmt.Sprintf("/api/maintenance?limit=%v&offset=%v&order=%v", v.limit, v.offset, v.order)
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
	deleteResult := testConnection.dbClient.Delete(createdRequests)
	if deleteResult.Error != nil {
		t.Fatalf("Couldn't clean up seeded maintenance requests: %v", deleteResult.Error)
	}
	deleteResult = testConnection.dbClient.Delete(createdTasks)
	if deleteResult.Error != nil {
		t.Fatalf("Couldn't clean up seeded tasks: %v", deleteResult.Error)
	}
	deleteResult = testConnection.dbClient.Delete(createdProperties)
	if deleteResult.Error != nil {
		t.Fatalf("Couldn't clean up seeded props: %v", deleteResult.Error)
	}
}

func TestMaintenanceController_Find(t *testing.T) {
	// Test setup
	// Create property
	propertyToCreate := &db.Property{
		Property_Name:    "Maintenance Property",
		Postcode:         80361,
		Suburb:           "Test Suburb",
		City:             "Test City",
		Street_Address_1: "Test Street Address 1",
		Bedrooms:         3,
		Bathrooms:        2,
		Description:      "Test Description",
		Managed:          true,
	}
	createdProperties := []db.Property{*propertyToCreate}
	createResult := testConnection.dbClient.Create(createdProperties)
	if createResult.Error != nil {
		t.Fatal("Failed to create properties for test: ", createResult.Error)
	}
	// Create task
	taskToCreate1 := &db.Task{
		TaskName: "Test Task",
		Type:     "Maintenance",
		Notes:    "Yohoo",
	}
	createdTasks := []db.Task{*taskToCreate1}
	createResult = testConnection.dbClient.Create(createdTasks)
	if createResult.Error != nil {
		t.Fatal("Failed to create tasks for test: ", createResult.Error)
	}
	// Create maintenance request
	requestToCreate1 := &db.MaintenanceRequest{
		Scale:          "Urgent",
		WorkDefinition: "Repair",
		Type:           "Electrical",
		Notes:          "Marketing team absolutely sucks",
		Property:       createdProperties[0],
		TaskID:         createdTasks[0].ID,
	}
	createdRequests := []db.MaintenanceRequest{*requestToCreate1}
	// Create task logs in db
	createResult = testConnection.dbClient.Create(createdRequests)
	if createResult.Error != nil {
		t.Fatal("Failed to create requests for test: ", createResult.Error)
	}

	// Create a request with an "id" URL parameter
	requestUrl := fmt.Sprintf("/api/maintenance/%v", createdRequests[0].ID)
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
	var body db.MaintenanceRequest
	json.Unmarshal(rr.Body.Bytes(), &body)
	checkMaintenanceRequestDetails(&body, &createdRequests[0], t, true)

	// Cleanup
	// Delete the created fixtures
	deleteResult := testConnection.dbClient.Delete(createdRequests)
	if deleteResult.Error != nil {
		t.Fatalf("Couldn't clean up seeded maintenance requests: %v", deleteResult.Error)
	}
	deleteResult = testConnection.dbClient.Delete(createdTasks)
	if deleteResult.Error != nil {
		t.Fatalf("Couldn't clean up seeded tasks: %v", deleteResult.Error)
	}
	deleteResult = testConnection.dbClient.Delete(createdProperties)
	if deleteResult.Error != nil {
		t.Fatalf("Couldn't clean up seeded props: %v", deleteResult.Error)
	}
}

func TestMaintenanceController_Delete(t *testing.T) {
	// Test setup
	// Create property
	propertyToCreate := &db.Property{
		Property_Name:    "Maintenance Property 2",
		Postcode:         80361,
		Suburb:           "Test Suburb",
		City:             "Test City",
		Street_Address_1: "Test Street Address 1",
		Bedrooms:         3,
		Bathrooms:        2,
		Description:      "Test Description",
		Managed:          true,
	}
	createdProperties := []db.Property{*propertyToCreate}
	createResult := testConnection.dbClient.Create(createdProperties)
	if createResult.Error != nil {
		t.Fatal("Failed to create properties for test: ", createResult.Error)
	}
	// Create task
	taskToCreate1 := &db.Task{
		TaskName: "Test Task",
		Type:     "Maintenance",
		Notes:    "Yohoo",
	}
	createdTasks := []db.Task{*taskToCreate1}
	createResult = testConnection.dbClient.Create(createdTasks)
	if createResult.Error != nil {
		t.Fatal("Failed to create tasks for test: ", createResult.Error)
	}
	// Create transactions
	requestToCreate1 := &db.MaintenanceRequest{
		Scale:          "Urgent",
		WorkDefinition: "Repair",
		Type:           "Electrical",
		Notes:          "Marketing team absolutely sucks",
		Property:       createdProperties[0],
		TaskID:         createdTasks[0].ID,
	}
	createdRequests := []db.MaintenanceRequest{*requestToCreate1}
	// Create task logs in db
	createResult = testConnection.dbClient.Create(createdRequests)
	if createResult.Error != nil {
		t.Fatal("Failed to create maintenance requests for test: ", createResult.Error)
	}

	// Create a request with an "id" URL parameter
	requestUrl := fmt.Sprintf("/api/maintenance/%v", createdRequests[0].ID)
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
		{testName: "Transaction basic user delete test", tokenToUse: testConnection.accounts.user.token, expectedResponseStatus: http.StatusForbidden},
		// Must be last
		// Tests of deletion success using admin privileges
		{testName: "Transaction admin delete test", tokenToUse: testConnection.accounts.admin.token, expectedResponseStatus: http.StatusOK},
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
	// Delete the created fixtures
	deleteResult := testConnection.dbClient.Delete(createdRequests)
	if deleteResult.Error != nil {
		t.Fatalf("Couldn't clean up seeded maintenance requests: %v", deleteResult.Error)
	}
	deleteResult = testConnection.dbClient.Delete(createdTasks)
	if deleteResult.Error != nil {
		t.Fatalf("Couldn't clean up seeded tasks: %v", deleteResult.Error)
	}
	deleteResult = testConnection.dbClient.Delete(createdProperties)
	if deleteResult.Error != nil {
		t.Fatalf("Couldn't clean up seeded props: %v", deleteResult.Error)
	}
}

func TestMaintenanceController_Update(t *testing.T) {
	// Test setup
	// Create property
	propertyToCreate := &db.Property{
		Property_Name:    "Maintenance Property 4",
		Postcode:         80361,
		Suburb:           "Test Suburb",
		City:             "Test City",
		Street_Address_1: "Test Street Address 1",
		Bedrooms:         3,
		Bathrooms:        2,
		Description:      "Test Description",
		Managed:          true,
	}
	createdProperties := []db.Property{*propertyToCreate}
	createResult := testConnection.dbClient.Create(createdProperties)
	if createResult.Error != nil {
		t.Fatal("Failed to create properties for test: ", createResult.Error)
	}
	// Create task
	taskToCreate1 := &db.Task{
		TaskName: "Test Task",
		Type:     "Maintenance",
		Notes:    "Yohoo",
	}
	createdTasks := []db.Task{*taskToCreate1}
	createResult = testConnection.dbClient.Create(createdTasks)
	if createResult.Error != nil {
		t.Fatal("Failed to create tasks for test: ", createResult.Error)
	}
	// Create maintenance request
	requestToCreate1 := &db.MaintenanceRequest{
		Scale:          "Urgent",
		WorkDefinition: "Repair",
		Type:           "Electrical",
		Notes:          "Marketing team absolutely sucks",
		Property:       createdProperties[0],
		TaskID:         createdTasks[0].ID,
	}
	createdRequests := []db.MaintenanceRequest{*requestToCreate1}
	createResult = testConnection.dbClient.Create(&createdRequests)
	if createResult.Error != nil {
		t.Fatal("Failed to create maintenance request for test: ", createResult.Error)
	}

	// Build test array
	var updateTests = []struct {
		data                   models.UpdateMaintenanceRequest
		tokenToUse             string
		expectedResponseStatus int
		checkDetails           bool
		testName               string
	}{
		// Test of update failure: basic user
		{models.UpdateMaintenanceRequest{
			Type: "HVAC",
		}, testConnection.accounts.user.token, http.StatusForbidden, false, "Basic user update test"},
		// Update should be allowed: admin
		{models.UpdateMaintenanceRequest{
			Type: "HVAC",
		}, testConnection.accounts.admin.token, http.StatusOK, true, "Admin update test"},
		// Update should be disallowed due to being invalid value for type
		{models.UpdateMaintenanceRequest{
			Type: "Sales",
		}, testConnection.accounts.admin.token, http.StatusBadRequest, false, "Invalid type update test"},
		// Update should be disallowed due to being invalid value for work definition
		{models.UpdateMaintenanceRequest{
			WorkDefinition: "Solitude",
		}, testConnection.accounts.admin.token, http.StatusBadRequest, false, "Invalid work definition update test"},
		// Update should be disallowed due to being invalid value for scale
		{models.UpdateMaintenanceRequest{
			Scale: "Beyond",
		}, testConnection.accounts.admin.token, http.StatusBadRequest, false, "Invalid scale update test"},
		// User should be forbidden before validating rather than Bad Request
		{models.UpdateMaintenanceRequest{
			Type: "Squalor",
		}, testConnection.accounts.user.token, http.StatusForbidden, false, "Invalid type update when basic user test"},
	}

	// t.Errorf("Created Requests: %v", createdRequests)
	// Create a request url with an "id" URL parameter
	requestUrl := fmt.Sprintf("/api/maintenance/%v", createdRequests[0].ID)
	fmt.Printf("Request URL: %v", requestUrl)

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
		var body db.MaintenanceRequest
		json.Unmarshal(rr.Body.Bytes(), &body)

		// Check response expected vs received
		if status := rr.Code; status != v.expectedResponseStatus {
			t.Errorf("Update test: got %v want %v. \nBody: %v", status, v.expectedResponseStatus, body.Scale)
		}

		// If need to check details
		if v.checkDetails == true {
			// Get maintenance request details from database
			var expected db.MaintenanceRequest
			findResult := testConnection.dbClient.Find(&expected, createdRequests[0].ID)
			if findResult.Error != nil {
				t.Errorf("Error finding updated maintenance request: %v", findResult.Error)
			}

			// Check task log details using updated object
			checkMaintenanceRequestDetails(&body, &expected, t, false)
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
		req, err := http.NewRequest("PUT", fmt.Sprint("/api/maintenance/"+v.urlExtension), buildReqBody(&db.MaintenanceRequest{
			Type: "HVAC",
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
	// Delete the created fixtures
	deleteResult := testConnection.dbClient.Delete(createdRequests)
	if deleteResult.Error != nil {
		t.Fatalf("Couldn't clean up seeded maintenance requests: %v", deleteResult.Error)
	}
	deleteResult = testConnection.dbClient.Delete(createdTasks)
	if deleteResult.Error != nil {
		t.Fatalf("Couldn't clean up seeded tasks: %v", deleteResult.Error)
	}
	deleteResult = testConnection.dbClient.Delete(createdProperties)
	if deleteResult.Error != nil {
		t.Fatalf("Couldn't clean up seeded props: %v", deleteResult.Error)
	}
}

func TestMaintenanceController_Create(t *testing.T) {
	// Setup
	//
	// Create property
	propertyToCreate := &db.Property{
		Property_Name:    "Maintenance Property 5",
		Postcode:         80361,
		Suburb:           "Test Suburb",
		City:             "Test City",
		Street_Address_1: "Test Street Address 1",
		Bedrooms:         3,
		Bathrooms:        2,
		Description:      "Test Description",
		Managed:          true,
	}
	createdProperties := []db.Property{*propertyToCreate}
	createResult := testConnection.dbClient.Create(createdProperties)
	if createResult.Error != nil {
		t.Fatal("Failed to create properties for test: ", createResult.Error)
	}

	var createTests = []struct {
		data                   models.CreateMaintenanceRequest
		expectedResponseStatus int
		tokenToUse             string
		testName               string
	}{
		// Should fail due to user role status of basic
		{models.CreateMaintenanceRequest{
			Scale:          "Urgent",
			WorkDefinition: "Repair",
			Type:           "Electrical",
			Notes:          "Marketing team absolutely sucks",
			Property:       createdProperties[0],
		}, http.StatusForbidden, testConnection.accounts.user.token, "basic user create"},
		// Should pass as user is admin
		{models.CreateMaintenanceRequest{
			Scale:          "Urgent",
			WorkDefinition: "Repair",
			Type:           "Electrical",
			Property:       createdProperties[0],
			Notes:          "Ridiculous things",
		}, http.StatusCreated, testConnection.accounts.admin.token, "admin create"},
		// Create should be disallowed due to invalid scale value
		{models.CreateMaintenanceRequest{
			Scale:          "MadeUpScale",
			WorkDefinition: "Repair",
			Type:           "Electrical",
			Notes:          "Waser team absolutely sucks",
			Property:       createdProperties[0],
		}, http.StatusBadRequest, testConnection.accounts.admin.token, "invalid scale create"},
		// Create should be disallowed due to invalid type value
		{models.CreateMaintenanceRequest{
			Scale:          "Urgent",
			WorkDefinition: "Repair",
			Type:           "Trains",
			Notes:          "Drover team absolutely sucks",
			Property:       createdProperties[0],
		}, http.StatusBadRequest, testConnection.accounts.admin.token, "invalid type create"},
		// Create should be disallowed due to invalid work definition type value
		{models.CreateMaintenanceRequest{
			Scale:          "Urgent",
			WorkDefinition: "Regulate",
			Type:           "Electrical",
			Notes:          "Marketing team absolutely sucks",
			Property:       createdProperties[0],
		}, http.StatusBadRequest, testConnection.accounts.admin.token, "invalid work definition create"},
		// Create should be disallowed due to notes being too long
		{models.CreateMaintenanceRequest{
			Scale:          "Urgent",
			WorkDefinition: "Repair",
			Type:           "Electrical",
			Notes:          "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Suspendisse vulputate, nunc sit amet efficitur bibendum, sapien odio auctor nisi, a interdum magna nisl ac purus. Fusce condimentum malesuada mi at eleifend. Sed laoreet varius risus, id mattis libero tristique nec. Sed eget malesuada magna. Morbi feugiat sapien euismod neque commodo suscipit. Vivamus vehicula euismod dui, id imperdiet elit lacinia non. Integer hendrerit, enim ac gravida malesuada, dolor leo dictum purus, nec bibendum velit est vel nulla. Nulla sagittis nulla non elit imperdiet convallis. Sed bibendum sollicitudin nunc, vel facilisis nulla convallis a. Nunc id ex feugiat, finibus magna sit amet, ultricies lacus.",
			Property:       createdProperties[0],
		}, http.StatusBadRequest, testConnection.accounts.admin.token, "invalid notes create"},
	}

	// Create a request url with an "id" URL parameter
	requestUrl := "/api/maintenance"

	for _, v := range createTests {
		// Each create test setup
		// Create task (for each test)
		taskToCreate1 := &db.Task{
			TaskName: "Test Task",
			Type:     "Maintenance",
			Notes:    "Yohoo",
		}
		createdTasks := []db.Task{*taskToCreate1}
		createResult = testConnection.dbClient.Create(createdTasks)
		if createResult.Error != nil {
			t.Fatal("Failed to create tasks for test: ", createResult.Error)
		}
		// Update v.data to include newly created task
		v.data.Task = createdTasks[0]

		// Make new request with maintenance request creation in body
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
			t.Errorf("response test (%v): got %v want %v. \nBody: %v\n", v.testName,
				status, v.expectedResponseStatus, rr.Body.String())
		}

		// Init body for response extraction
		var body db.MaintenanceRequest
		var foundRequest db.MaintenanceRequest
		// Grab ID from response body
		json.Unmarshal(rr.Body.Bytes(), &body)

		// Find the created transaction (to obtain full data with ID)
		testConnection.dbClient.Find(foundRequest, uint(body.ID))

		// Compare found details with those found in returned body
		checkMaintenanceRequestDetails(&body, &foundRequest, t, true)

		// If the task log was created successfully, check that it's deleted after test
		if v.expectedResponseStatus == http.StatusCreated {
			// Cleanup
			//
			// Delete the created task logs
			deleteMainResult := testConnection.dbClient.Delete(&db.MaintenanceRequest{}, uint(body.ID))
			if deleteMainResult.Error != nil {
				t.Fatalf("Couldn't clean up created maintenance requests: %v", deleteMainResult.Error)
			}
		}
		// Delete the created task regardless of response result
		deleteTaskResult := testConnection.dbClient.Delete(createdTasks)
		if deleteTaskResult.Error != nil {
			t.Fatalf("Couldn't clean up seeded tasks: %v", deleteTaskResult.Error)
		}
	}

	// cleanup

	deleteResult := testConnection.dbClient.Delete(createdProperties)
	if deleteResult.Error != nil {
		t.Fatalf("Couldn't clean up seeded props: %v", deleteResult.Error)
	}
}

// Check the transaction details
func checkMaintenanceRequestDetails(actual *db.MaintenanceRequest, expected *db.MaintenanceRequest, t *testing.T, checkId bool) {

	// Only check ID if parameter checkId is true
	if checkId == true {
		// Verify that the actual transaction id matches the created transaction id
		if actual.ID != expected.ID {
			t.Errorf("found transaction has incorrect ID: expected %d, got %d", expected.ID, actual.ID)
		}
	}

	// Verify that the actual details matches the expected details
	if actual.WorkDefinition != expected.WorkDefinition {
		t.Errorf("found maintenance request has incorrect work definition: expected %s, got %s", expected.WorkDefinition, actual.WorkDefinition)
	}
	if actual.Notes != expected.Notes {
		t.Errorf("found maintenance request has incorrect notes: expected %s, got %s", expected.Notes, actual.Notes)
	}
	if actual.Scale != expected.Scale {
		t.Errorf("found maintenance request has incorrect scale: expected %s, got %s", expected.Scale, actual.Scale)
	}
	if actual.Tax != expected.Tax {
		t.Errorf("found maintenance request has incorrect tax: expected %v, got %v", expected.Tax, actual.Tax)
	}
	if actual.TotalCost != expected.TotalCost {
		t.Errorf("found maintenance request has incorrect total cost: expected %v, got %v", expected.TotalCost, actual.TotalCost)
	}
	if actual.Type != expected.Type {
		t.Errorf("found maintenance request has incorrect type: expected %s, got %s", expected.Type, actual.Type)
	}

	// Relationships
	if actual.PropertyID != expected.PropertyID {
		t.Errorf("found maintenance request has incorrect property ID: expected %d, got %d", expected.Property.ID, actual.Property.ID)
	}
	if actual.TaskID != expected.TaskID {
		t.Errorf("found maintenance request has incorrect task ID: expected %d, got %d", expected.TaskID, actual.TaskID)
	}

}
