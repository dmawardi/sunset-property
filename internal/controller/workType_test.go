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

func TestWorkTypeController_FindAll(t *testing.T) {
	// Test setup
	typeToCreate1 := &db.WorkType{
		Name: "Test Type 1",
	}
	typeToCreate2 := &db.WorkType{
		Name: "Test Type 2",
	}
	createdWorkTypes := []db.WorkType{*typeToCreate1, *typeToCreate2}
	createResult := testConnection.dbClient.Create(createdWorkTypes)
	if createResult.Error != nil {
		t.Fatal("Failed to create work types for test: ", createResult.Error)
	}

	// Create a new request
	req, err := http.NewRequest("GET", "/api/work-types?limit=10&offset=0&order=", nil)
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
	var body []db.WorkType
	json.Unmarshal(rr.Body.Bytes(), &body)

	// Check length of array (should be two with seeded assets)
	if len(body) != len(createdWorkTypes) {
		t.Errorf("Array length check in findAll failed: expected %d, got %d", len(createdWorkTypes), len(body))
	}

	// Iterate through array received
	for _, actualWorkType := range body {
		// Iterate through prior created items to determine a match
		for _, createdWorkType := range createdWorkTypes {
			// If match found
			if actualWorkType.ID == createdWorkType.ID {
				// Check the details of the transaction match
				checkWorkTypeDetails(&actualWorkType, &createdWorkType, t, false)
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
		request := fmt.Sprintf("/api/work-types?limit=%v&offset=%v&order=%v", v.limit, v.offset, v.order)
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
	deleteResult := testConnection.dbClient.Delete(createdWorkTypes)
	if deleteResult.Error != nil {
		t.Fatalf("Couldn't clean up seeded work types: %v", deleteResult.Error)
	}
}

func TestWorkTypeController_Find(t *testing.T) {
	// Test setup
	typeToCreate1 := &db.WorkType{
		Name: "Test Type 3",
	}
	createdWorkTypes := []db.WorkType{*typeToCreate1}
	createResult := testConnection.dbClient.Create(createdWorkTypes)
	if createResult.Error != nil {
		t.Fatal("Failed to create work types for test: ", createResult.Error)
	}

	// Create a request with an "id" URL parameter
	requestUrl := fmt.Sprintf("/api/work-types/%v", createdWorkTypes[0].ID)
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
	var body db.WorkType
	json.Unmarshal(rr.Body.Bytes(), &body)
	checkWorkTypeDetails(&body, &createdWorkTypes[0], t, true)

	// Cleanup
	// Delete the created fixtures
	deleteResult := testConnection.dbClient.Delete(createdWorkTypes)
	if deleteResult.Error != nil {
		t.Fatalf("Couldn't clean up seeded work types: %v", deleteResult.Error)
	}
}

func TestWorkTypeController_Delete(t *testing.T) {
	// Test setup
	typeToCreate1 := &db.WorkType{
		Name: "Test Type 4",
	}
	createdWorkTypes := []db.WorkType{*typeToCreate1}
	createResult := testConnection.dbClient.Create(createdWorkTypes)
	if createResult.Error != nil {
		t.Fatal("Failed to create work types for test: ", createResult.Error)
	}

	// Create a request with an "id" URL parameter
	requestUrl := fmt.Sprintf("/api/work-types/%v", createdWorkTypes[0].ID)
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
	deleteResult := testConnection.dbClient.Delete(createdWorkTypes)
	if deleteResult.Error != nil {
		t.Fatalf("Couldn't clean up seeded work types: %v", deleteResult.Error)
	}
}

func TestWorkTypeController_Update(t *testing.T) {
	// Test setup
	typeToCreate1 := &db.WorkType{
		Name: "Test Type 5",
	}
	createdWorkTypes := []db.WorkType{*typeToCreate1}
	createResult := testConnection.dbClient.Create(createdWorkTypes)
	if createResult.Error != nil {
		t.Fatal("Failed to create work types for test: ", createResult.Error)
	}

	// Build test array
	var updateTests = []struct {
		data                   models.UpdateWorkType
		tokenToUse             string
		expectedResponseStatus int
		checkDetails           bool
		testName               string
	}{
		// Test of update failure: basic user
		{models.UpdateWorkType{
			Name: "HVAC",
		}, testConnection.accounts.user.token, http.StatusForbidden, false, "Basic user update test"},
		// Update should be allowed: admin
		{models.UpdateWorkType{
			Name: "HVAC",
		}, testConnection.accounts.admin.token, http.StatusOK, true, "Admin update test"},
		// Update should be disallowed due to being over valid length
		{models.UpdateWorkType{
			Name: "Loremipsum dolorsitametco dolorsitametco",
		}, testConnection.accounts.admin.token, http.StatusBadRequest, false, "Invalid type update test"},
		// User should be forbidden before validating rather than Bad Request
		{models.UpdateWorkType{
			Name: "HVAC",
		}, testConnection.accounts.user.token, http.StatusForbidden, false, "Invalid type update when basic user test"},
	}

	// Create a request url with an "id" URL parameter
	requestUrl := fmt.Sprintf("/api/work-types/%v", createdWorkTypes[0].ID)
	fmt.Printf("Request URL: %v", requestUrl)

	// Iterate through update tests
	for _, v := range updateTests {
		// Make new request with update in body
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
		var body db.WorkType
		json.Unmarshal(rr.Body.Bytes(), &body)

		// Check response expected vs received
		if status := rr.Code; status != v.expectedResponseStatus {
			t.Errorf("Update test: got %v want %v. \nBody: %v", status, v.expectedResponseStatus, rr.Body.String())
		}

		// If need to check details
		if v.checkDetails == true {
			// Get work type details from database
			var expected db.WorkType
			findResult := testConnection.dbClient.Find(&expected, createdWorkTypes[0].ID)
			if findResult.Error != nil {
				t.Errorf("Error finding updated work type: %v", findResult.Error)
			}

			// Check task log details using updated object
			checkWorkTypeDetails(&body, &expected, t, false)
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
		req, err := http.NewRequest("PUT", fmt.Sprint("/api/work-types/"+v.urlExtension), buildReqBody(&db.MaintenanceRequest{
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
	deleteResult := testConnection.dbClient.Delete(createdWorkTypes)
	if deleteResult.Error != nil {
		t.Fatalf("Couldn't clean up seeded work types: %v", deleteResult.Error)
	}
}

func TestWorkTypeController_Create(t *testing.T) {
	var createTests = []struct {
		data                   models.CreateWorkType
		expectedResponseStatus int
		tokenToUse             string
		testName               string
	}{
		// Should fail due to user role status of basic
		{models.CreateWorkType{
			Name: "Plumbing",
		}, http.StatusForbidden, testConnection.accounts.user.token, "basic user create"},
		// Should pass as user is admin
		{models.CreateWorkType{
			Name: "Plumbing",
		}, http.StatusCreated, testConnection.accounts.admin.token, "admin create"},
		// Create should be disallowed due to name being too long
		{models.CreateWorkType{
			Name: "Loremipsum dolorsitametco dolorsitametco",
		}, http.StatusBadRequest, testConnection.accounts.admin.token, "invalid name length create"},
	}

	// Create a request url
	requestUrl := "/api/work-types"

	for _, v := range createTests {

		// Make new request with work type creation in body
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
		var body db.WorkType
		var foundWorkType db.WorkType
		// Grab ID from response body
		json.Unmarshal(rr.Body.Bytes(), &body)

		// Find the created transaction (to obtain full data with ID)
		testConnection.dbClient.Find(foundWorkType, uint(body.ID))

		// Compare found details with those found in returned body
		checkWorkTypeDetails(&body, &foundWorkType, t, true)

		// If the task log was created successfully, check that it's deleted after test
		if v.expectedResponseStatus == http.StatusCreated {
			// Cleanup
			//
			// Delete the created work type
			deleteMainResult := testConnection.dbClient.Delete(&db.WorkType{}, uint(body.ID))
			if deleteMainResult.Error != nil {
				t.Fatalf("Couldn't clean up created work types: %v", deleteMainResult.Error)
			}
		}

	}
}

// Check the work type details
func checkWorkTypeDetails(actual *db.WorkType, expected *db.WorkType, t *testing.T, checkId bool) {
	// Only check ID if parameter checkId is true
	if checkId == true {
		// Verify that the actual id matches the created id
		if actual.ID != expected.ID {
			t.Errorf("found work type has incorrect ID: expected %d, got %d", expected.ID, actual.ID)
		}
	}

	// Check details
	if actual.Name != expected.Name {
		t.Errorf("found work type has incorrect name: expected %s, got %s", expected.Name, actual.Name)
	}
}
